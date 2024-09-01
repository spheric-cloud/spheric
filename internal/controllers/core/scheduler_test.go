// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package core_test

import (
	"fmt"
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	. "spheric.cloud/spheric/utils/testing"

	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
)

var _ = Describe("InstanceScheduler", func() {
	ns := SetupNamespace(k8sClient)
	instanceType := SetupInstanceType()

	It("should schedule instances on fleets", func(ctx SpecContext) {
		By("creating a fleet")
		fleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, fleet)).To(Succeed(), "failed to create fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleet, func() {
			fleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("creating a instance w/ the requested instance type")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image:           "my-image",
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create instance")

		By("waiting for the instance to be scheduled onto the fleet")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", Equal(corev1alpha1.NewLocalObjRef(fleet.Name))),
			HaveField("Status.State", Equal(corev1alpha1.InstanceStatePending)),
		))
	})

	It("should schedule schedule instances onto fleets if the fleet becomes available later than the instance", func(ctx SpecContext) {
		By("creating a instance w/ the requested instance type")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image:           "my-image",
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create instance")

		By("waiting for the instance to indicate it is pending")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", BeNil()),
			HaveField("Status.State", Equal(corev1alpha1.InstanceStatePending)),
		))

		By("creating a fleet")
		fleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, fleet)).To(Succeed(), "failed to create fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleet, func() {
			fleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("waiting for the instance to be scheduled onto the fleet")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", Equal(corev1alpha1.NewLocalObjRef(fleet.Name))),
			HaveField("Status.State", Equal(corev1alpha1.InstanceStatePending)),
		))
	})

	It("should schedule onto fleets with matching labels", func(ctx SpecContext) {
		By("creating a fleet w/o matching labels")
		fleetNoMatchingLabels := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, fleetNoMatchingLabels)).To(Succeed(), "failed to create fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleetNoMatchingLabels, func() {
			fleetNoMatchingLabels.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("creating a fleet w/ matching labels")
		fleetMatchingLabels := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
				Labels: map[string]string{
					"foo": "bar",
				},
			},
		}
		Expect(k8sClient.Create(ctx, fleetMatchingLabels)).To(Succeed(), "failed to create fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleetMatchingLabels, func() {
			fleetMatchingLabels.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("creating a instance w/ the requested instance type")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image: "my-image",
				FleetSelector: map[string]string{
					"foo": "bar",
				},
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create instance")

		By("waiting for the instance to be scheduled onto the fleet")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", Equal(corev1alpha1.NewLocalObjRef(fleetMatchingLabels.Name))),
			HaveField("Status.State", Equal(corev1alpha1.InstanceStatePending)),
		))
	})

	It("should schedule a instance with corresponding tolerations onto a fleet with taints", func(ctx SpecContext) {
		By("creating a fleet w/ taints")
		taintedFleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
			Spec: corev1alpha1.FleetSpec{
				Taints: []corev1alpha1.Taint{
					{
						Key:    "key",
						Value:  "value",
						Effect: corev1alpha1.TaintEffectNoSchedule,
					},
					{
						Key:    "key1",
						Effect: corev1alpha1.TaintEffectNoSchedule,
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, taintedFleet)).To(Succeed(), "failed to create the fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(taintedFleet, func() {
			taintedFleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("creating a instance")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image:           "my-image",
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create the instance")

		By("observing the instance isn't scheduled onto the fleet")
		Consistently(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", BeNil()),
		))

		By("patching the instance to contain only one of the corresponding tolerations")
		instanceBase := instance.DeepCopy()
		instance.Spec.Tolerations = append(instance.Spec.Tolerations, corev1alpha1.Toleration{
			Key:      "key",
			Value:    "value",
			Effect:   corev1alpha1.TaintEffectNoSchedule,
			Operator: corev1alpha1.TolerationOpEqual,
		})
		Expect(k8sClient.Patch(ctx, instance, client.MergeFrom(instanceBase))).To(Succeed(), "failed to patch the instance's spec")

		By("observing the instance isn't scheduled onto the fleet")
		Consistently(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", BeNil()),
		))

		By("patching the instance to contain all of the corresponding tolerations")
		instanceBase = instance.DeepCopy()
		instance.Spec.Tolerations = append(instance.Spec.Tolerations, corev1alpha1.Toleration{
			Key:      "key1",
			Effect:   corev1alpha1.TaintEffectNoSchedule,
			Operator: corev1alpha1.TolerationOpExists,
		})
		Expect(k8sClient.Patch(ctx, instance, client.MergeFrom(instanceBase))).To(Succeed(), "failed to patch the instance's spec")

		By("observing the instance is scheduled onto the fleet")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", Equal(corev1alpha1.NewLocalObjRef(taintedFleet.Name))),
		))
	})

	It("should schedule instance on fleet with most allocatable resources", func(ctx SpecContext) {
		By("creating a fleet")
		fleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, fleet)).To(Succeed(), "failed to create fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleet, func() {
			fleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("creating a second fleet")
		secondFleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "second-test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, secondFleet)).To(Succeed(), "failed to create the second fleet")

		By("creating a second instance type")
		secondInstanceType := &corev1alpha1.InstanceType{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "second-instance-type-",
			},
			Capabilities: corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:    resource.MustParse("1"),
				corev1alpha1.ResourceMemory: resource.MustParse("1Gi"),
			},
		}
		Expect(k8sClient.Create(ctx, secondInstanceType)).To(Succeed(), "failed to create second instance type")

		By("patching the second fleet status to contain a both instance typees")
		Eventually(UpdateStatus(secondFleet, func() {
			secondFleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name):       resource.MustParse("5"),
				corev1alpha1.ResourceInstanceType(secondInstanceType.Name): resource.MustParse("100"),
			}
		})).Should(Succeed())

		By("creating a instance")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image:           "my-image",
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create the instance")

		By("checking that the instance is scheduled onto the fleet")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef.Name", Equal(fleet.Name)),
		))
	})

	It("should schedule instances evenly on fleets", func(ctx SpecContext) {
		By("creating a fleet")
		fleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, fleet)).To(Succeed(), "failed to create fleet")

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleet, func() {
			fleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("50"),
			}
		})).Should(Succeed())

		By("creating a second fleet")
		secondFleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "second-test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, secondFleet)).To(Succeed(), "failed to create the second fleet")

		By("patching the second fleet status to contain a both instance typees")
		Eventually(UpdateStatus(secondFleet, func() {
			secondFleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("50"),
			}
		})).Should(Succeed())

		By("creating instances")
		var instances []*corev1alpha1.Instance
		for i := 0; i < 50; i++ {
			instance := &corev1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:    ns.Name,
					GenerateName: fmt.Sprintf("test-instance-%d-", i),
				},
				Spec: corev1alpha1.InstanceSpec{
					Image:           "my-image",
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create the instance")
			instances = append(instances, instance)
		}

		By("checking that every instance is scheduled onto a fleet")
		var numFleets1, numFleets2 int64
		for i := 0; i < 50; i++ {
			Eventually(Object(instances[i])).Should(SatisfyAll(
				HaveField("Spec.FleetRef", Not(BeNil())),
			))

			switch instances[i].Spec.FleetRef.Name {
			case fleet.Name:
				numFleets1++
			case secondFleet.Name:
				numFleets2++
			}
		}

		By("checking that instance are roughly distributed")
		Expect(math.Abs(float64(numFleets1 - numFleets2))).To(BeNumerically("<", 5))
	})

	It("should schedule a instances once the capacity is sufficient", func(ctx SpecContext) {
		By("creating a fleet")
		fleet := &corev1alpha1.Fleet{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-fleet-",
			},
		}
		Expect(k8sClient.Create(ctx, fleet)).To(Succeed(), "failed to create fleet")
		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleet, func() {})).Should(Succeed())

		By("creating a instance")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				Image:           "my-image",
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed(), "failed to create the instance")

		By("checking that the instance is scheduled onto the fleet")
		Consistently(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", BeNil()),
		))

		By("patching the fleet status to contain a instance type")
		Eventually(UpdateStatus(fleet, func() {
			fleet.Status.Allocatable = corev1alpha1.ResourceList{
				corev1alpha1.ResourceInstanceType(instanceType.Name): resource.MustParse("10"),
			}
		})).Should(Succeed())

		By("checking that the instance is scheduled onto the fleet")
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Spec.FleetRef", Equal(corev1alpha1.NewLocalObjRef(fleet.Name))),
		))
	})
})
