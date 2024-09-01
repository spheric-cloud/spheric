// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/registry/core/instancetype"
)

type InstanceTypeStorage struct {
	InstanceType *REST
}

type REST struct {
	*genericregistry.Store
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (InstanceTypeStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &core.InstanceType{}
		},
		NewListFunc: func() runtime.Object {
			return &core.InstanceTypeList{}
		},
		PredicateFunc:             instancetype.MatchInstanceType,
		DefaultQualifiedResource:  core.Resource("instancetypes"),
		SingularQualifiedResource: core.Resource("instancetype"),

		CreateStrategy: instancetype.Strategy,
		UpdateStrategy: instancetype.Strategy,
		DeleteStrategy: instancetype.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: instancetype.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return InstanceTypeStorage{}, err
	}

	return InstanceTypeStorage{
		InstanceType: &REST{store},
	}, nil
}
