// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"errors"
	"strconv"

	"spheric.cloud/spheric/actuo/list"
	"spheric.cloud/spheric/actuo/meta"
	"spheric.cloud/spheric/actuo/runtime"
	"spheric.cloud/spheric/actuo/watch"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// Update specifies how the existing object should be updated.
type Update[Object any] func(ctx context.Context, oldObj Object) (newObj Object, err error)

type Delete[Object any] func(ctx context.Context, oldObj Object) error

type Versioner[Object any, ObjectList runtime.List[Object]] interface {
	// UpdateObject sets storage metadata into an API object. Returns an error if the object
	// cannot be updated correctly. May return nil if the requested object does not need metadata
	// from database.
	UpdateObject(obj Object, resourceVersion uint64) error
	// ObjectResourceVersion returns the resource version (for persistence) of the specified object.
	// Should return an error if the specified object does not have a persistable version.
	ObjectResourceVersion(obj Object) (uint64, error)
	// PrepareObjectForStorage should set SelfLink and ResourceVersion to the empty value. Should
	// return an error if the specified object cannot be updated.
	PrepareObjectForStorage(obj Object) error
	// UpdateList sets the resource version into a list object.
	UpdateList(l ObjectList, resourceVersion uint64) error
}

type metaVersioner[Object interface {
	meta.Object
	runtime.Object
}, ObjectList interface {
	meta.List
	runtime.List[Object]
}] struct{}

func (metaVersioner[Object, ObjectList]) UpdateObject(obj Object, resourceVersion uint64) error {
	obj.SetResourceVersion(strconv.FormatUint(resourceVersion, 10))
	return nil
}

func (metaVersioner[Object, ObjectList]) ObjectResourceVersion(obj Object) (uint64, error) {
	return strconv.ParseUint(obj.GetResourceVersion(), 10, 64)
}

func (metaVersioner[Object, ObjectList]) PrepareObjectForStorage(obj Object) error {
	return nil // TODO: Do something here?
}

func (metaVersioner[Object, ObjectList]) UpdateList(l ObjectList, resourceVersion uint64) error {
	l.SetResourceVersion(strconv.FormatUint(resourceVersion, 10))
	return nil
}

func MetaVersioner[Object interface {
	meta.Object
	runtime.Object
}, ObjectList interface {
	meta.List
	runtime.List[Object]
}]() Versioner[Object, ObjectList] {
	return metaVersioner[Object, ObjectList]{}
}

func DefaultMetaVersioner[Object interface {
	meta.Object
	runtime.Object
	*ObjectVal
}, ObjectVal any]() Versioner[Object, *list.List[Object, ObjectVal]] {
	return MetaVersioner[Object, *list.List[Object, ObjectVal]]()
}

type nopVersioner[Object any, ObjectList runtime.List[Object]] struct{}

func (nopVersioner[Object, ObjectList]) UpdateObject(obj Object, resourceVersion uint64) error {
	return nil
}
func (nopVersioner[Object, ObjectList]) ObjectResourceVersion(obj Object) (uint64, error) {
	return 0, nil
}
func (nopVersioner[Object, ObjectList]) PrepareObjectForStorage(obj Object) error { return nil }
func (nopVersioner[Object, ObjectList]) UpdateList(l ObjectList, resourceVersion uint64) error {
	return nil
}

func NopVersioner[Object runtime.Object, ObjectList runtime.List[Object]]() Versioner[Object, ObjectList] {
	return nopVersioner[Object, ObjectList]{}
}

// Store is a generic key-value store of an object type.
// It forms the 'low-level' operations on which later on API-server-like storage is implemented.
type Store[Key, Object any] interface {
	// Create adds the object at the given Key to the store unless it already exists.
	Create(ctx context.Context, k Key, obj Object) (Object, error)
	// Get retrieves the object from the store if it exists.
	Get(ctx context.Context, k Key) (Object, error)
	Update(ctx context.Context, k Key, ignoreNotFound bool, update Update[Object]) (Object, error)
	Delete(ctx context.Context, k Key, del Delete[Object]) (Object, error)
	List(ctx context.Context, k Key) (runtime.List[Object], error)
	Watch(ctx context.Context, k Key) (watch.Watch[Object], error)
}

type Factory[Object any, ObjectList runtime.List[Object]] interface {
	New() Object
	NewList(len int) ObjectList
}

type defaultFactory[Object interface {
	*ObjectVal
	runtime.Object
}, ObjectVal any] struct{}

func DefaultFactory[Object interface {
	*ObjectVal
	runtime.Object
}, ObjectVal any]() Factory[Object, *list.List[Object, ObjectVal]] {
	return defaultFactory[Object, ObjectVal]{}
}

func (defaultFactory[Object, ObjectVal]) New() Object {
	return new(ObjectVal)
}

func (defaultFactory[Object, ObjectVal]) NewList(len int) *list.List[Object, ObjectVal] {
	return list.New[Object, ObjectVal](len)
}
