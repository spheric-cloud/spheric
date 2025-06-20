// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
)

// NetworkInterfaceApplyConfiguration represents a declarative configuration of the NetworkInterface type for use
// with apply.
type NetworkInterfaceApplyConfiguration struct {
	Name             *string                            `json:"name,omitempty"`
	SubnetRef        *SubnetReferenceApplyConfiguration `json:"subnetRef,omitempty"`
	IPFamilies       []v1.IPFamily                      `json:"ipFamilies,omitempty"`
	IPs              []string                           `json:"ips,omitempty"`
	AccessIPFamilies []v1.IPFamily                      `json:"accessIPFamilies,omitempty"`
	AccessIPs        []string                           `json:"accessIPs,omitempty"`
}

// NetworkInterfaceApplyConfiguration constructs a declarative configuration of the NetworkInterface type for use with
// apply.
func NetworkInterface() *NetworkInterfaceApplyConfiguration {
	return &NetworkInterfaceApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *NetworkInterfaceApplyConfiguration) WithName(value string) *NetworkInterfaceApplyConfiguration {
	b.Name = &value
	return b
}

// WithSubnetRef sets the SubnetRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SubnetRef field is set to the value of the last call.
func (b *NetworkInterfaceApplyConfiguration) WithSubnetRef(value *SubnetReferenceApplyConfiguration) *NetworkInterfaceApplyConfiguration {
	b.SubnetRef = value
	return b
}

// WithIPFamilies adds the given value to the IPFamilies field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the IPFamilies field.
func (b *NetworkInterfaceApplyConfiguration) WithIPFamilies(values ...v1.IPFamily) *NetworkInterfaceApplyConfiguration {
	for i := range values {
		b.IPFamilies = append(b.IPFamilies, values[i])
	}
	return b
}

// WithIPs adds the given value to the IPs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the IPs field.
func (b *NetworkInterfaceApplyConfiguration) WithIPs(values ...string) *NetworkInterfaceApplyConfiguration {
	for i := range values {
		b.IPs = append(b.IPs, values[i])
	}
	return b
}

// WithAccessIPFamilies adds the given value to the AccessIPFamilies field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the AccessIPFamilies field.
func (b *NetworkInterfaceApplyConfiguration) WithAccessIPFamilies(values ...v1.IPFamily) *NetworkInterfaceApplyConfiguration {
	for i := range values {
		b.AccessIPFamilies = append(b.AccessIPFamilies, values[i])
	}
	return b
}

// WithAccessIPs adds the given value to the AccessIPs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the AccessIPs field.
func (b *NetworkInterfaceApplyConfiguration) WithAccessIPs(values ...string) *NetworkInterfaceApplyConfiguration {
	for i := range values {
		b.AccessIPs = append(b.AccessIPs, values[i])
	}
	return b
}
