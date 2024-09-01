// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// DaemonEndpointApplyConfiguration represents an declarative configuration of the DaemonEndpoint type for use
// with apply.
type DaemonEndpointApplyConfiguration struct {
	Port *int32 `json:"port,omitempty"`
}

// DaemonEndpointApplyConfiguration constructs an declarative configuration of the DaemonEndpoint type for use with
// apply.
func DaemonEndpoint() *DaemonEndpointApplyConfiguration {
	return &DaemonEndpointApplyConfiguration{}
}

// WithPort sets the Port field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Port field is set to the value of the last call.
func (b *DaemonEndpointApplyConfiguration) WithPort(value int32) *DaemonEndpointApplyConfiguration {
	b.Port = &value
	return b
}