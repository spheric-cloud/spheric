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
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	internalinterfaces "spheric.cloud/spheric/client-go/informers/internalinterfaces"
	v1alpha1 "spheric.cloud/spheric/client-go/listers/core/v1alpha1"
	spheric "spheric.cloud/spheric/client-go/spheric"
)

// ResourceQuotaInformer provides access to a shared informer and lister for
// ResourceQuotas.
type ResourceQuotaInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ResourceQuotaLister
}

type resourceQuotaInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewResourceQuotaInformer constructs a new informer for ResourceQuota type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewResourceQuotaInformer(client spheric.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredResourceQuotaInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredResourceQuotaInformer constructs a new informer for ResourceQuota type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredResourceQuotaInformer(client spheric.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().ResourceQuotas(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().ResourceQuotas(namespace).Watch(context.TODO(), options)
			},
		},
		&corev1alpha1.ResourceQuota{},
		resyncPeriod,
		indexers,
	)
}

func (f *resourceQuotaInformer) defaultInformer(client spheric.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredResourceQuotaInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *resourceQuotaInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1alpha1.ResourceQuota{}, f.defaultInformer)
}

func (f *resourceQuotaInformer) Lister() v1alpha1.ResourceQuotaLister {
	return v1alpha1.NewResourceQuotaLister(f.Informer().GetIndexer())
}
