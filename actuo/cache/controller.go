// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-logr/logr"
)

type ResourceEventHandler[Object any] interface {
	OnAdd(obj Object, isInitialList bool)
	OnUpdate(oldObj, newObj Object)
	OnDelete(obj Object)
}

type ResourceEventHandlerFuncs[Object any] struct {
	AddFunc    func(obj Object, isInitialList bool)
	UpdateFunc func(oldObj, newObj Object)
	DeleteFunc func(obj Object)
}

func (f ResourceEventHandlerFuncs[Object]) OnAdd(obj Object, isInitialList bool) {
	if f.AddFunc != nil {
		f.AddFunc(obj, isInitialList)
	}
}

func (f ResourceEventHandlerFuncs[Object]) OnUpdate(oldObj, newObj Object) {
	if f.UpdateFunc != nil {
		f.UpdateFunc(oldObj, newObj)
	}
}

func (f ResourceEventHandlerFuncs[Object]) OnDelete(obj Object) {
	if f.DeleteFunc != nil {
		f.DeleteFunc(obj)
	}
}

type DeltaFIFOReconciler[Object any] interface {
	Reconcile(ctx context.Context, req ReconcileRequest[Object]) error
}

type ControllerOptions struct {
	Logger logr.Logger
	Resync <-chan struct{}
}

type controller[Key comparable, Object any] struct {
	lw      ListerWatcher[Object]
	fifo    *DeltaFIFO[Key, Object]
	store   Store[Key, Object]
	keyFunc KeyFunc[Key, Object]

	reconciler DeltaFIFOReconciler[Object]

	resync <-chan struct{}

	log logr.Logger
}

type Synced interface {
	Synced() <-chan struct{}
}

type Controller interface {
	Run(context.Context) error
}

func NewController[Key comparable, Object any](
	keyFunc KeyFunc[Key, Object],
	lw ListerWatcher[Object],
	store Store[Key, Object],
	reconciler DeltaFIFOReconciler[Object],
	opts ControllerOptions,
) Controller {
	logger := opts.Logger
	if logger.GetSink() == nil {
		logger = logr.Discard()
	}

	fifo := NewDeltaFIFO(keyFunc, DeltaFIFOOptions{
		Logger: logger.WithName("queue"),
	})

	return &controller[Key, Object]{
		lw:         lw,
		fifo:       fifo,
		store:      store,
		keyFunc:    keyFunc,
		reconciler: reconciler,
		resync:     opts.Resync,
		log:        logger,
	}
}

type DefaultDeltaFIFOReconciler[Key comparable, Object any] struct {
	keyFunc   KeyFunc[Key, Object]
	store     Store[Key, Object]
	handler   ResourceEventHandler[Object]
	closeOnce sync.Once
	synced    chan struct{}
}

func NewDeltaFIFOReconciler[Key comparable, Object any](
	keyFunc KeyFunc[Key, Object],
	store Store[Key, Object],
	handler ResourceEventHandler[Object],
) *DefaultDeltaFIFOReconciler[Key, Object] {
	return &DefaultDeltaFIFOReconciler[Key, Object]{
		keyFunc: keyFunc,
		store:   store,
		handler: handler,
	}
}

func processDeltas[Key, Object any](
	keyFunc KeyFunc[Key, Object],
	store Store[Key, Object],
	handler ResourceEventHandler[Object],
	deltas Deltas[Object],
) error {
	for _, d := range deltas {
		obj := d.Object
		key, err := keyFunc(obj)
		if err != nil {
			return err
		}

		switch d.Type {
		case Resync, Populate, Added, Updated:
			if old, err := store.Get(key); err == nil {
				if err := store.Set(obj); err != nil {
					return err
				}
				handler.OnUpdate(old, obj)
			} else {
				if err := store.Set(obj); err != nil {
					return err
				}
				handler.OnAdd(obj, d.Type == Populate)
			}
		case Deleted:
			if err := store.Delete(key); err != nil {
				return err
			}
			handler.OnDelete(obj)
		}
	}
	return nil
}

func (r *DefaultDeltaFIFOReconciler[Key, Object]) Reconcile(ctx context.Context, req ReconcileRequest[Object]) error {
	switch req.Type {
	case Populated:
		r.closeOnce.Do(func() { close(r.synced) })
		return nil
	case Resynced:
		return nil
	case NewDeltas:
		return processDeltas(r.keyFunc, r.store, r.handler, req.Deltas)
	default:
		return fmt.Errorf("unexpected event %q", req.Type)
	}
}

func (r *DefaultDeltaFIFOReconciler[Key, Object]) Synced() <-chan struct{} {
	return r.synced
}

type ReconcileRequest[Object any] struct {
	Type   DeltaFIFOEventType
	Deltas Deltas[Object]
}

func (c *controller[Key, Object]) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-c.fifo.Events():
			if !ok {
				c.log.Info("fifo closed")
				return
			}

			log := c.log.WithValues("EventType", event.Type)

			var deltas Deltas[Object]
			if event.Type == NewDeltas {
				var err error
				deltas, err = c.fifo.Pop(ctx, event.Key)
				if err != nil {
					if !errors.Is(err, ErrFIFOClosed) {
						log.Error(err, "Error popping deltas", "Key", event.Key)
						continue
					}

					log.Info("FIFO closed")
					return
				}
			}

			req := ReconcileRequest[Object]{Type: event.Type, Deltas: deltas}
			if err := c.reconciler.Reconcile(ctx, req); err != nil {
				log.Error(err, "Error reconciling")
				continue
			}
		}
	}
}

func (c *controller[Key, Object]) resyncLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-c.resync:
			if !ok {
				c.log.Info("Resync channel closed")
				return
			}

			if err := c.fifo.Resync(ctx, c.store.All()); err != nil {
				c.log.Error(err, "Error resyncing")
			}
		}
	}
}

func (c *controller[Key, Object]) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.fifo.Run(ctx); err != nil {
			c.log.Error(err, "Error running fifo")
		}
	}()

	reflector := NewReflector(c.lw, c.fifo, ReflectorOptions{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := reflector.Run(ctx); err != nil {
			c.log.Error(err, "Error running reflector")
		}
	}()

	if c.resync != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c.resyncLoop(ctx)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.loop(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
	return nil
}
