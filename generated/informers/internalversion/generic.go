/*
 * Copyright (c) 2021 by the OnMetal authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Code generated by informer-gen. DO NOT EDIT.

package internalversion

import (
	"fmt"

	compute "github.com/onmetal/onmetal-api/apis/compute"
	ipam "github.com/onmetal/onmetal-api/apis/ipam"
	networking "github.com/onmetal/onmetal-api/apis/networking"
	storage "github.com/onmetal/onmetal-api/apis/storage"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=compute.api.onmetal.de, Version=internalVersion
	case compute.SchemeGroupVersion.WithResource("machines"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Compute().InternalVersion().Machines().Informer()}, nil
	case compute.SchemeGroupVersion.WithResource("machineclasses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Compute().InternalVersion().MachineClasses().Informer()}, nil
	case compute.SchemeGroupVersion.WithResource("machinepools"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Compute().InternalVersion().MachinePools().Informer()}, nil

		// Group=ipam.api.onmetal.de, Version=internalVersion
	case ipam.SchemeGroupVersion.WithResource("prefixes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Ipam().InternalVersion().Prefixes().Informer()}, nil
	case ipam.SchemeGroupVersion.WithResource("prefixallocations"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Ipam().InternalVersion().PrefixAllocations().Informer()}, nil

		// Group=networking.api.onmetal.de, Version=internalVersion
	case networking.SchemeGroupVersion.WithResource("networkinterfaces"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Networking().InternalVersion().NetworkInterfaces().Informer()}, nil

		// Group=storage.api.onmetal.de, Version=internalVersion
	case storage.SchemeGroupVersion.WithResource("volumes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Storage().InternalVersion().Volumes().Informer()}, nil
	case storage.SchemeGroupVersion.WithResource("volumeclaims"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Storage().InternalVersion().VolumeClaims().Informer()}, nil
	case storage.SchemeGroupVersion.WithResource("volumeclasses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Storage().InternalVersion().VolumeClasses().Informer()}, nil
	case storage.SchemeGroupVersion.WithResource("volumepools"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Storage().InternalVersion().VolumePools().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
