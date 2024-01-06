// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/compute"
	machinepoolletclient "spheric.cloud/spheric/internal/machinepoollet/client"
	machinestorage "spheric.cloud/spheric/internal/registry/compute/machine/storage"
	machineclassstore "spheric.cloud/spheric/internal/registry/compute/machineclass/storage"
	machinepoolstorage "spheric.cloud/spheric/internal/registry/compute/machinepool/storage"
	sphericserializer "spheric.cloud/spheric/internal/serializer"

	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

type StorageProvider struct {
	MachinePoolletClientConfig machinepoolletclient.MachinePoolletClientConfig
}

func (p StorageProvider) GroupName() string {
	return compute.SchemeGroupVersion.Group
}

func (p StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(p.GroupName(), api.Scheme, api.ParameterCodec, api.Codecs)
	apiGroupInfo.NegotiatedSerializer = sphericserializer.DefaultSubsetNegotiatedSerializer(api.Codecs)

	storageMap, err := p.v1alpha1Storage(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, false, err
	}

	apiGroupInfo.VersionedResourcesStorageMap[computev1alpha1.SchemeGroupVersion.Version] = storageMap

	return apiGroupInfo, true, nil
}

func (p StorageProvider) v1alpha1Storage(restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storageMap := map[string]rest.Storage{}

	machineClassStorage, err := machineclassstore.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["machineclasses"] = machineClassStorage.MachineClass

	machinePoolStorage, err := machinepoolstorage.NewStorage(restOptionsGetter, p.MachinePoolletClientConfig)
	if err != nil {
		return storageMap, err
	}

	storageMap["machinepools"] = machinePoolStorage.MachinePool
	storageMap["machinepools/status"] = machinePoolStorage.Status

	machineStorage, err := machinestorage.NewStorage(restOptionsGetter, machinePoolStorage.MachinePoolletConnectionInfo)
	if err != nil {
		return storageMap, err
	}

	storageMap["machines"] = machineStorage.Machine
	storageMap["machines/status"] = machineStorage.Status
	storageMap["machines/exec"] = machineStorage.Exec

	return storageMap, nil
}
