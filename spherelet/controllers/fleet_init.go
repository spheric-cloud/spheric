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
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/spherelet/api/v1alpha1"
)

type FleetInit struct {
	client.Client

	FleetName  string
	ProviderID string

	// TODO: Remove OnInitialized / OnFailed as soon as the controller-runtime provides support for pre-start hooks:
	// https://github.com/kubernetes-sigs/controller-runtime/pull/2044

	OnInitialized func(ctx context.Context) error
	OnFailed      func(ctx context.Context, reason error) error
}

//+kubebuilder:rbac:groups=core.spheric.cloud,resources=fleets,verbs=get;list;create;update;patch

func (i *FleetInit) Start(ctx context.Context) error {
	log := ctrl.LoggerFrom(ctx).WithName("fleet").WithName("init")

	log.V(1).Info("Applying fleet")
	fleet := &corev1alpha1.Fleet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1alpha1.SchemeGroupVersion.String(),
			Kind:       "Fleet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: i.FleetName,
		},
		Spec: corev1alpha1.FleetSpec{
			ProviderID: i.ProviderID,
		},
	}
	if err := i.Patch(ctx, fleet, client.Apply, client.ForceOwnership, client.FieldOwner(v1alpha1.FieldOwner)); err != nil {
		if i.OnFailed != nil {
			log.V(1).Info("Failed applying, calling OnFailed callback", "Error", err)
			return i.OnFailed(ctx, err)
		}
		return fmt.Errorf("error applying fleet: %w", err)
	}

	log.V(1).Info("Successfully applied fleet")
	if i.OnInitialized != nil {
		log.V(1).Info("Calling OnInitialized callback")
		return i.OnInitialized(ctx)
	}
	return nil
}

func (i *FleetInit) SetupWithManager(mgr ctrl.Manager) error {
	return mgr.Add(i)
}
