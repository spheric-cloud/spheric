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
	commonv1alpha1 "github.com/onmetal/onmetal-api/apis/common/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// VolumeGK is a helper to easily access the GroupKind information of an Volume
var VolumeGK = schema.GroupKind{
	Group: SchemeGroupVersion.Group,
	Kind:  "Volume",
}

// VolumeSpec defines the desired state of Volume
type VolumeSpec struct {
	// VolumeClassRef is the VolumeClass of a volume
	VolumeClassRef corev1.LocalObjectReference `json:"volumeClassRef"`
	// VolumePoolSelector selects a suitable VolumePoolRef by the given labels.
	VolumePoolSelector map[string]string `json:"volumePoolSelector,omitempty"`
	// VolumePoolRef indicates which VolumePool to use for a volume.
	// If unset, the scheduler will figure out a suitable VolumePoolRef.
	VolumePoolRef corev1.LocalObjectReference `json:"volumePoolRef,omitempty"`
	// ClaimRef is the reference to the VolumeClaim used by the Volume.
	ClaimRef ClaimReference `json:"claimRef,omitempty"`
	// Resources is a description of the volume's resources and capacity.
	Resources corev1.ResourceList `json:"resources,omitempty"`
	// Tolerations define tolerations the Volume has. Only any VolumePool whose taints
	// covered by Tolerations will be considered to host the Volume.
	Tolerations []commonv1alpha1.Toleration `json:"tolerations,omitempty"`
}

// VolumeAccess represents information on how to access a volume.
type VolumeAccess struct {
	// SecretRef references the Secret containing the access credentials to consume a Volume.
	SecretRef corev1.LocalObjectReference `json:"secretRef,omitempty"`
	// Driver is the name of the drive to use for this volume. Required.
	Driver string `json:"driver"`
	// VolumeAttributes are attributes of the volume to use.
	VolumeAttributes map[string]string `json:"volumeAttributes,omitempty"`
}

// ClaimReference points to a referenced VolumeClaim.
type ClaimReference struct {
	// Name is the name of the referenced VolumeClaim.
	Name string `json:"name"`
	// UID is the UID of the referenced VolumeClaim.
	UID types.UID `json:"uid"`
}

// VolumeStatus defines the observed state of Volume
type VolumeStatus struct {
	// State represents the infrastructure state of a Volume.
	State VolumeState `json:"state,omitempty"`
	// Conditions represents different status aspects of a Volume.
	Conditions []VolumeCondition `json:"conditions,omitempty"`
	// Access specifies how to access a Volume.
	// This is set by the volume provider when the volume is provisioned.
	Access *VolumeAccess `json:"access,omitempty"`
}

const (
	// VolumeBoundReasonUnbound is used for any Volume that is not bound.
	VolumeBoundReasonUnbound = "Unbound"
	// VolumeBoundReasonPending is used for any Volume that is not available.
	VolumeBoundReasonPending = "Pending"
	// VolumeBoundReasonBound is used for any Volume that is bound.
	VolumeBoundReasonBound = "Bound"
)

// VolumePhase is the binding phase of a volume.
type VolumePhase string

const (
	// VolumePhaseUnknown is used for any Volume for which it is unknown whether it can be used for binding.
	VolumePhaseUnknown VolumePhase = "Unknown"
	// VolumePhaseUnbound is used for any Volume that not bound.
	VolumePhaseUnbound VolumePhase = "Unbound"
	// VolumePhasePending is used for any Volume that is currently awaiting binding.
	VolumePhasePending VolumePhase = "Pending"
	// VolumePhaseBound is used for any Volume that is properly bound.
	VolumePhaseBound VolumePhase = "Bound"
)

func FindVolumeCondition(conditions []VolumeCondition, conditionType VolumeConditionType) (VolumeCondition, int) {
	for i, condition := range conditions {
		if condition.Type == conditionType {
			return condition, i
		}
	}
	return VolumeCondition{}, -1
}

func VolumePhaseFromBoundStatusAndReason(status corev1.ConditionStatus, reason string) VolumePhase {
	switch {
	case status == corev1.ConditionFalse && reason == VolumeBoundReasonPending:
		return VolumePhasePending
	case status == corev1.ConditionFalse && reason == VolumeBoundReasonUnbound:
		return VolumePhaseUnbound
	case status == corev1.ConditionTrue:
		return VolumePhaseBound
	default:
		return VolumePhaseUnknown
	}
}

func GetVolumePhaseAndCondition(volume *Volume) (VolumePhase, VolumeCondition, int) {
	cond, idx := FindVolumeCondition(volume.Status.Conditions, VolumeBound)
	phase := VolumePhaseFromBoundStatusAndReason(cond.Status, cond.Reason)
	return phase, cond, idx
}

func GetVolumePhase(volume *Volume) VolumePhase {
	phase, _, _ := GetVolumePhaseAndCondition(volume)
	return phase
}

// VolumeState is a possible state a volume can be in.
type VolumeState string

const (
	// VolumeStateUnknown reports whether a Volume is in an unknown state.
	VolumeStateUnknown VolumeState = "Unknown"
	// VolumeStatePending reports whether a Volume is about to be ready.
	VolumeStatePending VolumeState = "Pending"
	// VolumeStateAvailable reports whether a Volume is available to be used.
	VolumeStateAvailable VolumeState = "Available"
	// VolumeStateError reports that a Volume is in an error state.
	VolumeStateError VolumeState = "Error"
)

// VolumeConditionType is a type a VolumeCondition can have.
type VolumeConditionType string

const (
	// VolumeSynced represents the condition of a volume being synced with its backing resources
	VolumeSynced VolumeConditionType = "Synced"

	// VolumeBound represents the binding state of a Volume.
	VolumeBound VolumeConditionType = "Bound"
)

// VolumeCondition is one of the conditions of a volume.
type VolumeCondition struct {
	// Type is the type of the condition.
	Type VolumeConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// Reason is a machine-readable indication of why the condition is in a certain state.
	Reason string `json:"reason"`
	// Message is a human-readable explanation of why the condition has a certain reason / state.
	Message string `json:"message"`
	// ObservedGeneration represents the .metadata.generation that the condition was set based upon.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// LastTransitionTime is the last time the status of a condition has transitioned from one state to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// Volume is the Schema for the volumes API
type Volume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeSpec   `json:"spec,omitempty"`
	Status VolumeStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeList contains a list of Volume
type VolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Volume `json:"items"`
}
