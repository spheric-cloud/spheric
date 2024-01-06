// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"spheric.cloud/spheric/internal/apis/storage"
	registrybucketclass "spheric.cloud/spheric/internal/registry/storage/bucketclass"
)

type BucketClassStorage struct {
	BucketClass *REST
}

type REST struct {
	*genericregistry.Store
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (BucketClassStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &storage.BucketClass{}
		},
		NewListFunc: func() runtime.Object {
			return &storage.BucketClassList{}
		},
		PredicateFunc:             registrybucketclass.MatchBucketClass,
		DefaultQualifiedResource:  storage.Resource("bucketclasses"),
		SingularQualifiedResource: storage.Resource("bucketclass"),

		CreateStrategy: registrybucketclass.Strategy,
		UpdateStrategy: registrybucketclass.Strategy,
		DeleteStrategy: registrybucketclass.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: registrybucketclass.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return BucketClassStorage{}, err
	}
	return BucketClassStorage{
		BucketClass: &REST{store},
	}, nil
}
