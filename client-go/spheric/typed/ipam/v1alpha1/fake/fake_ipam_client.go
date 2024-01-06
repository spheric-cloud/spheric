// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	v1alpha1 "spheric.cloud/spheric/client-go/spheric/typed/ipam/v1alpha1"
)

type FakeIpamV1alpha1 struct {
	*testing.Fake
}

func (c *FakeIpamV1alpha1) Prefixes(namespace string) v1alpha1.PrefixInterface {
	return &FakePrefixes{c, namespace}
}

func (c *FakeIpamV1alpha1) PrefixAllocations(namespace string) v1alpha1.PrefixAllocationInterface {
	return &FakePrefixAllocations{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeIpamV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}