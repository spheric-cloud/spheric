// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instanceevent

import sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"

type CreateEvent struct {
	Object *sri.Instance
}

type UpdateEvent struct {
	ObjectOld *sri.Instance
	ObjectNew *sri.Instance
}

type DeleteEvent struct {
	Object *sri.Instance
}

type GenericEvent struct {
	Object *sri.Instance
}

type Handler interface {
	Create(event CreateEvent)
	Update(event UpdateEvent)
	Delete(event DeleteEvent)
	Generic(event GenericEvent)
}

type HandlerFuncs struct {
	CreateFunc  func(event CreateEvent)
	UpdateFunc  func(event UpdateEvent)
	DeleteFunc  func(event DeleteEvent)
	GenericFunc func(event GenericEvent)
}

func (e HandlerFuncs) Create(event CreateEvent) {
	if e.CreateFunc != nil {
		e.CreateFunc(event)
	}
}

func (e HandlerFuncs) Update(event UpdateEvent) {
	if e.UpdateFunc != nil {
		e.UpdateFunc(event)
	}
}

func (e HandlerFuncs) Delete(event DeleteEvent) {
	if e.DeleteFunc != nil {
		e.DeleteFunc(event)
	}
}

func (e HandlerFuncs) Generic(event GenericEvent) {
	if e.GenericFunc != nil {
		e.GenericFunc(event)
	}
}

type HandlerRegistration interface {
	Remove() error
}

type Source interface {
	AddHandler(handler Handler) (HandlerRegistration, error)
}
