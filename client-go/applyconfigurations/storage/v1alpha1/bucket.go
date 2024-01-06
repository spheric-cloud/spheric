// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	managedfields "k8s.io/apimachinery/pkg/util/managedfields"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	internal "spheric.cloud/spheric/client-go/applyconfigurations/internal"
	v1 "spheric.cloud/spheric/client-go/applyconfigurations/meta/v1"
)

// BucketApplyConfiguration represents an declarative configuration of the Bucket type for use
// with apply.
type BucketApplyConfiguration struct {
	v1.TypeMetaApplyConfiguration    `json:",inline"`
	*v1.ObjectMetaApplyConfiguration `json:"metadata,omitempty"`
	Spec                             *BucketSpecApplyConfiguration   `json:"spec,omitempty"`
	Status                           *BucketStatusApplyConfiguration `json:"status,omitempty"`
}

// Bucket constructs an declarative configuration of the Bucket type for use with
// apply.
func Bucket(name, namespace string) *BucketApplyConfiguration {
	b := &BucketApplyConfiguration{}
	b.WithName(name)
	b.WithNamespace(namespace)
	b.WithKind("Bucket")
	b.WithAPIVersion("storage.spheric.cloud/v1alpha1")
	return b
}

// ExtractBucket extracts the applied configuration owned by fieldManager from
// bucket. If no managedFields are found in bucket for fieldManager, a
// BucketApplyConfiguration is returned with only the Name, Namespace (if applicable),
// APIVersion and Kind populated. It is possible that no managed fields were found for because other
// field managers have taken ownership of all the fields previously owned by fieldManager, or because
// the fieldManager never owned fields any fields.
// bucket must be a unmodified Bucket API object that was retrieved from the Kubernetes API.
// ExtractBucket provides a way to perform a extract/modify-in-place/apply workflow.
// Note that an extracted apply configuration will contain fewer fields than what the fieldManager previously
// applied if another fieldManager has updated or force applied any of the previously applied fields.
// Experimental!
func ExtractBucket(bucket *storagev1alpha1.Bucket, fieldManager string) (*BucketApplyConfiguration, error) {
	return extractBucket(bucket, fieldManager, "")
}

// ExtractBucketStatus is the same as ExtractBucket except
// that it extracts the status subresource applied configuration.
// Experimental!
func ExtractBucketStatus(bucket *storagev1alpha1.Bucket, fieldManager string) (*BucketApplyConfiguration, error) {
	return extractBucket(bucket, fieldManager, "status")
}

func extractBucket(bucket *storagev1alpha1.Bucket, fieldManager string, subresource string) (*BucketApplyConfiguration, error) {
	b := &BucketApplyConfiguration{}
	err := managedfields.ExtractInto(bucket, internal.Parser().Type("cloud.spheric.spheric.api.storage.v1alpha1.Bucket"), fieldManager, b, subresource)
	if err != nil {
		return nil, err
	}
	b.WithName(bucket.Name)
	b.WithNamespace(bucket.Namespace)

	b.WithKind("Bucket")
	b.WithAPIVersion("storage.spheric.cloud/v1alpha1")
	return b, nil
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithKind(value string) *BucketApplyConfiguration {
	b.Kind = &value
	return b
}

// WithAPIVersion sets the APIVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the APIVersion field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithAPIVersion(value string) *BucketApplyConfiguration {
	b.APIVersion = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithName(value string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithGenerateName(value string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithNamespace(value string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Namespace = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithUID(value types.UID) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.UID = &value
	return b
}

// WithResourceVersion sets the ResourceVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceVersion field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithResourceVersion(value string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ResourceVersion = &value
	return b
}

// WithGeneration sets the Generation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Generation field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithGeneration(value int64) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Generation = &value
	return b
}

// WithCreationTimestamp sets the CreationTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CreationTimestamp field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithCreationTimestamp(value metav1.Time) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.CreationTimestamp = &value
	return b
}

// WithDeletionTimestamp sets the DeletionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionTimestamp field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithDeletionTimestamp(value metav1.Time) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionTimestamp = &value
	return b
}

// WithDeletionGracePeriodSeconds sets the DeletionGracePeriodSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionGracePeriodSeconds field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithDeletionGracePeriodSeconds(value int64) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionGracePeriodSeconds = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *BucketApplyConfiguration) WithLabels(entries map[string]string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *BucketApplyConfiguration) WithAnnotations(entries map[string]string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithOwnerReferences adds the given value to the OwnerReferences field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OwnerReferences field.
func (b *BucketApplyConfiguration) WithOwnerReferences(values ...*v1.OwnerReferenceApplyConfiguration) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOwnerReferences")
		}
		b.OwnerReferences = append(b.OwnerReferences, *values[i])
	}
	return b
}

// WithFinalizers adds the given value to the Finalizers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Finalizers field.
func (b *BucketApplyConfiguration) WithFinalizers(values ...string) *BucketApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		b.Finalizers = append(b.Finalizers, values[i])
	}
	return b
}

func (b *BucketApplyConfiguration) ensureObjectMetaApplyConfigurationExists() {
	if b.ObjectMetaApplyConfiguration == nil {
		b.ObjectMetaApplyConfiguration = &v1.ObjectMetaApplyConfiguration{}
	}
}

// WithSpec sets the Spec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Spec field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithSpec(value *BucketSpecApplyConfiguration) *BucketApplyConfiguration {
	b.Spec = value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *BucketApplyConfiguration) WithStatus(value *BucketStatusApplyConfiguration) *BucketApplyConfiguration {
	b.Status = value
	return b
}
