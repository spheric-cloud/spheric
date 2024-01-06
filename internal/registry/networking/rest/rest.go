// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/networking"
	loadbalancerstorage "spheric.cloud/spheric/internal/registry/networking/loadbalancer/storage"
	loadbalancerroutingstorage "spheric.cloud/spheric/internal/registry/networking/loadbalancerrouting/storage"
	natgatewaystorage "spheric.cloud/spheric/internal/registry/networking/natgateway/storage"
	networkstorage "spheric.cloud/spheric/internal/registry/networking/network/storage"
	networkinterfacestorage "spheric.cloud/spheric/internal/registry/networking/networkinterface/storage"
	networkpolicystorage "spheric.cloud/spheric/internal/registry/networking/networkpolicy/storage"
	virtualipstorage "spheric.cloud/spheric/internal/registry/networking/virtualip/storage"
	sphericserializer "spheric.cloud/spheric/internal/serializer"

	"k8s.io/apiserver/pkg/server/storage"
)

type StorageProvider struct{}

func (p StorageProvider) GroupName() string {
	return networking.SchemeGroupVersion.Group
}

func (p StorageProvider) NewRESTStorage(apiResourceConfigSource storage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(p.GroupName(), api.Scheme, metav1.ParameterCodec, api.Codecs)
	apiGroupInfo.PrioritizedVersions = []schema.GroupVersion{networkingv1alpha1.SchemeGroupVersion}
	apiGroupInfo.NegotiatedSerializer = sphericserializer.DefaultSubsetNegotiatedSerializer(api.Codecs)

	storageMap, err := p.v1alpha1Storage(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, false, err
	}

	apiGroupInfo.VersionedResourcesStorageMap[networkingv1alpha1.SchemeGroupVersion.Version] = storageMap

	return apiGroupInfo, true, nil
}

func (p StorageProvider) v1alpha1Storage(restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storageMap := map[string]rest.Storage{}

	networkInterfaceStorage, err := networkinterfacestorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["networkinterfaces"] = networkInterfaceStorage.NetworkInterface
	storageMap["networkinterfaces/status"] = networkInterfaceStorage.Status

	networkStorage, err := networkstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["networks"] = networkStorage.Network
	storageMap["networks/status"] = networkStorage.Status

	networkPolicyStorage, err := networkpolicystorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["networkpolicies"] = networkPolicyStorage.NetworkPolicy
	storageMap["networkpolicies/status"] = networkPolicyStorage.Status

	virtualIPStorage, err := virtualipstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["virtualips"] = virtualIPStorage.VirtualIP
	storageMap["virtualips/status"] = virtualIPStorage.Status

	loadBalancerStorage, err := loadbalancerstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["loadbalancers"] = loadBalancerStorage.LoadBalancer
	storageMap["loadbalancers/status"] = loadBalancerStorage.Status

	loadBalancerRoutingStorage, err := loadbalancerroutingstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["loadbalancerroutings"] = loadBalancerRoutingStorage.LoadBalancerRouting

	natGatewayStorage, err := natgatewaystorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["natgateways"] = natGatewayStorage.NATGateway
	storageMap["natgateways/status"] = natGatewayStorage.Status

	return storageMap, nil
}
