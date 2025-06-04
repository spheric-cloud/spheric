// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0
package core_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("DiskReleaseReconciler", func() {
	ns := SetupNamespace(k8sClient)

	It("should release disks whose owner is gone", func(ctx SpecContext) {
		By("creating a disk referencing an owner that does not exist")
		disk := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "disk-",
			},
			Spec: corev1alpha1.DiskSpec{
				InstanceRef: &corev1alpha1.LocalUIDReference{
					Name: "should-not-exist",
					UID:  uuid.NewUUID(),
				},
			},
		}
		Expect(k8sClient.Create(ctx, disk)).To(Succeed())

		By("waiting for the disk to be released")
		Eventually(Object(disk)).Should(HaveField("Spec.InstanceRef", BeNil()))
	})
})
