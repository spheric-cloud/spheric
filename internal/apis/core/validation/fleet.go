// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/core"
)

func ValidateFleet(fleet *core.Fleet) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessor(fleet, false, validation.NameIsDNSLabel, field.NewPath("metadata"))...)

	return allErrs
}

func ValidateFleetUpdate(oldFleet, newFleet *core.Fleet) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessorUpdate(oldFleet, newFleet, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateFleet(newFleet)...)

	return allErrs
}
