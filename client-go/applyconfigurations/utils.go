// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfigurations

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
	v1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	corev1alpha1 "spheric.cloud/spheric/client-go/applyconfigurations/core/v1alpha1"
	internal "spheric.cloud/spheric/client-go/applyconfigurations/internal"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=core.spheric.cloud, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("AttachedDisk"):
		return &corev1alpha1.AttachedDiskApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AttachedDiskSource"):
		return &corev1alpha1.AttachedDiskSourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("AttachedDiskStatus"):
		return &corev1alpha1.AttachedDiskStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DaemonEndpoint"):
		return &corev1alpha1.DaemonEndpointApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Disk"):
		return &corev1alpha1.DiskApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DiskAccess"):
		return &corev1alpha1.DiskAccessApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DiskSpec"):
		return &corev1alpha1.DiskSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DiskStatus"):
		return &corev1alpha1.DiskStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DiskTemplateSpec"):
		return &corev1alpha1.DiskTemplateSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DiskType"):
		return &corev1alpha1.DiskTypeApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EFIVar"):
		return &corev1alpha1.EFIVarApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EmptyDiskSource"):
		return &corev1alpha1.EmptyDiskSourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("EphemeralDiskSource"):
		return &corev1alpha1.EphemeralDiskSourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Fleet"):
		return &corev1alpha1.FleetApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FleetAddress"):
		return &corev1alpha1.FleetAddressApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FleetCondition"):
		return &corev1alpha1.FleetConditionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FleetDaemonEndpoints"):
		return &corev1alpha1.FleetDaemonEndpointsApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FleetSpec"):
		return &corev1alpha1.FleetSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("FleetStatus"):
		return &corev1alpha1.FleetStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Instance"):
		return &corev1alpha1.InstanceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InstanceSpec"):
		return &corev1alpha1.InstanceSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InstanceStatus"):
		return &corev1alpha1.InstanceStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InstanceType"):
		return &corev1alpha1.InstanceTypeApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LocalObjectReference"):
		return &corev1alpha1.LocalObjectReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("LocalUIDReference"):
		return &corev1alpha1.LocalUIDReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Network"):
		return &corev1alpha1.NetworkApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NetworkInterface"):
		return &corev1alpha1.NetworkInterfaceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NetworkInterfaceStatus"):
		return &corev1alpha1.NetworkInterfaceStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("NetworkStatus"):
		return &corev1alpha1.NetworkStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SecretKeySelector"):
		return &corev1alpha1.SecretKeySelectorApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Subnet"):
		return &corev1alpha1.SubnetApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SubnetReference"):
		return &corev1alpha1.SubnetReferenceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SubnetSpec"):
		return &corev1alpha1.SubnetSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("SubnetStatus"):
		return &corev1alpha1.SubnetStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Taint"):
		return &corev1alpha1.TaintApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Toleration"):
		return &corev1alpha1.TolerationApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
