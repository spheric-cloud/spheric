// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FleetSpec defines the desired state of Fleet
type FleetSpec struct {
	// ProviderID identifies the Fleet on provider side.
	ProviderID string
	// Taints of the Fleet. Only Machines who tolerate all the taints
	// will land in the Fleet.
	Taints []Taint
}

// FleetStatus defines the observed state of Fleet
type FleetStatus struct {
	State           FleetState
	Conditions      []FleetCondition
	Addresses       []FleetAddress
	DaemonEndpoints FleetDaemonEndpoints
	// Capacity represents the total resources of a fleet.
	Capacity ResourceList
	// Allocatable represents the resources of a fleet that are available for scheduling.
	Allocatable ResourceList
}

// FleetDaemonEndpoints lists ports opened by daemons running on the Fleet.
type FleetDaemonEndpoints struct {
	// Endpoint on which spherelet is listening.
	// +optional
	SphereletEndpoint DaemonEndpoint
}

// DaemonEndpoint contains information about a single Daemon endpoint.
type DaemonEndpoint struct {
	// Port number of the given endpoint.
	Port int32
}

type FleetAddressType string

const (
	// FleetHostName identifies a name of the fleet. Although every fleet can be assumed
	// to have a FleetAddress of this type, its exact syntax and semantics are not
	// defined, and are not consistent between different clusters.
	FleetHostName FleetAddressType = "Hostname"

	// FleetInternalIP identifies an IP address which may not be visible to hosts outside the cluster.
	// By default, it is assumed that apiserver can reach fleet internal IPs, though it is possible
	// to configure clusters where this is not the case.
	//
	// FleetInternalIP is the default type of fleet IP, and does not necessarily imply
	// that the IP is ONLY reachable internally. If a fleet has multiple internal IPs,
	// no specific semantics are assigned to the additional IPs.
	FleetInternalIP FleetAddressType = "InternalIP"

	// FleetExternalIP identifies an IP address which is, in some way, intended to be more usable from outside
	// the cluster than an internal IP, though no specific semantics are defined.
	FleetExternalIP FleetAddressType = "ExternalIP"

	// FleetInternalDNS identifies a DNS name which resolves to an IP address which has
	// the characteristics of a FleetInternalIP. The IP it resolves to may or may not
	// be a listed FleetInternalIP address.
	FleetInternalDNS FleetAddressType = "InternalDNS"

	// FleetExternalDNS identifies a DNS name which resolves to an IP address which has the characteristics
	// of FleetExternalIP. The IP it resolves to may or may not be a listed MachineExternalIP address.
	FleetExternalDNS FleetAddressType = "ExternalDNS"
)

type FleetAddress struct {
	Type    FleetAddressType
	Address string
}

// FleetConditionType is a type a FleetCondition can have.
type FleetConditionType string

// FleetCondition is one of the conditions of a disk.
type FleetCondition struct {
	// Type is the type of the condition.
	Type FleetConditionType
	// Status is the status of the condition.
	Status corev1.ConditionStatus
	// Reason is a machine-readable indication of why the condition is in a certain state.
	Reason string
	// Message is a human-readable explanation of why the condition has a certain reason / state.
	Message string
	// ObservedGeneration represents the .metadata.generation that the condition was set based upon.
	ObservedGeneration int64
	// LastTransitionTime is the last time the status of a condition has transitioned from one state to another.
	LastTransitionTime metav1.Time
}

// FleetState is a state a Fleet can be in.
// +enum
type FleetState string

const (
	// FleetStateReady marks a Fleet as ready for accepting a Machine.
	FleetStateReady FleetState = "Ready"
	// FleetStatePending marks a Fleet as pending readiness.
	FleetStatePending FleetState = "Pending"
	// FleetStateError marks a Fleet in an error state.
	FleetStateError FleetState = "Error"
	// FleetStateOffline marks a Fleet as offline.
	FleetStateOffline FleetState = "Offline"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +genclient:nonNamespaced

// Fleet is the Schema for the fleets API
type Fleet struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   FleetSpec
	Status FleetStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FleetList contains a list of Fleet
type FleetList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Fleet
}
