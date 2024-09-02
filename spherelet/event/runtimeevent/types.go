// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package runtimeevent

import iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"

type UpdateResourcesEvent struct {
	ResourcesOld *iri.RuntimeResources
	ResourcesNew *iri.RuntimeResources
}

type Handler interface {
	UpdateResources(event *UpdateResourcesEvent)
}

type HandlerFuncs struct {
	UpdateResourcesFunc func(event *UpdateResourcesEvent)
}

func (f HandlerFuncs) UpdateResources(event *UpdateResourcesEvent) {
	if f.UpdateResourcesFunc != nil {
		f.UpdateResourcesFunc(event)
	}
}

type HandlerRegistration interface {
	Remove() error
}

type Source interface {
	AddHandler(handler Handler) (HandlerRegistration, error)
}
