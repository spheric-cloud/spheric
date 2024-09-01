// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/server/egressselector"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
)

// ConnectionInfo provides the information needed to connect to a spherelet.
type ConnectionInfo struct {
	Scheme                         string
	Hostname                       string
	Port                           string
	Transport                      http.RoundTripper
	InsecureSkipTLSVerifyTransport http.RoundTripper
}

type ConnectionInfoGetter interface {
	GetConnectionInfo(ctx context.Context, fleetName string) (*ConnectionInfo, error)
}

// FleetGetter defines an interface for looking up a node by name
type FleetGetter interface {
	Get(ctx context.Context, name string, options metav1.GetOptions) (*corev1alpha1.Fleet, error)
}

// FleetGetterFunc allows implementing FleetGetter with a function
type FleetGetterFunc func(ctx context.Context, name string, options metav1.GetOptions) (*corev1alpha1.Fleet, error)

// Get fetches information via FleetGetterFunc.
func (f FleetGetterFunc) Get(ctx context.Context, name string, options metav1.GetOptions) (*corev1alpha1.Fleet, error) {
	return f(ctx, name, options)
}

// FleetConnectionInfoGetter obtains connection info from the status of a Fleet API object
type FleetConnectionInfoGetter struct {
	// fleets is used to look up Fleet objects
	fleets FleetGetter
	// scheme is the scheme to use to connect to all spherelets.
	scheme string
	// defaultPort is the port to use if no spherelet endpoint port is recorded in the node status
	defaultPort int
	// transport is the transport to use to send a request to all spherelets.
	transport http.RoundTripper
	// insecureSkipTLSVerifyTransport is the transport to use if the kube-apiserver wants to skip verifying the TLS certificate of the spherelet
	insecureSkipTLSVerifyTransport http.RoundTripper
	// preferredAddressTypes specifies the preferred order to use to find a node address
	preferredAddressTypes []corev1alpha1.FleetAddressType
}

type SphereletClientConfig struct {
	// Port specifies the default port - used if no information about Spherelet port can be found in Node.NodeStatus.DaemonEndpoints.
	Port uint

	// ReadOnlyPort specifies the Port for ReadOnly communications.
	ReadOnlyPort uint

	// PreferredAddressTypes - used to select an address from Node.NodeStatus.Addresses
	PreferredAddressTypes []string

	// TLSClientConfig contains settings to enable transport layer security
	rest.TLSClientConfig

	// Server requires Bearer authentication
	BearerToken string `datapolicy:"token"`

	// HTTPTimeout is used by the client to timeout http requests to Spherelet.
	HTTPTimeout time.Duration

	// Dial is a custom dialer used for the client
	Dial utilnet.DialFunc

	// Lookup will give us a dialer if the egress selector is configured for it
	Lookup egressselector.Lookup
}

func (c *SphereletClientConfig) transportConfig() *transport.Config {
	cfg := &transport.Config{
		TLS: transport.TLSConfig{
			CAFile:     c.CAFile,
			CAData:     c.CAData,
			CertFile:   c.CertFile,
			CertData:   c.CertData,
			KeyFile:    c.KeyFile,
			KeyData:    c.KeyData,
			NextProtos: c.NextProtos,
		},
		BearerToken: c.BearerToken,
	}
	if !cfg.HasCA() {
		cfg.TLS.Insecure = true
	}
	return cfg
}

// MakeTransport creates a secure RoundTripper for HTTP Transport.
func MakeTransport(config *SphereletClientConfig) (http.RoundTripper, error) {
	return makeTransport(config, false)
}

// MakeInsecureTransport creates an insecure RoundTripper for HTTP Transport.
func MakeInsecureTransport(config *SphereletClientConfig) (http.RoundTripper, error) {
	return makeTransport(config, true)
}

// makeTransport creates a RoundTripper for HTTP Transport.
func makeTransport(config *SphereletClientConfig, insecureSkipTLSVerify bool) (http.RoundTripper, error) {
	// do the insecureSkipTLSVerify on the pre-transport *before* we go get a potentially cached connection.
	// transportConfig always produces a new struct pointer.
	preTLSConfig := config.transportConfig()
	if insecureSkipTLSVerify && preTLSConfig != nil {
		preTLSConfig.TLS.Insecure = true
		preTLSConfig.TLS.CAData = nil
		preTLSConfig.TLS.CAFile = ""
	}

	tlsConfig, err := transport.TLSConfigFor(preTLSConfig)
	if err != nil {
		return nil, err
	}

	rt := http.DefaultTransport
	dialer := config.Dial
	if dialer == nil && config.Lookup != nil {
		// Assuming EgressSelector if SSHTunnel is not turned on.
		// We will not get a dialer if egress selector is disabled.
		networkContext := egressselector.Cluster.AsNetworkContext()
		dialer, err = config.Lookup(networkContext)
		if err != nil {
			return nil, fmt.Errorf("failed to get context dialer for 'cluster': got %v", err)
		}
	}
	if dialer != nil || tlsConfig != nil {
		// If SSH Tunnel is turned on
		rt = utilnet.SetOldTransportDefaults(&http.Transport{
			DialContext:     dialer,
			TLSClientConfig: tlsConfig,
		})
	}

	return transport.HTTPWrappersForConfig(config.transportConfig(), rt)
}

// NoMatchError is a typed implementation of the error interface. It indicates a failure to get a matching Node.
type NoMatchError struct {
	addresses []corev1alpha1.FleetAddress
}

// Error is the implementation of the conventional interface for
// representing an error condition, with the nil value representing no error.
func (e *NoMatchError) Error() string {
	return fmt.Sprintf("no preferred addresses found; known addresses: %v", e.addresses)
}

// GetPreferredFleetAddress returns the address of the provided node, using the provided preference order.
// If none of the preferred address types are found, an error is returned.
func GetPreferredFleetAddress(fleet *corev1alpha1.Fleet, preferredAddressTypes []corev1alpha1.FleetAddressType) (string, error) {
	for _, addressType := range preferredAddressTypes {
		for _, address := range fleet.Status.Addresses {
			if address.Type == addressType {
				return address.Address, nil
			}
		}
	}
	return "", &NoMatchError{addresses: fleet.Status.Addresses}
}

// GetConnectionInfo retrieves connection info from the status of a Node API object.
func (k *FleetConnectionInfoGetter) GetConnectionInfo(ctx context.Context, fleetName string) (*ConnectionInfo, error) {
	fleet, err := k.fleets.Get(ctx, fleetName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Find a spherelet-reported address, using preferred address type
	host, err := GetPreferredFleetAddress(fleet, k.preferredAddressTypes)
	if err != nil {
		return nil, err
	}

	// Use the spherelet-reported port, if present
	port := int(fleet.Status.DaemonEndpoints.SphereletEndpoint.Port)
	if port <= 0 {
		port = k.defaultPort
	}

	return &ConnectionInfo{
		Scheme:                         k.scheme,
		Hostname:                       host,
		Port:                           strconv.Itoa(port),
		Transport:                      k.transport,
		InsecureSkipTLSVerifyTransport: k.insecureSkipTLSVerifyTransport,
	}, nil
}

func NewFleetConnectionInfoGetter(fleets FleetGetter, config SphereletClientConfig) (ConnectionInfoGetter, error) {
	transport, err := MakeTransport(&config)
	if err != nil {
		return nil, err
	}
	insecureSkipTLSVerifyTransport, err := MakeInsecureTransport(&config)
	if err != nil {
		return nil, err
	}

	var types []corev1alpha1.FleetAddressType
	for _, t := range config.PreferredAddressTypes {
		types = append(types, corev1alpha1.FleetAddressType(t))
	}

	return &FleetConnectionInfoGetter{
		fleets:                         fleets,
		scheme:                         "https",
		defaultPort:                    int(config.Port),
		transport:                      transport,
		insecureSkipTLSVerifyTransport: insecureSkipTLSVerifyTransport,

		preferredAddressTypes: types,
	}, nil
}
