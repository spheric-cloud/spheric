// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	commonv1alpha1 "spheric.cloud/spheric/api/common/v1alpha1"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	"spheric.cloud/spheric/poollet/machinepoollet/controllers/events"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/utils/claimmanager"
	utilslices "spheric.cloud/spheric/utils/slices"
)

type volumeClaimStrategy struct {
	client.Client
}

func (s *volumeClaimStrategy) ClaimState(claimer client.Object, obj client.Object) claimmanager.ClaimState {
	volume := obj.(*storagev1alpha1.Volume)
	if claimRef := volume.Spec.ClaimRef; claimRef != nil {
		if claimRef.UID == claimer.GetUID() {
			return claimmanager.ClaimStateClaimed
		}
		return claimmanager.ClaimStateTaken
	}
	return claimmanager.ClaimStateFree
}

func (s *volumeClaimStrategy) Adopt(ctx context.Context, claimer client.Object, obj client.Object) error {
	volume := obj.(*storagev1alpha1.Volume)
	base := volume.DeepCopy()
	volume.Spec.ClaimRef = commonv1alpha1.NewLocalObjUIDRef(claimer)
	return s.Patch(ctx, volume, client.StrategicMergeFrom(base))
}

func (s *volumeClaimStrategy) Release(ctx context.Context, claimer client.Object, obj client.Object) error {
	volume := obj.(*storagev1alpha1.Volume)
	base := volume.DeepCopy()
	volume.Spec.ClaimRef = nil
	return s.Patch(ctx, volume, client.StrategicMergeFrom(base))
}

func (r *MachineReconciler) volumeNameToMachineVolume(machine *computev1alpha1.Machine) map[string]computev1alpha1.Volume {
	sel := make(map[string]computev1alpha1.Volume)
	for _, machineVolume := range machine.Spec.Volumes {
		volumeName := computev1alpha1.MachineVolumeName(machine.Name, machineVolume)
		if volumeName == "" {
			// volume name is empty on empty disk volumes.
			continue
		}
		sel[volumeName] = machineVolume
	}
	return sel
}

func (r *MachineReconciler) machineVolumeSelector(machine *computev1alpha1.Machine) claimmanager.Selector {
	names := sets.New(computev1alpha1.MachineVolumeNames(machine)...)
	return claimmanager.SelectorFunc(func(obj client.Object) bool {
		volume := obj.(*storagev1alpha1.Volume)
		return names.Has(volume.Name)
	})
}

func (r *MachineReconciler) getVolumesForMachine(ctx context.Context, machine *computev1alpha1.Machine) ([]storagev1alpha1.Volume, error) {
	volumeList := &storagev1alpha1.VolumeList{}
	if err := r.List(ctx, volumeList,
		client.InNamespace(machine.Namespace),
	); err != nil {
		return nil, fmt.Errorf("error listing volumes: %w", err)
	}

	var (
		sel      = r.machineVolumeSelector(machine)
		claimMgr = claimmanager.New(machine, sel, &volumeClaimStrategy{r.Client})
		volumes  []storagev1alpha1.Volume
		errs     []error
	)
	for _, volume := range volumeList.Items {
		ok, err := claimMgr.Claim(ctx, &volume)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !ok {
			continue
		}

		if volume.Status.State != storagev1alpha1.VolumeStateAvailable || volume.Status.Access == nil {
			r.Eventf(machine, corev1.EventTypeNormal, events.VolumeNotReady, "Volume %s does not access information", volume.Name)
			continue
		}

		volumes = append(volumes, volume)
	}
	return volumes, errors.Join(errs...)
}

func (r *MachineReconciler) prepareRemoteSRIVolume(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	machineVolume *computev1alpha1.Volume,
	volume *storagev1alpha1.Volume,
) (*sri.Volume, bool, error) {
	access := volume.Status.Access
	if access == nil {
		r.Eventf(machine, corev1.EventTypeNormal, events.VolumeNotReady, "Volume %s does not report status access", volume.Name)
		return nil, false, nil
	}

	var secretData map[string][]byte
	if secretRef := access.SecretRef; secretRef != nil {
		secret := &corev1.Secret{}
		secretKey := client.ObjectKey{Namespace: volume.Namespace, Name: secretRef.Name}
		if err := r.Get(ctx, secretKey, secret); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, false, fmt.Errorf("error getting volume access secret %s: %w", secretKey.Name, err)
			}

			r.Eventf(machine, corev1.EventTypeNormal, events.VolumeNotReady,
				"Volume %s access secret %s not found",
				volume.Name,
				secretKey.Name,
			)
			return nil, false, nil
		}

		secretData = secret.Data
	}

	var encryptionData map[string][]byte
	if encryption := volume.Spec.Encryption; encryption != nil {
		secret := &corev1.Secret{}
		secretKey := client.ObjectKey{Namespace: volume.Namespace, Name: encryption.SecretRef.Name}
		if err := r.Get(ctx, secretKey, secret); err != nil {
			if !apierrors.IsNotFound(err) {
				return nil, false, fmt.Errorf("error getting volume encryption secret %s: %w", secretKey.Name, err)
			}

			r.Eventf(machine, corev1.EventTypeNormal, events.VolumeNotReady,
				"Volume %s encryption secret %s not found",
				volume.Name,
				secretKey.Name,
			)
			return nil, false, nil
		}

		encryptionData = secret.Data
	}

	return &sri.Volume{
		Name:   machineVolume.Name,
		Device: *machineVolume.Device,
		Connection: &sri.VolumeConnection{
			Driver:         access.Driver,
			Handle:         access.Handle,
			Attributes:     access.VolumeAttributes,
			SecretData:     secretData,
			EncryptionData: encryptionData,
		},
	}, true, nil
}

func (r *MachineReconciler) prepareEmptyDiskSRIVolume(machineVolume *computev1alpha1.Volume) *sri.Volume {
	var sizeBytes int64
	if sizeLimit := machineVolume.EmptyDisk.SizeLimit; sizeLimit != nil {
		sizeBytes = sizeLimit.Value()
	}
	return &sri.Volume{
		Name:   machineVolume.Name,
		Device: *machineVolume.Device,
		EmptyDisk: &sri.EmptyDisk{
			SizeBytes: sizeBytes,
		},
	}
}

func (r *MachineReconciler) prepareSRIVolumes(
	ctx context.Context,
	machine *computev1alpha1.Machine,
	volumes []storagev1alpha1.Volume,
) ([]*sri.Volume, bool, error) {
	var (
		volumeNameToMachineVolume = r.volumeNameToMachineVolume(machine)
		sriVolumes                []*sri.Volume
		errs                      []error
	)
	for _, volume := range volumes {
		machineVolume := volumeNameToMachineVolume[volume.Name]
		sriVolume, ok, err := r.prepareRemoteSRIVolume(ctx, machine, &machineVolume, &volume)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !ok {
			continue
		}

		sriVolumes = append(sriVolumes, sriVolume)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, false, err
	}

	for _, machineVolume := range machine.Spec.Volumes {
		if machineVolume.EmptyDisk == nil {
			continue
		}

		sriVolume := r.prepareEmptyDiskSRIVolume(&machineVolume)
		sriVolumes = append(sriVolumes, sriVolume)
	}

	if len(sriVolumes) != len(machine.Spec.Volumes) {
		expectedVolumeNames := utilslices.ToSetFunc(machine.Spec.Volumes, func(v computev1alpha1.Volume) string { return v.Name })
		actualVolumeNames := utilslices.ToSetFunc(sriVolumes, (*sri.Volume).GetName)
		missingVolumeNames := sets.List(expectedVolumeNames.Difference(actualVolumeNames))
		r.Eventf(machine, corev1.EventTypeNormal, events.VolumeNotReady, "Machine volumes are not ready: %v", missingVolumeNames)
		return nil, false, nil
	}
	return sriVolumes, true, nil
}

func (r *MachineReconciler) getExistingSRIVolumesForMachine(
	ctx context.Context,
	log logr.Logger,
	sriMachine *sri.Machine,
	desiredSRIVolumes []*sri.Volume,
) ([]*sri.Volume, error) {
	var (
		sriVolumes              []*sri.Volume
		desiredSRIVolumesByName = utilslices.ToMapByKey(desiredSRIVolumes, (*sri.Volume).GetName)
		errs                    []error
	)

	for _, sriVolume := range sriMachine.Spec.Volumes {
		log := log.WithValues("Volume", sriVolume.Name)

		desiredSRIVolume, ok := desiredSRIVolumesByName[sriVolume.Name]
		if ok && proto.Equal(desiredSRIVolume, sriVolume) {
			log.V(1).Info("Existing SRI volume is up-to-date")
			sriVolumes = append(sriVolumes, sriVolume)
			continue
		}

		log.V(1).Info("Detaching outdated SRI volume")
		_, err := r.MachineRuntime.DetachVolume(ctx, &sri.DetachVolumeRequest{
			MachineId: sriMachine.Metadata.Id,
			Name:      sriVolume.Name,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("[volume %s] %w", sriVolume.Name, err))
				continue
			}
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return sriVolumes, nil
}

func (r *MachineReconciler) getNewSRIVolumesForMachine(
	ctx context.Context,
	log logr.Logger,
	sriMachine *sri.Machine,
	desiredSRIVolumes, existingSRIVolumes []*sri.Volume,
) ([]*sri.Volume, error) {
	var (
		desiredNewSRIVolumes = FindNewSRIVolumes(desiredSRIVolumes, existingSRIVolumes)
		sriVolumes           []*sri.Volume
		errs                 []error
	)
	for _, newSRIVolume := range desiredNewSRIVolumes {
		log := log.WithValues("Volume", newSRIVolume.Name)
		log.V(1).Info("Attaching new volume")
		if _, err := r.MachineRuntime.AttachVolume(ctx, &sri.AttachVolumeRequest{
			MachineId: sriMachine.Metadata.Id,
			Volume:    newSRIVolume,
		}); err != nil {
			errs = append(errs, fmt.Errorf("[volume %s] %w", newSRIVolume.Name, err))
			continue
		}

		sriVolumes = append(sriVolumes, newSRIVolume)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return sriVolumes, nil
}

func (r *MachineReconciler) updateSRIVolumes(
	ctx context.Context,
	log logr.Logger,
	machine *computev1alpha1.Machine,
	sriMachine *sri.Machine,
	volumes []storagev1alpha1.Volume,
) error {
	desiredSRIVolumes, _, err := r.prepareSRIVolumes(ctx, machine, volumes)
	if err != nil {
		return fmt.Errorf("error preparing sri volumes: %w", err)
	}

	extistingSRIVolumes, err := r.getExistingSRIVolumesForMachine(ctx, log, sriMachine, desiredSRIVolumes)
	if err != nil {
		return fmt.Errorf("error getting existing sri volumes for machine: %w", err)
	}

	_, err = r.getNewSRIVolumesForMachine(ctx, log, sriMachine, desiredSRIVolumes, extistingSRIVolumes)
	if err != nil {
		return fmt.Errorf("error getting new sri volumes for machine: %w", err)
	}

	return nil
}

func (r *MachineReconciler) getVolumeStatusesForMachine(
	machine *computev1alpha1.Machine,
	sriMachine *sri.Machine,
	now metav1.Time,
) ([]computev1alpha1.VolumeStatus, error) {
	var (
		sriVolumeStatusByName        = utilslices.ToMapByKey(sriMachine.Status.Volumes, (*sri.VolumeStatus).GetName)
		existingVolumeStatusesByName = utilslices.ToMapByKey(machine.Status.Volumes, func(status computev1alpha1.VolumeStatus) string { return status.Name })
		volumeStatuses               []computev1alpha1.VolumeStatus
		errs                         []error
	)

	for _, machineVolume := range machine.Spec.Volumes {
		var (
			sriVolumeStatus, ok = sriVolumeStatusByName[machineVolume.Name]
			volumeStatusValues  computev1alpha1.VolumeStatus
		)
		if ok {
			var err error
			volumeStatusValues, err = r.convertSRIVolumeStatus(sriVolumeStatus)
			if err != nil {
				return nil, fmt.Errorf("[volume %s] %w", machineVolume.Name, err)
			}
		} else {
			volumeStatusValues = computev1alpha1.VolumeStatus{
				Name:  machineVolume.Name,
				State: computev1alpha1.VolumeStatePending,
			}
		}

		volumeStatus := existingVolumeStatusesByName[machineVolume.Name]
		r.addVolumeStatusValues(now, &volumeStatus, &volumeStatusValues)
		volumeStatuses = append(volumeStatuses, volumeStatus)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return volumeStatuses, nil
}

var sriVolumeStateToVolumeState = map[sri.VolumeState]computev1alpha1.VolumeState{
	sri.VolumeState_VOLUME_ATTACHED: computev1alpha1.VolumeStateAttached,
	sri.VolumeState_VOLUME_PENDING:  computev1alpha1.VolumeStatePending,
}

func (r *MachineReconciler) convertSRIVolumeState(sriState sri.VolumeState) (computev1alpha1.VolumeState, error) {
	if res, ok := sriVolumeStateToVolumeState[sriState]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown sri volume state %v", sriState)
}

func (r *MachineReconciler) convertSRIVolumeStatus(sriVolumeStatus *sri.VolumeStatus) (computev1alpha1.VolumeStatus, error) {
	state, err := r.convertSRIVolumeState(sriVolumeStatus.State)
	if err != nil {
		return computev1alpha1.VolumeStatus{}, err
	}

	return computev1alpha1.VolumeStatus{
		Name:   sriVolumeStatus.Name,
		Handle: sriVolumeStatus.Handle,
		State:  state,
	}, nil
}

func (r *MachineReconciler) addVolumeStatusValues(now metav1.Time, existing, newValues *computev1alpha1.VolumeStatus) {
	if existing.State != newValues.State {
		existing.LastStateTransitionTime = &now
	}
	existing.Name = newValues.Name
	existing.State = newValues.State
	existing.Handle = newValues.Handle
}
