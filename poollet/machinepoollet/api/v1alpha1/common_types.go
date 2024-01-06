// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MachineUIDLabel       = "machinepoollet.spheric.cloud/machine-uid"
	MachineNamespaceLabel = "machinepoollet.spheric.cloud/machine-namespace"
	MachineNameLabel      = "machinepoollet.spheric.cloud/machine-name"

	MachineGenerationAnnotation    = "machinepoollet.spheric.cloud/machine-generation"
	SRIMachineGenerationAnnotation = "machinepoollet.spheric.cloud/srimachine-generation"

	NetworkInterfaceMappingAnnotation = "machinepoollet.spheric.cloud/networkinterfacemapping"

	FieldOwner       = "machinepoollet.spheric.cloud/field-owner"
	MachineFinalizer = "machinepoollet.spheric.cloud/machine"

	// DownwardAPIPrefix is the prefix for any downward label.
	DownwardAPIPrefix = "downward-api.machinepoollet.spheric.cloud/"
)

// DownwardAPILabel makes a downward api label name from the given name.
func DownwardAPILabel(name string) string {
	return DownwardAPIPrefix + name
}

// DownwardAPIAnnotation makes a downward api annotation name from the given name.
func DownwardAPIAnnotation(name string) string {
	return DownwardAPIPrefix + name
}

// EncodeNetworkInterfaceMapping encodes the given network interface mapping to be used as an annotation.
func EncodeNetworkInterfaceMapping(nicMapping map[string]ObjectUIDRef) (string, error) {
	data, err := json.Marshal(nicMapping)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func DecodeNetworkInterfaceMapping(nicMappingString string) (map[string]ObjectUIDRef, error) {
	var nicMapping map[string]ObjectUIDRef
	if err := json.Unmarshal([]byte(nicMappingString), &nicMapping); err != nil {
		return nil, err
	}

	return nicMapping, nil
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
