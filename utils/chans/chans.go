// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package chans

import (
	"spheric.cloud/spheric/utils/constraints"
)

type Queue[T any] interface {
	Enqueue(v T)
	Dequeue() (T, bool)
}

type Toggle[Ch constraints.Channel[T], T any] struct {
	enabled bool
	ch      func() Ch
}

func (t *Toggle[Ch, T]) C() Ch {
	if t.enabled {
		return t.ch()
	}
	return nil
}

func (t *Toggle[Ch, T]) Disable() *Toggle[Ch, T] {
	t.enabled = false
	return t
}

func (t *Toggle[Ch, T]) Enable() *Toggle[Ch, T] {
	t.enabled = true
	return t
}

func ToggleOf[T any, Ch constraints.Channel[T]](ch Ch) *Toggle[Ch, T] {
	return &Toggle[Ch, T]{
		enabled: true,
		ch:      func() Ch { return ch },
	}
}

func ToggleOfFunc[T any, Ch constraints.Channel[T]](f func() Ch) *Toggle[Ch, T] {
	return &Toggle[Ch, T]{
		enabled: true,
		ch:      f,
	}
}

// BufferedReceive returns a receive-only channel that buffers items of the incoming chan using the given queue.
func BufferedReceive[T any](in <-chan T, q Queue[T]) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)

		var (
			next T
			nout chan T
		)
		for in != nil || nout != nil {
			select {
			case nout <- next:
				var ok bool
				next, ok = q.Dequeue()
				if !ok {
					nout = nil
				}
				continue
			default:
			}

			select {
			case v, ok := <-in:
				if !ok {
					in = nil
					continue
				}
				if nout == nil {
					nout = out
					next = v
				} else {
					q.Enqueue(v)
				}
			case nout <- next:
				var ok bool
				next, ok = q.Dequeue()
				if !ok {
					nout = nil
				}
			}
		}
	}()
	return out
}

func BufferedPumpExitOnClose[T any](out chan<- T, in <-chan T, q Queue[T]) {
	var (
		next T
		nout chan<- T
	)
	for {
		select {
		case nout <- next:
			var ok bool
			next, ok = q.Dequeue()
			if !ok {
				nout = nil
			}
			continue
		case v, ok := <-in:
			if !ok {
				return
			}

			if nout == nil {
				nout = out
				next = v
			} else {
				q.Enqueue(v)
			}
		}
	}
}
