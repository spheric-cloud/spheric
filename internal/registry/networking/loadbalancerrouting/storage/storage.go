// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"spheric.cloud/spheric/internal/apis/networking"
	"spheric.cloud/spheric/internal/registry/networking/loadbalancerrouting"
)

type LoadBalancerRoutingStorage struct {
	LoadBalancerRouting *REST
}

type REST struct {
	*genericregistry.Store
}

func (REST) ShortNames() []string {
	return []string{"lbr", "lbrt"}
}

func NewStorage(optsGetter generic.RESTOptionsGetter) (LoadBalancerRoutingStorage, error) {
	store := &genericregistry.Store{
		NewFunc: func() runtime.Object {
			return &networking.LoadBalancerRouting{}
		},
		NewListFunc: func() runtime.Object {
			return &networking.LoadBalancerRoutingList{}
		},
		PredicateFunc:             loadbalancerrouting.MatchLoadBalancerRouting,
		DefaultQualifiedResource:  networking.Resource("loadbalancerroutings"),
		SingularQualifiedResource: networking.Resource("loadbalancerrouting"),

		CreateStrategy: loadbalancerrouting.Strategy,
		UpdateStrategy: loadbalancerrouting.Strategy,
		DeleteStrategy: loadbalancerrouting.Strategy,

		TableConvertor: newTableConvertor(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: loadbalancerrouting.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return LoadBalancerRoutingStorage{}, err
	}

	return LoadBalancerRoutingStorage{
		LoadBalancerRouting: &REST{store},
	}, nil
}
