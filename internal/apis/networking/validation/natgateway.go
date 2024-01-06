// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	sphericvalidation "spheric.cloud/spheric/internal/api/validation"
	"spheric.cloud/spheric/internal/apis/networking"
)

// ValidateNATGateway validates an NATGateway object.
func ValidateNATGateway(natGateway *networking.NATGateway) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessor(natGateway, true, apivalidation.NameIsDNSLabel, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateNATGatewaySpec(&natGateway.Spec, field.NewPath("spec"))...)

	return allErrs
}

func validateNATGatewaySpec(spec *networking.NATGatewaySpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateNATGatewayType(spec.Type, fldPath.Child("type"))...)

	allErrs = append(allErrs, sphericvalidation.ValidateIPFamily(spec.IPFamily, fldPath.Child("ipFamily"))...)

	for _, msg := range apivalidation.NameIsDNSLabel(spec.NetworkRef.Name, false) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networkRef").Child("name"), spec.NetworkRef.Name, msg))
	}

	if spec.PortsPerNetworkInterface != nil {
		allErrs = append(allErrs, sphericvalidation.ValidatePowerOfTwo(int64(*spec.PortsPerNetworkInterface), fldPath.Child("portsPerNetworkInterface"))...)
	}

	return allErrs
}

var supportedNATGatewayTypes = sets.New(
	networking.NATGatewayTypePublic,
)

func validateNATGatewayType(natGatewayType networking.NATGatewayType, fldPath *field.Path) field.ErrorList {
	return sphericvalidation.ValidateEnum(supportedNATGatewayTypes, natGatewayType, fldPath, "must specify type")
}

// ValidateNATGatewayUpdate validates a NATGateway object before an update.
func ValidateNATGatewayUpdate(newNATGateway, oldNATGateway *networking.NATGateway) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessorUpdate(newNATGateway, oldNATGateway, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateNATGatewaySpecPrefixUpdate(&newNATGateway.Spec, &oldNATGateway.Spec, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidateNATGateway(newNATGateway)...)

	return allErrs
}

// validateNATGatewaySpecPrefixUpdate validates the spec of a natGateway object before an update.
func validateNATGatewaySpecPrefixUpdate(newSpec, oldSpec *networking.NATGatewaySpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, sphericvalidation.ValidateImmutableField(newSpec.NetworkRef, oldSpec.NetworkRef, fldPath.Child("networkRef"))...)

	return allErrs
}
