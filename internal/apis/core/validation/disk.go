// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/core"
)

func ValidateDisk(disk *core.Disk) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessor(disk, true, validation.NameIsDNSLabel, field.NewPath("metadata"))...)

	return allErrs
}

func ValidateDiskUpdate(oldDisk, newDisk *core.Disk) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessorUpdate(oldDisk, newDisk, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateDisk(newDisk)...)

	return allErrs
}
