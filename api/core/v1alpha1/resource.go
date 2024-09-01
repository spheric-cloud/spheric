// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
)

// ResourceName is the name of a resource, most often used alongside a resource.Quantity.
type ResourceName string

const (
	// ResourceCPU is the amount of cpu in cores.
	ResourceCPU ResourceName = "cpu"
	// ResourceMemory is the amount of memory in bytes.
	ResourceMemory ResourceName = "memory"
	// ResourceStorage is the amount of storage, in bytes.
	ResourceStorage ResourceName = "storage"
	// ResourceTPS defines max throughput per second. (e.g. 1Gi)
	ResourceTPS ResourceName = "tps"
	// ResourceIOPS defines max IOPS in input/output operations per second.
	ResourceIOPS ResourceName = "iops"

	// ResourceInstanceTypePrefix is the prefix for instance type resources.
	ResourceInstanceTypePrefix = "instance-type/"
)

// ResourceInstanceType is the resource for a specific instance type.
func ResourceInstanceType(name string) ResourceName {
	return ResourceName(ResourceInstanceTypePrefix + name)
}

// IsInstanceTypeResource determines whether the given resource name is for an instance type.
func IsInstanceTypeResource(name ResourceName) bool {
	return strings.HasPrefix(string(name), ResourceInstanceTypePrefix)
}

// ResourceList is a list of ResourceName alongside their resource.Quantity.
type ResourceList map[ResourceName]resource.Quantity

// Name returns the resource with name if specified, otherwise it returns a nil quantity with default format.
func (rl *ResourceList) Name(name ResourceName, defaultFormat resource.Format) *resource.Quantity {
	if val, ok := (*rl)[name]; ok {
		return &val
	}
	return &resource.Quantity{Format: defaultFormat}
}

// Storage is a shorthand for getting the quantity associated with ResourceStorage.
func (rl *ResourceList) Storage() *resource.Quantity {
	return rl.Name(ResourceStorage, resource.BinarySI)
}

// Memory is a shorthand for getting the quantity associated with ResourceMemory.
func (rl *ResourceList) Memory() *resource.Quantity {
	return rl.Name(ResourceMemory, resource.BinarySI)
}

// CPU is a shorthand for getting the quantity associated with ResourceCPU.
func (rl *ResourceList) CPU() *resource.Quantity {
	return rl.Name(ResourceCPU, resource.DecimalSI)
}

// TPS is a shorthand for getting the quantity associated with ResourceTPS.
func (rl *ResourceList) TPS() *resource.Quantity {
	return rl.Name(ResourceTPS, resource.DecimalSI)
}

// IOPS is a shorthand for getting the quantity associated with ResourceIOPS.
func (rl *ResourceList) IOPS() *resource.Quantity {
	return rl.Name(ResourceIOPS, resource.DecimalSI)
}
