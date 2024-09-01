// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

const InstanceSpecDiskNamesField = "Instance.spec-disk-names"

func SetupInstanceSpecDiskNamesFieldIndexer(ctx context.Context, idx client.FieldIndexer) error {
	return idx.IndexField(ctx, &corev1alpha1.Instance{}, InstanceSpecDiskNamesField, func(obj client.Object) []string {
		instance := obj.(*corev1alpha1.Instance)
		return corev1alpha1.InstanceDiskNames(instance)
	})
}

const InstanceSpecFleetRefNameField = "Instance.spec.fleetRef.name"

func SetupInstanceSpecFleetRefNameFieldIndexer(ctx context.Context, idx client.FieldIndexer) error {
	return idx.IndexField(ctx, &corev1alpha1.Instance{}, InstanceSpecFleetRefNameField, func(obj client.Object) []string {
		instance := obj.(*corev1alpha1.Instance)
		fleetRef := instance.Spec.FleetRef
		if fleetRef == nil {
			return []string{""}
		}

		return []string{fleetRef.Name}
	})
}

const InstanceSpecInstanceTypeRefNameField = corev1alpha1.InstanceInstanceTypeRefNameField

func SetupInstanceSpecInstanceTypeRefNameFieldIndexer(ctx context.Context, idx client.FieldIndexer) error {
	return idx.IndexField(ctx, &corev1alpha1.Instance{}, InstanceSpecInstanceTypeRefNameField, func(obj client.Object) []string {
		instance := obj.(*corev1alpha1.Instance)
		typeRef := instance.Spec.InstanceTypeRef
		return []string{typeRef.Name}
	})
}

const SubnetSpecNetworkRefNameField = "Subnet.spec.networkRef.name"

func SetupSubnetSpecNetworkRefNameField(ctx context.Context, idx client.FieldIndexer) error {
	return idx.IndexField(ctx, &corev1alpha1.Subnet{}, SubnetSpecNetworkRefNameField, func(obj client.Object) []string {
		subnet := obj.(*corev1alpha1.Subnet)
		networkRef := subnet.Spec.NetworkRef
		return []string{networkRef.Name}
	})
}
