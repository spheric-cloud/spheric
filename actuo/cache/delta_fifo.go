// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"sync"

	"github.com/go-logr/logr"
	"spheric.cloud/spheric/utils/chans"
	"spheric.cloud/spheric/utils/container/squeue"
	"spheric.cloud/spheric/utils/generic"
)

var ErrFIFOClosed = errors.New("fifo closed")

type DeltaType string

const (
	Added    DeltaType = "Added"
	Updated  DeltaType = "Updated"
	Deleted  DeltaType = "Deleted"
	Populate DeltaType = "Populate"
	Resync   DeltaType = "Resync"
)

type Delta[Object any] struct {
	Type   DeltaType
	Object Object
}

type Deltas[Object any] []Delta[Object]

func (d Deltas[Object]) Newest() *Delta[Object] {
	if n := len(d); n > 0 {
		return &d[n-1]
	}
	return nil
}

type deltaFIFOState uint8

const (
	deltaFIFOStateInitial deltaFIFOState = iota
	deltaFIFOStateRunning
	deltaFIFOStateStopping
	deltaFIFOStateStopped
)

type DeltaFIFO[Key comparable, Object any] struct {
	lock  sync.RWMutex
	state deltaFIFOState

	keyFunc KeyFunc[Key, Object]
	logger  logr.Logger

	in  chan deltaFIFOInput
	out chan DeltaFIFOEvent[Key]
}

type DeltaFIFOOptions struct {
	Logger logr.Logger
}

func NewDeltaFIFO[Key comparable, Object any](keyFunc KeyFunc[Key, Object], opts DeltaFIFOOptions) *DeltaFIFO[Key, Object] {
	logger := opts.Logger
	if logger.GetSink() == nil {
		logger = logr.Discard()
	}

	return &DeltaFIFO[Key, Object]{
		keyFunc: keyFunc,
		logger:  logger,

		in:  make(chan deltaFIFOInput),
		out: make(chan DeltaFIFOEvent[Key]),
	}
}

type deltaFIFOInput interface {
	deltaFIFOInput()
}

type deltaFIFOInputDelta[Object any] Delta[Object]

func (deltaFIFOInputDelta[Object]) deltaFIFOInput() {}

type deltaFIFOInputPopulate[Object any] iter.Seq[Object]

func (deltaFIFOInputPopulate[Object]) deltaFIFOInput() {}

type deltaFIFOInputResync[Object any] iter.Seq[Object]

func (deltaFIFOInputResync[Object]) deltaFIFOInput() {}

type deltaFIFOInputPop[Key, Object any] struct {
	key  Key
	resp chan<- Deltas[Object]
}

func (f deltaFIFOInputPop[Key, Object]) deltaFIFOInput() {}

func (f *DeltaFIFO[Key, Object]) submitInput(ctx context.Context, input deltaFIFOInput) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case f.in <- input:
		return nil
	}
}

func (f *DeltaFIFO[Key, Object]) Add(ctx context.Context, obj Object) error {
	return f.submitInput(ctx, deltaFIFOInputDelta[Object]{Added, obj})
}

func (f *DeltaFIFO[Key, Object]) Update(ctx context.Context, obj Object) error {
	return f.submitInput(ctx, deltaFIFOInputDelta[Object]{Updated, obj})
}

func (f *DeltaFIFO[Key, Object]) Delete(ctx context.Context, obj Object) error {
	return f.submitInput(ctx, deltaFIFOInputDelta[Object]{Deleted, obj})
}

func (f *DeltaFIFO[Key, Object]) Populate(ctx context.Context, objs iter.Seq[Object]) error {
	return f.submitInput(ctx, deltaFIFOInputPopulate[Object](objs))
}

func (f *DeltaFIFO[Key, Object]) Resync(ctx context.Context, objs iter.Seq[Object]) error {
	return f.submitInput(ctx, deltaFIFOInputResync[Object](objs))
}

func _[Key comparable, Object any]() ReflectorSink[Object] {
	return generic.Stub[*DeltaFIFO[Key, Object]]()
}

type deltaFIFOProcessor[Key comparable, Object any] struct {
	*DeltaFIFO[Key, Object]
	ctx     context.Context
	keyFunc KeyFunc[Key, Object]
	items   map[Key]Deltas[Object]

	inputs <-chan deltaFIFOInput
	events chan DeltaFIFOEvent[Key]
}

func (f *deltaFIFOProcessor[Key, Object]) submitEvent(event DeltaFIFOEvent[Key]) bool {
	select {
	case <-f.ctx.Done():
		return false
	case f.events <- event:
		return true
	}
}

func (f *deltaFIFOProcessor[Key, Object]) submitDelta(delta Delta[Object]) bool {
	key, _ := f.keyFunc(delta.Object)

	oldDeltas, ok := f.items[key]
	f.items[key] = appendDelta(oldDeltas, delta)
	if ok {
		// Only submit an event if none is in flight yet
		return true
	}

	return f.submitEvent(DeltaFIFOEvent[Key]{NewDeltas, key})
}

func (f *deltaFIFOProcessor[Key, Object]) loop() {
	defer close(f.events)

	defer func() {
		// Drain remaining inputs
		for in := range f.inputs {
			if p, ok := in.(deltaFIFOInputPop[Key, Object]); ok {
				close(p.resp)
			}
		}
	}()

	var populated bool
	for {
		select {
		case <-f.ctx.Done():
			return
		case in := <-f.inputs:
			switch in := in.(type) {
			case deltaFIFOInputPopulate[Object]:
				for obj := range in {
					if !f.submitDelta(Delta[Object]{Populate, obj}) {
						return
					}
				}
				if !populated {
					if !f.submitEvent(DeltaFIFOEvent[Key]{Type: Populated}) {
						return
					}
					populated = true
				}
			case deltaFIFOInputDelta[Object]:
				if !populated {
					if !f.submitEvent(DeltaFIFOEvent[Key]{Type: Populated}) {
						return
					}
					populated = true
				}
				if !f.submitDelta(Delta[Object](in)) {
					return
				}
			case deltaFIFOInputResync[Object]:
				for obj := range in {
					key, _ := f.keyFunc(obj)
					if _, exists := f.items[key]; exists {
						// Don't resync if there's already an event in-flight.
						continue
					}

					if !f.submitDelta(Delta[Object]{Resync, obj}) {
						return
					}
				}
				if !f.submitEvent(DeltaFIFOEvent[Key]{Type: Resynced}) {
					return
				}
			case deltaFIFOInputPop[Key, Object]:
				item := f.items[in.key]
				delete(f.items, in.key)
				in.resp <- item
			}
		}
	}
}

func (f *DeltaFIFO[Key, Object]) Events() <-chan DeltaFIFOEvent[Key] {
	return f.out
}

func (f *DeltaFIFO[Key, Object]) Pop(ctx context.Context, key Key) (Deltas[Object], error) {
	resp := make(chan Deltas[Object], 1)
	if err := f.submitInput(ctx, deltaFIFOInputPop[Key, Object]{key, resp}); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case deltas, ok := <-resp:
		if !ok {
			return nil, ErrFIFOClosed
		}
		return deltas, nil
	}
}

func (f *DeltaFIFO[Key, Object]) Run(ctx context.Context) error {
	f.lock.Lock()
	if f.state != deltaFIFOStateInitial {
		f.lock.Unlock()
		return fmt.Errorf("delta fifo already started")
	}

	var (
		wg sync.WaitGroup
	)
	func() {
		defer f.lock.Unlock()
		f.state = deltaFIFOStateRunning

		processor := &deltaFIFOProcessor[Key, Object]{
			ctx:     ctx,
			keyFunc: f.keyFunc,
			items:   make(map[Key]Deltas[Object]),
			inputs:  f.in,
			events:  make(chan DeltaFIFOEvent[Key]),
		}

		wg.Add(1)
		go func() {
			wg.Done()
			chans.BufferedPumpExitOnClose(f.out, processor.events, squeue.New[DeltaFIFOEvent[Key]](1024))
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			processor.loop()
		}()
	}()

	<-ctx.Done()

	f.lock.Lock()
	f.state = deltaFIFOStateStopping
	close(f.in)
	f.lock.Unlock()

	wg.Wait()

	f.lock.Lock()
	f.state = deltaFIFOStateStopped
	close(f.out)
	f.lock.Unlock()

	return nil
}

func appendDelta[Object any](into Deltas[Object], newDelta Delta[Object]) Deltas[Object] {
	if len(into) == 0 {
		return append(into, newDelta)
	}

	a := &into[len(into)-1]
	b := &newDelta
	if out := isDup(a, b); out != nil {
		into[len(into)-1] = *b
		return into
	}

	return append(into, newDelta)
}

func isDup[Object any](a, b *Delta[Object]) *Delta[Object] {
	if b.Type != Deleted || a.Type != Deleted {
		return nil
	}
	return b
}

type DeltaFIFOEventType string

const (
	Populated DeltaFIFOEventType = "Populated"
	Resynced  DeltaFIFOEventType = "Resynced"
	NewDeltas DeltaFIFOEventType = "NewDeltas"
)

type DeltaFIFOEvent[Key any] struct {
	Type DeltaFIFOEventType
	Key  Key
}
