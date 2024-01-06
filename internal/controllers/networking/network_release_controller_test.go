// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0
package networking

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("NetworkReleaseReconciler", func() {
	ns := SetupNamespace(&k8sClient)

	It("should release network interfaces whose owner is gone", func(ctx SpecContext) {
		By("creating a network having a peering claim that does not exist")
		network := &networkingv1alpha1.Network{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "network-",
			},
			Spec: networkingv1alpha1.NetworkSpec{
				PeeringClaimRefs: []networkingv1alpha1.NetworkPeeringClaimRef{
					{
						Name: "should-not-exist",
						UID:  uuid.NewUUID(),
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, network)).To(Succeed())

		By("waiting for the network to have the peering claim released")
		Eventually(Object(network)).Should(HaveField("Spec.PeeringClaimRefs", BeEmpty()))
	})
})
