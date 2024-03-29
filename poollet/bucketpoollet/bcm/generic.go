// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bcm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/exp/maps"
	"k8s.io/apimachinery/pkg/util/wait"
	ctrl "sigs.k8s.io/controller-runtime"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
)

type capabilities struct {
	tps  int64
	iops int64
}

func getCapabilities(sriCaps *sri.BucketClassCapabilities) capabilities {
	return capabilities{
		tps:  sriCaps.Tps,
		iops: sriCaps.Iops,
	}
}

type Generic struct {
	mu sync.RWMutex

	sync   bool
	synced chan struct{}

	bucketClassByName         map[string]*sri.BucketClass
	bucketClassByCapabilities map[capabilities][]*sri.BucketClass

	bucketRuntime sri.BucketRuntimeClient

	relistPeriod time.Duration
}

func (g *Generic) relist(ctx context.Context, log logr.Logger) error {
	log.V(1).Info("Relisting bucket classes")
	res, err := g.bucketRuntime.ListBucketClasses(ctx, &sri.ListBucketClassesRequest{})
	if err != nil {
		return fmt.Errorf("error listing bucket classes: %w", err)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	maps.Clear(g.bucketClassByName)
	maps.Clear(g.bucketClassByCapabilities)

	for _, bucketClass := range res.BucketClasses {
		caps := getCapabilities(bucketClass.Capabilities)
		g.bucketClassByName[bucketClass.Name] = bucketClass
		g.bucketClassByCapabilities[caps] = append(g.bucketClassByCapabilities[caps], bucketClass)
	}

	if !g.sync {
		g.sync = true
		close(g.synced)
	}

	return nil
}

func (g *Generic) Start(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx).WithName("vcm")
	wait.UntilWithContext(ctx, func(ctx context.Context) {
		if err := g.relist(ctx, log); err != nil {
			log.Error(err, "Error relisting")
		}
	}, g.relistPeriod)
	return nil
}

func (g *Generic) GetBucketClassFor(ctx context.Context, name string, caps *sri.BucketClassCapabilities) (*sri.BucketClass, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	expected := getCapabilities(caps)
	if byName, ok := g.bucketClassByName[name]; ok && getCapabilities(byName.Capabilities) == expected {
		return byName, nil
	}

	if byCaps, ok := g.bucketClassByCapabilities[expected]; ok {
		switch len(byCaps) {
		case 0:
			return nil, ErrNoMatchingBucketClass
		case 1:
			class := *byCaps[0]
			return &class, nil
		default:
			return nil, ErrAmbiguousMatchingBucketClass
		}
	}

	return nil, ErrNoMatchingBucketClass
}

func (g *Generic) WaitForSync(ctx context.Context) error {
	select {
	case <-g.synced:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type GenericOptions struct {
	RelistPeriod time.Duration
}

func setGenericOptionsDefaults(o *GenericOptions) {
	if o.RelistPeriod == 0 {
		o.RelistPeriod = 1 * time.Hour
	}
}

func NewGeneric(runtime sri.BucketRuntimeClient, opts GenericOptions) BucketClassMapper {
	setGenericOptionsDefaults(&opts)
	return &Generic{
		synced:                    make(chan struct{}),
		bucketClassByName:         map[string]*sri.BucketClass{},
		bucketClassByCapabilities: map[capabilities][]*sri.BucketClass{},
		bucketRuntime:             runtime,
		relistPeriod:              opts.RelistPeriod,
	}
}
