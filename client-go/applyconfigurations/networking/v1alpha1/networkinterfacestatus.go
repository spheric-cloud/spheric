// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	commonv1alpha1 "spheric.cloud/spheric/api/common/v1alpha1"
	v1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
)

// NetworkInterfaceStatusApplyConfiguration represents an declarative configuration of the NetworkInterfaceStatus type for use
// with apply.
type NetworkInterfaceStatusApplyConfiguration struct {
	State                   *v1alpha1.NetworkInterfaceState `json:"state,omitempty"`
	LastStateTransitionTime *v1.Time                        `json:"lastStateTransitionTime,omitempty"`
	IPs                     []commonv1alpha1.IP             `json:"ips,omitempty"`
	Prefixes                []commonv1alpha1.IPPrefix       `json:"prefixes,omitempty"`
	VirtualIP               *commonv1alpha1.IP              `json:"virtualIP,omitempty"`
}

// NetworkInterfaceStatusApplyConfiguration constructs an declarative configuration of the NetworkInterfaceStatus type for use with
// apply.
func NetworkInterfaceStatus() *NetworkInterfaceStatusApplyConfiguration {
	return &NetworkInterfaceStatusApplyConfiguration{}
}

// WithState sets the State field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the State field is set to the value of the last call.
func (b *NetworkInterfaceStatusApplyConfiguration) WithState(value v1alpha1.NetworkInterfaceState) *NetworkInterfaceStatusApplyConfiguration {
	b.State = &value
	return b
}

// WithLastStateTransitionTime sets the LastStateTransitionTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastStateTransitionTime field is set to the value of the last call.
func (b *NetworkInterfaceStatusApplyConfiguration) WithLastStateTransitionTime(value v1.Time) *NetworkInterfaceStatusApplyConfiguration {
	b.LastStateTransitionTime = &value
	return b
}

// WithIPs adds the given value to the IPs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the IPs field.
func (b *NetworkInterfaceStatusApplyConfiguration) WithIPs(values ...commonv1alpha1.IP) *NetworkInterfaceStatusApplyConfiguration {
	for i := range values {
		b.IPs = append(b.IPs, values[i])
	}
	return b
}

// WithPrefixes adds the given value to the Prefixes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Prefixes field.
func (b *NetworkInterfaceStatusApplyConfiguration) WithPrefixes(values ...commonv1alpha1.IPPrefix) *NetworkInterfaceStatusApplyConfiguration {
	for i := range values {
		b.Prefixes = append(b.Prefixes, values[i])
	}
	return b
}

// WithVirtualIP sets the VirtualIP field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VirtualIP field is set to the value of the last call.
func (b *NetworkInterfaceStatusApplyConfiguration) WithVirtualIP(value commonv1alpha1.IP) *NetworkInterfaceStatusApplyConfiguration {
	b.VirtualIP = &value
	return b
}
