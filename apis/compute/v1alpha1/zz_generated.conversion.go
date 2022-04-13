//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
 * Copyright (c) 2021 by the OnMetal authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	commonv1alpha1 "github.com/onmetal/onmetal-api/apis/common/v1alpha1"
	compute "github.com/onmetal/onmetal-api/apis/compute"
	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*EFIVar)(nil), (*compute.EFIVar)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_EFIVar_To_compute_EFIVar(a.(*EFIVar), b.(*compute.EFIVar), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.EFIVar)(nil), (*EFIVar)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_EFIVar_To_v1alpha1_EFIVar(a.(*compute.EFIVar), b.(*EFIVar), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Interface)(nil), (*compute.Interface)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Interface_To_compute_Interface(a.(*Interface), b.(*compute.Interface), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.Interface)(nil), (*Interface)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_Interface_To_v1alpha1_Interface(a.(*compute.Interface), b.(*Interface), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InterfaceStatus)(nil), (*compute.InterfaceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InterfaceStatus_To_compute_InterfaceStatus(a.(*InterfaceStatus), b.(*compute.InterfaceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.InterfaceStatus)(nil), (*InterfaceStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_InterfaceStatus_To_v1alpha1_InterfaceStatus(a.(*compute.InterfaceStatus), b.(*InterfaceStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Machine)(nil), (*compute.Machine)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Machine_To_compute_Machine(a.(*Machine), b.(*compute.Machine), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.Machine)(nil), (*Machine)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_Machine_To_v1alpha1_Machine(a.(*compute.Machine), b.(*Machine), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineClass)(nil), (*compute.MachineClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineClass_To_compute_MachineClass(a.(*MachineClass), b.(*compute.MachineClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachineClass)(nil), (*MachineClass)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachineClass_To_v1alpha1_MachineClass(a.(*compute.MachineClass), b.(*MachineClass), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineClassList)(nil), (*compute.MachineClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineClassList_To_compute_MachineClassList(a.(*MachineClassList), b.(*compute.MachineClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachineClassList)(nil), (*MachineClassList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachineClassList_To_v1alpha1_MachineClassList(a.(*compute.MachineClassList), b.(*MachineClassList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineCondition)(nil), (*compute.MachineCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineCondition_To_compute_MachineCondition(a.(*MachineCondition), b.(*compute.MachineCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachineCondition)(nil), (*MachineCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachineCondition_To_v1alpha1_MachineCondition(a.(*compute.MachineCondition), b.(*MachineCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineList)(nil), (*compute.MachineList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineList_To_compute_MachineList(a.(*MachineList), b.(*compute.MachineList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachineList)(nil), (*MachineList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachineList_To_v1alpha1_MachineList(a.(*compute.MachineList), b.(*MachineList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePool)(nil), (*compute.MachinePool)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachinePool_To_compute_MachinePool(a.(*MachinePool), b.(*compute.MachinePool), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachinePool)(nil), (*MachinePool)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachinePool_To_v1alpha1_MachinePool(a.(*compute.MachinePool), b.(*MachinePool), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePoolCondition)(nil), (*compute.MachinePoolCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachinePoolCondition_To_compute_MachinePoolCondition(a.(*MachinePoolCondition), b.(*compute.MachinePoolCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachinePoolCondition)(nil), (*MachinePoolCondition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachinePoolCondition_To_v1alpha1_MachinePoolCondition(a.(*compute.MachinePoolCondition), b.(*MachinePoolCondition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePoolList)(nil), (*compute.MachinePoolList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachinePoolList_To_compute_MachinePoolList(a.(*MachinePoolList), b.(*compute.MachinePoolList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachinePoolList)(nil), (*MachinePoolList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachinePoolList_To_v1alpha1_MachinePoolList(a.(*compute.MachinePoolList), b.(*MachinePoolList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePoolSpec)(nil), (*compute.MachinePoolSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachinePoolSpec_To_compute_MachinePoolSpec(a.(*MachinePoolSpec), b.(*compute.MachinePoolSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachinePoolSpec)(nil), (*MachinePoolSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachinePoolSpec_To_v1alpha1_MachinePoolSpec(a.(*compute.MachinePoolSpec), b.(*MachinePoolSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePoolStatus)(nil), (*compute.MachinePoolStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachinePoolStatus_To_compute_MachinePoolStatus(a.(*MachinePoolStatus), b.(*compute.MachinePoolStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachinePoolStatus)(nil), (*MachinePoolStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachinePoolStatus_To_v1alpha1_MachinePoolStatus(a.(*compute.MachinePoolStatus), b.(*MachinePoolStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineSpec)(nil), (*compute.MachineSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineSpec_To_compute_MachineSpec(a.(*MachineSpec), b.(*compute.MachineSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachineSpec)(nil), (*MachineSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachineSpec_To_v1alpha1_MachineSpec(a.(*compute.MachineSpec), b.(*MachineSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineStatus)(nil), (*compute.MachineStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineStatus_To_compute_MachineStatus(a.(*MachineStatus), b.(*compute.MachineStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.MachineStatus)(nil), (*MachineStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_MachineStatus_To_v1alpha1_MachineStatus(a.(*compute.MachineStatus), b.(*MachineStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*Volume)(nil), (*compute.Volume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Volume_To_compute_Volume(a.(*Volume), b.(*compute.Volume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.Volume)(nil), (*Volume)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_Volume_To_v1alpha1_Volume(a.(*compute.Volume), b.(*Volume), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*VolumeSource)(nil), (*compute.VolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_VolumeSource_To_compute_VolumeSource(a.(*VolumeSource), b.(*compute.VolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.VolumeSource)(nil), (*VolumeSource)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_VolumeSource_To_v1alpha1_VolumeSource(a.(*compute.VolumeSource), b.(*VolumeSource), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*VolumeStatus)(nil), (*compute.VolumeStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_VolumeStatus_To_compute_VolumeStatus(a.(*VolumeStatus), b.(*compute.VolumeStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*compute.VolumeStatus)(nil), (*VolumeStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_compute_VolumeStatus_To_v1alpha1_VolumeStatus(a.(*compute.VolumeStatus), b.(*VolumeStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_EFIVar_To_compute_EFIVar(in *EFIVar, out *compute.EFIVar, s conversion.Scope) error {
	out.Name = in.Name
	out.UUID = in.UUID
	out.Value = in.Value
	return nil
}

// Convert_v1alpha1_EFIVar_To_compute_EFIVar is an autogenerated conversion function.
func Convert_v1alpha1_EFIVar_To_compute_EFIVar(in *EFIVar, out *compute.EFIVar, s conversion.Scope) error {
	return autoConvert_v1alpha1_EFIVar_To_compute_EFIVar(in, out, s)
}

func autoConvert_compute_EFIVar_To_v1alpha1_EFIVar(in *compute.EFIVar, out *EFIVar, s conversion.Scope) error {
	out.Name = in.Name
	out.UUID = in.UUID
	out.Value = in.Value
	return nil
}

// Convert_compute_EFIVar_To_v1alpha1_EFIVar is an autogenerated conversion function.
func Convert_compute_EFIVar_To_v1alpha1_EFIVar(in *compute.EFIVar, out *EFIVar, s conversion.Scope) error {
	return autoConvert_compute_EFIVar_To_v1alpha1_EFIVar(in, out, s)
}

func autoConvert_v1alpha1_Interface_To_compute_Interface(in *Interface, out *compute.Interface, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha1_Interface_To_compute_Interface is an autogenerated conversion function.
func Convert_v1alpha1_Interface_To_compute_Interface(in *Interface, out *compute.Interface, s conversion.Scope) error {
	return autoConvert_v1alpha1_Interface_To_compute_Interface(in, out, s)
}

func autoConvert_compute_Interface_To_v1alpha1_Interface(in *compute.Interface, out *Interface, s conversion.Scope) error {
	return nil
}

// Convert_compute_Interface_To_v1alpha1_Interface is an autogenerated conversion function.
func Convert_compute_Interface_To_v1alpha1_Interface(in *compute.Interface, out *Interface, s conversion.Scope) error {
	return autoConvert_compute_Interface_To_v1alpha1_Interface(in, out, s)
}

func autoConvert_v1alpha1_InterfaceStatus_To_compute_InterfaceStatus(in *InterfaceStatus, out *compute.InterfaceStatus, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha1_InterfaceStatus_To_compute_InterfaceStatus is an autogenerated conversion function.
func Convert_v1alpha1_InterfaceStatus_To_compute_InterfaceStatus(in *InterfaceStatus, out *compute.InterfaceStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_InterfaceStatus_To_compute_InterfaceStatus(in, out, s)
}

func autoConvert_compute_InterfaceStatus_To_v1alpha1_InterfaceStatus(in *compute.InterfaceStatus, out *InterfaceStatus, s conversion.Scope) error {
	return nil
}

// Convert_compute_InterfaceStatus_To_v1alpha1_InterfaceStatus is an autogenerated conversion function.
func Convert_compute_InterfaceStatus_To_v1alpha1_InterfaceStatus(in *compute.InterfaceStatus, out *InterfaceStatus, s conversion.Scope) error {
	return autoConvert_compute_InterfaceStatus_To_v1alpha1_InterfaceStatus(in, out, s)
}

func autoConvert_v1alpha1_Machine_To_compute_Machine(in *Machine, out *compute.Machine, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_MachineSpec_To_compute_MachineSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_MachineStatus_To_compute_MachineStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Machine_To_compute_Machine is an autogenerated conversion function.
func Convert_v1alpha1_Machine_To_compute_Machine(in *Machine, out *compute.Machine, s conversion.Scope) error {
	return autoConvert_v1alpha1_Machine_To_compute_Machine(in, out, s)
}

func autoConvert_compute_Machine_To_v1alpha1_Machine(in *compute.Machine, out *Machine, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_compute_MachineSpec_To_v1alpha1_MachineSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_compute_MachineStatus_To_v1alpha1_MachineStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_compute_Machine_To_v1alpha1_Machine is an autogenerated conversion function.
func Convert_compute_Machine_To_v1alpha1_Machine(in *compute.Machine, out *Machine, s conversion.Scope) error {
	return autoConvert_compute_Machine_To_v1alpha1_Machine(in, out, s)
}

func autoConvert_v1alpha1_MachineClass_To_compute_MachineClass(in *MachineClass, out *compute.MachineClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Capabilities = *(*v1.ResourceList)(unsafe.Pointer(&in.Capabilities))
	return nil
}

// Convert_v1alpha1_MachineClass_To_compute_MachineClass is an autogenerated conversion function.
func Convert_v1alpha1_MachineClass_To_compute_MachineClass(in *MachineClass, out *compute.MachineClass, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineClass_To_compute_MachineClass(in, out, s)
}

func autoConvert_compute_MachineClass_To_v1alpha1_MachineClass(in *compute.MachineClass, out *MachineClass, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.Capabilities = *(*v1.ResourceList)(unsafe.Pointer(&in.Capabilities))
	return nil
}

// Convert_compute_MachineClass_To_v1alpha1_MachineClass is an autogenerated conversion function.
func Convert_compute_MachineClass_To_v1alpha1_MachineClass(in *compute.MachineClass, out *MachineClass, s conversion.Scope) error {
	return autoConvert_compute_MachineClass_To_v1alpha1_MachineClass(in, out, s)
}

func autoConvert_v1alpha1_MachineClassList_To_compute_MachineClassList(in *MachineClassList, out *compute.MachineClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]compute.MachineClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_MachineClassList_To_compute_MachineClassList is an autogenerated conversion function.
func Convert_v1alpha1_MachineClassList_To_compute_MachineClassList(in *MachineClassList, out *compute.MachineClassList, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineClassList_To_compute_MachineClassList(in, out, s)
}

func autoConvert_compute_MachineClassList_To_v1alpha1_MachineClassList(in *compute.MachineClassList, out *MachineClassList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]MachineClass)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_compute_MachineClassList_To_v1alpha1_MachineClassList is an autogenerated conversion function.
func Convert_compute_MachineClassList_To_v1alpha1_MachineClassList(in *compute.MachineClassList, out *MachineClassList, s conversion.Scope) error {
	return autoConvert_compute_MachineClassList_To_v1alpha1_MachineClassList(in, out, s)
}

func autoConvert_v1alpha1_MachineCondition_To_compute_MachineCondition(in *MachineCondition, out *compute.MachineCondition, s conversion.Scope) error {
	out.Type = compute.MachineConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	out.Reason = in.Reason
	out.Message = in.Message
	out.ObservedGeneration = in.ObservedGeneration
	out.LastUpdateTime = in.LastUpdateTime
	out.LastTransitionTime = in.LastTransitionTime
	return nil
}

// Convert_v1alpha1_MachineCondition_To_compute_MachineCondition is an autogenerated conversion function.
func Convert_v1alpha1_MachineCondition_To_compute_MachineCondition(in *MachineCondition, out *compute.MachineCondition, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineCondition_To_compute_MachineCondition(in, out, s)
}

func autoConvert_compute_MachineCondition_To_v1alpha1_MachineCondition(in *compute.MachineCondition, out *MachineCondition, s conversion.Scope) error {
	out.Type = MachineConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	out.Reason = in.Reason
	out.Message = in.Message
	out.ObservedGeneration = in.ObservedGeneration
	out.LastUpdateTime = in.LastUpdateTime
	out.LastTransitionTime = in.LastTransitionTime
	return nil
}

// Convert_compute_MachineCondition_To_v1alpha1_MachineCondition is an autogenerated conversion function.
func Convert_compute_MachineCondition_To_v1alpha1_MachineCondition(in *compute.MachineCondition, out *MachineCondition, s conversion.Scope) error {
	return autoConvert_compute_MachineCondition_To_v1alpha1_MachineCondition(in, out, s)
}

func autoConvert_v1alpha1_MachineList_To_compute_MachineList(in *MachineList, out *compute.MachineList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]compute.Machine)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_MachineList_To_compute_MachineList is an autogenerated conversion function.
func Convert_v1alpha1_MachineList_To_compute_MachineList(in *MachineList, out *compute.MachineList, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineList_To_compute_MachineList(in, out, s)
}

func autoConvert_compute_MachineList_To_v1alpha1_MachineList(in *compute.MachineList, out *MachineList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]Machine)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_compute_MachineList_To_v1alpha1_MachineList is an autogenerated conversion function.
func Convert_compute_MachineList_To_v1alpha1_MachineList(in *compute.MachineList, out *MachineList, s conversion.Scope) error {
	return autoConvert_compute_MachineList_To_v1alpha1_MachineList(in, out, s)
}

func autoConvert_v1alpha1_MachinePool_To_compute_MachinePool(in *MachinePool, out *compute.MachinePool, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_MachinePoolSpec_To_compute_MachinePoolSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_MachinePoolStatus_To_compute_MachinePoolStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_MachinePool_To_compute_MachinePool is an autogenerated conversion function.
func Convert_v1alpha1_MachinePool_To_compute_MachinePool(in *MachinePool, out *compute.MachinePool, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachinePool_To_compute_MachinePool(in, out, s)
}

func autoConvert_compute_MachinePool_To_v1alpha1_MachinePool(in *compute.MachinePool, out *MachinePool, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_compute_MachinePoolSpec_To_v1alpha1_MachinePoolSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_compute_MachinePoolStatus_To_v1alpha1_MachinePoolStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_compute_MachinePool_To_v1alpha1_MachinePool is an autogenerated conversion function.
func Convert_compute_MachinePool_To_v1alpha1_MachinePool(in *compute.MachinePool, out *MachinePool, s conversion.Scope) error {
	return autoConvert_compute_MachinePool_To_v1alpha1_MachinePool(in, out, s)
}

func autoConvert_v1alpha1_MachinePoolCondition_To_compute_MachinePoolCondition(in *MachinePoolCondition, out *compute.MachinePoolCondition, s conversion.Scope) error {
	out.Type = compute.MachinePoolConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	out.Reason = in.Reason
	out.Message = in.Message
	out.ObservedGeneration = in.ObservedGeneration
	out.LastUpdateTime = in.LastUpdateTime
	out.LastTransitionTime = in.LastTransitionTime
	return nil
}

// Convert_v1alpha1_MachinePoolCondition_To_compute_MachinePoolCondition is an autogenerated conversion function.
func Convert_v1alpha1_MachinePoolCondition_To_compute_MachinePoolCondition(in *MachinePoolCondition, out *compute.MachinePoolCondition, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachinePoolCondition_To_compute_MachinePoolCondition(in, out, s)
}

func autoConvert_compute_MachinePoolCondition_To_v1alpha1_MachinePoolCondition(in *compute.MachinePoolCondition, out *MachinePoolCondition, s conversion.Scope) error {
	out.Type = MachinePoolConditionType(in.Type)
	out.Status = v1.ConditionStatus(in.Status)
	out.Reason = in.Reason
	out.Message = in.Message
	out.ObservedGeneration = in.ObservedGeneration
	out.LastUpdateTime = in.LastUpdateTime
	out.LastTransitionTime = in.LastTransitionTime
	return nil
}

// Convert_compute_MachinePoolCondition_To_v1alpha1_MachinePoolCondition is an autogenerated conversion function.
func Convert_compute_MachinePoolCondition_To_v1alpha1_MachinePoolCondition(in *compute.MachinePoolCondition, out *MachinePoolCondition, s conversion.Scope) error {
	return autoConvert_compute_MachinePoolCondition_To_v1alpha1_MachinePoolCondition(in, out, s)
}

func autoConvert_v1alpha1_MachinePoolList_To_compute_MachinePoolList(in *MachinePoolList, out *compute.MachinePoolList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]compute.MachinePool)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_MachinePoolList_To_compute_MachinePoolList is an autogenerated conversion function.
func Convert_v1alpha1_MachinePoolList_To_compute_MachinePoolList(in *MachinePoolList, out *compute.MachinePoolList, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachinePoolList_To_compute_MachinePoolList(in, out, s)
}

func autoConvert_compute_MachinePoolList_To_v1alpha1_MachinePoolList(in *compute.MachinePoolList, out *MachinePoolList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]MachinePool)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_compute_MachinePoolList_To_v1alpha1_MachinePoolList is an autogenerated conversion function.
func Convert_compute_MachinePoolList_To_v1alpha1_MachinePoolList(in *compute.MachinePoolList, out *MachinePoolList, s conversion.Scope) error {
	return autoConvert_compute_MachinePoolList_To_v1alpha1_MachinePoolList(in, out, s)
}

func autoConvert_v1alpha1_MachinePoolSpec_To_compute_MachinePoolSpec(in *MachinePoolSpec, out *compute.MachinePoolSpec, s conversion.Scope) error {
	out.ProviderID = in.ProviderID
	out.Taints = *(*[]commonv1alpha1.Taint)(unsafe.Pointer(&in.Taints))
	return nil
}

// Convert_v1alpha1_MachinePoolSpec_To_compute_MachinePoolSpec is an autogenerated conversion function.
func Convert_v1alpha1_MachinePoolSpec_To_compute_MachinePoolSpec(in *MachinePoolSpec, out *compute.MachinePoolSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachinePoolSpec_To_compute_MachinePoolSpec(in, out, s)
}

func autoConvert_compute_MachinePoolSpec_To_v1alpha1_MachinePoolSpec(in *compute.MachinePoolSpec, out *MachinePoolSpec, s conversion.Scope) error {
	out.ProviderID = in.ProviderID
	out.Taints = *(*[]commonv1alpha1.Taint)(unsafe.Pointer(&in.Taints))
	return nil
}

// Convert_compute_MachinePoolSpec_To_v1alpha1_MachinePoolSpec is an autogenerated conversion function.
func Convert_compute_MachinePoolSpec_To_v1alpha1_MachinePoolSpec(in *compute.MachinePoolSpec, out *MachinePoolSpec, s conversion.Scope) error {
	return autoConvert_compute_MachinePoolSpec_To_v1alpha1_MachinePoolSpec(in, out, s)
}

func autoConvert_v1alpha1_MachinePoolStatus_To_compute_MachinePoolStatus(in *MachinePoolStatus, out *compute.MachinePoolStatus, s conversion.Scope) error {
	out.State = compute.MachinePoolState(in.State)
	out.Conditions = *(*[]compute.MachinePoolCondition)(unsafe.Pointer(&in.Conditions))
	out.AvailableMachineClasses = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.AvailableMachineClasses))
	return nil
}

// Convert_v1alpha1_MachinePoolStatus_To_compute_MachinePoolStatus is an autogenerated conversion function.
func Convert_v1alpha1_MachinePoolStatus_To_compute_MachinePoolStatus(in *MachinePoolStatus, out *compute.MachinePoolStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachinePoolStatus_To_compute_MachinePoolStatus(in, out, s)
}

func autoConvert_compute_MachinePoolStatus_To_v1alpha1_MachinePoolStatus(in *compute.MachinePoolStatus, out *MachinePoolStatus, s conversion.Scope) error {
	out.State = MachinePoolState(in.State)
	out.Conditions = *(*[]MachinePoolCondition)(unsafe.Pointer(&in.Conditions))
	out.AvailableMachineClasses = *(*[]v1.LocalObjectReference)(unsafe.Pointer(&in.AvailableMachineClasses))
	return nil
}

// Convert_compute_MachinePoolStatus_To_v1alpha1_MachinePoolStatus is an autogenerated conversion function.
func Convert_compute_MachinePoolStatus_To_v1alpha1_MachinePoolStatus(in *compute.MachinePoolStatus, out *MachinePoolStatus, s conversion.Scope) error {
	return autoConvert_compute_MachinePoolStatus_To_v1alpha1_MachinePoolStatus(in, out, s)
}

func autoConvert_v1alpha1_MachineSpec_To_compute_MachineSpec(in *MachineSpec, out *compute.MachineSpec, s conversion.Scope) error {
	out.MachineClassRef = in.MachineClassRef
	out.MachinePoolSelector = *(*map[string]string)(unsafe.Pointer(&in.MachinePoolSelector))
	out.MachinePoolRef = in.MachinePoolRef
	out.Image = in.Image
	out.Interfaces = *(*[]compute.Interface)(unsafe.Pointer(&in.Interfaces))
	out.Volumes = *(*[]compute.Volume)(unsafe.Pointer(&in.Volumes))
	out.IgnitionRef = (*commonv1alpha1.ConfigMapKeySelector)(unsafe.Pointer(in.IgnitionRef))
	out.EFIVars = *(*[]compute.EFIVar)(unsafe.Pointer(&in.EFIVars))
	out.Tolerations = *(*[]commonv1alpha1.Toleration)(unsafe.Pointer(&in.Tolerations))
	return nil
}

// Convert_v1alpha1_MachineSpec_To_compute_MachineSpec is an autogenerated conversion function.
func Convert_v1alpha1_MachineSpec_To_compute_MachineSpec(in *MachineSpec, out *compute.MachineSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineSpec_To_compute_MachineSpec(in, out, s)
}

func autoConvert_compute_MachineSpec_To_v1alpha1_MachineSpec(in *compute.MachineSpec, out *MachineSpec, s conversion.Scope) error {
	out.MachineClassRef = in.MachineClassRef
	out.MachinePoolSelector = *(*map[string]string)(unsafe.Pointer(&in.MachinePoolSelector))
	out.MachinePoolRef = in.MachinePoolRef
	out.Image = in.Image
	out.Interfaces = *(*[]Interface)(unsafe.Pointer(&in.Interfaces))
	out.Volumes = *(*[]Volume)(unsafe.Pointer(&in.Volumes))
	out.IgnitionRef = (*commonv1alpha1.ConfigMapKeySelector)(unsafe.Pointer(in.IgnitionRef))
	out.EFIVars = *(*[]EFIVar)(unsafe.Pointer(&in.EFIVars))
	out.Tolerations = *(*[]commonv1alpha1.Toleration)(unsafe.Pointer(&in.Tolerations))
	return nil
}

// Convert_compute_MachineSpec_To_v1alpha1_MachineSpec is an autogenerated conversion function.
func Convert_compute_MachineSpec_To_v1alpha1_MachineSpec(in *compute.MachineSpec, out *MachineSpec, s conversion.Scope) error {
	return autoConvert_compute_MachineSpec_To_v1alpha1_MachineSpec(in, out, s)
}

func autoConvert_v1alpha1_MachineStatus_To_compute_MachineStatus(in *MachineStatus, out *compute.MachineStatus, s conversion.Scope) error {
	out.State = compute.MachineState(in.State)
	out.Conditions = *(*[]compute.MachineCondition)(unsafe.Pointer(&in.Conditions))
	out.Interfaces = *(*[]compute.InterfaceStatus)(unsafe.Pointer(&in.Interfaces))
	out.VolumeAttachments = *(*[]compute.VolumeStatus)(unsafe.Pointer(&in.VolumeAttachments))
	return nil
}

// Convert_v1alpha1_MachineStatus_To_compute_MachineStatus is an autogenerated conversion function.
func Convert_v1alpha1_MachineStatus_To_compute_MachineStatus(in *MachineStatus, out *compute.MachineStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineStatus_To_compute_MachineStatus(in, out, s)
}

func autoConvert_compute_MachineStatus_To_v1alpha1_MachineStatus(in *compute.MachineStatus, out *MachineStatus, s conversion.Scope) error {
	out.State = MachineState(in.State)
	out.Conditions = *(*[]MachineCondition)(unsafe.Pointer(&in.Conditions))
	out.Interfaces = *(*[]InterfaceStatus)(unsafe.Pointer(&in.Interfaces))
	out.VolumeAttachments = *(*[]VolumeStatus)(unsafe.Pointer(&in.VolumeAttachments))
	return nil
}

// Convert_compute_MachineStatus_To_v1alpha1_MachineStatus is an autogenerated conversion function.
func Convert_compute_MachineStatus_To_v1alpha1_MachineStatus(in *compute.MachineStatus, out *MachineStatus, s conversion.Scope) error {
	return autoConvert_compute_MachineStatus_To_v1alpha1_MachineStatus(in, out, s)
}

func autoConvert_v1alpha1_Volume_To_compute_Volume(in *Volume, out *compute.Volume, s conversion.Scope) error {
	out.Name = in.Name
	if err := Convert_v1alpha1_VolumeSource_To_compute_VolumeSource(&in.VolumeSource, &out.VolumeSource, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_Volume_To_compute_Volume is an autogenerated conversion function.
func Convert_v1alpha1_Volume_To_compute_Volume(in *Volume, out *compute.Volume, s conversion.Scope) error {
	return autoConvert_v1alpha1_Volume_To_compute_Volume(in, out, s)
}

func autoConvert_compute_Volume_To_v1alpha1_Volume(in *compute.Volume, out *Volume, s conversion.Scope) error {
	out.Name = in.Name
	if err := Convert_compute_VolumeSource_To_v1alpha1_VolumeSource(&in.VolumeSource, &out.VolumeSource, s); err != nil {
		return err
	}
	return nil
}

// Convert_compute_Volume_To_v1alpha1_Volume is an autogenerated conversion function.
func Convert_compute_Volume_To_v1alpha1_Volume(in *compute.Volume, out *Volume, s conversion.Scope) error {
	return autoConvert_compute_Volume_To_v1alpha1_Volume(in, out, s)
}

func autoConvert_v1alpha1_VolumeSource_To_compute_VolumeSource(in *VolumeSource, out *compute.VolumeSource, s conversion.Scope) error {
	out.VolumeClaimRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.VolumeClaimRef))
	return nil
}

// Convert_v1alpha1_VolumeSource_To_compute_VolumeSource is an autogenerated conversion function.
func Convert_v1alpha1_VolumeSource_To_compute_VolumeSource(in *VolumeSource, out *compute.VolumeSource, s conversion.Scope) error {
	return autoConvert_v1alpha1_VolumeSource_To_compute_VolumeSource(in, out, s)
}

func autoConvert_compute_VolumeSource_To_v1alpha1_VolumeSource(in *compute.VolumeSource, out *VolumeSource, s conversion.Scope) error {
	out.VolumeClaimRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.VolumeClaimRef))
	return nil
}

// Convert_compute_VolumeSource_To_v1alpha1_VolumeSource is an autogenerated conversion function.
func Convert_compute_VolumeSource_To_v1alpha1_VolumeSource(in *compute.VolumeSource, out *VolumeSource, s conversion.Scope) error {
	return autoConvert_compute_VolumeSource_To_v1alpha1_VolumeSource(in, out, s)
}

func autoConvert_v1alpha1_VolumeStatus_To_compute_VolumeStatus(in *VolumeStatus, out *compute.VolumeStatus, s conversion.Scope) error {
	out.Name = in.Name
	out.DeviceID = in.DeviceID
	return nil
}

// Convert_v1alpha1_VolumeStatus_To_compute_VolumeStatus is an autogenerated conversion function.
func Convert_v1alpha1_VolumeStatus_To_compute_VolumeStatus(in *VolumeStatus, out *compute.VolumeStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_VolumeStatus_To_compute_VolumeStatus(in, out, s)
}

func autoConvert_compute_VolumeStatus_To_v1alpha1_VolumeStatus(in *compute.VolumeStatus, out *VolumeStatus, s conversion.Scope) error {
	out.Name = in.Name
	out.DeviceID = in.DeviceID
	return nil
}

// Convert_compute_VolumeStatus_To_v1alpha1_VolumeStatus is an autogenerated conversion function.
func Convert_compute_VolumeStatus_To_v1alpha1_VolumeStatus(in *compute.VolumeStatus, out *VolumeStatus, s conversion.Scope) error {
	return autoConvert_compute_VolumeStatus_To_v1alpha1_VolumeStatus(in, out, s)
}