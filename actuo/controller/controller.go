// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/client-go/util/workqueue"
	"spheric.cloud/spheric/actuo/reconcile"
	"spheric.cloud/spheric/actuo/source"
)

type Controller[Request comparable] interface {
	// Watch watches the provided Source.
	Watch(src source.Source[Request]) error

	// Start starts the controller.  Start blocks until the context is closed or a
	// controller has an error starting.
	Start(ctx context.Context) error
}

type controller[Request comparable] struct {
	mu      sync.Mutex
	started bool
	ctx     context.Context

	name                    string
	reconciler              reconcile.Reconciler[Request]
	queue                   workqueue.TypedRateLimitingInterface[Request]
	maxConcurrentReconciles int

	startWatches []source.Source[Request]
}

func New[Request comparable](name string, reconciler reconcile.Reconciler[Request]) (Controller[Request], error) {
	if name == "" {
		return nil, fmt.Errorf("must specify name")
	}
	if reconciler == nil {
		return nil, fmt.Errorf("must specify reconciler")
	}

	rateLimiter := workqueue.DefaultTypedControllerRateLimiter[Request]()
	queue := workqueue.NewTypedRateLimitingQueueWithConfig(
		rateLimiter,
		workqueue.TypedRateLimitingQueueConfig[Request]{
			Name: name,
		},
	)

	return &controller[Request]{
		name:       name,
		reconciler: reconciler,
		queue:      queue,
	}, nil
}

func (c *controller[Request]) reconcileHandler(ctx context.Context, req Request) {
	res, err := c.reconciler.Reconcile(ctx, req)
	switch {
	case err != nil:
		c.queue.AddRateLimited(req)
	case res.RequeueAfter > 0:
		c.queue.Forget(req)
		c.queue.AddAfter(req, res.RequeueAfter)
	case res.Requeue:
		c.queue.AddRateLimited(req)
	default:
		c.queue.Forget(req)
	}
}

func (c *controller[Request]) processNextWorkItem(ctx context.Context) bool {
	req, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(req)

	c.reconcileHandler(ctx, req)
	return true
}

func (c *controller[Request]) Watch(src source.Source[Request]) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.started {
		c.startWatches = append(c.startWatches, src)
		return nil
	}

	return src.Start(c.ctx, c.queue)
}

func (c *controller[Request]) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.started {
		return fmt.Errorf("controller %s already started", c.name)
	}

	c.ctx = ctx

	var wg sync.WaitGroup
	err := func() error {
		defer c.mu.Unlock()

		for _, src := range c.startWatches {
			if err := src.Start(ctx, c.queue); err != nil {
				return err
			}
		}
		for _, src := range c.startWatches {
			err := func() error {
				srcStartedCtx, cancel := context.WithCancel(ctx)
				defer cancel()

				select {
				case <-srcStartedCtx.Done():
					return fmt.Errorf("error waiting for source %v to start", src)
				case <-src.Started():
					return nil
				}
			}()
			if err != nil {
				return err
			}
		}

		for range c.maxConcurrentReconciles {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for c.processNextWorkItem(ctx) {
				}
			}()
		}

		c.started = true
		return nil
	}()
	if err != nil {
		return err
	}

	<-ctx.Done()
	wg.Wait()
	return nil
}
