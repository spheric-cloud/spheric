// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package debug

import (
	"context"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

type loggingQueue[T comparable] struct {
	mu sync.RWMutex

	done bool
	log  logr.Logger
	workqueue.TypedRateLimitingInterface[T]
}

func newLoggingQueue[T comparable](log logr.Logger, queue workqueue.TypedRateLimitingInterface[T]) *loggingQueue[T] {
	return &loggingQueue[T]{log: log, TypedRateLimitingInterface: queue}
}

func (q *loggingQueue[T]) Add(item T) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	q.log.Info("Add", "Item", item, "Done", q.done)
	q.TypedRateLimitingInterface.Add(item)
}

func (q *loggingQueue[T]) AddRateLimited(item T) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	q.log.Info("AddRateLimited", "Item", item, "Done", q.done)
	q.TypedRateLimitingInterface.AddRateLimited(item)
}

func (q *loggingQueue[T]) AddAfter(item T, duration time.Duration) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	q.log.Info("AddAfter", "Item", item, "Duration", duration, "Done", q.done)
	q.TypedRateLimitingInterface.AddAfter(item, duration)
}

func (q *loggingQueue[T]) Finish() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.done = true
}

type debugHandler[object any, request comparable] struct {
	log         logr.Logger
	handler     handler.TypedEventHandler[object, request]
	objectValue func(object) any
}

func (d *debugHandler[object, request]) Create(ctx context.Context, evt event.TypedCreateEvent[object], queue workqueue.TypedRateLimitingInterface[request]) {
	log := d.log.WithValues("Event", "Create", "Object", d.objectValue(evt.Object))
	log.Info("Handling Event")

	lQueue := newLoggingQueue(log.WithName("Queue"), queue)
	defer lQueue.Finish()

	d.handler.Create(ctx, evt, lQueue)
}

func (d *debugHandler[object, request]) Update(ctx context.Context, evt event.TypedUpdateEvent[object], queue workqueue.TypedRateLimitingInterface[request]) {
	log := d.log.WithValues("Event", "Update", "ObjectOld", d.objectValue(evt.ObjectOld), "ObjectNew", d.objectValue(evt.ObjectNew))
	log.Info("Handling Event")

	lQueue := newLoggingQueue(log.WithName("Queue"), queue)
	defer lQueue.Finish()

	d.handler.Update(ctx, evt, lQueue)
}

func (d *debugHandler[object, request]) Delete(ctx context.Context, evt event.TypedDeleteEvent[object], queue workqueue.TypedRateLimitingInterface[request]) {
	log := d.log.WithValues("Event", "Delete", "Object", d.objectValue(evt.Object))
	log.Info("Handling Event")

	lQueue := newLoggingQueue(log.WithName("Queue"), queue)
	defer lQueue.Finish()

	d.handler.Delete(ctx, evt, lQueue)
}

func (d *debugHandler[object, request]) Generic(ctx context.Context, evt event.TypedGenericEvent[object], queue workqueue.TypedRateLimitingInterface[request]) {
	log := d.log.WithValues("Event", "Generic", "Object", d.objectValue(evt.Object))
	log.Info("Handling Event")

	lQueue := newLoggingQueue(log.WithName("Queue"), queue)
	defer lQueue.Finish()

	d.handler.Generic(ctx, evt, lQueue)
}

// TypedHandler allows debugging a handler.EventHandler by wrapping it and logging each action it does.
//
// Caution: This has a heavy toll on runtime performance and should *not* be used in production code.
// Use only for debugging handlers and remove once done.
func TypedHandler[object any, request comparable](name string, handler handler.TypedEventHandler[object, request], opts ...TypedHandlerOption[object, request]) handler.TypedEventHandler[object, request] {
	o := (&TypedHandlerOptions[object, request]{}).ApplyOptions(opts)
	setTypedHandlerOptionsDefaults(o)

	return &debugHandler[object, request]{
		log:         o.Log.WithName(name),
		handler:     handler,
		objectValue: o.ObjectValue,
	}
}
