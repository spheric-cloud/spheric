// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "spheric.cloud/spheric/client-go/applyconfigurations/storage/v1alpha1"
)

// EphemeralVolumeSourceApplyConfiguration represents an declarative configuration of the EphemeralVolumeSource type for use
// with apply.
type EphemeralVolumeSourceApplyConfiguration struct {
	VolumeTemplate *v1alpha1.VolumeTemplateSpecApplyConfiguration `json:"volumeTemplate,omitempty"`
}

// EphemeralVolumeSourceApplyConfiguration constructs an declarative configuration of the EphemeralVolumeSource type for use with
// apply.
func EphemeralVolumeSource() *EphemeralVolumeSourceApplyConfiguration {
	return &EphemeralVolumeSourceApplyConfiguration{}
}

// WithVolumeTemplate sets the VolumeTemplate field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VolumeTemplate field is set to the value of the last call.
func (b *EphemeralVolumeSourceApplyConfiguration) WithVolumeTemplate(value *v1alpha1.VolumeTemplateSpecApplyConfiguration) *EphemeralVolumeSourceApplyConfiguration {
	b.VolumeTemplate = value
	return b
}
