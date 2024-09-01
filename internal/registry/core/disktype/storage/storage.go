// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/registry/core/disktype"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
)

type DiskTypeStorage struct {
	DiskType *REST
}

type REST struct {
	*genericregistry.Store
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (DiskTypeStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &core.DiskType{}
		},
		NewListFunc: func() runtime.Object {
			return &core.DiskTypeList{}
		},
		PredicateFunc:             disktype.MatchDiskType,
		DefaultQualifiedResource:  core.Resource("disktypes"),
		SingularQualifiedResource: core.Resource("disktype"),

		CreateStrategy: disktype.Strategy,
		UpdateStrategy: disktype.Strategy,
		DeleteStrategy: disktype.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{
		RESTOptions: optsGetter,
		AttrFunc:    disktype.GetAttrs,
		Indexers:    disktype.Indexers(),
	}
	if err := store.CompleteWithOptions(options); err != nil {
		return DiskTypeStorage{}, err
	}

	return DiskTypeStorage{
		DiskType: &REST{store},
	}, nil
}
