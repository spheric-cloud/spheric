/*
 * Copyright (c) 2022 by the OnMetal authors.
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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/onmetal/onmetal-api/api/storage/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// VolumePoolStatusApplyConfiguration represents an declarative configuration of the VolumePoolStatus type for use
// with apply.
type VolumePoolStatusApplyConfiguration struct {
	State                  *v1alpha1.VolumePoolState               `json:"state,omitempty"`
	Conditions             []VolumePoolConditionApplyConfiguration `json:"conditions,omitempty"`
	AvailableVolumeClasses []v1.LocalObjectReference               `json:"availableVolumeClasses,omitempty"`
	Available              *v1.ResourceList                        `json:"available,omitempty"`
	Used                   *v1.ResourceList                        `json:"used,omitempty"`
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

// WithAvailable sets the Available field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Available field is set to the value of the last call.
func (b *VolumePoolStatusApplyConfiguration) WithAvailable(value v1.ResourceList) *VolumePoolStatusApplyConfiguration {
	b.Available = &value
	return b
}

// WithUsed sets the Used field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Used field is set to the value of the last call.
func (b *VolumePoolStatusApplyConfiguration) WithUsed(value v1.ResourceList) *VolumePoolStatusApplyConfiguration {
	b.Used = &value
	return b
}
