// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	"spheric.cloud/spheric/utils/annotations"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("InstanceEphemeralDiskController", func() {
	ns := SetupNamespace(k8sClient)
	instanceClass := SetupInstanceType()

	It("should manage ephemeral disks for a instance", func(ctx SpecContext) {
		By("creating a disk that will be referenced by the instance")
		refDisk := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "ref-disk-",
			},
			Spec: corev1alpha1.DiskSpec{},
		}
		Expect(k8sClient.Create(ctx, refDisk)).To(Succeed())

		By("creating a instance")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
				Disks: []corev1alpha1.AttachedDisk{
					{
						Name: "ref-disk",
						AttachedDiskSource: corev1alpha1.AttachedDiskSource{
							DiskRef: corev1alpha1.NewLocalObjRef(refDisk.Name),
						},
					},
					{
						Name: "ephem-disk",
						AttachedDiskSource: corev1alpha1.AttachedDiskSource{
							Ephemeral: &corev1alpha1.EphemeralDiskSource{
								DiskTemplate: &corev1alpha1.DiskTemplateSpec{
									Spec: corev1alpha1.DiskSpec{},
								},
							},
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("creating an undesired controlled disk")
		undesiredControlledDiskClaimClaim := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "undesired-ctrl-disk-",
			},
			Spec: corev1alpha1.DiskSpec{},
		}
		annotations.SetDefaultEphemeralManagedBy(undesiredControlledDiskClaimClaim)
		_ = ctrl.SetControllerReference(instance, undesiredControlledDiskClaimClaim, k8sClient.Scheme())
		Expect(k8sClient.Create(ctx, undesiredControlledDiskClaimClaim)).To(Succeed())

		By("waiting for the ephemeral disk to exist")
		ephemDiskClaim := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns.Name,
				Name:      corev1alpha1.InstanceEphemeralDiskName(instance.Name, "ephem-disk"),
			},
		}
		Eventually(Get(ephemDiskClaim)).Should(Succeed())

		By("asserting the referenced disk still exists")
		Consistently(Get(refDisk)).Should(Succeed())

		By("waiting for the undesired controlled disk to be gone")
		Eventually(Get(undesiredControlledDiskClaimClaim)).Should(Satisfy(apierrors.IsNotFound))
	})

	It("should not delete externally managed disks for a instance", func(ctx SpecContext) {
		By("creating a instance")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("creating an undesired controlled disk")
		externalDiskClaim := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "external-disk-",
			},
			Spec: corev1alpha1.DiskSpec{},
		}
		_ = ctrl.SetControllerReference(instance, externalDiskClaim, k8sClient.Scheme())
		Expect(k8sClient.Create(ctx, externalDiskClaim)).To(Succeed())

		By("asserting that the external disk claim is not being deleted")
		Consistently(Object(externalDiskClaim)).Should(HaveField("DeletionTimestamp", BeNil()))
	})
})
