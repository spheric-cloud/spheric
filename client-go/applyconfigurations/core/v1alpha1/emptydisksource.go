// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	resource "k8s.io/apimachinery/pkg/api/resource"
)

// EmptyDiskSourceApplyConfiguration represents an declarative configuration of the EmptyDiskSource type for use
// with apply.
type EmptyDiskSourceApplyConfiguration struct {
	SizeLimit *resource.Quantity `json:"sizeLimit,omitempty"`
}

// EmptyDiskSourceApplyConfiguration constructs an declarative configuration of the EmptyDiskSource type for use with
// apply.
func EmptyDiskSource() *EmptyDiskSourceApplyConfiguration {
	return &EmptyDiskSourceApplyConfiguration{}
}

// WithSizeLimit sets the SizeLimit field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SizeLimit field is set to the value of the last call.
func (b *EmptyDiskSourceApplyConfiguration) WithSizeLimit(value resource.Quantity) *EmptyDiskSourceApplyConfiguration {
	b.SizeLimit = &value
	return b
}