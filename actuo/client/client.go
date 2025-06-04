// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"spheric.cloud/spheric/actuo/cache"
	"spheric.cloud/spheric/actuo/storage/store"
)

type Client[Key, Object any] interface {
	Get(ctx context.Context, key Key) (Object, error)
}

type client[Key comparable, Object any] struct {
	store    store.Store[Key, Object]
	informer cache.SharedInformer[Key, Object]
}
