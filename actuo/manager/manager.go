// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package manager

import (
	"context"
	"fmt"
	"sync"
)

type Runnable interface {
	Start(ctx context.Context) error
}

type Manager interface {
	Add(runnable Runnable) error
	Start(ctx context.Context) error
}

type manager struct {
	errChan   chan error
	runnables *runnableGroup
}

func New() Manager {
	errChan := make(chan error, 1)

	return &manager{
		errChan: errChan,
		runnables: newRunnableGroup(
			context.Background,
			func(err error) {
				errChan <- err
			},
		),
	}
}

func (m *manager) Add(runnable Runnable) error {
	if err := m.runnables.Add(runnable); err != nil {
		return err
	}
	return nil
}

func (m *manager) Start(ctx context.Context) error {
	if err := m.runnables.Start(); err != nil {
		return err
	}
	defer func() { _ = m.runnables.Stop(ctx) }()

	select {
	case <-ctx.Done():
		return nil
	case err := <-m.errChan:
		// Error starting or running a runnable
		return err
	}
}

type state uint8

const (
	stateInitial state = iota
	stateRunning
	stateStopped
)

type runnableGroup struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu         sync.RWMutex
	state      state
	startQueue []Runnable

	wg      sync.WaitGroup
	ch      chan Runnable
	onError func(error)
}

func newRunnableGroup(
	baseContext func() context.Context,
	onError func(error),
) *runnableGroup {
	ctx, cancel := context.WithCancel(baseContext())

	return &runnableGroup{
		ctx:     ctx,
		cancel:  cancel,
		ch:      make(chan Runnable),
		onError: onError,
	}
}

func (m *runnableGroup) reconcileLoop() {
	for runnable := range m.ch {
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			if err := runnable.Start(m.ctx); err != nil {
				if m.onError != nil {
					m.onError(err)
				}
			}
		}()
	}
}

func (m *runnableGroup) Start() error {
	m.mu.Lock()
	if m.state != stateInitial {
		m.mu.Unlock()
		return fmt.Errorf("cannot start multiple times")
	}
	defer m.mu.Unlock()

	m.state = stateRunning

	go m.reconcileLoop()

	for _, runnable := range m.startQueue {
		m.ch <- runnable
	}

	return nil
}

func (m *runnableGroup) Stop(ctx context.Context) error {
	m.mu.Lock()
	if m.state != stateRunning {
		m.mu.Unlock()
		return fmt.Errorf("cannot stop non-running group")
	}

	defer m.mu.Unlock()

	m.cancel()
	m.state = stateStopped

	done := make(chan struct{})
	go func() {
		defer close(done)
		m.wg.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *runnableGroup) Add(runnable Runnable) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch m.state {
	case stateStopped:
		return fmt.Errorf("group is stopped")
	case stateInitial:
		m.startQueue = append(m.startQueue, runnable)
		return nil
	case stateRunning:
		m.ch <- runnable
		return nil
	default:
		return fmt.Errorf("group is in invalid state %d", m.state)
	}
}
