// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"context"
	"errors"
	"sync"
)

type Group struct {
	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	shutdownHooksMu sync.Mutex
	shutdownHooks   []func()

	errsMu sync.Mutex
	errs   []error
}

func NewGroup(ctx context.Context) *Group {
	g := &Group{}
	g.ctx, g.cancel = context.WithCancel(ctx)

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		<-g.ctx.Done()

		g.shutdownHooksMu.Lock()
		defer g.shutdownHooksMu.Unlock()
		for _, shutdownHook := range g.shutdownHooks {
			shutdownHook()
		}
	}()

	return g
}

func (g *Group) OnShutdown(f func()) {
	g.shutdownHooksMu.Lock()
	defer g.shutdownHooksMu.Unlock()
	g.shutdownHooks = append(g.shutdownHooks, f)
}

func (g *Group) addErr(err error) {
	g.errsMu.Lock()
	defer g.errsMu.Unlock()
	g.errs = append(g.errs, err)
}

type OnError uint8

const (
	OnErrorContinue OnError = iota
	OnErrorStop
)

func (oe OnError) ApplyToStart(o *StartOptions) {
	o.OnError = oe
}

type StartOptions struct {
	Context context.Context
	OnError OnError
}

func (o *StartOptions) ApplyOptions(opts []StartOption) *StartOptions {
	for _, opt := range opts {
		opt.ApplyToStart(o)
	}
	return o
}

type StartOption interface {
	ApplyToStart(o *StartOptions)
}

func (g *Group) Start(f func(context.Context) error, opts ...StartOption) {
	o := (&StartOptions{}).ApplyOptions(opts)

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		if err := f(g.ctx); err != nil {
			g.addErr(err)

			switch o.OnError {
			case OnErrorContinue:
			case OnErrorStop:
				g.cancel()
			}
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return errors.Join(g.errs...)
}
