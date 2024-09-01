// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DiskSpec defines the desired state of Disk
type DiskSpec struct {
	// TypeRef references the DiskClass of the Disk.
	TypeRef LocalObjectReference `json:"typeRef"`
	// InstanceRef references the using instance of the Disk.
	InstanceRef *LocalUIDReference `json:"instanceRef,omitempty"`
	// Resources is a description of the Disk's resources and capacity.
	Resources ResourceList `json:"resources,omitempty"`
}

// DiskStatus defines the observed state of Disk
type DiskStatus struct {
	// State represents the infrastructure state of a Disk.
	State DiskState `json:"state,omitempty"`
	// Access contains information to access the Disk. Must be set when Disk is in DiskStateAvailable.
	Access *DiskAccess `json:"access,omitempty"`
	// LastStateTransitionTime is the last time the State transitioned between values.
	LastStateTransitionTime *metav1.Time `json:"lastStateTransitionTime,omitempty"`
}

type DiskAccess struct {
	// Driver is the name of the drive to use for this volume. Required.
	Driver string `json:"driver"`
	// Handle is the unique handle of the volume.
	Handle string `json:"handle"`
	// Attributes are attributes of the volume to use.
	Attributes map[string]string `json:"attributes,omitempty"`
	// SecretRef references the (optional) secret containing the data to access the Disk.
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`
}

// DiskState represents the infrastructure state of a Disk.
type DiskState string

const (
	// DiskStatePending reports whether a Disk is about to be ready.
	DiskStatePending DiskState = "Pending"
	// DiskStateAvailable reports whether a Disk is available to be used.
	DiskStateAvailable DiskState = "Available"
	// DiskStateError reports that a Disk is in an error state.
	DiskStateError DiskState = "Error"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// Disk is the Schema for the disks API
type Disk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DiskSpec   `json:"spec,omitempty"`
	Status DiskStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DiskList contains a list of Disk
type DiskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Disk `json:"items"`
}

// DiskTemplateSpec is the specification of a Disk template.
type DiskTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DiskSpec `json:"spec,omitempty"`
}
