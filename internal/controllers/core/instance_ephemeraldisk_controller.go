// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"
	"errors"
	"fmt"
	"maps"

	utilpredicate "spheric.cloud/spheric/spherelet/predicate"

	coreclient "spheric.cloud/spheric/internal/client/core"

	"spheric.cloud/spheric/utils/annotations"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

type InstanceEphemeralDiskReconciler struct {
	client.Client
}

//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instances,verbs=get;list;watch
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=disks,verbs=get;list;watch;create;update;delete

func (r *InstanceEphemeralDiskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	instance := &corev1alpha1.Instance{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcileExists(ctx, log, instance)
}

func (r *InstanceEphemeralDiskReconciler) reconcileExists(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance) (ctrl.Result, error) {
	if !instance.DeletionTimestamp.IsZero() {
		log.V(1).Info("Instance is deleting, nothing to do")
		return ctrl.Result{}, nil
	}

	return r.reconcile(ctx, log, instance)
}

func (r *InstanceEphemeralDiskReconciler) ephemeralInstanceDiskByName(instance *corev1alpha1.Instance) map[string]*corev1alpha1.Disk {
	res := make(map[string]*corev1alpha1.Disk)
	for _, attachedDisk := range instance.Spec.Disks {
		ephemeral := attachedDisk.Ephemeral
		if ephemeral == nil {
			continue
		}

		diskName := corev1alpha1.InstanceEphemeralDiskName(instance.Name, attachedDisk.Name)
		disk := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:   instance.Namespace,
				Name:        diskName,
				Labels:      ephemeral.DiskTemplate.Labels,
				Annotations: maps.Clone(ephemeral.DiskTemplate.Annotations),
			},
			Spec: ephemeral.DiskTemplate.Spec,
		}
		annotations.SetDefaultEphemeralManagedBy(disk)
		_ = ctrl.SetControllerReference(instance, disk, r.Scheme())
		res[diskName] = disk
	}
	return res
}

func (r *InstanceEphemeralDiskReconciler) handleExistingDisk(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	shouldManage bool,
	disk *corev1alpha1.Disk,
) error {
	if annotations.IsDefaultEphemeralControlledBy(disk, instance) {
		if shouldManage {
			log.V(1).Info("Ephemeral disk is present and controlled by instance")
			return nil
		}

		if !disk.DeletionTimestamp.IsZero() {
			log.V(1).Info("Undesired ephemeral disk is already deleting")
			return nil
		}

		log.V(1).Info("Deleting undesired ephemeral disk")
		if err := r.Delete(ctx, disk); client.IgnoreNotFound(err) != nil {
			return fmt.Errorf("error deleting disk %s: %w", disk.Name, err)
		}
		return nil
	}

	if shouldManage {
		log.V(1).Info("Won't adopt unmanaged disk")
	}
	return nil
}

func (r *InstanceEphemeralDiskReconciler) handleCreateDisk(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	disk *corev1alpha1.Disk,
) error {
	diskKey := client.ObjectKeyFromObject(disk)
	err := r.Create(ctx, disk)
	if err == nil {
		return nil
	}
	if !apierrors.IsAlreadyExists(err) {
		return err
	}

	// Due to a fast resync, we might get an already exists error.
	// In this case, try to fetch the disk again and, when successful, treat it as managing
	// an existing disk.
	if err := r.Get(ctx, diskKey, disk); err != nil {
		return fmt.Errorf("error getting disk %s after already exists: %w", diskKey.Name, err)
	}

	// Treat a retrieved disk as an existing we should manage.
	return r.handleExistingDisk(ctx, log, instance, true, disk)
}

func (r *InstanceEphemeralDiskReconciler) reconcile(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance) (ctrl.Result, error) {
	log.V(1).Info("Reconcile")

	log.V(1).Info("Listing disks")
	diskList := &corev1alpha1.DiskList{}
	if err := r.List(ctx, diskList,
		client.InNamespace(instance.Namespace),
	); err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing disks: %w", err)
	}

	var (
		ephemDiskByName = r.ephemeralInstanceDiskByName(instance)
		errs            []error
	)
	for _, disk := range diskList.Items {
		diskName := disk.Name
		_, shouldManage := ephemDiskByName[diskName]
		delete(ephemDiskByName, diskName)
		log := log.WithValues("Disk", klog.KObj(&disk), "ShouldManage", shouldManage)
		if err := r.handleExistingDisk(ctx, log, instance, shouldManage, &disk); err != nil {
			errs = append(errs, err)
		}
	}

	for _, disk := range ephemDiskByName {
		log := log.WithValues("Disk", klog.KObj(disk))
		if err := r.handleCreateDisk(ctx, log, instance, disk); err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return ctrl.Result{}, fmt.Errorf("error managing ephemeral disks: %w", err)
	}

	log.V(1).Info("Reconciled")
	return ctrl.Result{}, nil
}

func (r *InstanceEphemeralDiskReconciler) enqueueByDisk() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []ctrl.Request {
		disk := obj.(*corev1alpha1.Disk)
		log := ctrl.LoggerFrom(ctx)

		instanceList := &corev1alpha1.InstanceList{}
		if err := r.List(ctx, instanceList,
			client.InNamespace(disk.Namespace),
			client.MatchingFields{
				coreclient.InstanceSpecDiskNamesField: disk.Name,
			},
		); err != nil {
			log.Error(err, "Error listing instances")
			return nil
		}

		var reqs []ctrl.Request
		for _, instance := range instanceList.Items {
			if !instance.DeletionTimestamp.IsZero() {
				continue
			}

			reqs = append(reqs, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(&instance)})
		}
		return reqs
	})
}

func (r *InstanceEphemeralDiskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("instanceephemeraldisk").
		For(
			&corev1alpha1.Instance{},
			builder.WithPredicates(
				utilpredicate.IsNotInDeletionPredicate(),
			),
		).
		Owns(
			&corev1alpha1.Disk{},
		).
		Watches(
			&corev1alpha1.Disk{},
			r.enqueueByDisk(),
		).
		Complete(r)
}
