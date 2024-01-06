// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0
package networking

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	commonv1alpha1 "spheric.cloud/spheric/api/common/v1alpha1"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("VirtualIPReleaseReconciler", func() {
	ns := SetupNamespace(&k8sClient)

	It("should release virtual IPs whose owner is gone", func(ctx SpecContext) {
		By("creating a virtual IP referencing an owner that does not exist")
		nic := &networkingv1alpha1.VirtualIP{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "nic-",
			},
			Spec: networkingv1alpha1.VirtualIPSpec{
				Type:     networkingv1alpha1.VirtualIPTypePublic,
				IPFamily: corev1.IPv4Protocol,
				TargetRef: &commonv1alpha1.LocalUIDReference{
					Name: "should-not-exist",
					UID:  uuid.NewUUID(),
				},
			},
		}
		Expect(k8sClient.Create(ctx, nic)).To(Succeed())

		By("waiting for the virtual IP to be released")
		Eventually(Object(nic)).Should(HaveField("Spec.TargetRef", BeNil()))
	})
})
