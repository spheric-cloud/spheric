// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package core

const (
	// WatchLabel is a label that can be applied to any spheric resource.
	//
	// Provider controllers that allow for selective reconciliation may check this label and proceed
	// with reconciliation of the object only if this label and a configured value are present.
	WatchLabel = "spheric.cloud/watch-filter"

	// ReconcileRequestAnnotation is an annotation that requested a reconciliation at a specific time.
	ReconcileRequestAnnotation = "reconcile.spheric.cloud/requested-at"

	// ManagedByAnnotation is an annotation that can be applied to resources to signify that
	// some external system is managing the resource.
	ManagedByAnnotation = "spheric.cloud/managed-by"

	// EphemeralManagedByAnnotation is an annotation that can be applied to resources to signify that
	// some ephemeral controller is managing the resource.
	EphemeralManagedByAnnotation = "spheric.cloud/ephemeral-managed-by"

	// DefaultEphemeralManager is the default spheric ephemeral manager.
	DefaultEphemeralManager = "ephemeral-manager"

	FinalizerNetwork = "core.spheric.cloud/network"
)
