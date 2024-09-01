// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/utils/generic"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("FleetController", func() {
	instanceType := SetupInstanceType()

	It("should maintain the pool status", func(ctx SpecContext) {
		srv.SetStatus(&sri.RuntimeResources{
			CpuCount:    32,
			MemoryBytes: uint64(generic.Pointer(resource.MustParse("20Ti")).Value()),
			InstanceQuantities: map[string]int64{
				instanceType.Name: 10,
			},
		}, &sri.RuntimeResources{
			CpuCount:    24,
			MemoryBytes: uint64(generic.Pointer(resource.MustParse("10Ti")).Value()),
			InstanceQuantities: map[string]int64{
				instanceType.Name: 5,
			},
		})

		By("checking if the capacity is correct")
		Eventually(Object(&corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{Name: fleetName},
		})).Should(SatisfyAll(
			HaveField("Status.Capacity", EqualResources(corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:                             resource.MustParse("32"),
				corev1alpha1.ResourceMemory:                          resource.MustParse("20Ti"),
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			})),
			HaveField("Status.Allocatable", EqualResources(corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:                             resource.MustParse("24"),
				corev1alpha1.ResourceMemory:                          resource.MustParse("10Ti"),
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("5"),
			})),
		))
	})
})
