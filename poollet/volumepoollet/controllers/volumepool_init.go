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
	volumepoolletv1alpha1 "spheric.cloud/spheric/poollet/volumepoollet/api/v1alpha1"
)

type VolumePoolInit struct {
	client.Client

	VolumePoolName string
	ProviderID     string

	OnInitialized func(ctx context.Context) error
	OnFailed      func(ctx context.Context, reason error) error
}

//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=volumepools,verbs=get;list;create;update;patch

func (i *VolumePoolInit) Start(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx).WithName("volumepool").WithName("init")

	log.V(1).Info("Applying volume pool")
	volumePool := &storagev1alpha1.VolumePool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: storagev1alpha1.SchemeGroupVersion.String(),
			Kind:       "VolumePool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: i.VolumePoolName,
		},
		Spec: storagev1alpha1.VolumePoolSpec{
			ProviderID: i.ProviderID,
		},
	}
	if err := i.Patch(ctx, volumePool, client.Apply, client.ForceOwnership, client.FieldOwner(volumepoolletv1alpha1.FieldOwner)); err != nil {
		if i.OnFailed != nil {
			log.V(1).Info("Failed applying, calling OnFailed callback", "Error", err)
			return i.OnFailed(ctx, err)
		}
		return fmt.Errorf("error applying volume pool: %w", err)
	}

	log.V(1).Info("Successfully applied volume pool")
	if i.OnInitialized != nil {
		log.V(1).Info("Calling OnInitialized callback")
		return i.OnInitialized(ctx)
	}
	return nil
}

func (i *VolumePoolInit) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(i)
}
