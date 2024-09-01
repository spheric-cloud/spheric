// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("NetworkProtectionReconciler", func() {
	ns := SetupNamespace(k8sClient)

	var (
		network *corev1alpha1.Network
	)

	BeforeEach(func(ctx SpecContext) {
		By("creating a network")
		network = &corev1alpha1.Network{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "my-network-",
			},
		}
		Expect(k8sClient.Create(ctx, network)).To(Succeed())
	})

	It("should add and remove a finalizer for a network in use/not used by a network interface", func(ctx SpecContext) {
		By("creating a subnet referencing this network")
		subnet := &corev1alpha1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "subnet-",
			},
			Spec: corev1alpha1.SubnetSpec{
				NetworkRef: corev1alpha1.LocalObjRef(network.Name),
				CIDRs:      []string{"10.0.0.0/24"},
			},
		}
		Expect(k8sClient.Create(ctx, subnet)).To(Succeed())

		By("ensuring that the network finalizer has been set")
		networkKey := types.NamespacedName{
			Namespace: ns.Name,
			Name:      network.Name,
		}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(network.GetFinalizers()).To(ContainElement(corev1alpha1.FinalizerNetwork))
		}).Should(Succeed())

		By("deleting the subnet")
		Expect(k8sClient.Delete(ctx, subnet)).To(Succeed())

		By("deleting the network")
		Expect(k8sClient.Delete(ctx, network)).To(Succeed())

		By("ensuring that the network has been deleted")
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).To(Succeed())
		}).Should(Succeed())
	})

	It("should remove a finalizer for a network in deletion state once the reference network interface is deleted", func(ctx SpecContext) {
		By("creating a subnet referencing this network")
		subnet := &corev1alpha1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "subnet-",
			},
			Spec: corev1alpha1.SubnetSpec{
				NetworkRef: corev1alpha1.LocalObjRef(network.Name),
				CIDRs:      []string{"10.0.0.0/24"},
			},
		}
		Expect(k8sClient.Create(ctx, subnet)).To(Succeed())

		By("ensuring that the network finalizer has been set")
		networkKey := types.NamespacedName{
			Namespace: ns.Name,
			Name:      network.Name,
		}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(network.GetFinalizers()).To(ContainElement(corev1alpha1.FinalizerNetwork))
		}).Should(Succeed())

		By("deleting the network")
		Expect(k8sClient.Delete(ctx, network)).To(Succeed())

		By("ensuring that the network has a deletion timestamp set and the finalizer still present")
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(network.DeletionTimestamp.IsZero()).To(BeFalse())
			g.Expect(network.GetFinalizers()).To(ContainElement(corev1alpha1.FinalizerNetwork))
		}).Should(Succeed())

		By("deleting the subnet")
		Expect(k8sClient.Delete(ctx, subnet)).To(Succeed())

		By("ensuring that the network has been deleted")
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).To(Succeed())
		}).Should(Succeed())
	})

	It("should keep a finalizer if one of two network interfaces is removed", func(ctx SpecContext) {
		By("creating the first subnet referencing this network")
		subnet := &corev1alpha1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "subnet-",
			},
			Spec: corev1alpha1.SubnetSpec{
				NetworkRef: corev1alpha1.LocalObjRef(network.Name),
				CIDRs:      []string{"10.0.0.1/24"},
			},
		}
		Expect(k8sClient.Create(ctx, subnet)).To(Succeed())

		By("creating a second subnet referencing this network")
		subnet2 := &corev1alpha1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "subnet-",
			},
			Spec: corev1alpha1.SubnetSpec{
				NetworkRef: corev1alpha1.LocalObjRef(network.Name),
				CIDRs:      []string{"11.0.0.1/24"},
			},
		}
		Expect(k8sClient.Create(ctx, subnet2)).To(Succeed())

		By("ensuring that the network finalizer has been set")
		networkKey := types.NamespacedName{
			Namespace: ns.Name,
			Name:      network.Name,
		}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(network.GetFinalizers()).To(ContainElement(corev1alpha1.FinalizerNetwork))
		}).Should(Succeed())

		By("deleting the first subnet")
		Expect(k8sClient.Delete(ctx, subnet)).To(Succeed())

		By("deleting the network")
		Expect(k8sClient.Delete(ctx, network)).To(Succeed())

		By("ensuring that the network has a deletion timestamp set and finalizer still present")
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).NotTo(HaveOccurred())
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(network.DeletionTimestamp.IsZero()).To(BeFalse())
			g.Expect(network.GetFinalizers()).To(ContainElement(corev1alpha1.FinalizerNetwork))
		}).Should(Succeed())
	})

	It("should allow deletion of an unused network", func(ctx SpecContext) {
		By("deleting the network")
		Expect(k8sClient.Delete(ctx, network)).To(Succeed())

		By("ensuring that the network is not found")
		networkKey := types.NamespacedName{
			Namespace: ns.Name,
			Name:      network.Name,
		}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(ctx, networkKey, network)
			Expect(client.IgnoreNotFound(err)).To(Succeed())
		}).Should(Succeed())
	})
})
