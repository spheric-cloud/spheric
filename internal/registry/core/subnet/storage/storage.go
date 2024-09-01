// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	"spheric.cloud/spheric/internal/registry/core/subnet"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/apis/core"
)

type SubnetStorage struct {
	Subnet *REST
	Status *StatusREST
}

type REST struct {
	*genericregistry.Store
}

func (REST) ShortNames() []string {
	return []string{"net"}
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (SubnetStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &core.Subnet{}
		},
		NewListFunc: func() runtime.Object {
			return &core.SubnetList{}
		},
		PredicateFunc:             subnet.MatchSubnet,
		DefaultQualifiedResource:  core.Resource("subnets"),
		SingularQualifiedResource: core.Resource("subnet"),

		CreateStrategy: subnet.Strategy,
		UpdateStrategy: subnet.Strategy,
		DeleteStrategy: subnet.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: subnet.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return SubnetStorage{}, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = subnet.StatusStrategy
	statusStore.ResetFieldsStrategy = subnet.StatusStrategy

	return SubnetStorage{
		Subnet: &REST{store},
		Status: &StatusREST{&statusStore},
	}, nil
}

type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &core.Subnet{}
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
