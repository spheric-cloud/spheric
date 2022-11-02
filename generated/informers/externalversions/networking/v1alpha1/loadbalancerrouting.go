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

package v1alpha1

import (
	"context"
	time "time"

	networkingv1alpha1 "github.com/onmetal/onmetal-api/apis/networking/v1alpha1"
	versioned "github.com/onmetal/onmetal-api/generated/clientset/versioned"
	internalinterfaces "github.com/onmetal/onmetal-api/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/onmetal/onmetal-api/generated/listers/networking/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// LoadBalancerRoutingInformer provides access to a shared informer and lister for
// LoadBalancerRoutings.
type LoadBalancerRoutingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.LoadBalancerRoutingLister
}

type loadBalancerRoutingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewLoadBalancerRoutingInformer constructs a new informer for LoadBalancerRouting type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewLoadBalancerRoutingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredLoadBalancerRoutingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredLoadBalancerRoutingInformer constructs a new informer for LoadBalancerRouting type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredLoadBalancerRoutingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkingV1alpha1().LoadBalancerRoutings(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NetworkingV1alpha1().LoadBalancerRoutings(namespace).Watch(context.TODO(), options)
			},
		},
		&networkingv1alpha1.LoadBalancerRouting{},
		resyncPeriod,
		indexers,
	)
}

func (f *loadBalancerRoutingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredLoadBalancerRoutingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *loadBalancerRoutingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&networkingv1alpha1.LoadBalancerRouting{}, f.defaultInformer)
}

func (f *loadBalancerRoutingInformer) Lister() v1alpha1.LoadBalancerRoutingLister {
	return v1alpha1.NewLoadBalancerRoutingLister(f.Informer().GetIndexer())
}
