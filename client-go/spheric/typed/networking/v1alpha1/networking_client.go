// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"net/http"

	rest "k8s.io/client-go/rest"
	v1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	"spheric.cloud/spheric/client-go/spheric/scheme"
)

type NetworkingV1alpha1Interface interface {
	RESTClient() rest.Interface
	LoadBalancersGetter
	LoadBalancerRoutingsGetter
	NATGatewaysGetter
	NetworksGetter
	NetworkInterfacesGetter
	NetworkPoliciesGetter
	VirtualIPsGetter
}

// NetworkingV1alpha1Client is used to interact with features provided by the networking.spheric.cloud group.
type NetworkingV1alpha1Client struct {
	restClient rest.Interface
}

func (c *NetworkingV1alpha1Client) LoadBalancers(namespace string) LoadBalancerInterface {
	return newLoadBalancers(c, namespace)
}

func (c *NetworkingV1alpha1Client) LoadBalancerRoutings(namespace string) LoadBalancerRoutingInterface {
	return newLoadBalancerRoutings(c, namespace)
}

func (c *NetworkingV1alpha1Client) NATGateways(namespace string) NATGatewayInterface {
	return newNATGateways(c, namespace)
}

func (c *NetworkingV1alpha1Client) Networks(namespace string) NetworkInterface {
	return newNetworks(c, namespace)
}

func (c *NetworkingV1alpha1Client) NetworkInterfaces(namespace string) NetworkInterfaceInterface {
	return newNetworkInterfaces(c, namespace)
}

func (c *NetworkingV1alpha1Client) NetworkPolicies(namespace string) NetworkPolicyInterface {
	return newNetworkPolicies(c, namespace)
}

func (c *NetworkingV1alpha1Client) VirtualIPs(namespace string) VirtualIPInterface {
	return newVirtualIPs(c, namespace)
}

// NewForConfig creates a new NetworkingV1alpha1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*NetworkingV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new NetworkingV1alpha1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*NetworkingV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &NetworkingV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new NetworkingV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *NetworkingV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new NetworkingV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *NetworkingV1alpha1Client {
	return &NetworkingV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *NetworkingV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}