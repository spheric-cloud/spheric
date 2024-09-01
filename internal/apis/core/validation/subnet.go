// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/core"
)

func ValidateSubnet(subnet *core.Subnet) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessor(subnet, true, validation.NameIsDNSLabel, field.NewPath("metadata"))...)

	return allErrs
}

func ValidateSubnetUpdate(oldSubnet, newSubnet *core.Subnet) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessorUpdate(oldSubnet, newSubnet, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateSubnet(newSubnet)...)

	return allErrs
}
