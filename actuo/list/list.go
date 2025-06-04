// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package list

import (
	"iter"

	"spheric.cloud/spheric/actuo/meta"
	"spheric.cloud/spheric/actuo/runtime"
)

type List[O interface {
	runtime.Object
	*ObjectVal
}, ObjectVal any] struct {
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []ObjectVal
}

func New[O interface {
	runtime.Object
	*ObjectVal
}, ObjectVal any](len int) *List[O, ObjectVal] {
	return &List[O, ObjectVal]{
		Items: make([]ObjectVal, len),
	}
}

func (l *List[O, ObjectVal]) Item(idx int) O {
	return O(&l.Items[idx])
}

func (l *List[O, ObjectVal]) All() iter.Seq[O] {
	return func(yield func(O) bool) {
		for i := range l.Items {
			if !yield(O(&l.Items[i])) {
				return
			}
		}
	}
}

func (l *List[O, ObjectVal]) Len() int {
	return len(l.Items)
}
