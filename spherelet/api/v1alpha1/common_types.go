// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	InstanceUIDLabel       = "spherelet.spheric.cloud/instance-uid"
	InstanceNamespaceLabel = "spherelet.spheric.cloud/instance-namespace"
	InstanceNameLabel      = "spherelet.spheric.cloud/instance-name"

	InstanceGenerationAnnotation    = "spherelet.spheric.cloud/instance-generation"
	IRIInstanceGenerationAnnotation = "spherelet.spheric.cloud/iriinstance-generation"

	FieldOwner        = "spherelet.spheric.cloud/field-owner"
	InstanceFinalizer = "spherelet.spheric.cloud/instance"

	// DownwardAPIPrefix is the prefix for any downward label.
	DownwardAPIPrefix = "downward-api.spherelet.spheric.cloud/"
)

// DownwardAPILabel makes a downward api label name from the given name.
func DownwardAPILabel(name string) string {
	return DownwardAPIPrefix + name
}

// DownwardAPIAnnotation makes a downward api annotation name from the given name.
func DownwardAPIAnnotation(name string) string {
	return DownwardAPIPrefix + name
}

// ObjectUIDRef is a name-uid-reference to an object.
type ObjectUIDRef struct {
	Name string    `json:"name"`
	UID  types.UID `json:"uid"`
}

func ObjUID(obj client.Object) ObjectUIDRef {
	return ObjectUIDRef{
		Name: obj.GetName(),
		UID:  obj.GetUID(),
	}
}
