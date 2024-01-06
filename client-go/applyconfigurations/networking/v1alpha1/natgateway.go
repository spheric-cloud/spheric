// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	managedfields "k8s.io/apimachinery/pkg/util/managedfields"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	internal "spheric.cloud/spheric/client-go/applyconfigurations/internal"
	v1 "spheric.cloud/spheric/client-go/applyconfigurations/meta/v1"
)

// NATGatewayApplyConfiguration represents an declarative configuration of the NATGateway type for use
// with apply.
type NATGatewayApplyConfiguration struct {
	v1.TypeMetaApplyConfiguration    `json:",inline"`
	*v1.ObjectMetaApplyConfiguration `json:"metadata,omitempty"`
	Spec                             *NATGatewaySpecApplyConfiguration   `json:"spec,omitempty"`
	Status                           *NATGatewayStatusApplyConfiguration `json:"status,omitempty"`
}

// NATGateway constructs an declarative configuration of the NATGateway type for use with
// apply.
func NATGateway(name, namespace string) *NATGatewayApplyConfiguration {
	b := &NATGatewayApplyConfiguration{}
	b.WithName(name)
	b.WithNamespace(namespace)
	b.WithKind("NATGateway")
	b.WithAPIVersion("networking.spheric.cloud/v1alpha1")
	return b
}

// ExtractNATGateway extracts the applied configuration owned by fieldManager from
// nATGateway. If no managedFields are found in nATGateway for fieldManager, a
// NATGatewayApplyConfiguration is returned with only the Name, Namespace (if applicable),
// APIVersion and Kind populated. It is possible that no managed fields were found for because other
// field managers have taken ownership of all the fields previously owned by fieldManager, or because
// the fieldManager never owned fields any fields.
// nATGateway must be a unmodified NATGateway API object that was retrieved from the Kubernetes API.
// ExtractNATGateway provides a way to perform a extract/modify-in-place/apply workflow.
// Note that an extracted apply configuration will contain fewer fields than what the fieldManager previously
// applied if another fieldManager has updated or force applied any of the previously applied fields.
// Experimental!
func ExtractNATGateway(nATGateway *networkingv1alpha1.NATGateway, fieldManager string) (*NATGatewayApplyConfiguration, error) {
	return extractNATGateway(nATGateway, fieldManager, "")
}

// ExtractNATGatewayStatus is the same as ExtractNATGateway except
// that it extracts the status subresource applied configuration.
// Experimental!
func ExtractNATGatewayStatus(nATGateway *networkingv1alpha1.NATGateway, fieldManager string) (*NATGatewayApplyConfiguration, error) {
	return extractNATGateway(nATGateway, fieldManager, "status")
}

func extractNATGateway(nATGateway *networkingv1alpha1.NATGateway, fieldManager string, subresource string) (*NATGatewayApplyConfiguration, error) {
	b := &NATGatewayApplyConfiguration{}
	err := managedfields.ExtractInto(nATGateway, internal.Parser().Type("cloud.spheric.spheric.api.networking.v1alpha1.NATGateway"), fieldManager, b, subresource)
	if err != nil {
		return nil, err
	}
	b.WithName(nATGateway.Name)
	b.WithNamespace(nATGateway.Namespace)

	b.WithKind("NATGateway")
	b.WithAPIVersion("networking.spheric.cloud/v1alpha1")
	return b, nil
}

// WithKind sets the Kind field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kind field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithKind(value string) *NATGatewayApplyConfiguration {
	b.Kind = &value
	return b
}

// WithAPIVersion sets the APIVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the APIVersion field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithAPIVersion(value string) *NATGatewayApplyConfiguration {
	b.APIVersion = &value
	return b
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithName(value string) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithGenerateName(value string) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithNamespace(value string) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Namespace = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithUID(value types.UID) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.UID = &value
	return b
}

// WithResourceVersion sets the ResourceVersion field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ResourceVersion field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithResourceVersion(value string) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.ResourceVersion = &value
	return b
}

// WithGeneration sets the Generation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Generation field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithGeneration(value int64) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.Generation = &value
	return b
}

// WithCreationTimestamp sets the CreationTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CreationTimestamp field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithCreationTimestamp(value metav1.Time) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.CreationTimestamp = &value
	return b
}

// WithDeletionTimestamp sets the DeletionTimestamp field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionTimestamp field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithDeletionTimestamp(value metav1.Time) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionTimestamp = &value
	return b
}

// WithDeletionGracePeriodSeconds sets the DeletionGracePeriodSeconds field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DeletionGracePeriodSeconds field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithDeletionGracePeriodSeconds(value int64) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	b.DeletionGracePeriodSeconds = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *NATGatewayApplyConfiguration) WithLabels(entries map[string]string) *NATGatewayApplyConfiguration {
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
func (b *NATGatewayApplyConfiguration) WithAnnotations(entries map[string]string) *NATGatewayApplyConfiguration {
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
func (b *NATGatewayApplyConfiguration) WithOwnerReferences(values ...*v1.OwnerReferenceApplyConfiguration) *NATGatewayApplyConfiguration {
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
func (b *NATGatewayApplyConfiguration) WithFinalizers(values ...string) *NATGatewayApplyConfiguration {
	b.ensureObjectMetaApplyConfigurationExists()
	for i := range values {
		b.Finalizers = append(b.Finalizers, values[i])
	}
	return b
}

func (b *NATGatewayApplyConfiguration) ensureObjectMetaApplyConfigurationExists() {
	if b.ObjectMetaApplyConfiguration == nil {
		b.ObjectMetaApplyConfiguration = &v1.ObjectMetaApplyConfiguration{}
	}
}

// WithSpec sets the Spec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Spec field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithSpec(value *NATGatewaySpecApplyConfiguration) *NATGatewayApplyConfiguration {
	b.Spec = value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *NATGatewayApplyConfiguration) WithStatus(value *NATGatewayStatusApplyConfiguration) *NATGatewayApplyConfiguration {
	b.Status = value
	return b
}
