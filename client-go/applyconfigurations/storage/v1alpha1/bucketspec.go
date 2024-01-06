// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	v1alpha1 "spheric.cloud/spheric/client-go/applyconfigurations/common/v1alpha1"
)

// BucketSpecApplyConfiguration represents an declarative configuration of the BucketSpec type for use
// with apply.
type BucketSpecApplyConfiguration struct {
	BucketClassRef     *v1.LocalObjectReference                `json:"bucketClassRef,omitempty"`
	BucketPoolSelector map[string]string                       `json:"bucketPoolSelector,omitempty"`
	BucketPoolRef      *v1.LocalObjectReference                `json:"bucketPoolRef,omitempty"`
	Tolerations        []v1alpha1.TolerationApplyConfiguration `json:"tolerations,omitempty"`
}

// BucketSpecApplyConfiguration constructs an declarative configuration of the BucketSpec type for use with
// apply.
func BucketSpec() *BucketSpecApplyConfiguration {
	return &BucketSpecApplyConfiguration{}
}

// WithBucketClassRef sets the BucketClassRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BucketClassRef field is set to the value of the last call.
func (b *BucketSpecApplyConfiguration) WithBucketClassRef(value v1.LocalObjectReference) *BucketSpecApplyConfiguration {
	b.BucketClassRef = &value
	return b
}

// WithBucketPoolSelector puts the entries into the BucketPoolSelector field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the BucketPoolSelector field,
// overwriting an existing map entries in BucketPoolSelector field with the same key.
func (b *BucketSpecApplyConfiguration) WithBucketPoolSelector(entries map[string]string) *BucketSpecApplyConfiguration {
	if b.BucketPoolSelector == nil && len(entries) > 0 {
		b.BucketPoolSelector = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.BucketPoolSelector[k] = v
	}
	return b
}

// WithBucketPoolRef sets the BucketPoolRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BucketPoolRef field is set to the value of the last call.
func (b *BucketSpecApplyConfiguration) WithBucketPoolRef(value v1.LocalObjectReference) *BucketSpecApplyConfiguration {
	b.BucketPoolRef = &value
	return b
}

// WithTolerations adds the given value to the Tolerations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tolerations field.
func (b *BucketSpecApplyConfiguration) WithTolerations(values ...*v1alpha1.TolerationApplyConfiguration) *BucketSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTolerations")
		}
		b.Tolerations = append(b.Tolerations, *values[i])
	}
	return b
}
