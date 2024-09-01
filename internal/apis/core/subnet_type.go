// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubnetSpec defines the desired state of Subnet
type SubnetSpec struct {
	// NetworkRef references the network this subnet is part of.
	NetworkRef LocalObjectReference
	// CIDRs are the primary CIDR ranges of this Subnet.
	CIDRs []string
}

// SubnetStatus defines the observed state of Subnet
type SubnetStatus struct {
	// State is the state of the machine.
	State SubnetState
}

// SubnetState is the state of a network.
// +enum
type SubnetState string

const (
	// SubnetStatePending means the network is being provisioned.
	SubnetStatePending SubnetState = "Pending"
	// SubnetStateAvailable means the network is ready to use.
	SubnetStateAvailable SubnetState = "Available"
	// SubnetStateError means the network is in an error state.
	SubnetStateError SubnetState = "Error"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Subnet is the Schema for the network API
type Subnet struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   SubnetSpec
	Status SubnetStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SubnetList contains a list of Subnet
type SubnetList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Subnet
}
