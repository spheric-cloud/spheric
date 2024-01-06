// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package ipam

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	ipamv1alpha1 "spheric.cloud/spheric/api/ipam/v1alpha1"
)

const (
	PrefixAllocationSpecIPFamilyField      = "spec.ipFamily"
	PrefixAllocationSpecPrefixRefNameField = "spec.prefixRef.name"
)

func SetupPrefixAllocationSpecIPFamilyFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &ipamv1alpha1.PrefixAllocation{}, PrefixAllocationSpecIPFamilyField, func(obj client.Object) []string {
		prefixAllocation := obj.(*ipamv1alpha1.PrefixAllocation)
		return []string{string(prefixAllocation.Spec.IPFamily)}
	})
}

func SetupPrefixAllocationSpecPrefixRefNameField(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &ipamv1alpha1.PrefixAllocation{}, PrefixAllocationSpecPrefixRefNameField, func(obj client.Object) []string {
		allocation := obj.(*ipamv1alpha1.PrefixAllocation)
		prefixRef := allocation.Spec.PrefixRef
		if prefixRef == nil {
			return nil
		}
		return []string{prefixRef.Name}
	})
}
