// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

// Package networking contains API Schema definitions for the networking internal API group
// +groupName=networking.spheric.cloud
package networking

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "networking.spheric.cloud", Version: runtime.APIVersionInternal}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func Resource(name string) schema.GroupResource {
	return schema.GroupResource{
		Group:    SchemeGroupVersion.Group,
		Resource: name,
	}
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Network{},
		&NetworkList{},
		&NetworkPolicy{},
		&NetworkPolicyList{},
		&NetworkInterface{},
		&NetworkInterfaceList{},
		&VirtualIP{},
		&VirtualIPList{},
		&LoadBalancer{},
		&LoadBalancerList{},
		&LoadBalancerRouting{},
		&LoadBalancerRoutingList{},
		&NATGateway{},
		&NATGatewayList{},
	)
	return nil
}
