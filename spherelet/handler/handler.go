// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"context"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EnqueueRequestForName enqueues a reconcile.Request for a name without namespace.
type EnqueueRequestForName[object any] string

func (e EnqueueRequestForName[object]) enqueue(queue workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	queue.Add(reconcile.Request{NamespacedName: client.ObjectKey{Name: string(e)}})
}

// Create implements handler.EventHandler.
func (e EnqueueRequestForName[object]) Create(_ context.Context, _ event.TypedCreateEvent[object], queue workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueue(queue)
}

// Update implements handler.EventHandler.
func (e EnqueueRequestForName[object]) Update(_ context.Context, _ event.TypedUpdateEvent[object], queue workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueue(queue)
}

// Delete implements handler.EventHandler.
func (e EnqueueRequestForName[object]) Delete(_ context.Context, _ event.TypedDeleteEvent[object], queue workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueue(queue)
}

// Generic implements handler.EventHandler.
func (e EnqueueRequestForName[object]) Generic(_ context.Context, _ event.TypedGenericEvent[object], queue workqueue.TypedRateLimitingInterface[reconcile.Request]) {
	e.enqueue(queue)
}
