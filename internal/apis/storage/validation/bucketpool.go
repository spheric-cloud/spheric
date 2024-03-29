// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	sphericvalidation "spheric.cloud/spheric/internal/api/validation"
	"spheric.cloud/spheric/internal/apis/storage"
)

var ValidateBucketPoolName = apivalidation.NameIsDNSSubdomain

func ValidateBucketPool(bucketPool *storage.BucketPool) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessor(bucketPool, false, ValidateBucketPoolName, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateBucketPoolSpec(&bucketPool.Spec, field.NewPath("spec"))...)

	return allErrs
}

func validateBucketPoolSpec(bucketPoolSpec *storage.BucketPoolSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	return allErrs
}

func ValidateBucketPoolUpdate(newBucketPool, oldBucketPool *storage.BucketPool) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaAccessorUpdate(newBucketPool, oldBucketPool, field.NewPath("metadata"))...)
	allErrs = append(allErrs, validateBucketPoolSpecUpdate(&newBucketPool.Spec, &oldBucketPool.Spec, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidateBucketPool(newBucketPool)...)

	return allErrs
}

func validateBucketPoolSpecUpdate(newSpec, oldSpec *storage.BucketPoolSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if oldSpec.ProviderID != "" {
		allErrs = append(allErrs, sphericvalidation.ValidateImmutableField(newSpec.ProviderID, oldSpec.ProviderID, fldPath.Child("providerID"))...)
	}

	return allErrs
}
