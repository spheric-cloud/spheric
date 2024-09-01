// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

// InstanceSpecApplyConfiguration represents an declarative configuration of the InstanceSpec type for use
// with apply.
type InstanceSpecApplyConfiguration struct {
	InstanceTypeRef    *LocalObjectReferenceApplyConfiguration `json:"instanceTypeRef,omitempty"`
	FleetSelector      map[string]string                       `json:"fleetSelector,omitempty"`
	FleetRef           *LocalObjectReferenceApplyConfiguration `json:"fleetRef,omitempty"`
	Power              *corev1alpha1.Power                     `json:"power,omitempty"`
	Image              *string                                 `json:"image,omitempty"`
	ImagePullSecretRef *LocalObjectReferenceApplyConfiguration `json:"imagePullSecret,omitempty"`
	NetworkInterfaces  []NetworkInterfaceApplyConfiguration    `json:"networkInterfaces,omitempty"`
	Disks              []AttachedDiskApplyConfiguration        `json:"disks,omitempty"`
	IgnitionRef        *SecretKeySelectorApplyConfiguration    `json:"ignitionRef,omitempty"`
	EFIVars            []EFIVarApplyConfiguration              `json:"efiVars,omitempty"`
	Tolerations        []TolerationApplyConfiguration          `json:"tolerations,omitempty"`
}

// InstanceSpecApplyConfiguration constructs an declarative configuration of the InstanceSpec type for use with
// apply.
func InstanceSpec() *InstanceSpecApplyConfiguration {
	return &InstanceSpecApplyConfiguration{}
}

// WithInstanceTypeRef sets the InstanceTypeRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InstanceTypeRef field is set to the value of the last call.
func (b *InstanceSpecApplyConfiguration) WithInstanceTypeRef(value *LocalObjectReferenceApplyConfiguration) *InstanceSpecApplyConfiguration {
	b.InstanceTypeRef = value
	return b
}

// WithFleetSelector puts the entries into the FleetSelector field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the FleetSelector field,
// overwriting an existing map entries in FleetSelector field with the same key.
func (b *InstanceSpecApplyConfiguration) WithFleetSelector(entries map[string]string) *InstanceSpecApplyConfiguration {
	if b.FleetSelector == nil && len(entries) > 0 {
		b.FleetSelector = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.FleetSelector[k] = v
	}
	return b
}

// WithFleetRef sets the FleetRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FleetRef field is set to the value of the last call.
func (b *InstanceSpecApplyConfiguration) WithFleetRef(value *LocalObjectReferenceApplyConfiguration) *InstanceSpecApplyConfiguration {
	b.FleetRef = value
	return b
}

// WithPower sets the Power field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Power field is set to the value of the last call.
func (b *InstanceSpecApplyConfiguration) WithPower(value corev1alpha1.Power) *InstanceSpecApplyConfiguration {
	b.Power = &value
	return b
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *InstanceSpecApplyConfiguration) WithImage(value string) *InstanceSpecApplyConfiguration {
	b.Image = &value
	return b
}

// WithImagePullSecretRef sets the ImagePullSecretRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ImagePullSecretRef field is set to the value of the last call.
func (b *InstanceSpecApplyConfiguration) WithImagePullSecretRef(value *LocalObjectReferenceApplyConfiguration) *InstanceSpecApplyConfiguration {
	b.ImagePullSecretRef = value
	return b
}

// WithNetworkInterfaces adds the given value to the NetworkInterfaces field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NetworkInterfaces field.
func (b *InstanceSpecApplyConfiguration) WithNetworkInterfaces(values ...*NetworkInterfaceApplyConfiguration) *InstanceSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNetworkInterfaces")
		}
		b.NetworkInterfaces = append(b.NetworkInterfaces, *values[i])
	}
	return b
}

// WithDisks adds the given value to the Disks field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Disks field.
func (b *InstanceSpecApplyConfiguration) WithDisks(values ...*AttachedDiskApplyConfiguration) *InstanceSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithDisks")
		}
		b.Disks = append(b.Disks, *values[i])
	}
	return b
}

// WithIgnitionRef sets the IgnitionRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IgnitionRef field is set to the value of the last call.
func (b *InstanceSpecApplyConfiguration) WithIgnitionRef(value *SecretKeySelectorApplyConfiguration) *InstanceSpecApplyConfiguration {
	b.IgnitionRef = value
	return b
}

// WithEFIVars adds the given value to the EFIVars field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EFIVars field.
func (b *InstanceSpecApplyConfiguration) WithEFIVars(values ...*EFIVarApplyConfiguration) *InstanceSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithEFIVars")
		}
		b.EFIVars = append(b.EFIVars, *values[i])
	}
	return b
}

// WithTolerations adds the given value to the Tolerations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tolerations field.
func (b *InstanceSpecApplyConfiguration) WithTolerations(values ...*TolerationApplyConfiguration) *InstanceSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTolerations")
		}
		b.Tolerations = append(b.Tolerations, *values[i])
	}
	return b
}
