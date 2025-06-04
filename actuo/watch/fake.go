// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package watch

import (
	"sync"
)

type Fake[Object any] struct {
	mu      sync.Mutex
	stopped bool

	events chan Event[Object]
}

func NewFake[Object any]() *Fake[Object] {
	return &Fake[Object]{
		events: make(chan Event[Object]),
	}
}

func NewFakeWithCap[Object any](cap int) *Fake[Object] {
	return &Fake[Object]{
		events: make(chan Event[Object], cap),
	}
}

func (f *Fake[Object]) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.stopped {
		close(f.events)
		f.stopped = true
	}
}

func (f *Fake[Object]) Stopped() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.stopped
}

func (f *Fake[Object]) Events() <-chan Event[Object] {
	return f.events
}

func (f *Fake[Object]) Create(obj Object) {
	f.events <- Event[Object]{Type: EventTypeCreated, Object: obj}
}

func (f *Fake[Object]) Update(obj Object) {
	f.events <- Event[Object]{Type: EventTypeUpdated, Object: obj}
}

func (f *Fake[Object]) Delete(obj Object) {
	f.events <- Event[Object]{Type: EventTypeDeleted, Object: obj}
}
