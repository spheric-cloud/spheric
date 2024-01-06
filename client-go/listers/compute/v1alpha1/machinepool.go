// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
)

// MachinePoolLister helps list MachinePools.
// All objects returned here must be treated as read-only.
type MachinePoolLister interface {
	// List lists all MachinePools in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.MachinePool, err error)
	// Get retrieves the MachinePool from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.MachinePool, error)
	MachinePoolListerExpansion
}

// machinePoolLister implements the MachinePoolLister interface.
type machinePoolLister struct {
	indexer cache.Indexer
}

// NewMachinePoolLister returns a new MachinePoolLister.
func NewMachinePoolLister(indexer cache.Indexer) MachinePoolLister {
	return &machinePoolLister{indexer: indexer}
}

// List lists all MachinePools in the indexer.
func (s *machinePoolLister) List(selector labels.Selector) (ret []*v1alpha1.MachinePool, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.MachinePool))
	})
	return ret, err
}

// Get retrieves the MachinePool from the index for a given name.
func (s *machinePoolLister) Get(name string) (*v1alpha1.MachinePool, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("machinepool"), name)
	}
	return obj.(*v1alpha1.MachinePool), nil
}
