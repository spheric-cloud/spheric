// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/ironcore-dev/controller-utils/clientutils"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/kubectl/pkg/util/fieldpath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/spherelet/api/v1alpha1"
	sphereletclient "spheric.cloud/spheric/spherelet/client"
	"spheric.cloud/spheric/spherelet/controllers/events"
	iriinstance "spheric.cloud/spheric/spherelet/instance"
	utilclient "spheric.cloud/spheric/utils/client"
	"spheric.cloud/spheric/utils/predicates"
)

type InstanceReconciler struct {
	record.EventRecorder
	client.Client

	InstanceRuntime        iriinstance.RuntimeService
	InstanceRuntimeName    string
	InstanceRuntimeVersion string

	FleetName string

	DownwardAPILabels      map[string]string
	DownwardAPIAnnotations map[string]string

	WatchFilterValue string
}

func (r *InstanceReconciler) instanceKeyLabelSelector(instanceKey client.ObjectKey) map[string]string {
	return map[string]string{
		v1alpha1.InstanceNamespaceLabel: instanceKey.Namespace,
		v1alpha1.InstanceNameLabel:      instanceKey.Name,
	}
}

func (r *InstanceReconciler) instanceUIDLabelSelector(instanceUID types.UID) map[string]string {
	return map[string]string{
		v1alpha1.InstanceUIDLabel: string(instanceUID),
	}
}

//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instances,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=compute.spheric.cloud,resources=instances/finalizers,verbs=update
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=disks,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=networking.spheric.cloud,resources=networkinterfaces,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=networking.spheric.cloud,resources=networks,verbs=get;list;watch

func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	instance := &corev1alpha1.Instance{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("error getting instance %s: %w", req.NamespacedName, err)
		}
		return r.deleteGone(ctx, log, req.NamespacedName)
	}
	return r.reconcileExists(ctx, log, instance)
}

func (r *InstanceReconciler) getIRIInstancesForInstance(ctx context.Context, instance *corev1alpha1.Instance) ([]*iri.Instance, error) {
	res, err := r.InstanceRuntime.ListInstances(ctx, &iri.ListInstancesRequest{
		Filter: &iri.InstanceFilter{LabelSelector: r.instanceUIDLabelSelector(instance.UID)},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing instances by instance uid: %w", err)
	}
	return res.Instances, nil
}

func (r *InstanceReconciler) listInstancesByInstanceKey(ctx context.Context, instanceKey client.ObjectKey) ([]*iri.Instance, error) {
	res, err := r.InstanceRuntime.ListInstances(ctx, &iri.ListInstancesRequest{
		Filter: &iri.InstanceFilter{LabelSelector: r.instanceKeyLabelSelector(instanceKey)},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing instances by instance key: %w", err)
	}
	return res.Instances, nil
}

func (r *InstanceReconciler) getInstanceByID(ctx context.Context, id string) (*iri.Instance, error) {
	res, err := r.InstanceRuntime.ListInstances(ctx, &iri.ListInstancesRequest{
		Filter: &iri.InstanceFilter{Id: id},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing instances filtering by id: %w", err)
	}

	switch len(res.Instances) {
	case 0:
		return nil, status.Errorf(codes.NotFound, "instance %s not found", id)
	case 1:
		return res.Instances[0], nil
	default:
		return nil, fmt.Errorf("multiple instances found for id %s", id)
	}
}

func (r *InstanceReconciler) deleteInstances(ctx context.Context, log logr.Logger, instances []*iri.Instance) (bool, error) {
	var (
		errs        []error
		deletingIDs []string
	)
	for _, instance := range instances {
		instanceID := instance.Metadata.Id
		log := log.WithValues("InstanceID", instanceID)
		log.V(1).Info("Deleting matching instance")
		if _, err := r.InstanceRuntime.DeleteInstance(ctx, &iri.DeleteInstanceRequest{
			InstanceId: instanceID,
		}); err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("error deleting instance %s: %w", instanceID, err))
			} else {
				log.V(1).Info("Instance is already gone")
			}
		} else {
			log.V(1).Info("Issued instance deletion")
			deletingIDs = append(deletingIDs, instanceID)
		}
	}

	switch {
	case len(errs) > 0:
		return false, fmt.Errorf("error(s) deleting matching instance(s): %v", errs)
	case len(deletingIDs) > 0:
		log.V(1).Info("Instances are still deleting", "DeletingIDs", deletingIDs)
		return false, nil
	default:
		log.V(1).Info("No instance present")
		return true, nil
	}
}

func (r *InstanceReconciler) deleteGone(ctx context.Context, log logr.Logger, instanceKey client.ObjectKey) (ctrl.Result, error) {
	log.V(1).Info("Delete gone")

	log.V(1).Info("Listing instances by instance key")
	instances, err := r.listInstancesByInstanceKey(ctx, instanceKey)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing instances: %w", err)
	}

	ok, err := r.deleteInstances(ctx, log, instances)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting instances: %w", err)
	}
	if !ok {
		log.V(1).Info("Not all instances are gone yet, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}
	log.V(1).Info("Deleted gone")
	return ctrl.Result{}, nil
}

func (r *InstanceReconciler) reconcileExists(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance) (ctrl.Result, error) {
	if !instance.DeletionTimestamp.IsZero() {
		return r.delete(ctx, log, instance)
	}
	return r.reconcile(ctx, log, instance)
}

func (r *InstanceReconciler) delete(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance) (ctrl.Result, error) {
	log.V(1).Info("Delete")

	if !controllerutil.ContainsFinalizer(instance, v1alpha1.InstanceFinalizer) {
		log.V(1).Info("No finalizer present, nothing to do")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Finalizer present")

	log.V(1).Info("Deleting instances by UID")
	ok, err := r.deleteInstancesByInstanceUID(ctx, log, instance.UID)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting instances: %w", err)
	}
	if !ok {
		log.V(1).Info("Not all instances are gone, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Deleted iri instances by UID, removing finalizer")
	if err := clientutils.PatchRemoveFinalizer(ctx, r.Client, instance, v1alpha1.InstanceFinalizer); err != nil {
		return ctrl.Result{}, fmt.Errorf("error removing finalizer: %w", err)
	}

	log.V(1).Info("Deleted")
	return ctrl.Result{}, nil
}

func (r *InstanceReconciler) deleteInstancesByInstanceUID(ctx context.Context, log logr.Logger, instanceUID types.UID) (bool, error) {
	log.V(1).Info("Listing instances")
	res, err := r.InstanceRuntime.ListInstances(ctx, &iri.ListInstancesRequest{
		Filter: &iri.InstanceFilter{
			LabelSelector: map[string]string{
				v1alpha1.InstanceUIDLabel: string(instanceUID),
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("error listing instances: %w", err)
	}

	log.V(1).Info("Listed instances", "NoOfInstances", len(res.Instances))
	var (
		errs                []error
		deletingInstanceIDs []string
	)
	for _, instance := range res.Instances {
		instanceID := instance.Metadata.Id
		log := log.WithValues("InstanceID", instanceID)
		log.V(1).Info("Deleting instance")
		_, err := r.InstanceRuntime.DeleteInstance(ctx, &iri.DeleteInstanceRequest{
			InstanceId: instanceID,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("error deleting instance %s: %w", instanceID, err))
			} else {
				log.V(1).Info("Instance is already gone")
			}
		} else {
			log.V(1).Info("Issued instance deletion")
			deletingInstanceIDs = append(deletingInstanceIDs, instanceID)
		}
	}

	switch {
	case len(errs) > 0:
		return false, fmt.Errorf("error(s) deleting instance(s): %v", errs)
	case len(deletingInstanceIDs) > 0:
		log.V(1).Info("Instances are in deletion", "DeletingInstanceIDs", deletingInstanceIDs)
		return false, nil
	default:
		log.V(1).Info("All instances are gone")
		return true, nil
	}
}

func (r *InstanceReconciler) reconcile(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance) (ctrl.Result, error) {
	log.V(1).Info("Reconcile")

	log.V(1).Info("Ensuring finalizer")
	modified, err := clientutils.PatchEnsureFinalizer(ctx, r.Client, instance, v1alpha1.InstanceFinalizer)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring finalizer: %w", err)
	}
	if modified {
		log.V(1).Info("Added finalizer, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}
	log.V(1).Info("Finalizer is present")

	log.V(1).Info("Ensuring no reconcile annotation")
	modified, err = utilclient.PatchEnsureNoReconcileAnnotation(ctx, r.Client, instance)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring no reconcile annotation: %w", err)
	}
	if modified {
		log.V(1).Info("Removed reconcile annotation, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	disks, err := r.getDisksForInstance(ctx, instance)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting disks for instance: %w", err)
	}

	iriInstances, err := r.getIRIInstancesForInstance(ctx, instance)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting IRI instances for instance: %w", err)
	}

	switch len(iriInstances) {
	case 0:
		return r.create(ctx, log, instance, disks)
	case 1:
		iriInstance := iriInstances[0]
		return r.update(ctx, log, instance, iriInstance, disks)
	default:
		panic("unhandled: multiple IRI instances")
	}
}

func (r *InstanceReconciler) iriInstanceLabels(instance *corev1alpha1.Instance) (map[string]string, error) {
	annotations := map[string]string{
		v1alpha1.InstanceUIDLabel:       string(instance.UID),
		v1alpha1.InstanceNamespaceLabel: instance.Namespace,
		v1alpha1.InstanceNameLabel:      instance.Name,
	}

	for name, fieldPath := range r.DownwardAPILabels {
		value, err := fieldpath.ExtractFieldPathAsString(instance, fieldPath)
		if err != nil {
			return nil, fmt.Errorf("error extracting downward api label %q: %w", name, err)
		}

		annotations[v1alpha1.DownwardAPILabel(name)] = value
	}
	return annotations, nil
}

func (r *InstanceReconciler) iriInstanceAnnotations(
	instance *corev1alpha1.Instance,
	iriInstanceGeneration int64,
) (map[string]string, error) {
	annotations := map[string]string{
		v1alpha1.InstanceGenerationAnnotation:    strconv.FormatInt(instance.Generation, 10),
		v1alpha1.IRIInstanceGenerationAnnotation: strconv.FormatInt(iriInstanceGeneration, 10),
	}

	for name, fieldPath := range r.DownwardAPIAnnotations {
		value, err := fieldpath.ExtractFieldPathAsString(instance, fieldPath)
		if err != nil {
			return nil, fmt.Errorf("error extracting downward api annotation %q: %w", name, err)
		}

		annotations[v1alpha1.DownwardAPIAnnotation(name)] = value
	}

	return annotations, nil
}

func (r *InstanceReconciler) create(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	disks []corev1alpha1.Disk,
) (ctrl.Result, error) {
	log.V(1).Info("Create")

	log.V(1).Info("Getting instance config")
	iriInstance, ok, err := r.prepareIRIInstance(ctx, instance, disks)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error preparing iri instance: %w", err)
	}
	if !ok {
		log.V(1).Info("Instance is not yet ready")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Creating instance")
	res, err := r.InstanceRuntime.CreateInstance(ctx, &iri.CreateInstanceRequest{
		Instance: iriInstance,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error creating instance: %w", err)
	}
	log.V(1).Info("Created", "InstanceID", res.Instance.Metadata.Id)

	log.V(1).Info("Updating status")
	if err := r.updateStatus(ctx, log, instance, res.Instance); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating instance status: %w", err)
	}

	log.V(1).Info("Created")
	return ctrl.Result{}, nil
}

func (r *InstanceReconciler) getInstanceGeneration(iriInstance *iri.Instance) (int64, error) {
	return getAndParseFromStringMap(iriInstance.GetMetadata().GetAnnotations(),
		v1alpha1.InstanceGenerationAnnotation,
		parseInt64,
	)
}

func (r *InstanceReconciler) getIRIInstanceGeneration(iriInstance *iri.Instance) (int64, error) {
	return getAndParseFromStringMap(iriInstance.GetMetadata().GetAnnotations(),
		v1alpha1.IRIInstanceGenerationAnnotation,
		parseInt64,
	)
}

func (r *InstanceReconciler) updateStatus(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	iriInstance *iri.Instance,
) error {
	requiredIRIGeneration, err := r.getIRIInstanceGeneration(iriInstance)
	if err != nil {
		return err
	}

	iriGeneration := iriInstance.Metadata.Generation
	observedIRIGeneration := iriInstance.Status.ObservedGeneration

	if observedIRIGeneration < requiredIRIGeneration {
		log.V(1).Info("IRI instance was not observed at the latest generation",
			"IRIGeneration", iriGeneration,
			"ObservedIRIGeneration", observedIRIGeneration,
			"RequiredIRIGeneration", requiredIRIGeneration,
		)
		return nil
	}

	if err := r.updateInstanceStatus(ctx, instance, iriInstance); err != nil {
		return fmt.Errorf("error updating instance status: %w", err)
	}

	return nil
}

func (r *InstanceReconciler) updateInstanceStatus(ctx context.Context, instance *corev1alpha1.Instance, iriInstance *iri.Instance) error {
	now := metav1.Now()

	generation, err := r.getInstanceGeneration(iriInstance)
	if err != nil {
		return err
	}

	instanceID := iriinstance.MakeID(r.InstanceRuntimeName, iriInstance.Metadata.Id)

	state, err := r.convertIRIInstanceState(iriInstance.Status.State)
	if err != nil {
		return err
	}

	diskStatuses, err := r.getDiskStatusesForInstance(instance, iriInstance, now)
	if err != nil {
		return fmt.Errorf("error getting disk statuses: %w", err)
	}

	nicStatuses, err := r.getNetworkInterfaceStatusesForInstance(instance, iriInstance, now)
	if err != nil {
		return fmt.Errorf("error getting network interface statuses: %w", err)
	}

	base := instance.DeepCopy()

	instance.Status.State = state
	instance.Status.InstanceID = instanceID.String()
	instance.Status.ObservedGeneration = generation
	instance.Status.Disks = diskStatuses
	instance.Status.NetworkInterfaces = nicStatuses

	if err := r.Status().Patch(ctx, instance, client.MergeFrom(base)); err != nil {
		return fmt.Errorf("error patching status: %w", err)
	}
	return nil
}

func (r *InstanceReconciler) prepareIRIPower(power corev1alpha1.Power) (iri.Power, error) {
	switch power {
	case corev1alpha1.PowerOn:
		return iri.Power_POWER_ON, nil
	case corev1alpha1.PowerOff:
		return iri.Power_POWER_OFF, nil
	default:
		return 0, fmt.Errorf("unknown power %q", power)
	}
}

func (r *InstanceReconciler) updateIRIPower(ctx context.Context, log logr.Logger, instance *corev1alpha1.Instance, iriInstance *iri.Instance) error {
	actualPower := iriInstance.Spec.Power
	desiredPower, err := r.prepareIRIPower(instance.Spec.Power)
	if err != nil {
		return fmt.Errorf("error preparing iri power state: %w", err)
	}

	if actualPower == desiredPower {
		log.V(1).Info("Power is up-to-date", "Power", actualPower)
		return nil
	}

	if _, err := r.InstanceRuntime.UpdateInstancePower(ctx, &iri.UpdateInstancePowerRequest{
		InstanceId: iriInstance.Metadata.Id,
		Power:      desiredPower,
	}); err != nil {
		return fmt.Errorf("error updating instance power state: %w", err)
	}
	return nil
}

func (r *InstanceReconciler) update(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	iriInstance *iri.Instance,
	disks []corev1alpha1.Disk,
) (ctrl.Result, error) {
	log.V(1).Info("Updating existing instance")

	var errs []error

	log.V(1).Info("Updating network interfaces")
	if err := r.updateIRINetworkInterfaces(ctx, log, instance, iriInstance); err != nil {
		errs = append(errs, fmt.Errorf("error updating network interfaces: %w", err))
	}

	log.V(1).Info("Updating disks")
	if err := r.updateIRIDisks(ctx, log, instance, iriInstance, disks); err != nil {
		errs = append(errs, fmt.Errorf("error updating disks: %w", err))
	}

	log.V(1).Info("Updating power state")
	if err := r.updateIRIPower(ctx, log, instance, iriInstance); err != nil {
		errs = append(errs, fmt.Errorf("error updating power state: %w", err))
	}

	if len(errs) > 0 {
		return ctrl.Result{}, fmt.Errorf("error(s) updating instance: %v", errs)
	}

	log.V(1).Info("Updating annotations")
	if err := r.updateIRIAnnotations(ctx, log, instance, iriInstance); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating annotations: %w", err)
	}

	log.V(1).Info("Getting iri instance")
	iriInstance, err := r.getInstanceByID(ctx, iriInstance.Metadata.Id)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error getting iri instance: %w", err)
	}

	log.V(1).Info("Updating instance status")
	if err := r.updateStatus(ctx, log, instance, iriInstance); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating status: %w", err)
	}

	log.V(1).Info("Updated existing instance")
	return ctrl.Result{}, nil
}

func (r *InstanceReconciler) updateIRIAnnotations(
	ctx context.Context,
	log logr.Logger,
	instance *corev1alpha1.Instance,
	iriInstance *iri.Instance,
) error {
	desiredAnnotations, err := r.iriInstanceAnnotations(instance, iriInstance.GetMetadata().GetGeneration())
	if err != nil {
		return fmt.Errorf("error getting iri instance annotations: %w", err)
	}

	actualAnnotations := iriInstance.Metadata.Annotations

	if maps.Equal(desiredAnnotations, actualAnnotations) {
		log.V(1).Info("Annotations are up-to-date", "Annotations", desiredAnnotations)
		return nil
	}

	if _, err := r.InstanceRuntime.UpdateInstanceAnnotations(ctx, &iri.UpdateInstanceAnnotationsRequest{
		InstanceId:  iriInstance.Metadata.Id,
		Annotations: desiredAnnotations,
	}); err != nil {
		return fmt.Errorf("error updating instance annotations: %w", err)
	}
	return nil
}

var iriInstanceStateToInstanceState = map[iri.InstanceState]corev1alpha1.InstanceState{
	iri.InstanceState_INSTANCE_PENDING:    corev1alpha1.InstanceStatePending,
	iri.InstanceState_INSTANCE_RUNNING:    corev1alpha1.InstanceStateRunning,
	iri.InstanceState_INSTANCE_SUSPENDED:  corev1alpha1.InstanceStateShutdown,
	iri.InstanceState_INSTANCE_TERMINATED: corev1alpha1.InstanceStateTerminated,
}

func (r *InstanceReconciler) convertIRIInstanceState(state iri.InstanceState) (corev1alpha1.InstanceState, error) {
	if res, ok := iriInstanceStateToInstanceState[state]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown instance state %v", state)
}

func (r *InstanceReconciler) prepareIRIInstanceReqs(ctx context.Context, instance *corev1alpha1.Instance) (typeName string, cpu int64, memory uint64, ok bool, err error) {
	instanceType := &corev1alpha1.InstanceType{}
	instanceTypeKey := client.ObjectKey{Name: instance.Spec.InstanceTypeRef.Name}
	if err := r.Get(ctx, instanceTypeKey, instanceType); err != nil {
		if !apierrors.IsNotFound(err) {
			return "", 0, 0, false, fmt.Errorf("error getting instance type: %w", err)
		}

		r.Eventf(instance, corev1.EventTypeNormal, events.InstanceTypeNotReady, "Instance type %s is not ready: %v", instanceTypeKey.Name, err)
		return "", 0, 0, false, nil
	}

	typeName = instanceType.Name
	cpu = instanceType.Capabilities.CPU().Value()
	memory = uint64(instanceType.Capabilities.Memory().Value())
	return typeName, cpu, memory, true, nil
}

func (r *InstanceReconciler) prepareIRIIgnitionData(ctx context.Context, instance *corev1alpha1.Instance, ignitionRef *corev1alpha1.SecretKeySelector) ([]byte, bool, error) {
	ignitionSecret := &corev1.Secret{}
	ignitionSecretKey := client.ObjectKey{Namespace: instance.Namespace, Name: ignitionRef.Name}
	if err := r.Get(ctx, ignitionSecretKey, ignitionSecret); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, false, err
		}

		r.Eventf(instance, corev1.EventTypeNormal, events.IgnitionNotReady, "Ignition not ready: %v", err)
		return nil, false, nil
	}

	ignitionKey := ignitionRef.Key
	if ignitionKey == "" {
		ignitionKey = corev1alpha1.DefaultIgnitionKey
	}

	data, ok := ignitionSecret.Data[ignitionKey]
	if !ok {
		r.Eventf(instance, corev1.EventTypeNormal, events.IgnitionNotReady, "Ignition has no data at key %s", ignitionKey)
		return nil, false, nil
	}

	return data, true, nil
}

func (r *InstanceReconciler) prepareIRIInstance(
	ctx context.Context,
	instance *corev1alpha1.Instance,
	disks []corev1alpha1.Disk,
) (*iri.Instance, bool, error) {
	var (
		ok   = true
		errs []error
	)

	typeName, cpu, memory, reqsOk, err := r.prepareIRIInstanceReqs(ctx, instance)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing iri instance requirements: %w", err))
	case !reqsOk:
		ok = false
	}

	var imageSpec *iri.ImageSpec
	if image := instance.Spec.Image; image != "" {
		imageSpec = &iri.ImageSpec{
			Image: image,
		}
	}

	var ignitionData []byte
	if ignitionRef := instance.Spec.IgnitionRef; ignitionRef != nil {
		data, ignitionSpecOK, err := r.prepareIRIIgnitionData(ctx, instance, ignitionRef)
		switch {
		case err != nil:
			errs = append(errs, fmt.Errorf("error preparing iri ignition spec: %w", err))
		case !ignitionSpecOK:
			ok = false
		default:
			ignitionData = data
		}
	}

	instanceNics, instanceNicsOK, err := r.prepareIRINetworkInterfacesForInstance(ctx, instance)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing iri instance network interfaces: %w", err))
	case !instanceNicsOK:
		ok = false
	}

	instanceDisks, instanceDisksOK, err := r.prepareIRIDisks(ctx, instance, disks)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing iri instance disks: %w", err))
	case !instanceDisksOK:
		ok = false
	}

	labels, err := r.iriInstanceLabels(instance)
	if err != nil {
		errs = append(errs, fmt.Errorf("error preparing iri instance labels: %w", err))
	}

	annotations, err := r.iriInstanceAnnotations(instance, 1)
	if err != nil {
		errs = append(errs, fmt.Errorf("error preparing iri instance annotations: %w", err))
	}

	switch {
	case len(errs) > 0:
		return nil, false, fmt.Errorf("error(s) preparing instance: %v", errs)
	case !ok:
		return nil, false, nil
	default:
		return &iri.Instance{
			Metadata: &iri.ObjectMetadata{
				Labels:      labels,
				Annotations: annotations,
			},
			Spec: &iri.InstanceSpec{
				Image:             imageSpec,
				Type:              typeName,
				CpuCount:          cpu,
				MemoryBytes:       memory,
				IgnitionData:      ignitionData,
				Disks:             instanceDisks,
				NetworkInterfaces: instanceNics,
			},
		}, true, nil
	}
}

func InstanceRunsInFleet(instance *corev1alpha1.Instance, instancePoolName string) bool {
	instancePoolRef := instance.Spec.FleetRef
	if instancePoolRef == nil {
		return false
	}

	return instancePoolRef.Name == instancePoolName
}

func InstanceRunsInFleetPredicate(fleetName string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		instance := object.(*corev1alpha1.Instance)
		return InstanceRunsInFleet(instance, fleetName)
	})
}

func (r *InstanceReconciler) matchingWatchLabel() client.ListOption {
	var labels map[string]string
	if r.WatchFilterValue != "" {
		labels = map[string]string{
			corev1alpha1.WatchLabel: r.WatchFilterValue,
		}
	}
	return client.MatchingLabels(labels)
}

func (r *InstanceReconciler) enqueueInstancesReferencingDisk() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []ctrl.Request {
		disk := obj.(*corev1alpha1.Disk)
		log := ctrl.LoggerFrom(ctx)

		instanceList := &corev1alpha1.InstanceList{}
		if err := r.List(ctx, instanceList,
			client.InNamespace(disk.Namespace),
			client.MatchingFields{
				sphereletclient.InstanceSpecDiskNamesField: disk.Name,
			},
			r.matchingWatchLabel(),
		); err != nil {
			log.Error(err, "Error listing instances using disk", "DiskKey", client.ObjectKeyFromObject(disk))
			return nil
		}

		return utilclient.ReconcileRequestsFromObjectStructSlice[*corev1alpha1.Instance](instanceList.Items)
	})
}

func (r *InstanceReconciler) enqueueInstancesReferencingSecret() handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []ctrl.Request {
		secret := obj.(*corev1.Secret)
		log := ctrl.LoggerFrom(ctx)

		instanceList := &corev1alpha1.InstanceList{}
		if err := r.List(ctx, instanceList,
			client.InNamespace(secret.Namespace),
			client.MatchingFields{
				sphereletclient.InstanceSpecSecretNamesField: secret.Name,
			},
			r.matchingWatchLabel(),
		); err != nil {
			log.Error(err, "Error listing instances using secret", "SecretKey", client.ObjectKeyFromObject(secret))
			return nil
		}

		return utilclient.ReconcileRequestsFromObjectStructSlice[*corev1alpha1.Instance](instanceList.Items)
	})
}

func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log := ctrl.Log.WithName("spherelet")

	return ctrl.NewControllerManagedBy(mgr).
		For(
			&corev1alpha1.Instance{},
			builder.WithPredicates(
				InstanceRunsInFleetPredicate(r.FleetName),
				predicates.ResourceHasFilterLabel(log, r.WatchFilterValue),
				predicates.ResourceIsNotExternallyManaged(log),
			),
		).
		Watches(
			&corev1.Secret{},
			r.enqueueInstancesReferencingSecret(),
		).
		Watches(
			&corev1alpha1.Disk{},
			r.enqueueInstancesReferencingDisk(),
		).
		Complete(r)
}
