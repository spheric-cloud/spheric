// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache_test

import (
	"context"
	"iter"
	"slices"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "spheric.cloud/spheric/actuo/cache"
	"spheric.cloud/spheric/actuo/watch"
)

type fakeReflectorSinkItem[Object any] struct {
	populate []Object
	add      Object
	update   Object
	delete   Object
}

type fakeReflectorSink[Object any] struct {
	mu    sync.Mutex
	items []fakeReflectorSinkItem[Object]
}

func (s *fakeReflectorSink[Object]) add(item fakeReflectorSinkItem[Object]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = append(s.items, item)
}

func (s *fakeReflectorSink[Object]) Populate(ctx context.Context, objs iter.Seq[Object]) error {
	s.add(fakeReflectorSinkItem[Object]{
		populate: slices.Collect(objs),
	})
	return nil
}

func (s *fakeReflectorSink[Object]) Add(ctx context.Context, obj Object) error {
	s.add(fakeReflectorSinkItem[Object]{
		add: obj,
	})
	return nil
}

func (s *fakeReflectorSink[Object]) Update(ctx context.Context, obj Object) error {
	s.add(fakeReflectorSinkItem[Object]{
		update: obj,
	})
	return nil
}

func (s *fakeReflectorSink[Object]) Delete(ctx context.Context, obj Object) error {
	s.add(fakeReflectorSinkItem[Object]{
		delete: obj,
	})
	return nil
}

func (s *fakeReflectorSink[Object]) Items() []fakeReflectorSinkItem[Object] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return slices.Clone(s.items)
}

var _ = Describe("Reflector", func() {
	It("should reflect the events", func(ctx context.Context) {
		w := watch.NewFakeWithCap[string](10)
		lw := &fakeListerWatcher[string]{
			list:  []string{"foo", "bar"},
			watch: w,
		}
		sink := &fakeReflectorSink[string]{}

		w.Create("baz")
		w.Delete("bang")
		w.Update("qux")

		r := NewReflector(lw, sink, ReflectorOptions{})
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		done := make(chan error)
		go func() {
			defer close(done)
			done <- r.Run(ctx)
		}()

		Eventually(sink.Items).Should(Equal([]fakeReflectorSinkItem[string]{
			{populate: []string{"foo", "bar"}},
			{add: "baz"},
			{delete: "bang"},
			{update: "qux"},
		}))

		cancel()
		Eventually(done).Should(Receive(BeNil()))
	})
})
