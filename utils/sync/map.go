// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package sync

import (
	"iter"
	"sync"

	"golang.org/x/exp/maps"
)

type KLock[K comparable] struct {
	lock sync.RWMutex

	stateLock sync.Mutex

	pendingGlobalWrites uint32
	pendingGlobalReads  uint32

	activeGlobalReads uint32
	activeWrites      uint32

	stateCondGlobalWriteOk *sync.Cond
	stateCondGlobalReadOk  *sync.Cond
	stateCondReadOk        *sync.Cond
	stateCondWriteOk       *sync.Cond

	keyLocks map[K]*keyLockEntry
}

type keyLockEntry struct {
	waiting uint32
	lock    sync.RWMutex
}

func (m *KLock[K]) Lock() {
	m.lock.Lock()

	m.stateLock.Lock()
	if len(m.keyLocks) == 0 {
		return
	}

	m.pendingGlobalWrites++
	if m.stateCondGlobalWriteOk == nil {
		m.stateCondGlobalWriteOk = sync.NewCond(&m.stateLock)
	}
	m.stateCondGlobalWriteOk.Wait()
	m.pendingGlobalWrites--
}

func (m *KLock[K]) Unlock() {
	m.notify()
	m.stateLock.Unlock()
	m.lock.Unlock()
}

func (m *KLock[K]) RLock() {
	m.lock.RLock()

	m.stateLock.Lock()
	defer m.stateLock.Unlock()

	if m.activeWrites == 0 {
		m.activeGlobalReads++
		return
	}

	m.pendingGlobalReads++
	if m.stateCondGlobalReadOk == nil {
		m.stateCondGlobalReadOk = sync.NewCond(&m.stateLock)
	}
	m.stateCondGlobalReadOk.Wait()
	m.pendingGlobalReads--
}

func (m *KLock[K]) RUnlock() {
	m.stateLock.Lock()
	defer m.stateLock.Unlock()

	m.activeGlobalReads--
	m.lock.RUnlock()

	m.notify()
}

func (m *KLock[K]) addOrIncreasedEntry(k K) *keyLockEntry {
	if m.keyLocks == nil {
		m.keyLocks = make(map[K]*keyLockEntry)
	}

	entry := m.keyLocks[k]
	if entry == nil {
		entry = &keyLockEntry{}
		m.keyLocks[k] = entry
	} else {
		entry.waiting++
	}
	return entry
}

func (m *KLock[K]) deleteOrDecreasedEntry(k K) *keyLockEntry {
	if m.keyLocks == nil {
		m.keyLocks = make(map[K]*keyLockEntry)
	}

	entry := m.keyLocks[k]
	if entry.waiting == 0 {
		delete(m.keyLocks, k)
	} else {
		entry.waiting--
	}
	return entry
}

// notify must only be called while holding the state lock and after updating any
// relevant 'active' state values.
// It notifies relevant goroutines to continue.
func (m *KLock[K]) notify() {
	if m.pendingGlobalWrites > 0 {
		if len(m.keyLocks) == 0 {
			m.stateCondGlobalWriteOk.Signal()
		}
		return
	}

	if m.stateCondReadOk != nil {
		m.stateCondReadOk.Broadcast()
	}

	if m.pendingGlobalReads > 0 {
		if m.activeWrites == 0 {
			m.stateCondGlobalReadOk.Signal()
		}
		return
	}

	if m.stateCondWriteOk != nil {
		m.stateCondWriteOk.Broadcast()
	}
}

func (m *KLock[K]) canWriteAnyKey() bool {
	return m.activeGlobalReads == 0 && m.pendingGlobalReads == 0 && m.pendingGlobalWrites == 0
}

func (m *KLock[K]) LockKey(k K) {
	m.stateLock.Lock()
	defer m.stateLock.Unlock()

	for !m.canWriteAnyKey() {
		if m.stateCondWriteOk == nil {
			m.stateCondWriteOk = sync.NewCond(&m.stateLock)
		}
		m.stateCondWriteOk.Wait()
	}

	m.activeWrites++
	m.addOrIncreasedEntry(k).lock.Lock()
}

func (m *KLock[K]) UnlockKey(k K) {
	m.stateLock.Lock()
	defer m.stateLock.Unlock()

	entry := m.deleteOrDecreasedEntry(k)
	m.activeWrites--
	entry.lock.Unlock()

	m.notify()
}

func (m *KLock[K]) canReadAnyKey() bool {
	return m.pendingGlobalWrites == 0
}

func (m *KLock[K]) RLockKey(k K) {
	m.stateLock.Lock()
	defer m.stateLock.Unlock()

	for !m.canReadAnyKey() {
		if m.stateCondReadOk == nil {
			m.stateCondReadOk = sync.NewCond(&m.stateLock)
		}
		m.stateCondReadOk.Wait()
	}

	entry := m.addOrIncreasedEntry(k)
	entry.lock.RLock()
}

func (m *KLock[K]) RUnlockKey(k K) {
	m.stateLock.Lock()
	defer m.stateLock.Unlock()

	entry := m.deleteOrDecreasedEntry(k)
	entry.lock.RUnlock()

	m.notify()
}

type Map[K comparable, V any] struct {
	lock  KLock[K]
	items map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		items: make(map[K]V),
	}
}

func (m *Map[K, V]) GetOK(k K) (V, bool) {
	m.lock.RLockKey(k)
	defer m.lock.RUnlockKey(k)

	v, ok := m.items[k]
	return v, ok
}

func (m *Map[K, V]) Set(k K, v V) {
	m.lock.LockKey(k)
	defer m.lock.UnlockKey(k)

	m.items[k] = v
}

func (m *Map[K, V]) Delete(k K) {
	m.lock.LockKey(k)
	defer m.lock.UnlockKey(k)

	delete(m.items, k)
}

func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.lock.RLock()
		defer m.lock.RUnlock()

		for k, v := range m.items {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		m.lock.RLock()
		defer m.lock.RUnlock()

		for _, v := range m.items {
			if !yield(v) {
				return
			}
		}
	}
}

func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		m.lock.RLock()
		defer m.lock.RUnlock()

		for k := range m.items {
			if !yield(k) {
				return
			}
		}
	}
}

func (m *Map[K, V]) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.items)
}

func (m *Map[K, V]) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	maps.Clear(m.items)
}

func (m *Map[K, V]) Replace(values map[K]V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.items = values
}
