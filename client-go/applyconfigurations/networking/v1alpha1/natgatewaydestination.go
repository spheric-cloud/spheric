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
	types "k8s.io/apimachinery/pkg/types"
)

// NATGatewayDestinationApplyConfiguration represents an declarative configuration of the NATGatewayDestination type for use
// with apply.
type NATGatewayDestinationApplyConfiguration struct {
	Name *string                                     `json:"name,omitempty"`
	UID  *types.UID                                  `json:"uid,omitempty"`
	IPs  []NATGatewayDestinationIPApplyConfiguration `json:"ips,omitempty"`
}

// NATGatewayDestinationApplyConfiguration constructs an declarative configuration of the NATGatewayDestination type for use with
// apply.
func NATGatewayDestination() *NATGatewayDestinationApplyConfiguration {
	return &NATGatewayDestinationApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *NATGatewayDestinationApplyConfiguration) WithName(value string) *NATGatewayDestinationApplyConfiguration {
	b.Name = &value
	return b
}

// WithUID sets the UID field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UID field is set to the value of the last call.
func (b *NATGatewayDestinationApplyConfiguration) WithUID(value types.UID) *NATGatewayDestinationApplyConfiguration {
	b.UID = &value
	return b
}

// WithIPs adds the given value to the IPs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the IPs field.
func (b *NATGatewayDestinationApplyConfiguration) WithIPs(values ...*NATGatewayDestinationIPApplyConfiguration) *NATGatewayDestinationApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithIPs")
		}
		b.IPs = append(b.IPs, *values[i])
	}
	return b
}
