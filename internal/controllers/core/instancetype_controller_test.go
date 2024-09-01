// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("instancetype controller", func() {
	ns := SetupNamespace(k8sClient)

	It("removes the finalizer from instancetype only if there's no instance still using the instancetype", func(ctx SpecContext) {
		By("creating the instancetype consumed by the instance")
		instanceType := &corev1alpha1.InstanceType{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "instancetype-",
			},
			Capabilities: corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:    resource.MustParse("300m"),
				corev1alpha1.ResourceMemory: resource.MustParse("1Gi"),
			},
		}
		Expect(k8sClient.Create(ctx, instanceType)).Should(Succeed())
		DeferCleanup(DeleteIgnoreNotFound(k8sClient, instanceType))

		By("creating the instance")
		m := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image:           "my-image",
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, m)).Should(Succeed())

		By("checking the finalizer was added")
		instanceTypeKey := client.ObjectKeyFromObject(instanceType)
		Eventually(func(g Gomega) []string {
			err := k8sClient.Get(ctx, instanceTypeKey, instanceType)
			Expect(client.IgnoreNotFound(err)).To(Succeed(), "errors other than `not found` are not expected")
			g.Expect(err).NotTo(HaveOccurred())
			return instanceType.Finalizers
		}).Should(ContainElement(corev1alpha1.InstanceTypeFinalizer))

		By("checking the instancetype and its finalizer consistently exist upon deletion ")
		Expect(k8sClient.Delete(ctx, instanceType)).Should(Succeed())

		Consistently(func(g Gomega) []string {
			err := k8sClient.Get(ctx, instanceTypeKey, instanceType)
			Expect(client.IgnoreNotFound(err)).To(Succeed(), "errors other than `not found` are not expected")
			g.Expect(err).NotTo(HaveOccurred())
			return instanceType.Finalizers
		}).Should(ContainElement(corev1alpha1.InstanceTypeFinalizer))

		By("checking the instancetype is eventually gone after the deletion of the instance")
		Expect(k8sClient.Delete(ctx, m)).Should(Succeed())
		Eventually(func() bool {
			err := k8sClient.Get(ctx, instanceTypeKey, instanceType)
			Expect(client.IgnoreNotFound(err)).To(Succeed(), "errors other than `not found` are not expected")
			return apierrors.IsNotFound(err)
		}).Should(BeTrue(), "the error should be `not found`")
	})
})
