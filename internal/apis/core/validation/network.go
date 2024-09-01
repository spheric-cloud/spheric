// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/core"
)

func ValidateNetwork(network *core.Network) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessor(network, true, validation.NameIsDNSLabel, field.NewPath("metadata"))...)

	return allErrs
}

func ValidateNetworkUpdate(oldNetwork, newNetwork *core.Network) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validation.ValidateObjectMetaAccessorUpdate(oldNetwork, newNetwork, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateNetwork(newNetwork)...)

	return allErrs
}
