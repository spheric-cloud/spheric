// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core

import (
	"context"
	"fmt"

	"spheric.cloud/spheric/internal/controllers/core/scheduler"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	coreclient "spheric.cloud/spheric/internal/client/core"
)

const (
	outOfCapacity = "OutOfCapacity"
)

type InstanceScheduler struct {
	record.EventRecorder
	client.Client

	Cache    *scheduler.Cache
	snapshot *scheduler.Snapshot
}

//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instances,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=instances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.spheric.cloud,resources=fleets,verbs=get;list;watch

// Reconcile reconciles the desired with the actual state.
func (s *InstanceScheduler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	instance := &corev1alpha1.Instance{}
	if err := s.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if s.skipSchedule(log, instance) {
		log.V(1).Info("Skipping scheduling for instance")
		return ctrl.Result{}, nil
	}

	return s.reconcileExists(ctx, log, instance)
}

func (s *InstanceScheduler) skipSchedule(log logr.Logger, instance *corev1alpha1.Instance) bool {
	if !instance.DeletionTimestamp.IsZero() {
		return true
	}

	isAssumed, err := s.Cache.IsAssumedInstance(instance)
	if err != nil {
		log.Error(err, "Error checking whether instance has been assumed")
		return false
	}
	return isAssumed
}

func (s *InstanceScheduler) matchesLabels(ctx context.Context, info *scheduler.ContainerInfo, instance *corev1alpha1.Instance) bool {
	nodeLabels := labels.Set(info.Fleet().Labels)
	fleetSelector := labels.SelectorFromSet(instance.Spec.FleetSelector)

	return fleetSelector.Matches(nodeLabels)
}

func (s *InstanceScheduler) tolerateTaints(ctx context.Context, info *scheduler.ContainerInfo, instance *corev1alpha1.Instance) bool {
	return corev1alpha1.TolerateTaints(instance.Spec.Tolerations, info.Fleet().Spec.Taints)
}

func (s *InstanceScheduler) fitsFleet(ctx context.Context, info *scheduler.ContainerInfo, instance *corev1alpha1.Instance) bool {
	instanceClassName := instance.Spec.InstanceTypeRef.Name

	allocatable, ok := info.Fleet().Status.Allocatable[corev1alpha1.ResourceInstanceType(instanceClassName)]
	if !ok {
		return false
	}

	return allocatable.Cmp(*resource.NewQuantity(1, resource.DecimalSI)) >= 0
}

func (s *InstanceScheduler) reconcileExists(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance) (ctrl.Result, error) {
	s.updateSnapshot()

	fleets := s.snapshot.ListFleets()
	if len(fleets) == 0 {
		s.Event(instance, corev1.EventTypeNormal, outOfCapacity, "No fleets available to schedule instance on")
		return ctrl.Result{}, nil
	}

	var filteredFleets []*scheduler.ContainerInfo
	for _, fleet := range fleets {
		if !s.tolerateTaints(ctx, fleet, instance) {
			log.Info("fleet filtered", "reason", "taints do not match")
			continue
		}
		if !s.matchesLabels(ctx, fleet, instance) {
			log.Info("fleet filtered", "reason", "label do not match")
			continue
		}
		if !s.fitsFleet(ctx, fleet, instance) {
			log.Info("fleet filtered", "reason", "resources do not match")
			continue
		}

		filteredFleets = append(filteredFleets, fleet)
	}

	if len(filteredFleets) == 0 {
		s.Event(instance, corev1.EventTypeNormal, outOfCapacity, "No fleets available after filtering to schedule instance on")
		return ctrl.Result{}, nil
	}

	maxAllocatableFleet := filteredFleets[0]
	for _, fleet := range filteredFleets[1:] {
		if fleet.MaxAllocatable(instance.Spec.InstanceTypeRef.Name) > maxAllocatableFleet.MaxAllocatable(instance.Spec.InstanceTypeRef.Name) {
			maxAllocatableFleet = fleet
		}
	}
	log.V(1).Info("Determined fleet to schedule on", "FleetName", maxAllocatableFleet.Fleet().Name, "Instances", maxAllocatableFleet.NumInstances(), "Allocatable", maxAllocatableFleet.MaxAllocatable(instance.Spec.InstanceTypeRef.Name))

	log.V(1).Info("Assuming instance to be on fleet")
	if err := s.assume(instance, maxAllocatableFleet.Fleet().Name); err != nil {
		return ctrl.Result{}, err
	}

	log.V(1).Info("Running binding asynchronously")
	go func() {
		if err := s.bindingCycle(ctx, log, instance); err != nil {
			if err := s.Cache.ForgetInstance(instance); err != nil {
				log.Error(err, "Error forgetting instance")
			}
		}
	}()
	return ctrl.Result{}, nil
}

func (s *InstanceScheduler) updateSnapshot() {
	if s.snapshot == nil {
		s.snapshot = s.Cache.Snapshot()
	} else {
		s.snapshot.Update()
	}
}

func (s *InstanceScheduler) assume(assumed *corev1alpha1.Instance, nodeName string) error {
	assumed.Spec.FleetRef = corev1alpha1.NewLocalObjRef(nodeName)
	if err := s.Cache.AssumeInstance(assumed); err != nil {
		return err
	}
	return nil
}

func (s *InstanceScheduler) bindingCycle(ctx context.Context, log logr.Logger, assumedInstance *corev1alpha1.Instance) error {
	if err := s.bind(ctx, log, assumedInstance); err != nil {
		return fmt.Errorf("error binding: %w", err)
	}
	log.V(1).Info("Bound instance to fleet")
	return nil
}

func (s *InstanceScheduler) bind(ctx context.Context, log logr.Logger, assumed *corev1alpha1.Instance) error {
	defer func() {
		if err := s.Cache.FinishBinding(assumed); err != nil {
			log.Error(err, "Error finishing cache binding")
		}
	}()

	nonAssumed := assumed.DeepCopy()
	nonAssumed.Spec.FleetRef = nil

	if err := s.Patch(ctx, assumed, client.MergeFrom(nonAssumed)); err != nil {
		return fmt.Errorf("error patching instance: %w", err)
	}
	return nil
}

func (s *InstanceScheduler) enqueueUnscheduledInstances(ctx context.Context, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
	log := ctrl.LoggerFrom(ctx)
	instanceList := &corev1alpha1.InstanceList{}
	if err := s.List(ctx, instanceList, client.MatchingFields{coreclient.InstanceSpecFleetRefNameField: ""}); err != nil {
		log.Error(fmt.Errorf("could not list instances w/o fleet: %w", err), "Error listing fleets")
		return
	}

	for _, instance := range instanceList.Items {
		if !instance.DeletionTimestamp.IsZero() {
			continue
		}
		if instance.Spec.FleetRef != nil {
			continue
		}
		queue.Add(ctrl.Request{NamespacedName: client.ObjectKeyFromObject(&instance)})
	}
}

func (s *InstanceScheduler) isInstanceAssigned() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		instance := obj.(*corev1alpha1.Instance)
		return instance.Spec.FleetRef != nil
	})
}

func (s *InstanceScheduler) isInstanceNotAssigned() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		instance := obj.(*corev1alpha1.Instance)
		return instance.Spec.FleetRef == nil
	})
}

func (s *InstanceScheduler) handleInstance() handler.EventHandler {
	return handler.Funcs{
		CreateFunc: func(ctx context.Context, evt event.CreateEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
			instance := evt.Object.(*corev1alpha1.Instance)
			log := ctrl.LoggerFrom(ctx)

			if err := s.Cache.AddInstance(instance); err != nil {
				log.Error(err, "Error adding instance to cache")
			}
		},
		UpdateFunc: func(ctx context.Context, evt event.UpdateEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
			log := ctrl.LoggerFrom(ctx)

			oldInstance := evt.ObjectOld.(*corev1alpha1.Instance)
			newInstance := evt.ObjectNew.(*corev1alpha1.Instance)
			if err := s.Cache.UpdateInstance(oldInstance, newInstance); err != nil {
				log.Error(err, "Error updating instance in cache")
			}
		},
		DeleteFunc: func(ctx context.Context, evt event.DeleteEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
			log := ctrl.LoggerFrom(ctx)

			instance := evt.Object.(*corev1alpha1.Instance)
			if err := s.Cache.RemoveInstance(instance); err != nil {
				log.Error(err, "Error adding instance to cache")
			}
		},
	}
}

func (s *InstanceScheduler) handleFleet() handler.EventHandler {
	return handler.Funcs{
		CreateFunc: func(ctx context.Context, evt event.CreateEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
			info := evt.Object.(*corev1alpha1.Fleet)
			s.Cache.AddContainer(info)
			s.enqueueUnscheduledInstances(ctx, queue)
		},
		UpdateFunc: func(ctx context.Context, evt event.UpdateEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
			oldFleet := evt.ObjectOld.(*corev1alpha1.Fleet)
			newFleet := evt.ObjectNew.(*corev1alpha1.Fleet)
			s.Cache.UpdateContainer(oldFleet, newFleet)
			s.enqueueUnscheduledInstances(ctx, queue)
		},
		DeleteFunc: func(ctx context.Context, evt event.DeleteEvent, queue workqueue.TypedRateLimitingInterface[ctrl.Request]) {
			log := ctrl.LoggerFrom(ctx)

			info := evt.Object.(*corev1alpha1.Fleet)
			if err := s.Cache.RemoveContainer(info); err != nil {
				log.Error(err, "Error removing fleet from cache")
			}
		},
	}
}

func (s *InstanceScheduler) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("instance-scheduler").
		WithOptions(controller.Options{
			// Only a single concurrent reconcile since it is serialized on the scheduling algorithm's node fitting.
			MaxConcurrentReconciles: 1,
		}).
		// Enqueue unscheduled instances.
		For(&corev1alpha1.Instance{},
			builder.WithPredicates(
				s.isInstanceNotAssigned(),
			),
		).
		Watches(
			&corev1alpha1.Instance{},
			s.handleInstance(),
			builder.WithPredicates(
				s.isInstanceAssigned(),
			),
		).
		// Enqueue unscheduled instances if a fleet w/ required instance classes becomes available.
		Watches(
			&corev1alpha1.Fleet{},
			s.handleFleet(),
		).
		Complete(s)
}
