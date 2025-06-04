// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"context"

	"k8s.io/client-go/util/workqueue"
	cache2 "spheric.cloud/spheric/actuo/cache"
	"spheric.cloud/spheric/actuo/event"
	"spheric.cloud/spheric/actuo/handler"
)

type Source[Request comparable] interface {
	// Start is internal and should be called only by the Controller to start the source.
	// Start must be non-blocking.
	Start(context.Context, workqueue.TypedRateLimitingInterface[Request]) error
	// Started returns any startup errors. If there were no startup errors, started should
	// be closed.
	Started() <-chan error
}

type informer[Key comparable, Object any, Request comparable] struct {
	sharedInformer cache2.SharedInformer[Key, Object]
	handler        handler.EventHandler[Object, Request]
	started        chan error
}

func NewInformer[Key comparable, Object any, Request comparable](
	inf cache2.SharedInformer[Key, Object],
	handler handler.EventHandler[Object, Request],
) Source[Request] {
	return &informer[Key, Object, Request]{
		sharedInformer: inf,
		handler:        handler,
		started:        make(chan error, 1),
	}
}

func (s *informer[Key, Object, Request]) Start(ctx context.Context, q workqueue.TypedRateLimitingInterface[Request]) error {
	go func() {
		handle, err := s.sharedInformer.AddEventHandler(cache2.ResourceEventHandlerFuncs[Object]{
			AddFunc: func(obj Object, isInitialList bool) {
				s.handler.Create(ctx, event.CreateEvent[Object]{Object: obj}, q)
			},
			UpdateFunc: func(oldObj, newObj Object) {
				s.handler.Update(ctx, event.UpdateEvent[Object]{ObjectOld: oldObj, ObjectNew: newObj}, q)
			},
			DeleteFunc: func(obj Object) {
				s.handler.Delete(ctx, event.DeleteEvent[Object]{Object: obj}, q)
			},
		}, cache2.EventHandlerOptions{})
		if err != nil {
			s.started <- err
			return
		}
		defer func() { _ = s.sharedInformer.RemoveEventHandler(handle) }()

		select {
		case <-ctx.Done():
			s.started <- ctx.Err()
			return
		case <-handle.Synced():
			close(s.started)
		}

		<-ctx.Done()
	}()
	return nil
}

func (s *informer[Key, Object, Request]) Started() <-chan error {
	return s.started
}
