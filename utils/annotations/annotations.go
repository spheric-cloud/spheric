// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package annotations

import (
	"time"

	"github.com/ironcore-dev/controller-utils/metautils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

func HasReconcileAnnotation(o metav1.Object) bool {
	return metautils.HasAnnotation(o, corev1alpha1.ReconcileRequestAnnotation)
}

func SetReconcileAnnotation(o metav1.Object) {
	metautils.SetAnnotation(o, corev1alpha1.ReconcileRequestAnnotation, time.Now().Format(time.RFC3339Nano))
}

func RemoveReconcileAnnotation(o metav1.Object) {
	metautils.DeleteAnnotation(o, corev1alpha1.ReconcileRequestAnnotation)
}

func IsExternallyManaged(o metav1.Object) bool {
	return metautils.HasAnnotation(o, corev1alpha1.ManagedByAnnotation)
}

func IsEphemeralManagedBy(o metav1.Object, manager string) bool {
	actual, ok := o.GetAnnotations()[corev1alpha1.EphemeralManagedByAnnotation]
	return ok && actual == manager
}

func IsDefaultEphemeralControlledBy(o metav1.Object, owner metav1.Object) bool {
	return metav1.IsControlledBy(o, owner) && IsEphemeralManagedBy(o, corev1alpha1.DefaultEphemeralManager)
}

func SetDefaultEphemeralManagedBy(o metav1.Object) {
	metautils.SetAnnotation(o, corev1alpha1.EphemeralManagedByAnnotation, corev1alpha1.DefaultEphemeralManager)
}

func SetExternallyMangedBy(o metav1.Object, manager string) {
	metautils.SetAnnotation(o, corev1alpha1.ManagedByAnnotation, manager)
}

func RemoveExternallyManaged(o metav1.Object) {
	metautils.DeleteAnnotation(o, corev1alpha1.ManagedByAnnotation)
}
