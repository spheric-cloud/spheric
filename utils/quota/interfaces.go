// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package quota

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

type Evaluator interface {
	Type() client.Object
	MatchesResourceName(resourceName corev1alpha1.ResourceName) bool
	MatchesResourceScopeSelectorRequirement(item client.Object, req corev1alpha1.ResourceScopeSelectorRequirement) (bool, error)
	Usage(ctx context.Context, item client.Object) (corev1alpha1.ResourceList, error)
}

type Registry interface {
	// Add to registry
	Add(e Evaluator) error
	// Remove from registry
	Remove(obj client.Object) error
	// Get by group resource
	Get(obj client.Object) (Evaluator, error)
	// List from registry
	List() []Evaluator
}

type UsageCalculator interface {
	CalculateUsage(ctx context.Context, quota *corev1alpha1.ResourceQuota) (corev1alpha1.ResourceList, error)
}
