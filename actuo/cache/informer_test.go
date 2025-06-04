// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache_test

import (
	"context"
	"slices"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "spheric.cloud/spheric/actuo/cache"
	watch2 "spheric.cloud/spheric/actuo/watch"
)

func IdentityKeyFunc[ObjectAndKey any](objAndKey ObjectAndKey) (ObjectAndKey, error) {
	return objAndKey, nil
}

type fakeListerWatcher[Object any] struct {
	list  []Object
	watch watch2.Watch[Object]
}

func (f *fakeListerWatcher[Object]) List(ctx context.Context) ([]Object, error) {
	return f.list, nil
}

func (f *fakeListerWatcher[Object]) Watch(ctx context.Context) (watch2.Watch[Object], error) {
	return f.watch, nil
}

type fakeResourceEvent[Object any] struct {
	add           Object
	isInitialList bool
	update        Object
	delete        Object
}

type fakeResourceEventHandler[Object any] struct {
	mu     sync.Mutex
	events []fakeResourceEvent[Object]
}

func (h *fakeResourceEventHandler[Object]) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = h.events[:0]
}

func (h *fakeResourceEventHandler[Object]) add(event fakeResourceEvent[Object]) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, event)
}

func (h *fakeResourceEventHandler[Object]) OnAdd(obj Object, isInitialList bool) {
	h.add(fakeResourceEvent[Object]{add: obj, isInitialList: isInitialList})
}

func (h *fakeResourceEventHandler[Object]) OnUpdate(oldObj, newObj Object) {
	h.add(fakeResourceEvent[Object]{update: oldObj})
}

func (h *fakeResourceEventHandler[Object]) OnDelete(obj Object) {
	h.add(fakeResourceEvent[Object]{delete: obj})
}

func (h *fakeResourceEventHandler[Object]) Events() []fakeResourceEvent[Object] {
	h.mu.Lock()
	defer h.mu.Unlock()
	return slices.Clone(h.events)
}

var _ = Describe("Informer", func() {
	Context("SharedInformer", func() {
		Describe("AddEventHandler", func() {
			It("should add an event handler to the informer", MustPassRepeatedly(100), func(ctx context.Context) {
				w := watch2.NewFakeWithCap[string](10)
				lw := &fakeListerWatcher[string]{
					list:  []string{"foo", "bar"},
					watch: w,
				}

				sharedInformer := NewSharedInformer(IdentityKeyFunc, lw, SharedInformerOptions{
					Logger: GinkgoLogr.WithName("shared-informer"),
				})

				ctx, cancel := context.WithCancel(ctx)
				defer cancel()
				done := make(chan error, 1)
				go func() {
					done <- sharedInformer.Run(ctx)
				}()

				By("waiting until the informer is synced")
				Eventually(sharedInformer.Synced()).Should(BeClosed())

				By("adding a handler to the informer")
				handler := &fakeResourceEventHandler[string]{}
				handle, err := sharedInformer.AddEventHandler(handler, EventHandlerOptions{})
				Expect(err).NotTo(HaveOccurred())

				By("waiting for the handler to be synced")
				Eventually(handle.Synced()).Should(BeClosed())

				By("waiting for the event handler to be called")
				Eventually(handler.Events).Should(ConsistOf(
					fakeResourceEvent[string]{add: "foo", isInitialList: true},
					fakeResourceEvent[string]{add: "bar", isInitialList: true},
				))

				By("clearing the handler events")
				handler.Clear()

				By("adding sme watch events")
				w.Create("baz")
				w.Delete("foo")
				w.Update("bar")

				By("waiting for the event handler to be called")
				Eventually(handler.Events).Should(ConsistOf(
					fakeResourceEvent[string]{add: "baz"},
					fakeResourceEvent[string]{delete: "foo"},
					fakeResourceEvent[string]{update: "bar"},
				))

				By("cancelling the informer")
				cancel()

				By("waiting for the informer to stop without error")
				Eventually(done).Should(Receive(BeNil()))
			})
		})
	})
})
