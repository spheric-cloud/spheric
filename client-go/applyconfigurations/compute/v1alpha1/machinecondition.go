/*
 * Copyright (c) 2022 by the OnMetal authors.
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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/onmetal/onmetal-api/api/compute/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineConditionApplyConfiguration represents an declarative configuration of the MachineCondition type for use
// with apply.
type MachineConditionApplyConfiguration struct {
	Type               *v1alpha1.MachineConditionType `json:"type,omitempty"`
	Status             *v1.ConditionStatus            `json:"status,omitempty"`
	Reason             *string                        `json:"reason,omitempty"`
	Message            *string                        `json:"message,omitempty"`
	ObservedGeneration *int64                         `json:"observedGeneration,omitempty"`
	LastTransitionTime *metav1.Time                   `json:"lastTransitionTime,omitempty"`
}

// MachineConditionApplyConfiguration constructs an declarative configuration of the MachineCondition type for use with
// apply.
func MachineCondition() *MachineConditionApplyConfiguration {
	return &MachineConditionApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *MachineConditionApplyConfiguration) WithType(value v1alpha1.MachineConditionType) *MachineConditionApplyConfiguration {
	b.Type = &value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *MachineConditionApplyConfiguration) WithStatus(value v1.ConditionStatus) *MachineConditionApplyConfiguration {
	b.Status = &value
	return b
}

// WithReason sets the Reason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reason field is set to the value of the last call.
func (b *MachineConditionApplyConfiguration) WithReason(value string) *MachineConditionApplyConfiguration {
	b.Reason = &value
	return b
}

// WithMessage sets the Message field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Message field is set to the value of the last call.
func (b *MachineConditionApplyConfiguration) WithMessage(value string) *MachineConditionApplyConfiguration {
	b.Message = &value
	return b
}

// WithObservedGeneration sets the ObservedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObservedGeneration field is set to the value of the last call.
func (b *MachineConditionApplyConfiguration) WithObservedGeneration(value int64) *MachineConditionApplyConfiguration {
	b.ObservedGeneration = &value
	return b
}

// WithLastTransitionTime sets the LastTransitionTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastTransitionTime field is set to the value of the last call.
func (b *MachineConditionApplyConfiguration) WithLastTransitionTime(value metav1.Time) *MachineConditionApplyConfiguration {
	b.LastTransitionTime = &value
	return b
}
