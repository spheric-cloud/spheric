// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"fmt"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/internal/registry/core/fleet"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/apis/core/v1alpha1"
	"spheric.cloud/spheric/internal/spherelet/client"
)

type REST struct {
	*genericregistry.Store
}

type FleetStorage struct {
	Fleet                   *REST
	Status                  *StatusREST
	SphereletConnectionInfo client.ConnectionInfoGetter
}

func NewStorage(optsGetter generic.RESTOptionsGetter, sphereletClientConfig client.SphereletClientConfig) (FleetStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &core.Fleet{}
		},
		NewListFunc: func() runtime.Object {
			return &core.FleetList{}
		},
		PredicateFunc:             fleet.MatchFleet,
		DefaultQualifiedResource:  core.Resource("fleets"),
		SingularQualifiedResource: core.Resource("fleet"),

		CreateStrategy: fleet.Strategy,
		UpdateStrategy: fleet.Strategy,
		DeleteStrategy: fleet.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: fleet.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return FleetStorage{}, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = fleet.StatusStrategy
	statusStore.ResetFieldsStrategy = fleet.StatusStrategy

	fleetRest := &REST{store}
	statusRest := &StatusREST{&statusStore}

	// Build a FleetGetter that looks up nodes using the REST handler
	fleetGetter := client.FleetGetterFunc(func(ctx context.Context, fleetName string, options metav1.GetOptions) (*corev1alpha1.Fleet, error) {
		obj, err := fleetRest.Get(ctx, fleetName, &options)
		if err != nil {
			return nil, err
		}
		fleet, ok := obj.(*core.Fleet)
		if !ok {
			return nil, fmt.Errorf("unexpected type %T", obj)
		}
		// TODO: Remove the conversion. Consider only return the FleetAddresses
		externalFleet := &corev1alpha1.Fleet{}
		if err := v1alpha1.Convert_core_Fleet_To_v1alpha1_Fleet(fleet, externalFleet, nil); err != nil {
			return nil, fmt.Errorf("failed to convert to v1alpha1.Fleet: %v", err)
		}
		return externalFleet, nil
	})
	connectionInfoGetter, err := client.NewFleetConnectionInfoGetter(fleetGetter, sphereletClientConfig)
	if err != nil {
		return FleetStorage{}, err
	}

	return FleetStorage{
		Fleet:                   fleetRest,
		Status:                  statusRest,
		SphereletConnectionInfo: connectionInfoGetter,
	}, nil
}

type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &core.Fleet{}
}

func (r *StatusREST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	return r.store.Get(ctx, name, options)
}

func (r *StatusREST) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	return r.store.Update(ctx, name, objInfo, createValidation, updateValidation, false, options)
}

func (r *StatusREST) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return r.store.GetResetFields()
}

func (r *StatusREST) Destroy() {}
