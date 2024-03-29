// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"spheric.cloud/spheric/internal/apis/compute"
	"spheric.cloud/spheric/internal/registry/compute/machineclass"
)

type MachineClassStorage struct {
	MachineClass *REST
}

type REST struct {
	*genericregistry.Store
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (MachineClassStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &compute.MachineClass{}
		},
		NewListFunc: func() runtime.Object {
			return &compute.MachineClassList{}
		},
		PredicateFunc:             machineclass.MatchMachineClass,
		DefaultQualifiedResource:  compute.Resource("machineclasses"),
		SingularQualifiedResource: compute.Resource("machineclass"),

		CreateStrategy: machineclass.Strategy,
		UpdateStrategy: machineclass.Strategy,
		DeleteStrategy: machineclass.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: machineclass.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return MachineClassStorage{}, err
	}

	return MachineClassStorage{
		MachineClass: &REST{store},
	}, nil
}
