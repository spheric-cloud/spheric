// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/ironcore-dev/controller-utils/clientutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	volumepoolletv1alpha1 "spheric.cloud/spheric/poollet/volumepoollet/api/v1alpha1"
	"spheric.cloud/spheric/poollet/volumepoollet/controllers/events"
	"spheric.cloud/spheric/poollet/volumepoollet/vcm"
	srimeta "spheric.cloud/spheric/sri/apis/meta/v1alpha1"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
	"spheric.cloud/spheric/utils/predicates"
)

type VolumeReconciler struct {
	record.EventRecorder
	client.Client
	Scheme *runtime.Scheme

	VolumeRuntime sri.VolumeRuntimeClient

	VolumeClassMapper vcm.VolumeClassMapper

	VolumePoolName   string
	WatchFilterValue string
}

func (r *VolumeReconciler) sriVolumeLabels(volume *storagev1alpha1.Volume) map[string]string {
	return map[string]string{
		volumepoolletv1alpha1.VolumeUIDLabel:       string(volume.UID),
		volumepoolletv1alpha1.VolumeNamespaceLabel: volume.Namespace,
		volumepoolletv1alpha1.VolumeNameLabel:      volume.Name,
	}
}

func (r *VolumeReconciler) sriVolumeAnnotations(_ *storagev1alpha1.Volume) map[string]string {
	return map[string]string{}
}

func (r *VolumeReconciler) listSRIVolumesByKey(ctx context.Context, volumeKey client.ObjectKey) ([]*sri.Volume, error) {
	res, err := r.VolumeRuntime.ListVolumes(ctx, &sri.ListVolumesRequest{
		Filter: &sri.VolumeFilter{
			LabelSelector: map[string]string{
				volumepoolletv1alpha1.VolumeNamespaceLabel: volumeKey.Namespace,
				volumepoolletv1alpha1.VolumeNameLabel:      volumeKey.Name,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing volumes by key: %w", err)
	}
	volumes := res.Volumes
	return volumes, nil
}

func (r *VolumeReconciler) listSRIVolumesByUID(ctx context.Context, volumeUID types.UID) ([]*sri.Volume, error) {
	res, err := r.VolumeRuntime.ListVolumes(ctx, &sri.ListVolumesRequest{
		Filter: &sri.VolumeFilter{
			LabelSelector: map[string]string{
				volumepoolletv1alpha1.VolumeUIDLabel: string(volumeUID),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing volumes by uid: %w", err)
	}
	return res.Volumes, nil
}

//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=volumes,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=volumes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=volumes/finalizers,verbs=update

func (r *VolumeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	volume := &storagev1alpha1.Volume{}
	if err := r.Get(ctx, req.NamespacedName, volume); err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("error getting volume %s: %w", req.NamespacedName, err)
		}
		return r.deleteGone(ctx, log, req.NamespacedName)
	}
	return r.reconcileExists(ctx, log, volume)
}

func (r *VolumeReconciler) deleteGone(ctx context.Context, log logr.Logger, volumeKey client.ObjectKey) (ctrl.Result, error) {
	log.V(1).Info("Delete gone")

	log.V(1).Info("Listing sri volumes by key")
	volumes, err := r.listSRIVolumesByKey(ctx, volumeKey)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing sri volumes by key: %w", err)
	}

	ok, err := r.deleteSRIVolumes(ctx, log, volumes)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting sri volumes: %w", err)
	}
	if !ok {
		log.V(1).Info("Not all sri volumes are gone, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Deleted gone")
	return ctrl.Result{}, nil
}

func (r *VolumeReconciler) deleteSRIVolumes(ctx context.Context, log logr.Logger, volumes []*sri.Volume) (bool, error) {
	var (
		errs                 []error
		deletingSRIVolumeIDs []string
	)

	for _, volume := range volumes {
		sriVolumeID := volume.Metadata.Id
		log := log.WithValues("SRIVolumeID", sriVolumeID)
		log.V(1).Info("Deleting sri volume")
		_, err := r.VolumeRuntime.DeleteVolume(ctx, &sri.DeleteVolumeRequest{
			VolumeId: sriVolumeID,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("error deleting sri volume %s: %w", sriVolumeID, err))
			} else {
				log.V(1).Info("SRI Volume is already gone")
			}
		} else {
			log.V(1).Info("Issued sri volume deletion")
			deletingSRIVolumeIDs = append(deletingSRIVolumeIDs, sriVolumeID)
		}
	}

	switch {
	case len(errs) > 0:
		return false, fmt.Errorf("error(s) deleting sri volume(s): %v", errs)
	case len(deletingSRIVolumeIDs) > 0:
		log.V(1).Info("Volumes are in deletion", "DeletingSRIVolumeIDs", deletingSRIVolumeIDs)
		return false, nil
	default:
		log.V(1).Info("No sri volumes present")
		return true, nil
	}
}

func (r *VolumeReconciler) reconcileExists(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume) (ctrl.Result, error) {
	if !volume.DeletionTimestamp.IsZero() {
		return r.delete(ctx, log, volume)
	}
	return r.reconcile(ctx, log, volume)
}

func (r *VolumeReconciler) delete(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume) (ctrl.Result, error) {
	log.V(1).Info("Delete")

	if !controllerutil.ContainsFinalizer(volume, volumepoolletv1alpha1.VolumeFinalizer) {
		log.V(1).Info("No finalizer present, nothing to do")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Finalizer present")

	log.V(1).Info("Listing volumes")
	volumes, err := r.listSRIVolumesByUID(ctx, volume.UID)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing volumes by uid: %w", err)
	}

	ok, err := r.deleteSRIVolumes(ctx, log, volumes)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting sri volumes: %w", err)
	}
	if !ok {
		log.V(1).Info("Not all sri volumes are gone, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Deleted all sri volumes, removing finalizer")
	if err := clientutils.PatchRemoveFinalizer(ctx, r.Client, volume, volumepoolletv1alpha1.VolumeFinalizer); err != nil {
		return ctrl.Result{}, fmt.Errorf("error removing finalizer: %w", err)
	}

	log.V(1).Info("Deleted")
	return ctrl.Result{}, nil
}

func getSRIVolumeClassCapabilities(volumeClass *storagev1alpha1.VolumeClass) (*sri.VolumeClassCapabilities, error) {
	tps := volumeClass.Capabilities.TPS()
	iops := volumeClass.Capabilities.IOPS()

	return &sri.VolumeClassCapabilities{
		Tps:  tps.Value(),
		Iops: iops.Value(),
	}, nil
}

func (r *VolumeReconciler) prepareSRIVolumeMetadata(volume *storagev1alpha1.Volume) *srimeta.ObjectMetadata {
	return &srimeta.ObjectMetadata{
		Labels:      r.sriVolumeLabels(volume),
		Annotations: r.sriVolumeAnnotations(volume),
	}
}

func (r *VolumeReconciler) prepareSRIVolumeClass(ctx context.Context, volume *storagev1alpha1.Volume, volumeClassName string) (string, bool, error) {
	volumeClass := &storagev1alpha1.VolumeClass{}
	volumeClassKey := client.ObjectKey{Name: volumeClassName}
	if err := r.Get(ctx, volumeClassKey, volumeClass); err != nil {
		err = fmt.Errorf("error getting volume class %s: %w", volumeClassKey, err)
		if !apierrors.IsNotFound(err) {
			return "", false, fmt.Errorf("error getting volume class %s: %w", volumeClassName, err)
		}

		r.Eventf(volume, corev1.EventTypeNormal, events.VolumeClassNotReady, "Volume class %s not found", volumeClassName)
		return "", false, nil
	}

	caps, err := getSRIVolumeClassCapabilities(volumeClass)
	if err != nil {
		return "", false, fmt.Errorf("error getting sri volume class capabilities: %w", err)
	}

	class, _, err := r.VolumeClassMapper.GetVolumeClassFor(ctx, volumeClassName, caps)
	if err != nil {
		return "", false, fmt.Errorf("error getting matching volume class: %w", err)
	}
	return class.Name, true, nil
}

func (r *VolumeReconciler) prepareSRIVolumeEncryption(ctx context.Context, volume *storagev1alpha1.Volume) (*sri.EncryptionSpec, bool, error) {
	encryption := volume.Spec.Encryption
	if encryption == nil {
		return nil, true, nil
	}

	encryptionSecret := &corev1.Secret{}
	encryptionSecretKey := client.ObjectKey{Name: encryption.SecretRef.Name, Namespace: volume.Namespace}
	if err := r.Get(ctx, encryptionSecretKey, encryptionSecret); err != nil {
		err = fmt.Errorf("error getting volume encryption secret %s: %w", encryptionSecretKey, err)
		if !apierrors.IsNotFound(err) {
			return nil, false, fmt.Errorf("error getting volume encryption secret %s: %w", encryption.SecretRef.Name, err)
		}

		r.Eventf(volume, corev1.EventTypeNormal, events.VolumeEncryptionSecretNotReady, "Volume encryption secret %s not found", encryption.SecretRef.Name)
		return nil, false, nil
	}

	return &sri.EncryptionSpec{
		SecretData: encryptionSecret.Data,
	}, true, nil
}

func (r *VolumeReconciler) prepareSRIVolumeResources(_ context.Context, _ *storagev1alpha1.Volume, resources corev1alpha1.ResourceList) (*sri.VolumeResources, bool, error) {
	storageBytes := resources.Storage().Value()

	return &sri.VolumeResources{
		StorageBytes: storageBytes,
	}, true, nil
}

func (r *VolumeReconciler) prepareSRIVolume(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume) (*sri.Volume, bool, error) {
	var (
		ok   = true
		errs []error
	)

	log.V(1).Info("Getting volume class")
	class, classOK, err := r.prepareSRIVolumeClass(ctx, volume, volume.Spec.VolumeClassRef.Name)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing sri volume class: %w", err))
	case !classOK:
		ok = false
	}

	log.V(1).Info("Getting encryption secret")
	encryption, encryptionOK, err := r.prepareSRIVolumeEncryption(ctx, volume)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing sri volume class: %w", err))
	case !encryptionOK:
		ok = false
	}

	resources, resourcesOK, err := r.prepareSRIVolumeResources(ctx, volume, volume.Spec.Resources)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing sri volume resources: %w", err))
	case !resourcesOK:
		ok = false
	}

	metadata := r.prepareSRIVolumeMetadata(volume)

	if len(errs) > 0 {
		return nil, false, fmt.Errorf("error(s) preparing sri volume: %v", errs)
	}
	if !ok {
		return nil, false, nil
	}

	return &sri.Volume{
		Metadata: metadata,
		Spec: &sri.VolumeSpec{
			Image:      volume.Spec.Image,
			Class:      class,
			Resources:  resources,
			Encryption: encryption,
		},
	}, true, nil
}

func (r *VolumeReconciler) reconcile(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume) (ctrl.Result, error) {
	log.V(1).Info("Reconcile")

	log.V(1).Info("Ensuring finalizer")
	modified, err := clientutils.PatchEnsureFinalizer(ctx, r.Client, volume, volumepoolletv1alpha1.VolumeFinalizer)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring finalizer: %w", err)
	}
	if modified {
		log.V(1).Info("Added finalizer, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}
	log.V(1).Info("Finalizer is present")

	log.V(1).Info("Ensuring no reconcile annotation")
	modified, err = sphericclient.PatchEnsureNoReconcileAnnotation(ctx, r.Client, volume)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring no reconcile annotation: %w", err)
	}
	if modified {
		log.V(1).Info("Removed reconcile annotation, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Listing volumes")
	res, err := r.VolumeRuntime.ListVolumes(ctx, &sri.ListVolumesRequest{
		Filter: &sri.VolumeFilter{
			LabelSelector: map[string]string{
				volumepoolletv1alpha1.VolumeUIDLabel: string(volume.UID),
			},
		},
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing volumes: %w", err)
	}

	switch len(res.Volumes) {
	case 0:
		return r.create(ctx, log, volume)
	case 1:
		sriVolume := res.Volumes[0]
		if err := r.update(ctx, log, volume, sriVolume); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating volume: %w", err)
		}

		if err := r.updateStatus(ctx, log, volume, sriVolume); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating volume status: %w", err)
		}
		return ctrl.Result{}, nil
	default:
		panic("unhandled multiple volumes")
	}
}

func (r *VolumeReconciler) create(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume) (ctrl.Result, error) {
	log.V(1).Info("Create")

	log.V(1).Info("Preparing sri volume")
	sriVolume, ok, err := r.prepareSRIVolume(ctx, log, volume)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error preparing sri volume: %w", err)
	}
	if !ok {
		log.V(1).Info("SRI volume is not yet ready to be prepared")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Creating volume")
	res, err := r.VolumeRuntime.CreateVolume(ctx, &sri.CreateVolumeRequest{
		Volume: sriVolume,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error creating volume: %w", err)
	}

	sriVolume = res.Volume

	volumeID := sriVolume.Metadata.Id
	log = log.WithValues("VolumeID", volumeID)
	log.V(1).Info("Created")

	log.V(1).Info("Updating status")
	if err := r.updateStatus(ctx, log, volume, sriVolume); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating volume status: %w", err)
	}

	log.V(1).Info("Created")
	return ctrl.Result{}, nil
}

func (r *VolumeReconciler) update(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume, sriVolume *sri.Volume) error {
	storageBytes := volume.Spec.Resources.Storage().Value()
	oldStorageBytes := sriVolume.Spec.Resources.StorageBytes
	if storageBytes != oldStorageBytes {
		log.V(1).Info("Expanding volume", "StorageBytes", storageBytes, "OldStorageBytes", oldStorageBytes)
		if _, err := r.VolumeRuntime.ExpandVolume(ctx, &sri.ExpandVolumeRequest{
			VolumeId: sriVolume.Metadata.Id,
			Resources: &sri.VolumeResources{
				StorageBytes: storageBytes,
			},
		}); err != nil {
			return fmt.Errorf("failed to expand volume: %w", err)
		}
	}

	return nil
}

func (r *VolumeReconciler) volumeSecretName(volumeName string, volumeHandle string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s/%s", volumeName, volumeHandle)))
	return hex.EncodeToString(sum[:])[:63]
}

var sriVolumeStateToVolumeState = map[sri.VolumeState]storagev1alpha1.VolumeState{
	sri.VolumeState_VOLUME_PENDING:   storagev1alpha1.VolumeStatePending,
	sri.VolumeState_VOLUME_AVAILABLE: storagev1alpha1.VolumeStateAvailable,
	sri.VolumeState_VOLUME_ERROR:     storagev1alpha1.VolumeStateError,
}

func (r *VolumeReconciler) convertSRIVolumeState(sriState sri.VolumeState) (storagev1alpha1.VolumeState, error) {
	if res, ok := sriVolumeStateToVolumeState[sriState]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown volume state %v", sriState)
}

func (r *VolumeReconciler) updateStatus(ctx context.Context, log logr.Logger, volume *storagev1alpha1.Volume, sriVolume *sri.Volume) error {
	var access *storagev1alpha1.VolumeAccess

	if sriVolume.Status.State == sri.VolumeState_VOLUME_AVAILABLE {
		if sriAccess := sriVolume.Status.Access; sriAccess != nil {
			var secretRef *corev1.LocalObjectReference

			if sriAccess.SecretData != nil {
				log.V(1).Info("Applying volume secret")
				volumeSecret := &corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						APIVersion: corev1.SchemeGroupVersion.String(),
						Kind:       "Secret",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: volume.Namespace,
						Name:      r.volumeSecretName(volume.Name, sriAccess.Handle),
						Labels: map[string]string{
							volumepoolletv1alpha1.VolumeUIDLabel: string(volume.UID),
						},
					},
					Data: sriAccess.SecretData,
				}
				_ = ctrl.SetControllerReference(volume, volumeSecret, r.Scheme)
				if err := r.Patch(ctx, volumeSecret, client.Apply, client.FieldOwner(volumepoolletv1alpha1.FieldOwner)); err != nil {
					return fmt.Errorf("error applying volume secret: %w", err)
				}
				secretRef = &corev1.LocalObjectReference{Name: volumeSecret.Name}
			} else {
				log.V(1).Info("Deleting any corresponding volume secret")
				if err := r.DeleteAllOf(ctx, &corev1.Secret{},
					client.InNamespace(volume.Namespace),
					client.MatchingLabels{
						volumepoolletv1alpha1.VolumeUIDLabel: string(volume.UID),
					},
				); err != nil {
					return fmt.Errorf("error deleting any corresponding volume secret: %w", err)
				}
			}

			access = &storagev1alpha1.VolumeAccess{
				SecretRef:        secretRef,
				Driver:           sriAccess.Driver,
				Handle:           sriAccess.Handle,
				VolumeAttributes: sriAccess.Attributes,
			}
		}

	}

	base := volume.DeepCopy()
	now := metav1.Now()

	volume.Status.Access = access
	newState, err := r.convertSRIVolumeState(sriVolume.Status.State)
	if err != nil {
		return err
	}
	if newState != volume.Status.State {
		volume.Status.LastStateTransitionTime = &now
	}
	volume.Status.State = newState

	if err := r.Status().Patch(ctx, volume, client.MergeFrom(base)); err != nil {
		return fmt.Errorf("error patching volume status: %w", err)
	}
	return nil
}

func (r *VolumeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log := ctrl.Log.WithName("volumepoollet")

	return ctrl.NewControllerManagedBy(mgr).
		For(
			&storagev1alpha1.Volume{},
			builder.WithPredicates(
				VolumeRunsInVolumePoolPredicate(r.VolumePoolName),
				predicates.ResourceHasFilterLabel(log, r.WatchFilterValue),
				predicates.ResourceIsNotExternallyManaged(log),
			),
		).
		Complete(r)
}

func VolumeRunsInVolumePool(volume *storagev1alpha1.Volume, volumePoolName string) bool {
	volumePoolRef := volume.Spec.VolumePoolRef
	if volumePoolRef == nil {
		return false
	}

	return volumePoolRef.Name == volumePoolName
}

func VolumeRunsInVolumePoolPredicate(volumePoolName string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		volume := object.(*storagev1alpha1.Volume)
		return VolumeRunsInVolumePool(volume, volumePoolName)
	})
}
