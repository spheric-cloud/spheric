// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// InstanceTypeRef references the instance type of the instance.
	InstanceTypeRef LocalObjectReference
	// FleetSelector selects a suitable Fleet by the given labels.
	FleetSelector map[string]string
	// FleetRef defines the fleet to run the instance in.
	// If empty, a scheduler will figure out an appropriate pool to run the instance in.
	FleetRef *LocalObjectReference
	// Power is the desired instance power state.
	// Defaults to PowerOn.
	Power Power
	// Image is the optional URL providing the operating system image of the instance.
	// +optional
	Image string
	// ImagePullSecretRef is an optional secret for pulling the image of a instance.
	ImagePullSecretRef *LocalObjectReference
	// NetworkInterfaces define a list of network interfaces present on the instance
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	NetworkInterfaces []NetworkInterface
	// Disks are the disks attached to this instance.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	Disks []AttachedDisk
	// IgnitionRef is a reference to a secret containing the ignition YAML for the instance to boot up.
	// If key is empty, DefaultIgnitionKey will be used as fallback.
	IgnitionRef *SecretKeySelector
	// EFIVars are variables to pass to EFI while booting up.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	EFIVars []EFIVar
	// Tolerations define tolerations the Instance has. Only fleets whose taints
	// covered by Tolerations will be considered to run the Instance.
	Tolerations []Toleration
}

// Power is the desired power state of a Instance.
type Power string

const (
	// PowerOn indicates that a Instance should be powered on.
	PowerOn Power = "On"
	// PowerOff indicates that a Instance should be powered off.
	PowerOff Power = "Off"
)

// EFIVar is a variable to pass to EFI while booting up.
type EFIVar struct {
	// Name is the name of the EFIVar.
	Name string
	// UUID is the uuid of the EFIVar.
	UUID string
	// Value is the value of the EFIVar.
	Value string
}

// DefaultIgnitionKey is the default key for InstanceSpec.UserData.
const DefaultIgnitionKey = "ignition.yaml"

type SubnetReference struct {
	// NetworkName is the name of the referenced network.
	NetworkName string
	// Name of the referenced subnet.
	Name string
}

// NetworkInterface is the definition of a single interface
type NetworkInterface struct {
	// Name is the name of the network interface.
	Name string
	// SubnetRef references the Subnet this NetworkInterface is connected to
	SubnetRef SubnetReference
	// IPFamilies defines which IPFamilies this NetworkInterface is supporting
	IPFamilies []corev1.IPFamily
	// IPs are the literal requested IPs for this NetworkInterface.
	IPs []string
	// AccessIPFamilies are the access configuration IP families.
	AccessIPFamilies []corev1.IPFamily
	// AccessIPs are the literal request access IPs.
	AccessIPs []string
}

// AttachedDisk defines a disk attached to a instance.
type AttachedDisk struct {
	// Name is the name of the disk.
	Name string
	// Device is the device name where the disk should be attached.
	// Pointer to distinguish between explicit zero and not specified.
	// If empty, an unused device name will be determined if possible.
	Device *string
	// AttachedDiskSource is the source where the storage for the disk resides at.
	AttachedDiskSource
}

// AttachedDiskSource specifies the source to use for a disk.
type AttachedDiskSource struct {
	// DiskRef instructs to use the specified Disk as source for the attachment.
	DiskRef *LocalObjectReference
	// EmptyDisk instructs to use a disk offered by the fleet provider.
	EmptyDisk *EmptyDiskSource
	// Ephemeral instructs to create an ephemeral (i.e. coupled to the lifetime of the surrounding object)
	// disk to use.
	Ephemeral *EphemeralDiskSource
}

// EphemeralDiskSource is a definition for an ephemeral (i.e. coupled to the lifetime of the surrounding object)
// disk.
type EphemeralDiskSource struct {
	// DiskTemplate is the template definition of a Disk.
	DiskTemplate *DiskTemplateSpec
}

// EmptyDiskSource is a disk that's offered by the fleet provider.
// Usually ephemeral (i.e. deleted when the surrounding entity is deleted), with
// varying performance characteristics. Potentially not recoverable.
type EmptyDiskSource struct {
	// SizeLimit is the total amount of local storage required for this EmptyDisk disk.
	// The default is nil which means that the limit is undefined.
	SizeLimit *resource.Quantity
}

// NetworkInterfaceStatus reports the status of an NetworkInterfaceSource.
type NetworkInterfaceStatus struct {
	// Name is the name of the NetworkInterface to whom the status belongs to.
	Name string
	// IPs are the ips allocated for the network interface.
	IPs []string
	// AccessIPs are the allocated access IPs for the network interface.
	AccessIPs []string
	// State represents the attachment state of a NetworkInterface.
	State NetworkInterfaceState
	// LastStateTransitionTime is the last time the State transitioned.
	LastStateTransitionTime *metav1.Time
}

// NetworkInterfaceState is the infrastructure attachment state a NetworkInterface can be in.
type NetworkInterfaceState string

const (
	// NetworkInterfaceStatePending indicates that the attachment of a network interface is pending.
	NetworkInterfaceStatePending NetworkInterfaceState = "Pending"
	// NetworkInterfaceStateAttached indicates that a network interface has been successfully attached.
	NetworkInterfaceStateAttached NetworkInterfaceState = "Attached"
)

// AttachedDiskStatus is the status of a disk.
type AttachedDiskStatus struct {
	// Name is the name of the attached disk.
	Name string
	// State represents the attachment state of a disk.
	State AttachedDiskState
	// LastStateTransitionTime is the last time the State transitioned.
	LastStateTransitionTime *metav1.Time
}

// AttachedDiskState is the infrastructure attachment state a disk can be in.
type AttachedDiskState string

const (
	// AttachedDiskStatePending indicates that the attachment of a disk is pending.
	AttachedDiskStatePending AttachedDiskState = "Pending"
	// AttachedDiskStateAttached indicates that a disk has been successfully attached.
	AttachedDiskStateAttached AttachedDiskState = "Attached"
)

// InstanceStatus defines the observed state of Instance
type InstanceStatus struct {
	// InstanceID is the provider specific instance ID in the format '<type>://<instance_id>'.
	InstanceID string
	// ObservedGeneration is the last generation the Fleet observed of the Instance.
	ObservedGeneration int64
	// State is the infrastructure state of the instance.
	State InstanceState
	// NetworkInterfaces is the list of network interface states for the instance.
	NetworkInterfaces []NetworkInterfaceStatus
	// Disks is the list of disk states for the instance.
	Disks []AttachedDiskStatus
}

// InstanceState is the state of a instance.
// +enum
type InstanceState string

const (
	// InstanceStatePending means the Instance has been accepted by the system, but not yet completely started.
	// This includes time before being bound to a Fleet, as well as time spent setting up the Instance on that
	// Fleet.
	InstanceStatePending InstanceState = "Pending"
	// InstanceStateRunning means the instance is running on a Fleet.
	InstanceStateRunning InstanceState = "Running"
	// InstanceStateShutdown means the instance is shut down.
	InstanceStateShutdown InstanceState = "Shutdown"
	// InstanceStateTerminated means the instance has been permanently stopped and cannot be started.
	InstanceStateTerminated InstanceState = "Terminated"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   InstanceSpec
	Status InstanceStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Instance
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:conversion-gen:explicit-from=net/url.Values

// InstanceExecOptions is the query options to a Instance's remote exec call
type InstanceExecOptions struct {
	metav1.TypeMeta
	InsecureSkipTLSVerifyBackend bool
}
