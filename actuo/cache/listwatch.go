// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"

	"spheric.cloud/spheric/actuo/watch"
)

type ListerWatcher[Object any] interface {
	List(ctx context.Context) ([]Object, error)
	Watch(context.Context) (watch.Watch[Object], error)
}
