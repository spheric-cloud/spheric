// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

// InstanceLister helps list Instances.
// All objects returned here must be treated as read-only.
type InstanceLister interface {
	// List lists all Instances in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Instance, err error)
	// Instances returns an object that can list and get Instances.
	Instances(namespace string) InstanceNamespaceLister
	InstanceListerExpansion
}

// instanceLister implements the InstanceLister interface.
type instanceLister struct {
	indexer cache.Indexer
}

// NewInstanceLister returns a new InstanceLister.
func NewInstanceLister(indexer cache.Indexer) InstanceLister {
	return &instanceLister{indexer: indexer}
}

// List lists all Instances in the indexer.
func (s *instanceLister) List(selector labels.Selector) (ret []*v1alpha1.Instance, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Instance))
	})
	return ret, err
}

// Instances returns an object that can list and get Instances.
func (s *instanceLister) Instances(namespace string) InstanceNamespaceLister {
	return instanceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// InstanceNamespaceLister helps list and get Instances.
// All objects returned here must be treated as read-only.
type InstanceNamespaceLister interface {
	// List lists all Instances in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Instance, err error)
	// Get retrieves the Instance from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Instance, error)
	InstanceNamespaceListerExpansion
}

// instanceNamespaceLister implements the InstanceNamespaceLister
// interface.
type instanceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Instances in the indexer for a given namespace.
func (s instanceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Instance, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Instance))
	})
	return ret, err
}

// Get retrieves the Instance from the indexer for a given namespace and name.
func (s instanceNamespaceLister) Get(name string) (*v1alpha1.Instance, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("instance"), name)
	}
	return obj.(*v1alpha1.Instance), nil
}
