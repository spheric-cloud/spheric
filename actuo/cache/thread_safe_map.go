// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"iter"

	"spheric.cloud/spheric/utils/sync"
)

type Cache[Key comparable, Object any] struct {
	keyFunc KeyFunc[Key, Object]
	store   *sync.Map[Key, Object]
}

func NewCache[Key comparable, Object any](keyFunc KeyFunc[Key, Object]) *Cache[Key, Object] {
	return &Cache[Key, Object]{
		keyFunc: keyFunc,
		store:   sync.NewMap[Key, Object](),
	}
}

func (c *Cache[Key, Object]) Set(obj Object) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return err
	}

	c.store.Set(key, obj)
	return nil
}

func (c *Cache[Key, Object]) Delete(key Key) error {
	c.store.Delete(key)
	return nil
}

func (c *Cache[Key, Object]) Get(key Key) (Object, error) {
	obj, ok := c.store.GetOK(key)
	if !ok {
		return obj, ErrNotFound
	}
	return obj, nil
}

func (c *Cache[Key, Object]) All() iter.Seq[Object] {
	return c.store.Values()
}
