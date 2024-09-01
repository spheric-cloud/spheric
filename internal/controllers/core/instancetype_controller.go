// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"
	"fmt"
	"sort"

	"spheric.cloud/spheric/api/core/v1alpha1"

	"github.com/go-logr/logr"
	"github.com/ironcore-dev/controller-utils/clientutils"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	coreclient "spheric.cloud/spheric/internal/client/core"
	"spheric.cloud/spheric/utils/slices"
)

// InstanceTypeReconciler reconciles a InstanceTypeRef object
type InstanceTypeReconciler struct {
	client.Client
	APIReader client.Reader
}

//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instancetypes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instancetypes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instancetypes/finalizers,verbs=update
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instances,verbs=get;list;watch

// Reconcile moves the current state of the cluster closer to the desired state
func (r *InstanceTypeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	instanceType := &v1alpha1.InstanceType{}
	if err := r.Get(ctx, req.NamespacedName, instanceType); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return r.reconcileExists(ctx, log, instanceType)
}

func (r *InstanceTypeReconciler) listReferencingInstancesWithReader(
	ctx context.Context,
	rd client.Reader,
	instanceType *v1alpha1.InstanceType,
) ([]v1alpha1.Instance, error) {
	instanceList := &v1alpha1.InstanceList{}
	if err := rd.List(ctx, instanceList,
		client.InNamespace(instanceType.Namespace),
		client.MatchingFields{coreclient.InstanceSpecInstanceTypeRefNameField: instanceType.Name},
	); err != nil {
		return nil, fmt.Errorf("error listing the instances using the instance type: %w", err)
	}

	return instanceList.Items, nil
}

func (r *InstanceTypeReconciler) collectInstanceNames(instances []v1alpha1.Instance) []string {
	instanceNames := slices.MapRef(instances, func(instance *v1alpha1.Instance) string {
		return instance.Name
	})
	sort.Strings(instanceNames)
	return instanceNames
}

func (r *InstanceTypeReconciler) delete(ctx context.Context, log logr.Logger, instanceType *v1alpha1.InstanceType) (ctrl.Result, error) {
	if !controllerutil.ContainsFinalizer(instanceType, v1alpha1.InstanceTypeFinalizer) {
		return ctrl.Result{}, nil
	}

	instances, err := r.listReferencingInstancesWithReader(ctx, r.Client, instanceType)
	if err != nil {
		return ctrl.Result{}, err
	}
	if len(instances) > 0 {
		log.V(1).Info("Instance type is still in use", "ReferencingInstanceNames", r.collectInstanceNames(instances))
		return ctrl.Result{Requeue: true}, nil
	}

	instances, err = r.listReferencingInstancesWithReader(ctx, r.APIReader, instanceType)
	if err != nil {
		return ctrl.Result{}, err
	}
	if len(instances) > 0 {
		log.V(1).Info("Instance type is still in use", "ReferencingInstanceNames", r.collectInstanceNames(instances))
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Instance type is not in use anymore, removing finalizer")
	if err := clientutils.PatchRemoveFinalizer(ctx, r.Client, instanceType, v1alpha1.InstanceTypeFinalizer); err != nil {
		return ctrl.Result{}, err
	}

	log.V(1).Info("Successfully removed finalizer")
	return ctrl.Result{}, nil
}

func (r *InstanceTypeReconciler) reconcile(ctx context.Context, log logr.Logger, instanceType *v1alpha1.InstanceType) (ctrl.Result, error) {
	log.V(1).Info("Ensuring finalizer")
	if modified, err := clientutils.PatchEnsureFinalizer(ctx, r.Client, instanceType, v1alpha1.InstanceTypeFinalizer); err != nil || modified {
		return ctrl.Result{}, err
	}

	log.V(1).Info("Finalizer is present")
	return ctrl.Result{}, nil
}

func (r *InstanceTypeReconciler) reconcileExists(ctx context.Context, log logr.Logger, instanceType *v1alpha1.InstanceType) (ctrl.Result, error) {
	if !instanceType.DeletionTimestamp.IsZero() {
		return r.delete(ctx, log, instanceType)
	}
	return r.reconcile(ctx, log, instanceType)
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceTypeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.InstanceType{}).
		Watches(
			&v1alpha1.Instance{},
			handler.Funcs{
				DeleteFunc: func(ctx context.Context, event event.DeleteEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
					instance := event.Object.(*v1alpha1.Instance)
					queue.Add(ctrl.Request{NamespacedName: types.NamespacedName{Name: instance.Spec.InstanceTypeRef.Name}})
				},
			},
		).
		Complete(r)
}
