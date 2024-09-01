// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"spheric.cloud/spheric/api/core/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func instanceIsOnFleet(instance *v1alpha1.Instance, FleetName string) bool {
	FleetRef := instance.Spec.FleetRef
	if FleetRef == nil {
		return false
	}

	return FleetRef.Name == FleetName
}

const InstanceSpecDiskNamesField = "instance-spec-disks"

func SetupInstanceSpecDiskNamesField(ctx context.Context, indexer client.FieldIndexer, FleetName string) error {
	return indexer.IndexField(
		ctx,
		&v1alpha1.Instance{},
		InstanceSpecDiskNamesField,
		func(object client.Object) []string {
			instance := object.(*v1alpha1.Instance)
			if !instanceIsOnFleet(instance, FleetName) {
				return nil
			}
			return v1alpha1.InstanceDiskNames(instance)
		},
	)
}

const InstanceSpecSecretNamesField = "instance-spec-secrets"

func SetupInstanceSpecSecretNamesField(ctx context.Context, indexer client.FieldIndexer, FleetName string) error {
	return indexer.IndexField(
		ctx,
		&v1alpha1.Instance{},
		InstanceSpecSecretNamesField,
		func(object client.Object) []string {
			instance := object.(*v1alpha1.Instance)
			if !instanceIsOnFleet(instance, FleetName) {
				return nil
			}

			return v1alpha1.InstanceSecretNames(instance)
		},
	)
}
