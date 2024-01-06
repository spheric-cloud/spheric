// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package srievent

import (
	srimeta "spheric.cloud/spheric/sri/apis/meta/v1alpha1"
)

type CreateEvent[O srimeta.Object] struct {
	Object O
}

type UpdateEvent[O srimeta.Object] struct {
	ObjectOld O
	ObjectNew O
}

type DeleteEvent[O srimeta.Object] struct {
	Object O
}

type GenericEvent[O srimeta.Object] struct {
	Object O
}

type Handler[O srimeta.Object] interface {
	Create(event CreateEvent[O])
	Update(event UpdateEvent[O])
	Delete(event DeleteEvent[O])
	Generic(event GenericEvent[O])
}

type HandlerFuncs[O srimeta.Object] struct {
	CreateFunc  func(event CreateEvent[O])
	UpdateFunc  func(event UpdateEvent[O])
	DeleteFunc  func(event DeleteEvent[O])
	GenericFunc func(event GenericEvent[O])
}

func (e HandlerFuncs[O]) Create(event CreateEvent[O]) {
	if e.CreateFunc != nil {
		e.CreateFunc(event)
	}
}

func (e HandlerFuncs[O]) Update(event UpdateEvent[O]) {
	if e.UpdateFunc != nil {
		e.UpdateFunc(event)
	}
}

func (e HandlerFuncs[O]) Delete(event DeleteEvent[O]) {
	if e.DeleteFunc != nil {
		e.DeleteFunc(event)
	}
}

func (e HandlerFuncs[O]) Generic(event GenericEvent[O]) {
	if e.GenericFunc != nil {
		e.GenericFunc(event)
	}
}

type HandlerRegistration interface {
	Remove() error
}

type Source[O srimeta.Object] interface {
	AddHandler(handler Handler[O]) (HandlerRegistration, error)
}
