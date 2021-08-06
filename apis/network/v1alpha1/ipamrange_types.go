/*
 * Copyright (c) 2021 by the OnMetal authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	common "github.com/onmetal/onmetal-api/apis/common/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var IPAMRangeGK = schema.GroupKind{
	Group: GroupVersion.Group,
	Kind:  "IPAMRange",
}

// IPAMRangeSpec defines the desired state of IPAMRange
// Either parent and size or a give CIDR must be specified. If parent is specified,
// the effective range of the given size is allocated from the parent IP range. If parent and CIDR
// is defined, the given CIDR must be in the parent range and unused. It will be allocated if possible.
// Otherwise the status of the object will be set to "Failed".
type IPAMRangeSpec struct {
	Parent *common.ScopedReference `json:"parent,omitempty"`
	Size   string                  `json:"size,omitempty"`

	CIDR string `json:"cidr,omitempty"`
}

// IPAMRangeStatus defines the observed state of IPAMRange
type IPAMRangeStatus struct {
	common.StateFields `json:",inline"`
	Bound              *common.ScopedKindReference `json:"bound,omitempty"`
	CIDR               string                      `json:"cidr,omitempty"`
	FreeBlocks         []string                    `json:"freeBlocks,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=ipr
//+kubebuilder:printcolumn:name="RequestedCIDR",type=string,JSONPath=`.spec.cidr`
//+kubebuilder:printcolumn:name="RequestedSize",type=string,JSONPath=`.spec.size`
//+kubebuilder:printcolumn:name="EffectiveCIDR",type=string,JSONPath=`.status.cidr`
//+kubebuilder:printcolumn:name="Parent",type=string,JSONPath=`.spec.parent.name`
//+kubebuilder:printcolumn:name="BoundKind",type=string,JSONPath=`.status.bound.kind`
//+kubebuilder:printcolumn:name="BoundName",type=string,JSONPath=`.status.bound.kind`
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// IPAMRange is the Schema for the ipamranges API
type IPAMRange struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPAMRangeSpec   `json:"spec,omitempty"`
	Status IPAMRangeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IPAMRangeList contains a list of IPAMRange
type IPAMRangeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPAMRange `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPAMRange{}, &IPAMRangeList{})
}
