// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package runtimeevent

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/server/healthz"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
)

type Generator interface {
	Source
	healthz.HealthChecker
	manager.Runnable
}

type event struct {
	UpdateResources *UpdateResourcesEvent
}

type handler struct {
	Handler Handler
}

type generator struct {
	mu sync.RWMutex

	eventChannel chan *event

	handlers sets.Set[*handler]

	// relistPeriod is the period for relisting.
	relistPeriod time.Duration
	// relistThreshold is the maximum threshold between two relists to become unhealthy.
	relistThreshold time.Duration
	// relistTime is the last time a relist happened.
	relistTime atomic.Pointer[time.Time]
	// firstListTime is the first time a relist happened.
	firstListTime time.Time

	resources *sri.RuntimeResources

	getResources func(ctx context.Context) (*sri.RuntimeResources, error)
}

type GeneratorOptions struct {
	ChannelCapacity int
	RelistPeriod    time.Duration
	RelistThreshold time.Duration
}

func setGeneratorOptionsDefaults(o *GeneratorOptions) {
	if o.ChannelCapacity == 0 {
		o.ChannelCapacity = 1024
	}
	if o.RelistPeriod <= 0 {
		o.RelistPeriod = 1 * time.Second
	}
	if o.RelistThreshold <= 0 {
		o.RelistThreshold = 3 * time.Minute
	}
}

func NewGenerator(getResources func(ctx context.Context) (*sri.RuntimeResources, error), opts GeneratorOptions) Generator {
	setGeneratorOptionsDefaults(&opts)

	return &generator{
		eventChannel:    make(chan *event, opts.ChannelCapacity),
		relistPeriod:    opts.RelistPeriod,
		relistThreshold: opts.RelistThreshold,
		relistTime:      atomic.Pointer[time.Time]{},
		firstListTime:   time.Time{},
		getResources:    getResources,
		handlers:        sets.New[*handler](),
	}
}

func (g *generator) Name() string {
	return "runtime-event-generator"
}

func (g *generator) Check(_ *http.Request) error {
	relistTime := g.relistTime.Load()
	if relistTime == nil {
		return fmt.Errorf("reg did not relist yet")
	}

	elapsed := time.Since(*relistTime)
	if elapsed > g.relistThreshold {
		return fmt.Errorf("reg was last seen active %v ago, threshold is %v", elapsed, g.relistThreshold)
	}
	return nil
}

func (g *generator) readHandlers() []*handler {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.handlers.UnsortedList()
}

func (g *generator) Start(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx).WithName("event-generator")

	go func() {
		for evt := range g.eventChannel {
			handlers := g.readHandlers()

			for _, handler := range handlers {
				switch {
				case evt.UpdateResources != nil:
					handler.Handler.UpdateResources(evt.UpdateResources)
				}
			}
		}
	}()

	go func() {
		defer close(g.eventChannel)

		t := time.NewTicker(g.relistPeriod)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := g.relist(ctx, log); err != nil {
					log.Error(err, "Error relisting")
				}
			}
		}
	}()

	return nil
}

func (g *generator) relist(ctx context.Context, log logr.Logger) error {
	timestamp := time.Now()
	newResources, err := g.getResources(ctx)
	if err != nil {
		return fmt.Errorf("error listing: %w", err)
	}
	g.relistTime.Store(&timestamp)
	g.firstListTime = timestamp

	oldResources := g.resources
	g.resources = newResources

	if proto.Equal(newResources, oldResources) {
		return nil
	}

	evt := &event{UpdateResources: &UpdateResourcesEvent{
		ResourcesOld: oldResources,
		ResourcesNew: newResources,
	}}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case g.eventChannel <- evt:
	default:
		log.Info("Event channel is full, discarding event")
	}

	return nil
}

type handlerRegistration struct {
	generator *generator
	handler   *handler
}

func (r *handlerRegistration) Remove() error {
	r.generator.mu.Lock()
	defer r.generator.mu.Unlock()

	r.generator.handlers.Delete(r.handler)
	return nil
}

func (g *generator) AddHandler(hdl Handler) (HandlerRegistration, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	h := &handler{Handler: hdl}

	g.handlers.Insert(h)
	return &handlerRegistration{
		generator: g,
		handler:   h,
	}, nil
}
