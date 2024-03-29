// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	sphericvalidation "spheric.cloud/spheric/internal/api/validation"
	"spheric.cloud/spheric/internal/apis/compute"
)

var ValidateMachinePoolName = apivalidation.NameIsDNSSubdomain

func ValidateMachinePool(machinePool *compute.MachinePool) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessor(machinePool, false, ValidateMachinePoolName, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateMachinePoolSpec(&machinePool.Spec, field.NewPath("spec"))...)

	return allErrs
}

func validateMachinePoolSpec(machinePoolSpec *compute.MachinePoolSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	return allErrs
}

func ValidateMachinePoolUpdate(newMachinePool, oldMachinePool *compute.MachinePool) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessorUpdate(newMachinePool, oldMachinePool, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateMachinePoolSpecUpdate(&newMachinePool.Spec, &oldMachinePool.Spec, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidateMachinePool(newMachinePool)...)

	return allErrs
}

func validateMachinePoolSpecUpdate(newSpec, oldSpec *compute.MachinePoolSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if oldSpec.ProviderID != "" {
		allErrs = append(allErrs, sphericvalidation.ValidateImmutableField(newSpec.ProviderID, oldSpec.ProviderID, fldPath.Child("providerID"))...)
	}

	return allErrs
}
