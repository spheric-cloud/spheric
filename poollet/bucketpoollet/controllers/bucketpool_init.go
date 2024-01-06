// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	bucketpoolletv1alpha1 "spheric.cloud/spheric/poollet/bucketpoollet/api/v1alpha1"
)

type BucketPoolInit struct {
	client.Client

	BucketPoolName string
	ProviderID     string

	OnInitialized func(ctx context.Context) error
	OnFailed      func(ctx context.Context, reason error) error
}

//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=bucketpools,verbs=get;list;create;update;patch

func (i *BucketPoolInit) Start(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx).WithName("bucketpool").WithName("init")

	log.V(1).Info("Applying bucket pool")
	bucketPool := &storagev1alpha1.BucketPool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: storagev1alpha1.SchemeGroupVersion.String(),
			Kind:       "BucketPool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: i.BucketPoolName,
		},
		Spec: storagev1alpha1.BucketPoolSpec{
			ProviderID: i.ProviderID,
		},
	}
	if err := i.Patch(ctx, bucketPool, client.Apply, client.ForceOwnership, client.FieldOwner(bucketpoolletv1alpha1.FieldOwner)); err != nil {
		if i.OnFailed != nil {
			log.V(1).Info("Failed applying, calling OnFailed callback", "Error", err)
			return i.OnFailed(ctx, err)
		}
		return fmt.Errorf("error applying bucket pool: %w", err)
	}

	log.V(1).Info("Successfully applied bucket pool")
	if i.OnInitialized != nil {
		log.V(1).Info("Calling OnInitialized callback")
		return i.OnInitialized(ctx)
	}
	return nil
}

func (i *BucketPoolInit) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(i)
}
