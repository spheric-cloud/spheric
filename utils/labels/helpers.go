// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package labels

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

// HasWatchLabel returns true if the object has a label with the WatchLabel key matching the given value.
func HasWatchLabel(o metav1.Object, labelValue string) bool {
	val, ok := o.GetLabels()[corev1alpha1.WatchLabel]
	if !ok {
		return false
	}
	return val == labelValue
}
