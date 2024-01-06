// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	internalinterfaces "spheric.cloud/spheric/client-go/informers/internalinterfaces"
	v1alpha1 "spheric.cloud/spheric/client-go/listers/compute/v1alpha1"
	spheric "spheric.cloud/spheric/client-go/spheric"
)

// MachinePoolInformer provides access to a shared informer and lister for
// MachinePools.
type MachinePoolInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.MachinePoolLister
}

type machinePoolInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewMachinePoolInformer constructs a new informer for MachinePool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewMachinePoolInformer(client spheric.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredMachinePoolInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredMachinePoolInformer constructs a new informer for MachinePool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredMachinePoolInformer(client spheric.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ComputeV1alpha1().MachinePools().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ComputeV1alpha1().MachinePools().Watch(context.TODO(), options)
			},
		},
		&computev1alpha1.MachinePool{},
		resyncPeriod,
		indexers,
	)
}

func (f *machinePoolInformer) defaultInformer(client spheric.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredMachinePoolInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *machinePoolInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&computev1alpha1.MachinePool{}, f.defaultInformer)
}

func (f *machinePoolInformer) Lister() v1alpha1.MachinePoolLister {
	return v1alpha1.NewMachinePoolLister(f.Informer().GetIndexer())
}
