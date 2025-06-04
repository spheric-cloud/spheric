// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"fmt"
	"net/http"
	"path"

	ahttp "spheric.cloud/spheric/actuo/http/server"
	"spheric.cloud/spheric/actuo/meta"
	"spheric.cloud/spheric/actuo/resource"
	"spheric.cloud/spheric/actuo/storage/store"
	"spheric.cloud/spheric/actuo/types"
)

type createOp[Key, Object any] struct {
}

type getOp[Key, Object any] struct {
}

type listOp[Object any] struct {
}

type updateOp[Key, Object any] struct {
}

type deleteOp[Key, Object any] struct {
}

type watchOp[Object any] struct {
}

type Builder[StoreKey, Key, Object any] struct {
	namespacePath string
	resource      string
	itemPath      string
	requestKeyer  ahttp.RequestKeyer[Key]
	objectKeyer   resource.ObjectKeyer[Key, Object]
	storeKeyer    resource.StoreKeyer[StoreKey, Key]
	objectFactory ahttp.ObjectFactory[Object]
	store         store.Store[StoreKey, Object]

	create *createOp[StoreKey, Object]
	get    *getOp[StoreKey, Object]
	list   *listOp[Object]
	update *updateOp[StoreKey, Object]
	delete *deleteOp[StoreKey, Object]
	watch  *watchOp[Object]

	creater resource.Creater[Object]
	getter  resource.Getter[Key, Object]
	lister  resource.Lister[Object]
	updater resource.Updater[Key, Object]
	deleter resource.Deleter[Key, Object]
	watcher resource.Watcher[Object]
}

func New[StoreKey, Key, Object any](
	namespacePath string,
	resource string,
	itemPath string,
	requestKeyer ahttp.RequestKeyer[Key],
	objectKeyer resource.ObjectKeyer[Key, Object],
	storeKeyer resource.StoreKeyer[StoreKey, Key],
	objectFactory ahttp.ObjectFactory[Object],
	store store.Store[StoreKey, Object],
) *Builder[StoreKey, Key, Object] {
	return &Builder[StoreKey, Key, Object]{
		namespacePath: namespacePath,
		resource:      resource,
		itemPath:      itemPath,
		requestKeyer:  requestKeyer,
		objectKeyer:   objectKeyer,
		storeKeyer:    storeKeyer,
		objectFactory: objectFactory,
		store:         store,
	}
}

func NewNamespaced[Object any](
	res string,
	objectKeyer resource.ObjectKeyer[types.NamespacedName, Object],
	objectFactory ahttp.ObjectFactory[Object],
	store store.Store[string, Object],
) *Builder[string, types.NamespacedName, Object] {
	return New(
		ahttp.DefaultNamespacePath,
		res,
		ahttp.DefaultItemPath,
		ahttp.DefaultNamespacedNameRequestKeyer,
		objectKeyer,
		resource.NewNamespacedStoreKeyer(path.Join("/", res)),
		objectFactory,
		store,
	)
}

func NewNamespacedMeta[Object interface {
	*ObjectVal
	meta.Object
}, ObjectVal any](
	res string,
	store store.Store[string, Object],
) *Builder[string, types.NamespacedName, Object] {
	return NewNamespaced(
		res,
		resource.NamespacedObjectKeyer[Object](),
		ahttp.ObjectValFactory[Object, ObjectVal](),
		store,
	)
}

func (b *Builder[StoreKey, Key, Object]) Create() *Builder[StoreKey, Key, Object] {
	b.create = &createOp[StoreKey, Object]{}
	return b
}

func (b *Builder[StoreKey, Key, Object]) Get() *Builder[StoreKey, Key, Object] {
	b.get = &getOp[StoreKey, Object]{}
	return b
}

func (b *Builder[StoreKey, Key, Object]) List() *Builder[StoreKey, Key, Object] {
	b.list = &listOp[Object]{}
	return b
}

func (b *Builder[StoreKey, Key, Object]) Update() *Builder[StoreKey, Key, Object] {
	b.update = &updateOp[StoreKey, Object]{}
	return b
}

func (b *Builder[StoreKey, Key, Object]) Delete() *Builder[StoreKey, Key, Object] {
	b.delete = &deleteOp[StoreKey, Object]{}
	return b
}

func (b *Builder[StoreKey, Key, Object]) Watch() *Builder[StoreKey, Key, Object] {
	b.watch = &watchOp[Object]{}
	return b
}

func (b *Builder[StoreKey, Key, Object]) CRUD() *Builder[StoreKey, Key, Object] {
	return b.
		Create().
		Get().
		List().
		Update().
		Delete().
		Watch()
}

func (b *Builder[StoreKey, Key, Object]) Build() (http.Handler, error) {
	mux := ahttp.NewServeMux()

	if b.create != nil {
		b.creater = resource.NewCreater(b.storeKeyer, b.objectKeyer, b.store)
		pattern := path.Join("/", b.namespacePath, b.resource)
		handler := ahttp.CreateHandler(b.objectFactory, b.creater)
		mux.Handle(fmt.Sprintf("POST %s", pattern), handler)
	}

	if b.get != nil {
		b.getter = resource.NewGetter(b.storeKeyer, b.store)
		pattern := path.Join("/", b.namespacePath, b.resource, b.itemPath)
		handler := ahttp.GetHandler(b.requestKeyer, b.getter)
		mux.Handle(fmt.Sprintf("GET %s", pattern), handler)
	}

	if b.list != nil {
		b.lister = resource.NewLister(b.storeKeyer, b.store)
		handler := ahttp.ListHandler(b.requestKeyer, b.lister)
		patterns := []string{path.Join("/", b.namespacePath, b.resource)}
		if b.namespacePath != "" {
			patterns = append(patterns, path.Join("/", b.resource))
		}
		for _, pattern := range patterns {
			mux.Handle(fmt.Sprintf("GET %s", pattern), handler)
		}
	}

	if b.update != nil {
		b.updater = resource.NewUpdater(b.storeKeyer, b.store)
		pattern := path.Join("/", b.namespacePath, b.resource, b.itemPath)
		mux.Handle(fmt.Sprintf("PUT %s", pattern), ahttp.UpdateHandler(b.requestKeyer, b.objectFactory, b.updater))
	}

	if b.delete != nil {
		b.deleter = resource.NewDeleter(b.storeKeyer, b.store)
		pattern := path.Join("/", b.namespacePath, b.resource, b.itemPath)
		handler := ahttp.DeleteHandler(b.requestKeyer, b.deleter)
		mux.Handle(fmt.Sprintf("DELETE %s", pattern), handler)
	}

	if b.watch != nil {
		b.watcher = resource.NewWatcher(b.storeKeyer, b.store)
		handler := ahttp.WatchHandler(b.requestKeyer, b.watcher)
		patterns := []string{path.Join("/", b.namespacePath, b.resource)}
		if b.namespacePath != "" {
			patterns = append(patterns, path.Join("/", b.resource))
		}
		for _, pattern := range patterns {
			mux.Handle(fmt.Sprintf("WATCH %s", pattern), handler)
		}
	}

	return mux, nil
}
