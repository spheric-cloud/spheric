// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/core"
)

func ValidateInstanceType(instanceType *core.InstanceType) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessor(instanceType, false, validation.NameIsDNSLabel, field.NewPath("metadata"))...)

	return allErrs
}

func ValidateInstanceTypeUpdate(oldInstanceType, newInstanceType *core.InstanceType) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessorUpdate(oldInstanceType, newInstanceType, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateInstanceType(newInstanceType)...)

	return allErrs
}
