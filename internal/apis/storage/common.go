// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import corev1 "k8s.io/api/core/v1"

const (
	VolumeVolumePoolRefNameField  = "spec.volumePoolRef.name"
	VolumeVolumeClassRefNameField = "spec.volumeClassRef.name"

	BucketBucketPoolRefNameField  = "spec.bucketPoolRef.name"
	BucketBucketClassRefNameField = "spec.bucketClassRef.name"

	// VolumePoolsGroup is the system rbac group all volume pools are in.
	VolumePoolsGroup = "storage.spheric.cloud:system:volumepools"

	// VolumePoolUserNamePrefix is the prefix all volume pool users should have.
	VolumePoolUserNamePrefix = "storage.spheric.cloud:system:volumepool:"

	// BucketPoolsGroup is the system rbac group all bucket pools are in.
	BucketPoolsGroup = "storage.spheric.cloud:system:bucketpools"

	// BucketPoolUserNamePrefix is the prefix all bucket pool users should have.
	BucketPoolUserNamePrefix = "storage.spheric.cloud:system:bucketpool:"

	SecretTypeVolumeAuth = corev1.SecretType("storage.spheric.cloud/volume-auth")
)

// VolumePoolCommonName constructs the common name for a certificate of a volume pool user.
func VolumePoolCommonName(name string) string {
	return VolumePoolUserNamePrefix + name
}

// BucketPoolCommonName constructs the common name for a certificate of a bucket pool user.
func BucketPoolCommonName(name string) string {
	return BucketPoolUserNamePrefix + name
}
