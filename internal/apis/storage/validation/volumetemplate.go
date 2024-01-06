// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	sphericvalidation "spheric.cloud/spheric/internal/api/validation"
	"spheric.cloud/spheric/internal/apis/storage"
)

var allowedVolumeTemplateObjectMetaFields = sets.New(
	"Annotations",
	"Labels",
)

func validateVolumeTemplateSpecMetadata(objMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, metav1validation.ValidateLabels(objMeta.Labels, fldPath.Child("labels"))...)
	allErrs = append(allErrs, apivalidation.ValidateAnnotations(objMeta.Annotations, fldPath.Child("annotations"))...)
	allErrs = append(allErrs, sphericvalidation.ValidateFieldAllowList(*objMeta, allowedVolumeTemplateObjectMetaFields, "cannot be set for a volume template", fldPath)...)

	return allErrs
}

// ValidateVolumeTemplateSpec validates the spec of a volume template.
func ValidateVolumeTemplateSpec(spec *storage.VolumeTemplateSpec, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateVolumeTemplateSpecMetadata(&spec.ObjectMeta, fldPath.Child("metadata"))...)
	allErrs = append(allErrs, validateVolumeSpec(&spec.Spec, fldPath.Child("spec"))...)

	return allErrs
}
