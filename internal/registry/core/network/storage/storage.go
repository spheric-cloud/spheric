// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	"spheric.cloud/spheric/internal/registry/core/network"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/apis/core"
)

type NetworkStorage struct {
	Network *REST
	Status  *StatusREST
}

type REST struct {
	*genericregistry.Store
}

func (REST) ShortNames() []string {
	return []string{"net"}
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (NetworkStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &core.Network{}
		},
		NewListFunc: func() runtime.Object {
			return &core.NetworkList{}
		},
		PredicateFunc:             network.MatchNetwork,
		DefaultQualifiedResource:  core.Resource("networks"),
		SingularQualifiedResource: core.Resource("network"),

		CreateStrategy: network.Strategy,
		UpdateStrategy: network.Strategy,
		DeleteStrategy: network.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: network.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return NetworkStorage{}, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = network.StatusStrategy
	statusStore.ResetFieldsStrategy = network.StatusStrategy

	return NetworkStorage{
		Network: &REST{store},
		Status:  &StatusREST{&statusStore},
	}, nil
}

type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &core.Network{}
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
