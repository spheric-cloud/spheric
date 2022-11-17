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

// ReferencedVolumeStatusApplyConfiguration represents an declarative configuration of the ReferencedVolumeStatus type for use
// with apply.
type ReferencedVolumeStatusApplyConfiguration struct {
	Driver *string `json:"driver,omitempty"`
	Handle *string `json:"handle,omitempty"`
}

// ReferencedVolumeStatusApplyConfiguration constructs an declarative configuration of the ReferencedVolumeStatus type for use with
// apply.
func ReferencedVolumeStatus() *ReferencedVolumeStatusApplyConfiguration {
	return &ReferencedVolumeStatusApplyConfiguration{}
}

// WithDriver sets the Driver field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Driver field is set to the value of the last call.
func (b *ReferencedVolumeStatusApplyConfiguration) WithDriver(value string) *ReferencedVolumeStatusApplyConfiguration {
	b.Driver = &value
	return b
}

// WithHandle sets the Handle field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Handle field is set to the value of the last call.
func (b *ReferencedVolumeStatusApplyConfiguration) WithHandle(value string) *ReferencedVolumeStatusApplyConfiguration {
	b.Handle = &value
	return b
}
