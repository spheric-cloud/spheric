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
	"context"
	time "time"

	networking "github.com/onmetal/onmetal-api/apis/networking"
	clientsetinternalversion "github.com/onmetal/onmetal-api/generated/clientset/internalversion"
	internalinterfaces "github.com/onmetal/onmetal-api/generated/informers/internalversion/internalinterfaces"
	internalversion "github.com/onmetal/onmetal-api/generated/listers/networking/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// LoadBalancerInformer provides access to a shared informer and lister for
// LoadBalancers.
type LoadBalancerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.LoadBalancerLister
}

type loadBalancerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewLoadBalancerInformer constructs a new informer for LoadBalancer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewLoadBalancerInformer(client clientsetinternalversion.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredLoadBalancerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredLoadBalancerInformer constructs a new informer for LoadBalancer type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredLoadBalancerInformer(client clientsetinternalversion.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Networking().LoadBalancers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Networking().LoadBalancers(namespace).Watch(context.TODO(), options)
			},
		},
		&networking.LoadBalancer{},
		resyncPeriod,
		indexers,
	)
}

func (f *loadBalancerInformer) defaultInformer(client clientsetinternalversion.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredLoadBalancerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *loadBalancerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&networking.LoadBalancer{}, f.defaultInformer)
}

func (f *loadBalancerInformer) Lister() internalversion.LoadBalancerLister {
	return internalversion.NewLoadBalancerLister(f.Informer().GetIndexer())
}