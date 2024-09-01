// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// InstanceTypeFinalizer is the finalizer for InstanceType.
	InstanceTypeFinalizer = GroupName + "/instancetype"
)

// InstanceTypeClass denotes the type of InstanceType.
type InstanceTypeClass string

const (
	InstanceTypeContinuous InstanceTypeClass = "Continuous"
	InstanceTypeDiscrete   InstanceTypeClass = "Discrete"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +genclient:nonNamespaced

// InstanceType is the Schema for the instancetypes API
type InstanceType struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Class specifies the class of the InstanceType.
	// Can either be 'Continuous' or 'Discrete'.
	Class InstanceTypeClass

	// Capabilities are the capabilities of the instance type.
	Capabilities ResourceList
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceTypeList contains a list of InstanceType
type InstanceTypeList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []InstanceType
}
