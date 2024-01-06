//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	commonv1alpha1 "spheric.cloud/spheric/api/common/v1alpha1"
	v1alpha1 "spheric.cloud/spheric/api/ipam/v1alpha1"
	ipam "spheric.cloud/spheric/internal/apis/ipam"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1alpha1.Prefix)(nil), (*ipam.Prefix)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Prefix_To_ipam_Prefix(a.(*v1alpha1.Prefix), b.(*ipam.Prefix), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.Prefix)(nil), (*v1alpha1.Prefix)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_Prefix_To_v1alpha1_Prefix(a.(*ipam.Prefix), b.(*v1alpha1.Prefix), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixAllocation)(nil), (*ipam.PrefixAllocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixAllocation_To_ipam_PrefixAllocation(a.(*v1alpha1.PrefixAllocation), b.(*ipam.PrefixAllocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixAllocation)(nil), (*v1alpha1.PrefixAllocation)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixAllocation_To_v1alpha1_PrefixAllocation(a.(*ipam.PrefixAllocation), b.(*v1alpha1.PrefixAllocation), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixAllocationList)(nil), (*ipam.PrefixAllocationList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixAllocationList_To_ipam_PrefixAllocationList(a.(*v1alpha1.PrefixAllocationList), b.(*ipam.PrefixAllocationList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixAllocationList)(nil), (*v1alpha1.PrefixAllocationList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixAllocationList_To_v1alpha1_PrefixAllocationList(a.(*ipam.PrefixAllocationList), b.(*v1alpha1.PrefixAllocationList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixAllocationSpec)(nil), (*ipam.PrefixAllocationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixAllocationSpec_To_ipam_PrefixAllocationSpec(a.(*v1alpha1.PrefixAllocationSpec), b.(*ipam.PrefixAllocationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixAllocationSpec)(nil), (*v1alpha1.PrefixAllocationSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixAllocationSpec_To_v1alpha1_PrefixAllocationSpec(a.(*ipam.PrefixAllocationSpec), b.(*v1alpha1.PrefixAllocationSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixAllocationStatus)(nil), (*ipam.PrefixAllocationStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixAllocationStatus_To_ipam_PrefixAllocationStatus(a.(*v1alpha1.PrefixAllocationStatus), b.(*ipam.PrefixAllocationStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixAllocationStatus)(nil), (*v1alpha1.PrefixAllocationStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixAllocationStatus_To_v1alpha1_PrefixAllocationStatus(a.(*ipam.PrefixAllocationStatus), b.(*v1alpha1.PrefixAllocationStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixList)(nil), (*ipam.PrefixList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixList_To_ipam_PrefixList(a.(*v1alpha1.PrefixList), b.(*ipam.PrefixList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixList)(nil), (*v1alpha1.PrefixList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixList_To_v1alpha1_PrefixList(a.(*ipam.PrefixList), b.(*v1alpha1.PrefixList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixSpec)(nil), (*ipam.PrefixSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec(a.(*v1alpha1.PrefixSpec), b.(*ipam.PrefixSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixSpec)(nil), (*v1alpha1.PrefixSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec(a.(*ipam.PrefixSpec), b.(*v1alpha1.PrefixSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixStatus)(nil), (*ipam.PrefixStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixStatus_To_ipam_PrefixStatus(a.(*v1alpha1.PrefixStatus), b.(*ipam.PrefixStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixStatus)(nil), (*v1alpha1.PrefixStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixStatus_To_v1alpha1_PrefixStatus(a.(*ipam.PrefixStatus), b.(*v1alpha1.PrefixStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1alpha1.PrefixTemplateSpec)(nil), (*ipam.PrefixTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_PrefixTemplateSpec_To_ipam_PrefixTemplateSpec(a.(*v1alpha1.PrefixTemplateSpec), b.(*ipam.PrefixTemplateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ipam.PrefixTemplateSpec)(nil), (*v1alpha1.PrefixTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_ipam_PrefixTemplateSpec_To_v1alpha1_PrefixTemplateSpec(a.(*ipam.PrefixTemplateSpec), b.(*v1alpha1.PrefixTemplateSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Prefix_To_ipam_Prefix(in *v1alpha1.Prefix, out *ipam.Prefix, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_PrefixStatus_To_ipam_PrefixStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Prefix_To_ipam_Prefix is an autogenerated conversion function.
func Convert_v1alpha1_Prefix_To_ipam_Prefix(in *v1alpha1.Prefix, out *ipam.Prefix, s conversion.Scope) error {
	return autoConvert_v1alpha1_Prefix_To_ipam_Prefix(in, out, s)
}

func autoConvert_ipam_Prefix_To_v1alpha1_Prefix(in *ipam.Prefix, out *v1alpha1.Prefix, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_ipam_PrefixStatus_To_v1alpha1_PrefixStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_ipam_Prefix_To_v1alpha1_Prefix is an autogenerated conversion function.
func Convert_ipam_Prefix_To_v1alpha1_Prefix(in *ipam.Prefix, out *v1alpha1.Prefix, s conversion.Scope) error {
	return autoConvert_ipam_Prefix_To_v1alpha1_Prefix(in, out, s)
}

func autoConvert_v1alpha1_PrefixAllocation_To_ipam_PrefixAllocation(in *v1alpha1.PrefixAllocation, out *ipam.PrefixAllocation, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_PrefixAllocationSpec_To_ipam_PrefixAllocationSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_PrefixAllocationStatus_To_ipam_PrefixAllocationStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_PrefixAllocation_To_ipam_PrefixAllocation is an autogenerated conversion function.
func Convert_v1alpha1_PrefixAllocation_To_ipam_PrefixAllocation(in *v1alpha1.PrefixAllocation, out *ipam.PrefixAllocation, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixAllocation_To_ipam_PrefixAllocation(in, out, s)
}

func autoConvert_ipam_PrefixAllocation_To_v1alpha1_PrefixAllocation(in *ipam.PrefixAllocation, out *v1alpha1.PrefixAllocation, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_ipam_PrefixAllocationSpec_To_v1alpha1_PrefixAllocationSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_ipam_PrefixAllocationStatus_To_v1alpha1_PrefixAllocationStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_ipam_PrefixAllocation_To_v1alpha1_PrefixAllocation is an autogenerated conversion function.
func Convert_ipam_PrefixAllocation_To_v1alpha1_PrefixAllocation(in *ipam.PrefixAllocation, out *v1alpha1.PrefixAllocation, s conversion.Scope) error {
	return autoConvert_ipam_PrefixAllocation_To_v1alpha1_PrefixAllocation(in, out, s)
}

func autoConvert_v1alpha1_PrefixAllocationList_To_ipam_PrefixAllocationList(in *v1alpha1.PrefixAllocationList, out *ipam.PrefixAllocationList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ipam.PrefixAllocation, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_PrefixAllocation_To_ipam_PrefixAllocation(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_PrefixAllocationList_To_ipam_PrefixAllocationList is an autogenerated conversion function.
func Convert_v1alpha1_PrefixAllocationList_To_ipam_PrefixAllocationList(in *v1alpha1.PrefixAllocationList, out *ipam.PrefixAllocationList, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixAllocationList_To_ipam_PrefixAllocationList(in, out, s)
}

func autoConvert_ipam_PrefixAllocationList_To_v1alpha1_PrefixAllocationList(in *ipam.PrefixAllocationList, out *v1alpha1.PrefixAllocationList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1alpha1.PrefixAllocation, len(*in))
		for i := range *in {
			if err := Convert_ipam_PrefixAllocation_To_v1alpha1_PrefixAllocation(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_ipam_PrefixAllocationList_To_v1alpha1_PrefixAllocationList is an autogenerated conversion function.
func Convert_ipam_PrefixAllocationList_To_v1alpha1_PrefixAllocationList(in *ipam.PrefixAllocationList, out *v1alpha1.PrefixAllocationList, s conversion.Scope) error {
	return autoConvert_ipam_PrefixAllocationList_To_v1alpha1_PrefixAllocationList(in, out, s)
}

func autoConvert_v1alpha1_PrefixAllocationSpec_To_ipam_PrefixAllocationSpec(in *v1alpha1.PrefixAllocationSpec, out *ipam.PrefixAllocationSpec, s conversion.Scope) error {
	out.IPFamily = v1.IPFamily(in.IPFamily)
	out.Prefix = (*commonv1alpha1.IPPrefix)(unsafe.Pointer(in.Prefix))
	out.PrefixLength = in.PrefixLength
	out.PrefixRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.PrefixRef))
	out.PrefixSelector = (*metav1.LabelSelector)(unsafe.Pointer(in.PrefixSelector))
	return nil
}

// Convert_v1alpha1_PrefixAllocationSpec_To_ipam_PrefixAllocationSpec is an autogenerated conversion function.
func Convert_v1alpha1_PrefixAllocationSpec_To_ipam_PrefixAllocationSpec(in *v1alpha1.PrefixAllocationSpec, out *ipam.PrefixAllocationSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixAllocationSpec_To_ipam_PrefixAllocationSpec(in, out, s)
}

func autoConvert_ipam_PrefixAllocationSpec_To_v1alpha1_PrefixAllocationSpec(in *ipam.PrefixAllocationSpec, out *v1alpha1.PrefixAllocationSpec, s conversion.Scope) error {
	out.IPFamily = v1.IPFamily(in.IPFamily)
	out.Prefix = (*commonv1alpha1.IPPrefix)(unsafe.Pointer(in.Prefix))
	out.PrefixLength = in.PrefixLength
	out.PrefixRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.PrefixRef))
	out.PrefixSelector = (*metav1.LabelSelector)(unsafe.Pointer(in.PrefixSelector))
	return nil
}

// Convert_ipam_PrefixAllocationSpec_To_v1alpha1_PrefixAllocationSpec is an autogenerated conversion function.
func Convert_ipam_PrefixAllocationSpec_To_v1alpha1_PrefixAllocationSpec(in *ipam.PrefixAllocationSpec, out *v1alpha1.PrefixAllocationSpec, s conversion.Scope) error {
	return autoConvert_ipam_PrefixAllocationSpec_To_v1alpha1_PrefixAllocationSpec(in, out, s)
}

func autoConvert_v1alpha1_PrefixAllocationStatus_To_ipam_PrefixAllocationStatus(in *v1alpha1.PrefixAllocationStatus, out *ipam.PrefixAllocationStatus, s conversion.Scope) error {
	out.Prefix = (*commonv1alpha1.IPPrefix)(unsafe.Pointer(in.Prefix))
	out.Phase = ipam.PrefixAllocationPhase(in.Phase)
	out.LastPhaseTransitionTime = (*metav1.Time)(unsafe.Pointer(in.LastPhaseTransitionTime))
	return nil
}

// Convert_v1alpha1_PrefixAllocationStatus_To_ipam_PrefixAllocationStatus is an autogenerated conversion function.
func Convert_v1alpha1_PrefixAllocationStatus_To_ipam_PrefixAllocationStatus(in *v1alpha1.PrefixAllocationStatus, out *ipam.PrefixAllocationStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixAllocationStatus_To_ipam_PrefixAllocationStatus(in, out, s)
}

func autoConvert_ipam_PrefixAllocationStatus_To_v1alpha1_PrefixAllocationStatus(in *ipam.PrefixAllocationStatus, out *v1alpha1.PrefixAllocationStatus, s conversion.Scope) error {
	out.Prefix = (*commonv1alpha1.IPPrefix)(unsafe.Pointer(in.Prefix))
	out.LastPhaseTransitionTime = (*metav1.Time)(unsafe.Pointer(in.LastPhaseTransitionTime))
	out.Phase = v1alpha1.PrefixAllocationPhase(in.Phase)
	return nil
}

// Convert_ipam_PrefixAllocationStatus_To_v1alpha1_PrefixAllocationStatus is an autogenerated conversion function.
func Convert_ipam_PrefixAllocationStatus_To_v1alpha1_PrefixAllocationStatus(in *ipam.PrefixAllocationStatus, out *v1alpha1.PrefixAllocationStatus, s conversion.Scope) error {
	return autoConvert_ipam_PrefixAllocationStatus_To_v1alpha1_PrefixAllocationStatus(in, out, s)
}

func autoConvert_v1alpha1_PrefixList_To_ipam_PrefixList(in *v1alpha1.PrefixList, out *ipam.PrefixList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]ipam.Prefix)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_PrefixList_To_ipam_PrefixList is an autogenerated conversion function.
func Convert_v1alpha1_PrefixList_To_ipam_PrefixList(in *v1alpha1.PrefixList, out *ipam.PrefixList, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixList_To_ipam_PrefixList(in, out, s)
}

func autoConvert_ipam_PrefixList_To_v1alpha1_PrefixList(in *ipam.PrefixList, out *v1alpha1.PrefixList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]v1alpha1.Prefix)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_ipam_PrefixList_To_v1alpha1_PrefixList is an autogenerated conversion function.
func Convert_ipam_PrefixList_To_v1alpha1_PrefixList(in *ipam.PrefixList, out *v1alpha1.PrefixList, s conversion.Scope) error {
	return autoConvert_ipam_PrefixList_To_v1alpha1_PrefixList(in, out, s)
}

func autoConvert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec(in *v1alpha1.PrefixSpec, out *ipam.PrefixSpec, s conversion.Scope) error {
	out.IPFamily = v1.IPFamily(in.IPFamily)
	out.Prefix = (*commonv1alpha1.IPPrefix)(unsafe.Pointer(in.Prefix))
	out.PrefixLength = in.PrefixLength
	out.ParentRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.ParentRef))
	out.ParentSelector = (*metav1.LabelSelector)(unsafe.Pointer(in.ParentSelector))
	return nil
}

// Convert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec is an autogenerated conversion function.
func Convert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec(in *v1alpha1.PrefixSpec, out *ipam.PrefixSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec(in, out, s)
}

func autoConvert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec(in *ipam.PrefixSpec, out *v1alpha1.PrefixSpec, s conversion.Scope) error {
	out.IPFamily = v1.IPFamily(in.IPFamily)
	out.Prefix = (*commonv1alpha1.IPPrefix)(unsafe.Pointer(in.Prefix))
	out.PrefixLength = in.PrefixLength
	out.ParentRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.ParentRef))
	out.ParentSelector = (*metav1.LabelSelector)(unsafe.Pointer(in.ParentSelector))
	return nil
}

// Convert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec is an autogenerated conversion function.
func Convert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec(in *ipam.PrefixSpec, out *v1alpha1.PrefixSpec, s conversion.Scope) error {
	return autoConvert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec(in, out, s)
}

func autoConvert_v1alpha1_PrefixStatus_To_ipam_PrefixStatus(in *v1alpha1.PrefixStatus, out *ipam.PrefixStatus, s conversion.Scope) error {
	out.Phase = ipam.PrefixPhase(in.Phase)
	out.LastPhaseTransitionTime = (*metav1.Time)(unsafe.Pointer(in.LastPhaseTransitionTime))
	out.Used = *(*[]commonv1alpha1.IPPrefix)(unsafe.Pointer(&in.Used))
	return nil
}

// Convert_v1alpha1_PrefixStatus_To_ipam_PrefixStatus is an autogenerated conversion function.
func Convert_v1alpha1_PrefixStatus_To_ipam_PrefixStatus(in *v1alpha1.PrefixStatus, out *ipam.PrefixStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixStatus_To_ipam_PrefixStatus(in, out, s)
}

func autoConvert_ipam_PrefixStatus_To_v1alpha1_PrefixStatus(in *ipam.PrefixStatus, out *v1alpha1.PrefixStatus, s conversion.Scope) error {
	out.Phase = v1alpha1.PrefixPhase(in.Phase)
	out.LastPhaseTransitionTime = (*metav1.Time)(unsafe.Pointer(in.LastPhaseTransitionTime))
	out.Used = *(*[]commonv1alpha1.IPPrefix)(unsafe.Pointer(&in.Used))
	return nil
}

// Convert_ipam_PrefixStatus_To_v1alpha1_PrefixStatus is an autogenerated conversion function.
func Convert_ipam_PrefixStatus_To_v1alpha1_PrefixStatus(in *ipam.PrefixStatus, out *v1alpha1.PrefixStatus, s conversion.Scope) error {
	return autoConvert_ipam_PrefixStatus_To_v1alpha1_PrefixStatus(in, out, s)
}

func autoConvert_v1alpha1_PrefixTemplateSpec_To_ipam_PrefixTemplateSpec(in *v1alpha1.PrefixTemplateSpec, out *ipam.PrefixTemplateSpec, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_PrefixSpec_To_ipam_PrefixSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_PrefixTemplateSpec_To_ipam_PrefixTemplateSpec is an autogenerated conversion function.
func Convert_v1alpha1_PrefixTemplateSpec_To_ipam_PrefixTemplateSpec(in *v1alpha1.PrefixTemplateSpec, out *ipam.PrefixTemplateSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_PrefixTemplateSpec_To_ipam_PrefixTemplateSpec(in, out, s)
}

func autoConvert_ipam_PrefixTemplateSpec_To_v1alpha1_PrefixTemplateSpec(in *ipam.PrefixTemplateSpec, out *v1alpha1.PrefixTemplateSpec, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_ipam_PrefixSpec_To_v1alpha1_PrefixSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

// Convert_ipam_PrefixTemplateSpec_To_v1alpha1_PrefixTemplateSpec is an autogenerated conversion function.
func Convert_ipam_PrefixTemplateSpec_To_v1alpha1_PrefixTemplateSpec(in *ipam.PrefixTemplateSpec, out *v1alpha1.PrefixTemplateSpec, s conversion.Scope) error {
	return autoConvert_ipam_PrefixTemplateSpec_To_v1alpha1_PrefixTemplateSpec(in, out, s)
}
