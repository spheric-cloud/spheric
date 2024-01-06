// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
)

const (
	MachineMachinePoolRefNameField  = "spec.machinePoolRef.name"
	MachineMachineClassRefNameField = "spec.machineClassRef.name"

	// MachinePoolsGroup is the system rbac group all machine pools are in.
	MachinePoolsGroup = "compute.spheric.cloud:system:machinepools"

	// MachinePoolUserNamePrefix is the prefix all machine pool users should have.
	MachinePoolUserNamePrefix = "compute.spheric.cloud:system:machinepool:"

	SecretTypeIgnition = corev1.SecretType("compute.spheric.cloud/ignition")
)

// MachinePoolCommonName constructs the common name for a certificate of a machine pool user.
func MachinePoolCommonName(name string) string {
	return MachinePoolUserNamePrefix + name
}

// EphemeralNetworkInterfaceSource is a definition for an ephemeral (i.e. coupled to the lifetime of the surrounding
// object) networking.NetworkInterface.
type EphemeralNetworkInterfaceSource struct {
	// NetworkInterfaceTemplate is the template definition of the networking.NetworkInterface.
	NetworkInterfaceTemplate *networkingv1alpha1.NetworkInterfaceTemplateSpec `json:"networkInterfaceTemplate,omitempty"`
}

// EphemeralVolumeSource is a definition for an ephemeral (i.e. coupled to the lifetime of the surrounding object)
// storage.Volume.
type EphemeralVolumeSource struct {
	// VolumeTemplate is the template definition of the storage.Volume.
	VolumeTemplate *storagev1alpha1.VolumeTemplateSpec `json:"volumeTemplate,omitempty"`
}
