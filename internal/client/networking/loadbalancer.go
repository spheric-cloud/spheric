// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package networking

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
)

const (
	LoadBalancerPrefixNamesField = "loadbalancer-prefix-names"

	LoadBalancerNetworkNameField = "loadbalancer-network-name"
)

func SetupLoadBalancerPrefixNamesFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &networkingv1alpha1.LoadBalancer{}, LoadBalancerPrefixNamesField, func(obj client.Object) []string {
		loadBalancer := obj.(*networkingv1alpha1.LoadBalancer)
		return networkingv1alpha1.LoadBalancerPrefixNames(loadBalancer)
	})
}

func SetupLoadBalancerNetworkNameFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &networkingv1alpha1.LoadBalancer{}, LoadBalancerNetworkNameField, func(obj client.Object) []string {
		loadBalancer := obj.(*networkingv1alpha1.LoadBalancer)
		return []string{loadBalancer.Spec.NetworkRef.Name}
	})
}
