// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instanceevent

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"

	"github.com/go-logr/logr"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/server/healthz"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Generator interface {
	Source
	healthz.HealthChecker
	manager.Runnable
}

type event struct {
	Create  *CreateEvent
	Update  *UpdateEvent
	Delete  *DeleteEvent
	Generic *GenericEvent
}

type oldNewMapEntry struct {
	Old     *sri.Instance
	Current *sri.Instance
}

type oldNewMap map[string]*oldNewMapEntry

func (m oldNewMap) id(obj *sri.Instance) string {
	return obj.GetMetadata().GetId()
}

func (m oldNewMap) setCurrent(current []*sri.Instance) {
	for _, v := range m {
		v.Current = nil
	}

	for _, item := range current {
		item := item
		id := m.id(item)
		if r, ok := m[id]; ok {
			r.Current = item
		} else {
			m[id] = &oldNewMapEntry{
				Current: item,
			}
		}
	}
}

func (m oldNewMap) getCurrent(id string) (*sri.Instance, bool) {
	r, ok := m[id]
	if ok && r.Current != nil {
		return r.Current, true
	}
	return nil, false
}

func (m oldNewMap) getOld(id string) (*sri.Instance, bool) {
	r, ok := m[id]
	if ok && r.Old != nil {
		return r.Old, true
	}
	return nil, false
}

func (m oldNewMap) update(id string) {
	r, ok := m[id]
	if !ok {
		return
	}

	if r.Current == nil {
		delete(m, id)
		return
	}

	r.Old = r.Current
	r.Current = nil
}

type handler struct {
	Handler
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

	items oldNewMap

	list func(ctx context.Context) ([]*sri.Instance, error)
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

func NewGenerator(list func(ctx context.Context) ([]*sri.Instance, error), opts GeneratorOptions) Generator {
	setGeneratorOptionsDefaults(&opts)

	return &generator{
		eventChannel:    make(chan *event, opts.ChannelCapacity),
		relistPeriod:    opts.RelistPeriod,
		relistThreshold: opts.RelistThreshold,
		relistTime:      atomic.Pointer[time.Time]{},
		firstListTime:   time.Time{},
		items:           make(oldNewMap),
		list:            list,
		handlers:        sets.New[*handler](),
	}
}

func (g *generator) Name() string {
	return "instance-event-generator"
}

func (g *generator) Check(_ *http.Request) error {
	relistTime := g.relistTime.Load()
	if relistTime == nil {
		return fmt.Errorf("mleg did not relist yet")
	}

	elapsed := time.Since(*relistTime)
	if elapsed > g.relistThreshold {
		return fmt.Errorf("mleg was last seen active %v ago, threshold is %v", elapsed, g.relistThreshold)
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
				case evt.Create != nil:
					handler.Create(*evt.Create)
				case evt.Update != nil:
					handler.Update(*evt.Update)
				case evt.Delete != nil:
					handler.Delete(*evt.Delete)
				case evt.Generic != nil:
					handler.Generic(*evt.Generic)
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
	objects, err := g.list(ctx)
	if err != nil {
		return fmt.Errorf("error listing: %w", err)
	}
	g.relistTime.Store(&timestamp)
	g.firstListTime = timestamp

	g.items.setCurrent(objects)

	eventsByKey := make(map[string][]*event)
	for key := range g.items {
		itemOld, oldOK := g.items.getOld(key)
		itemNew, newOK := g.items.getCurrent(key)
		switch {
		case !oldOK && newOK:
			createdAt := time.Unix(0, itemNew.GetMetadata().CreatedAt)
			if createdAt.Before(g.firstListTime) {
				eventsByKey[key] = []*event{{Create: &CreateEvent{Object: itemNew}}}
			} else {
				eventsByKey[key] = []*event{{Generic: &GenericEvent{Object: itemNew}}}
			}
		case oldOK && !newOK:
			eventsByKey[key] = []*event{{Delete: &DeleteEvent{Object: itemOld}}}
		case oldOK && newOK:
			if !proto.Equal(itemOld, itemNew) {
				eventsByKey[key] = []*event{{Update: &UpdateEvent{ObjectOld: itemOld, ObjectNew: itemNew}}}
			}
		}
	}

	for machineID, events := range eventsByKey {
		g.items.update(machineID)
		for i := range events {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case g.eventChannel <- events[i]:
			default:
				log.Info("Event channel is full, discarding event", "InstanceID", machineID)
			}
		}
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
