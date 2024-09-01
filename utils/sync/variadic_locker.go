// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package sync

import "sync"

func WithLocker0To0(l sync.Locker, f func()) {
	l.Lock()
	defer l.Unlock()
	f()
}

func WithLocker0To1[F1 any](l sync.Locker, f func() F1) F1 {
	l.Lock()
	defer l.Unlock()
	return f()
}

func WithLocker1To0[E1 any](l sync.Locker, f func(E1), e1 E1) {
	l.Lock()
	defer l.Unlock()
	f(e1)
}

func WithLocker1To1[E1, F1 any](l sync.Locker, f func(E1) F1, e1 E1) F1 {
	l.Lock()
	defer l.Unlock()
	return f(e1)
}

func WithLocker0To2[F1, F2 any](l sync.Locker, f func() (F1, F2)) (F1, F2) {
	l.Lock()
	defer l.Unlock()
	return f()
}

func WithLocker1To2[E1, F1, F2 any](l sync.Locker, f func(E1) (F1, F2), e1 E1) (F1, F2) {
	l.Lock()
	defer l.Unlock()
	return f(e1)
}

func WithLocker2To0[E1, E2 any](l sync.Locker, f func(E1, E2), e1 E1, e2 E2) {
	l.Lock()
	defer l.Unlock()
	f(e1, e2)
}

func WithLocker2To1[E1, E2, F1 any](l sync.Locker, f func(E1, E2) F1, e1 E1, e2 E2) F1 {
	l.Lock()
	defer l.Unlock()
	return f(e1, e2)
}

func WithLocker2To2[E1, E2, F1, F2 any](l sync.Locker, f func(E1, E2) (F1, F2), e1 E1, e2 E2) (F1, F2) {
	l.Lock()
	defer l.Unlock()
	return f(e1, e2)
}
