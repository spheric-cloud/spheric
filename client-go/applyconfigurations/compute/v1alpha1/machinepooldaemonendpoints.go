// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// MachinePoolDaemonEndpointsApplyConfiguration represents an declarative configuration of the MachinePoolDaemonEndpoints type for use
// with apply.
type MachinePoolDaemonEndpointsApplyConfiguration struct {
	MachinepoolletEndpoint *DaemonEndpointApplyConfiguration `json:"machinepoolletEndpoint,omitempty"`
}

// MachinePoolDaemonEndpointsApplyConfiguration constructs an declarative configuration of the MachinePoolDaemonEndpoints type for use with
// apply.
func MachinePoolDaemonEndpoints() *MachinePoolDaemonEndpointsApplyConfiguration {
	return &MachinePoolDaemonEndpointsApplyConfiguration{}
}

// WithMachinepoolletEndpoint sets the MachinepoolletEndpoint field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the MachinepoolletEndpoint field is set to the value of the last call.
func (b *MachinePoolDaemonEndpointsApplyConfiguration) WithMachinepoolletEndpoint(value *DaemonEndpointApplyConfiguration) *MachinePoolDaemonEndpointsApplyConfiguration {
	b.MachinepoolletEndpoint = value
	return b
}
