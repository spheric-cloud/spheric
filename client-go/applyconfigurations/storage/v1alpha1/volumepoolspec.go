// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "spheric.cloud/spheric/client-go/applyconfigurations/common/v1alpha1"
)

// VolumePoolSpecApplyConfiguration represents an declarative configuration of the VolumePoolSpec type for use
// with apply.
type VolumePoolSpecApplyConfiguration struct {
	ProviderID *string                            `json:"providerID,omitempty"`
	Taints     []v1alpha1.TaintApplyConfiguration `json:"taints,omitempty"`
}

// VolumePoolSpecApplyConfiguration constructs an declarative configuration of the VolumePoolSpec type for use with
// apply.
func VolumePoolSpec() *VolumePoolSpecApplyConfiguration {
	return &VolumePoolSpecApplyConfiguration{}
}

// WithProviderID sets the ProviderID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ProviderID field is set to the value of the last call.
func (b *VolumePoolSpecApplyConfiguration) WithProviderID(value string) *VolumePoolSpecApplyConfiguration {
	b.ProviderID = &value
	return b
}

// WithTaints adds the given value to the Taints field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Taints field.
func (b *VolumePoolSpecApplyConfiguration) WithTaints(values ...*v1alpha1.TaintApplyConfiguration) *VolumePoolSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTaints")
		}
		b.Taints = append(b.Taints, *values[i])
	}
	return b
}
