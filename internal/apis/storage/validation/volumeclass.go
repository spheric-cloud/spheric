// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/api/resource"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	sphericvalidation "spheric.cloud/spheric/internal/api/validation"
	"spheric.cloud/spheric/internal/apis/core"
	"spheric.cloud/spheric/internal/apis/storage"
)

func ValidateVolumeClass(volumeClass *storage.VolumeClass) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessor(volumeClass, false, apivalidation.NameIsDNSLabel, field.NewPath("metadata"))...)

	allErrs = append(allErrs, validateVolumeClassCapabilities(volumeClass.Capabilities, field.NewPath("capabilities"))...)

	allErrs = append(allErrs, validateVolumeClassResizePolicy(volumeClass.ResizePolicy, field.NewPath("resizePolicy"))...)

	return allErrs
}

func validateVolumeClassResizePolicy(policy storage.ResizePolicy, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, sphericvalidation.ValidateResizePolicy(policy, fldPath)...)

	return allErrs
}

func validateVolumeClassCapabilities(capabilities core.ResourceList, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	tps := capabilities.Name(core.ResourceTPS, resource.DecimalSI)
	allErrs = append(allErrs, sphericvalidation.ValidatePositiveQuantity(*tps, fldPath.Key(string(core.ResourceTPS)))...)

	iops := capabilities.Name(core.ResourceIOPS, resource.DecimalSI)
	allErrs = append(allErrs, sphericvalidation.ValidatePositiveQuantity(*iops, fldPath.Key(string(core.ResourceIOPS)))...)

	return allErrs
}

func ValidateVolumeClassUpdate(newVolumeClass, oldVolumeClass *storage.VolumeClass) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessorUpdate(newVolumeClass, oldVolumeClass, field.NewPath("metadata"))...)
	allErrs = append(allErrs, sphericvalidation.ValidateImmutableField(newVolumeClass.Capabilities, oldVolumeClass.Capabilities, field.NewPath("capabilities"))...)
	allErrs = append(allErrs, ValidateVolumeClass(newVolumeClass)...)

	return allErrs
}
