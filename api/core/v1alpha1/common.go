// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	InstanceFleetRefNameField        = "spec.fleetRef.name"
	InstanceInstanceTypeRefNameField = "spec.instanceTypeRef.name"

	// FleetsGroup is the system rbac group all fleets are in.
	FleetsGroup = "core.spheric.cloud:system:fleets"

	// FleetUserNamePrefix is the prefix all fleet users should have.
	FleetUserNamePrefix = "core.spheric.cloud:system:fleet:"

	SecretTypeIgnition = corev1.SecretType("core.spheric.cloud/ignition")
)

// FleetCommonName constructs the common name for a certificate of a fleet user.
func FleetCommonName(name string) string {
	return FleetUserNamePrefix + name
}
