// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	"spheric.cloud/spheric/spherelet/instance"
	"spheric.cloud/spheric/utils/generic"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/resource"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
)

type FleetReconciler struct {
	client.Client

	// FleetName is the name of the computev1alpha1.Fleet to report / update.
	FleetName string
	// Addresses are the addresses the spherelet server is available on.
	Addresses []corev1alpha1.FleetAddress
	// Port is the port the spherelet server is available on.
	Port int32

	InstanceRuntime instance.RuntimeService
}

//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instancepools,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instancepools/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instanceclasses,verbs=get;list;watch

func (r *FleetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	fleet := &corev1alpha1.Fleet{}
	if err := r.Get(ctx, req.NamespacedName, fleet); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcileExists(ctx, log, fleet)
}

func (r *FleetReconciler) reconcileExists(ctx context.Context, log logr.Logger, fleet *corev1alpha1.Fleet) (ctrl.Result, error) {
	if !fleet.DeletionTimestamp.IsZero() {
		return r.delete(ctx, log, fleet)
	}
	return r.reconcile(ctx, log, fleet)
}

func (r *FleetReconciler) delete(ctx context.Context, log logr.Logger, fleet *corev1alpha1.Fleet) (ctrl.Result, error) {
	log.V(1).Info("Delete")
	log.V(1).Info("Deleted")
	return ctrl.Result{}, nil
}

func getFleetResources(iriRes *iri.RuntimeResources) corev1alpha1.ResourceList {
	res := make(corev1alpha1.ResourceList)
	res[corev1alpha1.ResourceCPU] = *resource.NewQuantity(iriRes.CpuCount, resource.DecimalSI)
	res[corev1alpha1.ResourceMemory] = *resource.NewQuantity(int64(iriRes.MemoryBytes), resource.BinarySI)
	for name, count := range iriRes.InstanceQuantities {
		res[corev1alpha1.ResourceInstanceType(name)] = *resource.NewQuantity(count, resource.DecimalSI)
	}
	return res
}

func (r *FleetReconciler) calculateCapacity(
	ctx context.Context,
	log logr.Logger,
) (capacity, allocatable corev1alpha1.ResourceList, err error) {
	log.V(1).Info("Determining supported instance classes, capacity and allocatable")

	res, err := r.InstanceRuntime.Status(ctx, &iri.StatusRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("error getting instance status: %w", err)
	}

	capacity = getFleetResources(generic.PointerOrNew(res.Capacity))
	allocatable = getFleetResources(generic.PointerOrNew(res.Allocatable))
	return capacity, allocatable, nil
}

func (r *FleetReconciler) updateStatus(ctx context.Context, log logr.Logger, fleet *corev1alpha1.Fleet) error {
	capacity, allocatable, err := r.calculateCapacity(ctx, log)
	if err != nil {
		return fmt.Errorf("error calculating pool resources:%w", err)
	}

	base := fleet.DeepCopy()
	fleet.Status.State = corev1alpha1.FleetStateReady
	fleet.Status.Addresses = r.Addresses
	fleet.Status.Capacity = capacity
	fleet.Status.Allocatable = allocatable
	fleet.Status.DaemonEndpoints.SphereletEndpoint.Port = r.Port

	if err := r.Status().Patch(ctx, fleet, client.MergeFrom(base)); err != nil {
		return fmt.Errorf("error patching instance pool status: %w", err)
	}

	return nil
}

func (r *FleetReconciler) reconcile(ctx context.Context, log logr.Logger, fleet *corev1alpha1.Fleet) (ctrl.Result, error) {
	log.V(1).Info("Reconcile")

	log.V(1).Info("Ensuring no reconcile annotation")
	modified, err := sphericclient.PatchEnsureNoReconcileAnnotation(ctx, r.Client, fleet)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring no reconcile annotation: %w", err)
	}
	if modified {
		log.V(1).Info("Removed reconcile annotation, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Updating instance pool status")
	if err := r.updateStatus(ctx, log, fleet); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating status: %w", err)
	}

	log.V(1).Info("Reconciled")
	return ctrl.Result{}, nil
}

func (r *FleetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(
			&corev1alpha1.Fleet{},
			builder.WithPredicates(
				predicate.NewPredicateFuncs(func(obj client.Object) bool {
					return obj.GetName() == r.FleetName
				}),
			),
		).
		Watches(
			&corev1alpha1.InstanceType{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []ctrl.Request {
				return []ctrl.Request{{NamespacedName: client.ObjectKey{Name: r.FleetName}}}
			}),
		).
		Watches(
			&corev1alpha1.Instance{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []ctrl.Request {
				return []ctrl.Request{{NamespacedName: client.ObjectKey{Name: r.FleetName}}}
			}),
			builder.WithPredicates(
				InstanceRunsInFleetPredicate(r.FleetName),
			),
		).
		Complete(r)
}
