// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package watch

type EventType string

const (
	EventTypeCreated EventType = "Created"
	EventTypeUpdated EventType = "Updated"
	EventTypeDeleted EventType = "Deleted"
)

type Event[Object any] struct {
	Type   EventType
	Object Object
}

type Watch[Object any] interface {
	Stop()
	Events() <-chan Event[Object]
}
