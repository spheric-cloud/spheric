// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	"spheric.cloud/spheric/client-go/informers"
	"spheric.cloud/spheric/client-go/spheric"
	"spheric.cloud/spheric/internal/quota/evaluator/generic"
	utilsgeneric "spheric.cloud/spheric/utils/generic"
	"spheric.cloud/spheric/utils/quota"
	"spheric.cloud/spheric/utils/quota/resourceaccess"
)

func NewEvaluators(volumeClassCapabilities, bucketClassCapabilities generic.CapabilitiesReader) []quota.Evaluator {
	return []quota.Evaluator{
		NewVolumeEvaluator(volumeClassCapabilities),
		NewBucketEvaluator(bucketClassCapabilities),
	}
}

func extractVolumeClassCapabilities(volumeClass *storagev1alpha1.VolumeClass) corev1alpha1.ResourceList {
	return volumeClass.Capabilities
}

func NewClientVolumeCapabilitiesReader(c client.Client) generic.CapabilitiesReader {
	getter := resourceaccess.NewTypedClientGetter[storagev1alpha1.VolumeClass](c)
	return generic.NewGetterCapabilitiesReader(getter,
		extractVolumeClassCapabilities,
		func(s string) client.ObjectKey { return client.ObjectKey{Name: s} },
	)
}

func NewPrimeLRUVolumeClassCapabilitiesReader(c spheric.Interface, f informers.SharedInformerFactory) generic.CapabilitiesReader {
	getter := resourceaccess.NewPrimeLRUGetter[*storagev1alpha1.VolumeClass, string](
		func(ctx context.Context, className string) (*storagev1alpha1.VolumeClass, error) {
			return c.StorageV1alpha1().VolumeClasses().Get(ctx, className, metav1.GetOptions{})
		},
		func(ctx context.Context, className string) (*storagev1alpha1.VolumeClass, error) {
			return f.Storage().V1alpha1().VolumeClasses().Lister().Get(className)
		},
	)
	return generic.NewGetterCapabilitiesReader(getter, extractVolumeClassCapabilities, utilsgeneric.Identity[string])
}

func extractBucketClassCapabilities(bucketClass *storagev1alpha1.BucketClass) corev1alpha1.ResourceList {
	return bucketClass.Capabilities
}

func NewClientBucketCapabilitiesReader(c client.Client) generic.CapabilitiesReader {
	getter := resourceaccess.NewTypedClientGetter[storagev1alpha1.BucketClass](c)
	return generic.NewGetterCapabilitiesReader(getter,
		extractBucketClassCapabilities,
		func(s string) client.ObjectKey { return client.ObjectKey{Name: s} },
	)
}

func NewPrimeLRUBucketClassCapabilitiesReader(c spheric.Interface, f informers.SharedInformerFactory) generic.CapabilitiesReader {
	getter := resourceaccess.NewPrimeLRUGetter[*storagev1alpha1.BucketClass, string](
		func(ctx context.Context, className string) (*storagev1alpha1.BucketClass, error) {
			return c.StorageV1alpha1().BucketClasses().Get(ctx, className, metav1.GetOptions{})
		},
		func(ctx context.Context, className string) (*storagev1alpha1.BucketClass, error) {
			return f.Storage().V1alpha1().BucketClasses().Lister().Get(className)
		},
	)
	return generic.NewGetterCapabilitiesReader(getter, extractBucketClassCapabilities, utilsgeneric.Identity[string])
}
