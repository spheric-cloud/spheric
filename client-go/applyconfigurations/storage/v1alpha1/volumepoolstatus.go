// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	v1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
)

// VolumePoolStatusApplyConfiguration represents an declarative configuration of the VolumePoolStatus type for use
// with apply.
type VolumePoolStatusApplyConfiguration struct {
	State                  *v1alpha1.VolumePoolState               `json:"state,omitempty"`
	Conditions             []VolumePoolConditionApplyConfiguration `json:"conditions,omitempty"`
	AvailableVolumeClasses []v1.LocalObjectReference               `json:"availableVolumeClasses,omitempty"`
	Capacity               *corev1alpha1.ResourceList              `json:"capacity,omitempty"`
	Allocatable            *corev1alpha1.ResourceList              `json:"allocatable,omitempty"`
}

// VolumePoolStatusApplyConfiguration constructs an declarative configuration of the VolumePoolStatus type for use with
// apply.
func VolumePoolStatus() *VolumePoolStatusApplyConfiguration {
	return &VolumePoolStatusApplyConfiguration{}
}

// WithState sets the State field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the State field is set to the value of the last call.
func (b *VolumePoolStatusApplyConfiguration) WithState(value v1alpha1.VolumePoolState) *VolumePoolStatusApplyConfiguration {
	b.State = &value
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
func (b *VolumePoolStatusApplyConfiguration) WithConditions(values ...*VolumePoolConditionApplyConfiguration) *VolumePoolStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}

// WithAvailableVolumeClasses adds the given value to the AvailableVolumeClasses field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the AvailableVolumeClasses field.
func (b *VolumePoolStatusApplyConfiguration) WithAvailableVolumeClasses(values ...v1.LocalObjectReference) *VolumePoolStatusApplyConfiguration {
	for i := range values {
		b.AvailableVolumeClasses = append(b.AvailableVolumeClasses, values[i])
	}
	return b
}

// WithCapacity sets the Capacity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Capacity field is set to the value of the last call.
func (b *VolumePoolStatusApplyConfiguration) WithCapacity(value corev1alpha1.ResourceList) *VolumePoolStatusApplyConfiguration {
	b.Capacity = &value
	return b
}

// WithAllocatable sets the Allocatable field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Allocatable field is set to the value of the last call.
func (b *VolumePoolStatusApplyConfiguration) WithAllocatable(value corev1alpha1.ResourceList) *VolumePoolStatusApplyConfiguration {
	b.Allocatable = &value
	return b
}
