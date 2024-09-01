// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"fmt"
)

// InstanceEphemeralDiskName returns the name of a Disk for an ephemeral instance disk.
func InstanceEphemeralDiskName(instanceName, instanceDiskName string) string {
	return fmt.Sprintf("%s-%s", instanceName, instanceDiskName)
}

// InstanceDiskName returns the name of the DiskClaim for a instance disk.
func InstanceDiskName(instanceName string, disk AttachedDisk) string {
	switch {
	case disk.DiskRef != nil:
		return disk.DiskRef.Name
	case disk.Ephemeral != nil:
		return InstanceEphemeralDiskName(instanceName, disk.Name)
	default:
		return ""
	}
}

// InstanceDiskNames returns all Disk names of a instance.
func InstanceDiskNames(instance *Instance) []string {
	var names []string
	for _, disk := range instance.Spec.Disks {
		if name := InstanceDiskName(instance.Name, disk); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// InstanceSecretNames returns all secret names of a instance.
func InstanceSecretNames(instance *Instance) []string {
	var names []string

	if imagePullSecretRef := instance.Spec.ImagePullSecretRef; imagePullSecretRef != nil {
		names = append(names, imagePullSecretRef.Name)
	}

	if ignitionRef := instance.Spec.IgnitionRef; ignitionRef != nil {
		names = append(names, ignitionRef.Name)
	}

	return names
}
