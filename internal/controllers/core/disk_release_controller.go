// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"
	"fmt"

	utilpredicate "spheric.cloud/spheric/spherelet/predicate"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/lru"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type DiskReleaseReconciler struct {
	client.Client
	APIReader client.Reader

	AbsenceCache *lru.Cache
}

//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=disks,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instances,verbs=get;list;watch

func (r *DiskReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	disk := &corev1alpha1.Disk{}
	if err := r.Get(ctx, req.NamespacedName, disk); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcileExists(ctx, log, disk)
}

func (r *DiskReleaseReconciler) reconcileExists(ctx context.Context, log logr.Logger, disk *corev1alpha1.Disk) (ctrl.Result, error) {
	if !disk.DeletionTimestamp.IsZero() {
		log.V(1).Info("Disk is already deleting, nothing to do")
		return ctrl.Result{}, nil
	}

	return r.reconcile(ctx, log, disk)
}

func (r *DiskReleaseReconciler) instanceExists(ctx context.Context, disk *corev1alpha1.Disk) (bool, error) {
	instanceRef := disk.Spec.InstanceRef
	if _, ok := r.AbsenceCache.Get(instanceRef.UID); ok {
		return false, nil
	}

	instance := &metav1.PartialObjectMetadata{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1alpha1.SchemeGroupVersion.String(),
			Kind:       "Instance",
		},
	}
	instanceKey := client.ObjectKey{Namespace: disk.Namespace, Name: instanceRef.Name}
	if err := r.APIReader.Get(ctx, instanceKey, instance); err != nil {
		if !apierrors.IsNotFound(err) {
			return false, fmt.Errorf("error getting instance %s: %w", instanceRef.Name, err)
		}

		r.AbsenceCache.Add(instanceRef.UID, nil)
		return false, nil
	}
	return true, nil
}

func (r *DiskReleaseReconciler) releaseDisk(ctx context.Context, disk *corev1alpha1.Disk) error {
	baseNic := disk.DeepCopy()
	disk.Spec.InstanceRef = nil
	if err := r.Patch(ctx, disk, client.StrategicMergeFrom(baseNic, client.MergeFromWithOptimisticLock{})); err != nil {
		return fmt.Errorf("error patching disk: %w", err)
	}
	return nil
}

func (r *DiskReleaseReconciler) reconcile(ctx context.Context, log logr.Logger, disk *corev1alpha1.Disk) (ctrl.Result, error) {
	log.V(1).Info("Reconcile")

	if disk.Spec.InstanceRef == nil {
		log.V(1).Info("Disk is not claimed, nothing to do")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Checking whether instance exists")
	ok, err := r.instanceExists(ctx, disk)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error checking whether instance exists: %w", err)
	}
	if ok {
		log.V(1).Info("Instance is still present")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Instance does not exist, releasing disk")
	if err := r.releaseDisk(ctx, disk); err != nil {
		if !apierrors.IsConflict(err) {
			return ctrl.Result{}, fmt.Errorf("error releasing disk: %w", err)
		}
		log.V(1).Info("Disk was updated, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Reconciled")
	return ctrl.Result{}, nil
}

func (r *DiskReleaseReconciler) diskClaimedPredicate() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		disk := obj.(*corev1alpha1.Disk)
		return disk.Spec.InstanceRef != nil
	})
}

func (r *DiskReleaseReconciler) enqueueByInstance() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []ctrl.Request {
		instance := obj.(*corev1alpha1.Instance)
		log := ctrl.LoggerFrom(ctx)

		diskList := &corev1alpha1.DiskList{}
		if err := r.List(ctx, diskList,
			client.InNamespace(instance.Namespace),
		); err != nil {
			log.Error(err, "Error listing disks")
			return nil
		}

		var reqs []ctrl.Request
		for _, disk := range diskList.Items {
			instanceRef := disk.Spec.InstanceRef
			if instanceRef == nil {
				continue
			}

			if instanceRef.UID != instance.UID {
				continue
			}

			reqs = append(reqs, ctrl.Request{NamespacedName: client.ObjectKeyFromObject(&disk)})
		}
		return reqs
	})
}

func (r *DiskReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("diskrelease").
		For(
			&corev1alpha1.Disk{},
			builder.WithPredicates(r.diskClaimedPredicate()),
		).
		Watches(
			&corev1alpha1.Instance{},
			r.enqueueByInstance(),
			builder.WithPredicates(utilpredicate.IsInDeletionPredicate()),
		).
		Complete(r)
}
