// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"
	"errors"
	"iter"
)

var ErrNotFound = errors.New("not found")

type ReflectorSink[Object any] interface {
	Populate(ctx context.Context, objs iter.Seq[Object]) error
	Add(ctx context.Context, obj Object) error
	Update(ctx context.Context, obj Object) error
	Delete(ctx context.Context, obj Object) error
}

type Store[Key, Object any] interface {
	Set(obj Object) error
	Delete(key Key) error
	Get(key Key) (Object, error)
	All() iter.Seq[Object]
}

type KeyFunc[Key, Object any] func(Object) (Key, error)
