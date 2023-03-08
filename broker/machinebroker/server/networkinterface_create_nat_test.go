// Copyright 2022 OnMetal authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server_test

import (
	commonv1alpha1 "github.com/onmetal/onmetal-api/api/common/v1alpha1"
	machinebrokerv1alpha1 "github.com/onmetal/onmetal-api/broker/machinebroker/api/v1alpha1"
	ori "github.com/onmetal/onmetal-api/ori/apis/machine/v1alpha1"
	orimeta "github.com/onmetal/onmetal-api/ori/apis/meta/v1alpha1"
	. "github.com/onmetal/onmetal-api/utils/testing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NetworkInterfaceCreateNAT", func() {
	ctx := SetupContext()
	_, srv := SetupTest(ctx)

	It("should correctly create nat for a network interface", func() {
		By("creating a network interface")
		res, err := srv.CreateNetworkInterface(ctx, &ori.CreateNetworkInterfaceRequest{
			NetworkInterface: &ori.NetworkInterface{
				Metadata: &orimeta.ObjectMetadata{},
				Spec: &ori.NetworkInterfaceSpec{
					Network: &ori.NetworkSpec{Handle: "foo"},
					Ips:     []string{"192.168.178.1"},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("inspecting the created network interface")
		networkInterface := res.NetworkInterface
		Expect(networkInterface.Spec.Nats).To(BeEmpty())

		By("creating a nat gateway target for the network interface")
		_, err = srv.CreateNetworkInterfaceNAT(ctx, &ori.CreateNetworkInterfaceNATRequest{
			NetworkInterfaceId: networkInterface.Metadata.Id,
			Nat: &ori.NATSpec{
				Ip:      "10.0.0.1",
				Port:    80,
				EndPort: 8080,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("listing the nat gateways for the network interface")
		nats, err := srv.NATGateways().ListByDependent(ctx, networkInterface.Metadata.Id)
		Expect(err).NotTo(HaveOccurred())

		Expect(nats).To(ConsistOf(machinebrokerv1alpha1.NATGateway{
			NetworkHandle: "foo",
			IP:            commonv1alpha1.MustParseIP("10.0.0.1"),
			Destinations: []machinebrokerv1alpha1.NATGatewayDestination{
				{
					ID:      networkInterface.Metadata.Id,
					Port:    80,
					EndPort: 8080,
				},
			},
		}))
	})
})