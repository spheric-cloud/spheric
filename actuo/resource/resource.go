// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"context"
	"path"

	"spheric.cloud/spheric/actuo/meta"
	"spheric.cloud/spheric/actuo/runtime"
	"spheric.cloud/spheric/actuo/storage/store"
	"spheric.cloud/spheric/actuo/types"
	"spheric.cloud/spheric/actuo/watch"
	"spheric.cloud/spheric/utils/generic"
)

type Creater[Object any] interface {
	Create(ctx context.Context, obj Object) (Object, error)
}

type Getter[Key, Object any] interface {
	Get(ctx context.Context, key Key) (Object, error)
}

type Lister[Object any] interface {
	List(ctx context.Context) (runtime.List[Object], error)
}

type Deleter[Key, Object any] interface {
	Delete(ctx context.Context, key Key) (Object, error)
}

type Updater[Key, Object any] interface {
	Update(ctx context.Context, key Key, obj Object) (Object, error)
}

type Watcher[Object any] interface {
	Watch(ctx context.Context) (watch.Watch[Object], error)
}

type ObjectKeyer[Key, Object any] interface {
	ObjectKey(Object) Key
}

type globalMetaObjectKeyer[Object meta.Object] struct{}

func (globalMetaObjectKeyer[Object]) ObjectKey(obj Object) string {
	return obj.GetName()
}

func GlobalMetaObjectKeyer[Object meta.Object]() ObjectKeyer[string, Object] {
	return globalMetaObjectKeyer[Object]{}
}

type namespacedObjectKeyer[Object meta.Object] struct{}

func (namespacedObjectKeyer[Object]) ObjectKey(obj Object) types.NamespacedName {
	return types.NamespacedName{
		Namespace: obj.GetNamespace(),
	}
}

func NamespacedObjectKeyer[Object meta.Object]() ObjectKeyer[types.NamespacedName, Object] {
	return namespacedObjectKeyer[Object]{}
}

type StoreKeyer[StoreKey, Key any] interface {
	ObjectStoreKey(ctx context.Context, key Key) (StoreKey, error)
	ListStoreKey(ctx context.Context) (StoreKey, error)
}

type namespacedStoreKeyer struct {
	rootKey string
}

func NewNamespacedStoreKeyer(rootKey string) StoreKeyer[string, types.NamespacedName] {
	return &namespacedStoreKeyer{rootKey: rootKey}
}

func (k *namespacedStoreKeyer) ObjectStoreKey(ctx context.Context, key types.NamespacedName) (string, error) {
	return path.Join(k.rootKey, key.Namespace, key.Name), nil
}

func (k *namespacedStoreKeyer) ListStoreKey(ctx context.Context) (string, error) {
	return path.Join(k.rootKey, types.NamespaceValue(ctx)), nil
}

type SimpleStoreKeyer struct {
	rootKey    string
	contextKey func(ctx context.Context) (string, bool)
}

type SimpleStoreKeyerOptions struct {
	ContextKey func(ctx context.Context) (string, bool)
}

func (o *SimpleStoreKeyerOptions) ApplyOptions(opts []SimpleStoreKeyerOption) *SimpleStoreKeyerOptions {
	for _, opt := range opts {
		opt.ApplyToSimpleStoreKeyer(o)
	}
	return o
}

type SimpleStoreKeyerOption interface {
	ApplyToSimpleStoreKeyer(*SimpleStoreKeyerOptions)
}

type WithContextKey func(ctx context.Context) (context.Context, bool)

var _ StoreKeyer[string, string] = (*SimpleStoreKeyer)(nil)

func NewSimpleStoreKeyer(rootKey string, opts ...SimpleStoreKeyerOption) *SimpleStoreKeyer {
	o := (&SimpleStoreKeyerOptions{}).ApplyOptions(opts)
	return &SimpleStoreKeyer{
		rootKey:    rootKey,
		contextKey: o.ContextKey,
	}
}

func (s *SimpleStoreKeyer) ObjectStoreKey(ctx context.Context, key string) (string, error) {
	if s.contextKey != nil {
		if ctxKey, ok := s.contextKey(ctx); ok {
			return path.Join(s.rootKey, key, ctxKey), nil
		}
	}
	return path.Join(s.rootKey, key), nil
}

func (s *SimpleStoreKeyer) ListStoreKey(ctx context.Context) (string, error) {
	if s.contextKey != nil {
		if ctxKey, ok := s.contextKey(ctx); ok {
			return path.Join(s.rootKey, ctxKey), nil
		}
	}
	return s.rootKey, nil
}

type creater[StoreKey, Key, Object any] struct {
	storeKeyer StoreKeyer[StoreKey, Key]
	keyer      ObjectKeyer[Key, Object]
	store      store.Store[StoreKey, Object]
}

func NewCreater[StoreKey, Key, Object any](
	storeKeyer StoreKeyer[StoreKey, Key],
	keyer ObjectKeyer[Key, Object],
	store store.Store[StoreKey, Object],
) Creater[Object] {
	return &creater[StoreKey, Key, Object]{
		storeKeyer: storeKeyer,
		keyer:      keyer,
		store:      store,
	}
}

func (c *creater[StoreKey, Key, Object]) Create(ctx context.Context, obj Object) (Object, error) {
	key := c.keyer.ObjectKey(obj)
	storeKey, err := c.storeKeyer.ObjectStoreKey(ctx, key)
	if err != nil {
		return generic.Zero[Object](), err
	}

	created, err := c.store.Create(ctx, storeKey, obj)
	if err != nil {
		return generic.Zero[Object](), err
	}

	return created, nil
}

type deleter[StoreKey, Key, Object any] struct {
	storeKeyer StoreKeyer[StoreKey, Key]
	store      store.Store[StoreKey, Object]
}

func NewDeleter[StoreKey, Key, Object any](
	storeKeyer StoreKeyer[StoreKey, Key],
	store store.Store[StoreKey, Object],
) Deleter[Key, Object] {
	return &deleter[StoreKey, Key, Object]{
		storeKeyer: storeKeyer,
		store:      store,
	}
}

func (d *deleter[StoreKey, Key, Object]) Delete(ctx context.Context, key Key) (Object, error) {
	storeKey, err := d.storeKeyer.ObjectStoreKey(ctx, key)
	if err != nil {
		return generic.Zero[Object](), err
	}

	return d.store.Delete(ctx, storeKey, func(ctx context.Context, oldObj Object) error {
		return nil
	})
}

type getter[StoreKey, Key, Object any] struct {
	storeKeyer StoreKeyer[StoreKey, Key]
	store      store.Store[StoreKey, Object]
}

func NewGetter[StoreKey, Key, Object any](
	storeKeyer StoreKeyer[StoreKey, Key],
	store store.Store[StoreKey, Object],
) Getter[Key, Object] {
	return &getter[StoreKey, Key, Object]{
		storeKeyer: storeKeyer,
		store:      store,
	}
}

func (g *getter[StoreKey, Key, Object]) Get(ctx context.Context, key Key) (Object, error) {
	storeKey, err := g.storeKeyer.ObjectStoreKey(ctx, key)
	if err != nil {
		return generic.Zero[Object](), err
	}

	return g.store.Get(ctx, storeKey)
}

type lister[StoreKey, Key, Object any] struct {
	storeKeyer StoreKeyer[StoreKey, Key]
	store      store.Store[StoreKey, Object]
}

func NewLister[StoreKey, Key, Object any](
	storeKeyer StoreKeyer[StoreKey, Key],
	store store.Store[StoreKey, Object],
) Lister[Object] {
	return &lister[StoreKey, Key, Object]{
		storeKeyer: storeKeyer,
		store:      store,
	}
}

func (l *lister[StoreKey, Key, Object]) List(ctx context.Context) (runtime.List[Object], error) {
	listStoreKey, err := l.storeKeyer.ListStoreKey(ctx)
	if err != nil {
		return nil, err
	}
	return l.store.List(ctx, listStoreKey)
}

type updater[StoreKey, Key, Object any] struct {
	storeKeyer StoreKeyer[StoreKey, Key]
	store      store.Store[StoreKey, Object]
}

func NewUpdater[StoreKey, Key, Object any](
	storeKeyer StoreKeyer[StoreKey, Key],
	store store.Store[StoreKey, Object],
) Updater[Key, Object] {
	return &updater[StoreKey, Key, Object]{
		storeKeyer: storeKeyer,
		store:      store,
	}
}

func (u *updater[StoreKey, Key, Object]) Update(ctx context.Context, key Key, obj Object) (Object, error) {
	storeKey, err := u.storeKeyer.ObjectStoreKey(ctx, key)
	if err != nil {
		return generic.Zero[Object](), err
	}

	return u.store.Update(ctx, storeKey, false, func(ctx context.Context, oldObj Object) (newObj Object, err error) {
		return obj, nil
	})
}

type watcher[StoreKey, Key, Object any] struct {
	storeKeyer StoreKeyer[StoreKey, Key]
	store      store.Store[StoreKey, Object]
}

func NewWatcher[StoreKey, Key, Object any](
	storeKeyer StoreKeyer[StoreKey, Key],
	store store.Store[StoreKey, Object],
) Watcher[Object] {
	return &watcher[StoreKey, Key, Object]{
		storeKeyer: storeKeyer,
		store:      store,
	}
}

func (w *watcher[StoreKey, Key, Object]) Watch(ctx context.Context) (watch.Watch[Object], error) {
	listStoreKey, err := w.storeKeyer.ListStoreKey(ctx)
	if err != nil {
		return nil, err
	}
	return w.store.Watch(ctx, listStoreKey)
}
