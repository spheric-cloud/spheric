// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"errors"
	"fmt"

	"spheric.cloud/spheric/api/core/v1alpha1"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/spherelet/controllers/events"
	"spheric.cloud/spheric/utils/claimmanager"
	utilslices "spheric.cloud/spheric/utils/slices"
)

type diskClaimStrategy struct {
	client.Client
}

func (s *diskClaimStrategy) ClaimState(claimer client.Object, obj client.Object) claimmanager.ClaimState {
	disk := obj.(*corev1alpha1.Disk)
	if claimRef := disk.Spec.InstanceRef; claimRef != nil {
		if claimRef.UID == claimer.GetUID() {
			return claimmanager.ClaimStateClaimed
		}
		return claimmanager.ClaimStateTaken
	}
	return claimmanager.ClaimStateFree
}

func (s *diskClaimStrategy) Adopt(ctx context.Context, claimer client.Object, obj client.Object) error {
	disk := obj.(*corev1alpha1.Disk)
	base := disk.DeepCopy()
	disk.Spec.InstanceRef = corev1alpha1.NewLocalObjUIDRef(claimer)
	return s.Patch(ctx, disk, client.StrategicMergeFrom(base))
}

func (s *diskClaimStrategy) Release(ctx context.Context, claimer client.Object, obj client.Object) error {
	disk := obj.(*corev1alpha1.Disk)
	base := disk.DeepCopy()
	disk.Spec.InstanceRef = nil
	return s.Patch(ctx, disk, client.StrategicMergeFrom(base))
}

func (r *InstanceReconciler) diskNameToAttachedDisk(instance *v1alpha1.Instance) map[string]v1alpha1.AttachedDisk {
	sel := make(map[string]v1alpha1.AttachedDisk)
	for _, instanceDisk := range instance.Spec.Disks {
		diskClaimName := v1alpha1.InstanceDiskName(instance.Name, instanceDisk)
		if diskClaimName == "" {
			// disk claim name is empty on empty disk disks.
			continue
		}
		sel[diskClaimName] = instanceDisk
	}
	return sel
}

func (r *InstanceReconciler) instanceClaimSelector(instance *v1alpha1.Instance) claimmanager.Selector {
	names := sets.New(v1alpha1.InstanceDiskNames(instance)...)
	return claimmanager.SelectorFunc(func(obj client.Object) bool {
		disk := obj.(*corev1alpha1.Disk)
		return names.Has(disk.Name)
	})
}

func (r *InstanceReconciler) getDisksForInstance(ctx context.Context, instance *v1alpha1.Instance) ([]corev1alpha1.Disk, error) {
	diskList := &corev1alpha1.DiskList{}
	if err := r.List(ctx, diskList,
		client.InNamespace(instance.Namespace),
	); err != nil {
		return nil, fmt.Errorf("error listing disks: %w", err)
	}

	var (
		sel      = r.instanceClaimSelector(instance)
		claimMgr = claimmanager.New(instance, sel, &diskClaimStrategy{r.Client})
		disks    []corev1alpha1.Disk
		errs     []error
	)
	for _, disk := range diskList.Items {
		ok, err := claimMgr.Claim(ctx, &disk)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !ok {
			continue
		}

		if disk.Status.State != corev1alpha1.DiskStateAvailable || disk.Status.Access == nil {
			r.Eventf(instance, corev1.EventTypeNormal, events.DiskNotReady, "Disk %s is not yet available", disk.Name)
			continue
		}

		disks = append(disks, disk)
	}
	return disks, errors.Join(errs...)
}

func (r *InstanceReconciler) prepareRemoteIRIDisk(
	ctx context.Context,
	instance *v1alpha1.Instance,
	instanceDisk *v1alpha1.AttachedDisk,
	disk *corev1alpha1.Disk,
) (*iri.Disk, bool, error) {
	access := disk.Status.Access

	var secretData map[string][]byte
	if secretRef := access.SecretRef; secretRef != nil {
		secret := &corev1.Secret{}
		secretKey := client.ObjectKey{Namespace: disk.Namespace, Name: secretRef.Name}
		if err := r.Get(ctx, secretKey, secret); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, false, fmt.Errorf("error getting disk access secret %s: %w", secretKey.Name, err)
			}

			r.Eventf(instance, corev1.EventTypeNormal, events.DiskNotReady,
				"Disk %s access secret %s not found",
				disk.Name,
				secretKey.Name,
			)
			return nil, false, nil
		}

		secretData = secret.Data
	}

	return &iri.Disk{
		Name:   instanceDisk.Name,
		Device: *instanceDisk.Device,
		Connection: &iri.DiskConnection{
			Driver:     access.Driver,
			Handle:     access.Handle,
			Attributes: access.Attributes,
			SecretData: secretData,
		},
	}, true, nil
}

func (r *InstanceReconciler) prepareEmptyDiskIRIDisk(instanceDisk *v1alpha1.AttachedDisk) *iri.Disk {
	var sizeBytes int64
	if sizeLimit := instanceDisk.EmptyDisk.SizeLimit; sizeLimit != nil {
		sizeBytes = sizeLimit.Value()
	}
	return &iri.Disk{
		Name:   instanceDisk.Name,
		Device: *instanceDisk.Device,
		EmptyDisk: &iri.EmptyDisk{
			SizeBytes: sizeBytes,
		},
	}
}

func (r *InstanceReconciler) prepareIRIDisks(
	ctx context.Context,
	instance *v1alpha1.Instance,
	disks []corev1alpha1.Disk,
) ([]*iri.Disk, bool, error) {
	var (
		diskClaimNameToAttachedDisk = r.diskNameToAttachedDisk(instance)
		iriDisks                    []*iri.Disk
		errs                        []error
	)
	for _, disk := range disks {
		instanceDisk := diskClaimNameToAttachedDisk[disk.Name]
		iriDisk, ok, err := r.prepareRemoteIRIDisk(ctx, instance, &instanceDisk, &disk)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !ok {
			continue
		}

		iriDisks = append(iriDisks, iriDisk)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, false, err
	}

	for _, instanceDisk := range instance.Spec.Disks {
		if instanceDisk.EmptyDisk == nil {
			continue
		}

		iriDisk := r.prepareEmptyDiskIRIDisk(&instanceDisk)
		iriDisks = append(iriDisks, iriDisk)
	}

	if len(iriDisks) != len(instance.Spec.Disks) {
		expectedDiskNames := utilslices.ToSetFunc(instance.Spec.Disks, func(v v1alpha1.AttachedDisk) string { return v.Name })
		actualDiskNames := utilslices.ToSetFunc(iriDisks, (*iri.Disk).GetName)
		missingDiskNames := sets.List(expectedDiskNames.Difference(actualDiskNames))
		r.Eventf(instance, corev1.EventTypeNormal, events.DiskNotReady, "Instance disks are not ready: %v", missingDiskNames)
		return nil, false, nil
	}
	return iriDisks, true, nil
}

func (r *InstanceReconciler) getExistingIRIDisksForInstance(
	ctx context.Context,
	log logr.Logger,
	iriInstance *iri.Instance,
	desiredIRIDisks []*iri.Disk,
) ([]*iri.Disk, error) {
	var (
		iriDisks              []*iri.Disk
		desiredIRIDisksByName = utilslices.ToMapByKey(desiredIRIDisks, (*iri.Disk).GetName)
		errs                  []error
	)

	for _, iriDisk := range iriInstance.Spec.Disks {
		log := log.WithValues("Disk", iriDisk.Name)

		desiredIRIDisk, ok := desiredIRIDisksByName[iriDisk.Name]
		if ok && proto.Equal(desiredIRIDisk, iriDisk) {
			log.V(1).Info("Existing IRI disk is up-to-date")
			iriDisks = append(iriDisks, iriDisk)
			continue
		}

		log.V(1).Info("Detaching outdated IRI disk")
		_, err := r.InstanceRuntime.DetachDisk(ctx, &iri.DetachDiskRequest{
			InstanceId: iriInstance.Metadata.Id,
			Name:       iriDisk.Name,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("[disk %s] %w", iriDisk.Name, err))
				continue
			}
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return iriDisks, nil
}

func (r *InstanceReconciler) getNewIRIDisksForInstance(
	ctx context.Context,
	log logr.Logger,
	iriInstance *iri.Instance,
	desiredIRIDisks, existingIRIDisks []*iri.Disk,
) ([]*iri.Disk, error) {
	var (
		desiredNewIRIDisks = FindNewIRIDisks(desiredIRIDisks, existingIRIDisks)
		iriDisks           []*iri.Disk
		errs               []error
	)
	for _, newIRIDisk := range desiredNewIRIDisks {
		log := log.WithValues("Disk", newIRIDisk.Name)
		log.V(1).Info("Attaching new disk")
		if _, err := r.InstanceRuntime.AttachDisk(ctx, &iri.AttachDiskRequest{
			InstanceId: iriInstance.Metadata.Id,
			Disk:       newIRIDisk,
		}); err != nil {
			errs = append(errs, fmt.Errorf("[disk %s] %w", newIRIDisk.Name, err))
			continue
		}

		iriDisks = append(iriDisks, newIRIDisk)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return iriDisks, nil
}

func (r *InstanceReconciler) updateIRIDisks(
	ctx context.Context,
	log logr.Logger,
	instance *v1alpha1.Instance,
	iriInstance *iri.Instance,
	disks []corev1alpha1.Disk,
) error {
	desiredIRIDisks, _, err := r.prepareIRIDisks(ctx, instance, disks)
	if err != nil {
		return fmt.Errorf("error preparing iri disks: %w", err)
	}

	extistingIRIDisks, err := r.getExistingIRIDisksForInstance(ctx, log, iriInstance, desiredIRIDisks)
	if err != nil {
		return fmt.Errorf("error getting existing iri disks for instance: %w", err)
	}

	_, err = r.getNewIRIDisksForInstance(ctx, log, iriInstance, desiredIRIDisks, extistingIRIDisks)
	if err != nil {
		return fmt.Errorf("error getting new iri disks for instance: %w", err)
	}

	return nil
}

func (r *InstanceReconciler) getDiskStatusesForInstance(
	instance *v1alpha1.Instance,
	iriInstance *iri.Instance,
	now metav1.Time,
) ([]v1alpha1.AttachedDiskStatus, error) {
	var (
		iriDiskStatusByName        = utilslices.ToMapByKey(iriInstance.Status.Disks, (*iri.DiskStatus).GetName)
		existingDiskStatusesByName = utilslices.ToMapByKey(instance.Status.Disks, func(status v1alpha1.AttachedDiskStatus) string { return status.Name })
		diskStatuses               []v1alpha1.AttachedDiskStatus
		errs                       []error
	)

	for _, instanceDisk := range instance.Spec.Disks {
		var (
			iriDiskStatus, ok = iriDiskStatusByName[instanceDisk.Name]
			diskStatusValues  v1alpha1.AttachedDiskStatus
		)
		if ok {
			var err error
			diskStatusValues, err = r.convertIRIDiskStatus(iriDiskStatus)
			if err != nil {
				return nil, fmt.Errorf("[disk %s] %w", instanceDisk.Name, err)
			}
		} else {
			diskStatusValues = v1alpha1.AttachedDiskStatus{
				Name:  instanceDisk.Name,
				State: v1alpha1.AttachedDiskStatePending,
			}
		}

		diskStatus := existingDiskStatusesByName[instanceDisk.Name]
		r.addDiskStatusValues(now, &diskStatus, &diskStatusValues)
		diskStatuses = append(diskStatuses, diskStatus)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return diskStatuses, nil
}

var iriDiskStateToDiskState = map[iri.DiskState]v1alpha1.AttachedDiskState{
	iri.DiskState_DISK_ATTACHED: v1alpha1.AttachedDiskStateAttached,
	iri.DiskState_DISK_PENDING:  v1alpha1.AttachedDiskStatePending,
}

func (r *InstanceReconciler) convertIRIDiskState(iriState iri.DiskState) (v1alpha1.AttachedDiskState, error) {
	if res, ok := iriDiskStateToDiskState[iriState]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown iri disk state %v", iriState)
}

func (r *InstanceReconciler) convertIRIDiskStatus(iriDiskStatus *iri.DiskStatus) (v1alpha1.AttachedDiskStatus, error) {
	state, err := r.convertIRIDiskState(iriDiskStatus.State)
	if err != nil {
		return v1alpha1.AttachedDiskStatus{}, err
	}

	return v1alpha1.AttachedDiskStatus{
		Name:  iriDiskStatus.Name,
		State: state,
	}, nil
}

func (r *InstanceReconciler) addDiskStatusValues(now metav1.Time, existing, newValues *v1alpha1.AttachedDiskStatus) {
	if existing.State != newValues.State {
		existing.LastStateTransitionTime = &now
	}
	existing.Name = newValues.Name
	existing.State = newValues.State
}
