// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"spheric.cloud/spheric/internal/apis/storage"
	"spheric.cloud/spheric/internal/registry/storage/bucketpool"
)

type REST struct {
	*genericregistry.Store
}

type BucketPoolStorage struct {
	BucketPool *REST
	Status     *StatusREST
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (BucketPoolStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &storage.BucketPool{}
		},
		NewListFunc: func() runtime.Object {
			return &storage.BucketPoolList{}
		},
		PredicateFunc:             bucketpool.MatchBucketPool,
		DefaultQualifiedResource:  storage.Resource("bucketpools"),
		SingularQualifiedResource: storage.Resource("bucketpool"),

		CreateStrategy: bucketpool.Strategy,
		UpdateStrategy: bucketpool.Strategy,
		DeleteStrategy: bucketpool.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: bucketpool.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return BucketPoolStorage{}, err
	}

	statusStore := *store
	statusStore.UpdateStrategy = bucketpool.StatusStrategy
	statusStore.ResetFieldsStrategy = bucketpool.StatusStrategy

	return BucketPoolStorage{
		BucketPool: &REST{store},
		Status:     &StatusREST{&statusStore},
	}, nil
}

type StatusREST struct {
	store *genericregistry.Store
}

func (r *StatusREST) New() runtime.Object {
	return &storage.BucketPool{}
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
