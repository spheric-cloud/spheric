// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"context"

	"k8s.io/client-go/util/workqueue"
	"spheric.cloud/spheric/actuo/event"
)

type EventHandler[Object any, Request comparable] interface {
	// Create is called in response to a create event.
	Create(context.Context, event.CreateEvent[Object], workqueue.TypedRateLimitingInterface[Request])

	// Update is called in response to an update event.
	Update(context.Context, event.UpdateEvent[Object], workqueue.TypedRateLimitingInterface[Request])

	// Delete is called in response to a delete event.
	Delete(context.Context, event.DeleteEvent[Object], workqueue.TypedRateLimitingInterface[Request])

	// Generic is called in response to an event of an unknown type or a synthetic event.
	Generic(context.Context, event.GenericEvent[Object], workqueue.TypedRateLimitingInterface[Request])
}
