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
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	bucketpoolletv1alpha1 "spheric.cloud/spheric/poollet/bucketpoollet/api/v1alpha1"
	"spheric.cloud/spheric/poollet/bucketpoollet/bcm"
	"spheric.cloud/spheric/poollet/bucketpoollet/controllers/events"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	srimeta "spheric.cloud/spheric/sri/apis/meta/v1alpha1"
	sphericclient "spheric.cloud/spheric/utils/client"
	"spheric.cloud/spheric/utils/predicates"
)

type BucketReconciler struct {
	record.EventRecorder
	client.Client
	Scheme *runtime.Scheme

	BucketRuntime sri.BucketRuntimeClient

	BucketClassMapper bcm.BucketClassMapper

	BucketPoolName   string
	WatchFilterValue string
}

func (r *BucketReconciler) sriBucketLabels(bucket *storagev1alpha1.Bucket) map[string]string {
	return map[string]string{
		bucketpoolletv1alpha1.BucketUIDLabel:       string(bucket.UID),
		bucketpoolletv1alpha1.BucketNamespaceLabel: bucket.Namespace,
		bucketpoolletv1alpha1.BucketNameLabel:      bucket.Name,
	}
}

func (r *BucketReconciler) sriBucketAnnotations(_ *storagev1alpha1.Bucket) map[string]string {
	return map[string]string{}
}

func (r *BucketReconciler) listSRIBucketsByKey(ctx context.Context, bucketKey client.ObjectKey) ([]*sri.Bucket, error) {
	res, err := r.BucketRuntime.ListBuckets(ctx, &sri.ListBucketsRequest{
		Filter: &sri.BucketFilter{
			LabelSelector: map[string]string{
				bucketpoolletv1alpha1.BucketNamespaceLabel: bucketKey.Namespace,
				bucketpoolletv1alpha1.BucketNameLabel:      bucketKey.Name,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing buckets by key: %w", err)
	}
	buckets := res.Buckets
	return buckets, nil
}

func (r *BucketReconciler) listSRIBucketsByUID(ctx context.Context, bucketUID types.UID) ([]*sri.Bucket, error) {
	res, err := r.BucketRuntime.ListBuckets(ctx, &sri.ListBucketsRequest{
		Filter: &sri.BucketFilter{
			LabelSelector: map[string]string{
				bucketpoolletv1alpha1.BucketUIDLabel: string(bucketUID),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error listing buckets by uid: %w", err)
	}
	return res.Buckets, nil
}

//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=buckets,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=buckets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=storage.spheric.cloud,resources=buckets/finalizers,verbs=update

func (r *BucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	bucket := &storagev1alpha1.Bucket{}
	if err := r.Get(ctx, req.NamespacedName, bucket); err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("error getting bucket %s: %w", req.NamespacedName, err)
		}
		return r.deleteGone(ctx, log, req.NamespacedName)
	}
	return r.reconcileExists(ctx, log, bucket)
}

func (r *BucketReconciler) deleteGone(ctx context.Context, log logr.Logger, bucketKey client.ObjectKey) (ctrl.Result, error) {
	log.V(1).Info("Delete gone")

	log.V(1).Info("Listing sri buckets by key")
	buckets, err := r.listSRIBucketsByKey(ctx, bucketKey)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing sri buckets by key: %w", err)
	}

	ok, err := r.deleteSRIBuckets(ctx, log, buckets)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting sri buckets: %w", err)
	}
	if !ok {
		log.V(1).Info("Not all sri buckets are gone, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Deleted gone")
	return ctrl.Result{}, nil
}

func (r *BucketReconciler) deleteSRIBuckets(ctx context.Context, log logr.Logger, buckets []*sri.Bucket) (bool, error) {
	var (
		errs                 []error
		deletingSRIBucketIDs []string
	)

	for _, bucket := range buckets {
		sriBucketID := bucket.Metadata.Id
		log := log.WithValues("SRIBucketID", sriBucketID)
		log.V(1).Info("Deleting sri bucket")
		_, err := r.BucketRuntime.DeleteBucket(ctx, &sri.DeleteBucketRequest{
			BucketId: sriBucketID,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound {
				errs = append(errs, fmt.Errorf("error deleting sri bucket %s: %w", sriBucketID, err))
			} else {
				log.V(1).Info("SRI Bucket is already gone")
			}
		} else {
			log.V(1).Info("Issued sri bucket deletion")
			deletingSRIBucketIDs = append(deletingSRIBucketIDs, sriBucketID)
		}
	}

	switch {
	case len(errs) > 0:
		return false, fmt.Errorf("error(s) deleting sri bucket(s): %v", errs)
	case len(deletingSRIBucketIDs) > 0:
		log.V(1).Info("Buckets are in deletion", "DeletingSRIBucketIDs", deletingSRIBucketIDs)
		return false, nil
	default:
		log.V(1).Info("No sri buckets present")
		return true, nil
	}
}

func (r *BucketReconciler) reconcileExists(ctx context.Context, log logr.Logger, bucket *storagev1alpha1.Bucket) (ctrl.Result, error) {
	if !bucket.DeletionTimestamp.IsZero() {
		return r.delete(ctx, log, bucket)
	}
	return r.reconcile(ctx, log, bucket)
}

func (r *BucketReconciler) delete(ctx context.Context, log logr.Logger, bucket *storagev1alpha1.Bucket) (ctrl.Result, error) {
	log.V(1).Info("Delete")

	if !controllerutil.ContainsFinalizer(bucket, bucketpoolletv1alpha1.BucketFinalizer) {
		log.V(1).Info("No finalizer present, nothing to do")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Finalizer present")

	log.V(1).Info("Listing buckets")
	buckets, err := r.listSRIBucketsByUID(ctx, bucket.UID)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing buckets by uid: %w", err)
	}

	ok, err := r.deleteSRIBuckets(ctx, log, buckets)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error deleting sri buckets: %w", err)
	}
	if !ok {
		log.V(1).Info("Not all sri buckets are gone, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Deleted all sri buckets, removing finalizer")
	if err := clientutils.PatchRemoveFinalizer(ctx, r.Client, bucket, bucketpoolletv1alpha1.BucketFinalizer); err != nil {
		return ctrl.Result{}, fmt.Errorf("error removing finalizer: %w", err)
	}

	log.V(1).Info("Deleted")
	return ctrl.Result{}, nil
}

func getSRIBucketClassCapabilities(bucketClass *storagev1alpha1.BucketClass) (*sri.BucketClassCapabilities, error) {
	tps := bucketClass.Capabilities.TPS()
	iops := bucketClass.Capabilities.IOPS()

	return &sri.BucketClassCapabilities{
		Tps:  tps.Value(),
		Iops: iops.Value(),
	}, nil
}

func (r *BucketReconciler) prepareSRIBucketMetadata(bucket *storagev1alpha1.Bucket) *srimeta.ObjectMetadata {
	return &srimeta.ObjectMetadata{
		Labels:      r.sriBucketLabels(bucket),
		Annotations: r.sriBucketAnnotations(bucket),
	}
}

func (r *BucketReconciler) prepareSRIBucketClass(ctx context.Context, bucket *storagev1alpha1.Bucket, bucketClassName string) (string, bool, error) {
	bucketClass := &storagev1alpha1.BucketClass{}
	bucketClassKey := client.ObjectKey{Name: bucketClassName}
	if err := r.Get(ctx, bucketClassKey, bucketClass); err != nil {
		err = fmt.Errorf("error getting bucket class %s: %w", bucketClassKey, err)
		if !apierrors.IsNotFound(err) {
			return "", false, fmt.Errorf("error getting bucket class %s: %w", bucketClassName, err)
		}

		r.Eventf(bucket, corev1.EventTypeNormal, events.BucketClassNotReady, "Bucket class %s not found", bucketClassName)
		return "", false, nil
	}

	caps, err := getSRIBucketClassCapabilities(bucketClass)
	if err != nil {
		return "", false, fmt.Errorf("error getting sri bucket class capabilities: %w", err)
	}

	class, err := r.BucketClassMapper.GetBucketClassFor(ctx, bucketClassName, caps)
	if err != nil {
		return "", false, fmt.Errorf("error getting matching bucket class: %w", err)
	}
	return class.Name, true, nil
}

func (r *BucketReconciler) prepareSRIBucket(ctx context.Context, log logr.Logger, bucket *storagev1alpha1.Bucket) (*sri.Bucket, bool, error) {
	var (
		ok   = true
		errs []error
	)

	log.V(1).Info("Getting bucket class")
	class, classOK, err := r.prepareSRIBucketClass(ctx, bucket, bucket.Spec.BucketClassRef.Name)
	switch {
	case err != nil:
		errs = append(errs, fmt.Errorf("error preparing sri bucket class: %w", err))
	case !classOK:
		ok = false
	}

	metadata := r.prepareSRIBucketMetadata(bucket)

	if len(errs) > 0 {
		return nil, false, fmt.Errorf("error(s) preparing sri bucket: %v", errs)
	}
	if !ok {
		return nil, false, nil
	}

	return &sri.Bucket{
		Metadata: metadata,
		Spec: &sri.BucketSpec{
			Class: class,
		},
	}, true, nil
}

func (r *BucketReconciler) reconcile(ctx context.Context, log logr.Logger, bucket *storagev1alpha1.Bucket) (ctrl.Result, error) {
	log.V(1).Info("Reconcile")

	log.V(1).Info("Ensuring finalizer")
	modified, err := clientutils.PatchEnsureFinalizer(ctx, r.Client, bucket, bucketpoolletv1alpha1.BucketFinalizer)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring finalizer: %w", err)
	}
	if modified {
		log.V(1).Info("Added finalizer, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}
	log.V(1).Info("Finalizer is present")

	log.V(1).Info("Ensuring no reconcile annotation")
	modified, err = sphericclient.PatchEnsureNoReconcileAnnotation(ctx, r.Client, bucket)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error ensuring no reconcile annotation: %w", err)
	}
	if modified {
		log.V(1).Info("Removed reconcile annotation, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	log.V(1).Info("Listing buckets")
	res, err := r.BucketRuntime.ListBuckets(ctx, &sri.ListBucketsRequest{
		Filter: &sri.BucketFilter{
			LabelSelector: map[string]string{
				bucketpoolletv1alpha1.BucketUIDLabel: string(bucket.UID),
			},
		},
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error listing buckets: %w", err)
	}

	switch len(res.Buckets) {
	case 0:
		return r.create(ctx, log, bucket)
	case 1:
		sriBucket := res.Buckets[0]
		if err := r.updateStatus(ctx, log, bucket, sriBucket); err != nil {
			return ctrl.Result{}, fmt.Errorf("error updating bucket status: %w", err)
		}
		return ctrl.Result{}, nil
	default:
		panic("unhandled multiple buckets")
	}
}

func (r *BucketReconciler) create(ctx context.Context, log logr.Logger, bucket *storagev1alpha1.Bucket) (ctrl.Result, error) {
	log.V(1).Info("Create")

	log.V(1).Info("Preparing sri bucket")
	sriBucket, ok, err := r.prepareSRIBucket(ctx, log, bucket)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error preparing sri bucket: %w", err)
	}
	if !ok {
		log.V(1).Info("SRI bucket is not yet ready to be prepared")
		return ctrl.Result{}, nil
	}

	log.V(1).Info("Creating bucket")
	res, err := r.BucketRuntime.CreateBucket(ctx, &sri.CreateBucketRequest{
		Bucket: sriBucket,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("error creating bucket: %w", err)
	}

	sriBucket = res.Bucket

	bucketID := sriBucket.Metadata.Id
	log = log.WithValues("BucketID", bucketID)
	log.V(1).Info("Created")

	log.V(1).Info("Updating status")
	if err := r.updateStatus(ctx, log, bucket, sriBucket); err != nil {
		return ctrl.Result{}, fmt.Errorf("error updating bucket status: %w", err)
	}

	log.V(1).Info("Created")
	return ctrl.Result{}, nil
}

func (r *BucketReconciler) bucketSecretName(bucketName string) string {
	sum := sha256.Sum256([]byte(bucketName))
	return hex.EncodeToString(sum[:])[:63]
}

var sriBucketStateToBucketState = map[sri.BucketState]storagev1alpha1.BucketState{
	sri.BucketState_BUCKET_PENDING:   storagev1alpha1.BucketStatePending,
	sri.BucketState_BUCKET_AVAILABLE: storagev1alpha1.BucketStateAvailable,
	sri.BucketState_BUCKET_ERROR:     storagev1alpha1.BucketStateError,
}

func (r *BucketReconciler) convertSRIBucketState(sriState sri.BucketState) (storagev1alpha1.BucketState, error) {
	if res, ok := sriBucketStateToBucketState[sriState]; ok {
		return res, nil
	}
	return "", fmt.Errorf("unknown bucket state %v", sriState)
}

func (r *BucketReconciler) updateStatus(ctx context.Context, log logr.Logger, bucket *storagev1alpha1.Bucket, sriBucket *sri.Bucket) error {
	var access *storagev1alpha1.BucketAccess

	if sriBucket.Status.State == sri.BucketState_BUCKET_AVAILABLE {
		if sriAccess := sriBucket.Status.Access; sriAccess != nil {
			var secretRef *corev1.LocalObjectReference

			if sriAccess.SecretData != nil {
				log.V(1).Info("Applying bucket secret")
				bucketSecret := &corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						APIVersion: corev1.SchemeGroupVersion.String(),
						Kind:       "Secret",
					},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: bucket.Namespace,
						Name:      r.bucketSecretName(bucket.Name),
						Labels: map[string]string{
							bucketpoolletv1alpha1.BucketUIDLabel: string(bucket.UID),
						},
					},
					Data: sriAccess.SecretData,
				}
				_ = ctrl.SetControllerReference(bucket, bucketSecret, r.Scheme)
				if err := r.Patch(ctx, bucketSecret, client.Apply, client.FieldOwner(bucketpoolletv1alpha1.FieldOwner)); err != nil {
					return fmt.Errorf("error applying bucket secret: %w", err)
				}
				secretRef = &corev1.LocalObjectReference{Name: bucketSecret.Name}
			} else {
				log.V(1).Info("Deleting any corresponding bucket secret")
				if err := r.DeleteAllOf(ctx, &corev1.Secret{},
					client.InNamespace(bucket.Namespace),
					client.MatchingLabels{
						bucketpoolletv1alpha1.BucketUIDLabel: string(bucket.UID),
					},
				); err != nil {
					return fmt.Errorf("error deleting any corresponding bucket secret: %w", err)
				}
			}

			access = &storagev1alpha1.BucketAccess{
				SecretRef: secretRef,
				Endpoint:  sriAccess.Endpoint,
			}
		}

	}

	base := bucket.DeepCopy()
	now := metav1.Now()

	bucket.Status.Access = access
	newState, err := r.convertSRIBucketState(sriBucket.Status.State)
	if err != nil {
		return err
	}
	if newState != bucket.Status.State {
		bucket.Status.LastStateTransitionTime = &now
	}
	bucket.Status.State = newState

	if err := r.Status().Patch(ctx, bucket, client.MergeFrom(base)); err != nil {
		return fmt.Errorf("error patching bucket status: %w", err)
	}
	return nil
}

func (r *BucketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	log := ctrl.Log.WithName("bucketpoollet")

	return ctrl.NewControllerManagedBy(mgr).
		For(
			&storagev1alpha1.Bucket{},
			builder.WithPredicates(
				BucketRunsInBucketPoolPredicate(r.BucketPoolName),
				predicates.ResourceHasFilterLabel(log, r.WatchFilterValue),
				predicates.ResourceIsNotExternallyManaged(log),
			),
		).
		Complete(r)
}

func BucketRunsInBucketPool(bucket *storagev1alpha1.Bucket, bucketPoolName string) bool {
	bucketPoolRef := bucket.Spec.BucketPoolRef
	if bucketPoolRef == nil {
		return false
	}

	return bucketPoolRef.Name == bucketPoolName
}

func BucketRunsInBucketPoolPredicate(bucketPoolName string) predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		bucket := object.(*storagev1alpha1.Bucket)
		return BucketRunsInBucketPool(bucket, bucketPoolName)
	})
}
