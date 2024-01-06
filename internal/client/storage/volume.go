// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
)

const (
	VolumeSpecVolumeClassRefNameField = storagev1alpha1.VolumeVolumeClassRefNameField
	VolumeSpecVolumePoolRefNameField  = storagev1alpha1.VolumeVolumePoolRefNameField
)

func SetupVolumeSpecVolumeClassRefNameFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &storagev1alpha1.Volume{}, VolumeSpecVolumeClassRefNameField, func(obj client.Object) []string {
		volume := obj.(*storagev1alpha1.Volume)
		volumeClassRef := volume.Spec.VolumeClassRef
		if volumeClassRef == nil {
			return []string{""}
		}
		return []string{volumeClassRef.Name}
	})
}

func SetupVolumeSpecVolumePoolRefNameFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &storagev1alpha1.Volume{}, VolumeSpecVolumePoolRefNameField, func(obj client.Object) []string {
		volume := obj.(*storagev1alpha1.Volume)
		volumePoolRef := volume.Spec.VolumePoolRef
		if volumePoolRef == nil {
			return []string{""}
		}
		return []string{volumePoolRef.Name}
	})
}
