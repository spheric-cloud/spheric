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
	PrefixSpecIPFamilyField      = "spec.ipFamily"
	PrefixSpecParentRefNameField = "spec.parentRef.name"
)

func SetupPrefixSpecIPFamilyFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &ipamv1alpha1.Prefix{}, PrefixSpecIPFamilyField, func(obj client.Object) []string {
		prefix := obj.(*ipamv1alpha1.Prefix)
		return []string{string(prefix.Spec.IPFamily)}
	})
}

func SetupPrefixSpecParentRefFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &ipamv1alpha1.Prefix{}, PrefixSpecParentRefNameField, func(obj client.Object) []string {
		prefix := obj.(*ipamv1alpha1.Prefix)
		parentRef := prefix.Spec.ParentRef
		if parentRef == nil {
			return []string{""}
		}
		return []string{parentRef.Name}
	})
}
