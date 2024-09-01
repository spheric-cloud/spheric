// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

// NetworkStatusApplyConfiguration represents an declarative configuration of the NetworkStatus type for use
// with apply.
type NetworkStatusApplyConfiguration struct {
	State *v1alpha1.NetworkState `json:"state,omitempty"`
}

// NetworkStatusApplyConfiguration constructs an declarative configuration of the NetworkStatus type for use with
// apply.
func NetworkStatus() *NetworkStatusApplyConfiguration {
	return &NetworkStatusApplyConfiguration{}
}

// WithState sets the State field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the State field is set to the value of the last call.
func (b *NetworkStatusApplyConfiguration) WithState(value v1alpha1.NetworkState) *NetworkStatusApplyConfiguration {
	b.State = &value
	return b
}