// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/internal/api"
	"spheric.cloud/spheric/internal/apis/core"
	diskstorage "spheric.cloud/spheric/internal/registry/core/disk/storage"
	disktypestorage "spheric.cloud/spheric/internal/registry/core/disktype/storage"
	fleetstorage "spheric.cloud/spheric/internal/registry/core/fleet/storage"
	instancestorage "spheric.cloud/spheric/internal/registry/core/instance/storage"
	instancetypestorage "spheric.cloud/spheric/internal/registry/core/instancetype/storage"
	networkstorage "spheric.cloud/spheric/internal/registry/core/network/storage"
	subnetstorage "spheric.cloud/spheric/internal/registry/core/subnet/storage"
	sphericserializer "spheric.cloud/spheric/internal/serializer"
	sphereletclient "spheric.cloud/spheric/internal/spherelet/client"

	serverstorage "k8s.io/apiserver/pkg/server/storage"
)

type StorageProvider struct {
	SphereletClientConfig sphereletclient.SphereletClientConfig
}

func (p StorageProvider) GroupName() string {
	return core.SchemeGroupVersion.Group
}

func (p StorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(p.GroupName(), api.Scheme, api.ParameterCodec, api.Codecs)
	apiGroupInfo.NegotiatedSerializer = sphericserializer.DefaultSubsetNegotiatedSerializer(api.Codecs)

	storageMap, err := p.v1alpha1Storage(restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, false, err
	}

	apiGroupInfo.VersionedResourcesStorageMap[corev1alpha1.SchemeGroupVersion.Version] = storageMap

	return apiGroupInfo, true, nil
}

func (p StorageProvider) v1alpha1Storage(restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storageMap := map[string]rest.Storage{}

	diskStorage, err := diskstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["disks"] = diskStorage.Disk
	storageMap["disks/status"] = diskStorage.Status

	diskTypeStorage, err := disktypestorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["disktypes"] = diskTypeStorage.DiskType

	fleetStorage, err := fleetstorage.NewStorage(restOptionsGetter, p.SphereletClientConfig)
	if err != nil {
		return storageMap, err
	}

	storageMap["fleets"] = fleetStorage.Fleet
	storageMap["fleets/status"] = fleetStorage.Status

	instanceStorage, err := instancestorage.NewStorage(restOptionsGetter, fleetStorage.SphereletConnectionInfo)
	if err != nil {
		return storageMap, err
	}

	storageMap["instances"] = instanceStorage.Instance
	storageMap["instances/status"] = instanceStorage.Status
	storageMap["instances/exec"] = instanceStorage.Exec

	instanceTypeStorage, err := instancetypestorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["instancetypes"] = instanceTypeStorage.InstanceType

	networkStorage, err := networkstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["networks"] = networkStorage.Network
	storageMap["networks/status"] = networkStorage.Status

	subnetStorage, err := subnetstorage.NewStorage(restOptionsGetter)
	if err != nil {
		return storageMap, err
	}

	storageMap["subnets"] = subnetStorage.Subnet
	storageMap["subnets/status"] = subnetStorage.Status

	return storageMap, nil
}
