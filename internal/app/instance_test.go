// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package app_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	. "spheric.cloud/spheric/utils/testing"
)

var _ = Describe("Core", func() {
	var (
		ctx           = SetupContext()
		ns            = SetupTest(ctx)
		instanceClass = &corev1alpha1.InstanceType{}
	)

	const (
		fieldOwner = client.FieldOwner("fieldowner.test.spheric.cloud/apiserver")
	)

	BeforeEach(func() {
		*instanceClass = corev1alpha1.InstanceType{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "instance-class-",
			},
			Capabilities: corev1alpha1.ResourceList{
				corev1alpha1.ResourceCPU:    resource.MustParse("1"),
				corev1alpha1.ResourceMemory: resource.MustParse("1Gi"),
			},
		}
		Expect(k8sClient.Create(ctx, instanceClass)).To(Succeed(), "failed to create test instance class")
		DeferCleanup(k8sClient.Delete, ctx, instanceClass)
	})

	Context("Instance", func() {
		It("should correctly apply instances with volumes and default devices", func() {
			By("applying a instance with volumes")
			instance := &corev1alpha1.Instance{
				TypeMeta: metav1.TypeMeta{
					APIVersion: corev1alpha1.SchemeGroupVersion.String(),
					Kind:       "Instance",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "my-instance",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
					Disks: []corev1alpha1.AttachedDisk{
						{
							Name: "foo",
							AttachedDiskSource: corev1alpha1.AttachedDiskSource{
								EmptyDisk: &corev1alpha1.EmptyDiskSource{},
							},
						},
					},
				},
			}
			Expect(k8sClient.Patch(ctx, instance, client.Apply, fieldOwner)).To(Succeed())

			By("inspecting the instance's volumes")
			Expect(instance.Spec.Disks).To(Equal([]corev1alpha1.AttachedDisk{
				{
					Name:   "foo",
					Device: ptr.To("oda"),
					AttachedDiskSource: corev1alpha1.AttachedDiskSource{
						EmptyDisk: &corev1alpha1.EmptyDiskSource{},
					},
				},
			}))

			By("applying a changed instance with a second disk")
			instance = &corev1alpha1.Instance{
				TypeMeta: metav1.TypeMeta{
					APIVersion: corev1alpha1.SchemeGroupVersion.String(),
					Kind:       "Instance",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "my-instance",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
					Disks: []corev1alpha1.AttachedDisk{
						{
							Name: "foo",
							AttachedDiskSource: corev1alpha1.AttachedDiskSource{
								EmptyDisk: &corev1alpha1.EmptyDiskSource{},
							},
						},
						{
							Name: "bar",
							AttachedDiskSource: corev1alpha1.AttachedDiskSource{
								EmptyDisk: &corev1alpha1.EmptyDiskSource{},
							},
						},
					},
				},
			}
			Expect(k8sClient.Patch(ctx, instance, client.Apply, fieldOwner)).To(Succeed())

			By("inspecting the instance's volumes")
			Expect(instance.Spec.Disks).To(Equal([]corev1alpha1.AttachedDisk{
				{
					Name:   "foo",
					Device: ptr.To("oda"),
					AttachedDiskSource: corev1alpha1.AttachedDiskSource{
						EmptyDisk: &corev1alpha1.EmptyDiskSource{},
					},
				},
				{
					Name:   "bar",
					Device: ptr.To("odb"),
					AttachedDiskSource: corev1alpha1.AttachedDiskSource{
						EmptyDisk: &corev1alpha1.EmptyDiskSource{},
					},
				},
			}))

			By("applying a changed instance with the first disk removed")
			instance = &corev1alpha1.Instance{
				TypeMeta: metav1.TypeMeta{
					APIVersion: corev1alpha1.SchemeGroupVersion.String(),
					Kind:       "Instance",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      "my-instance",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
					Disks: []corev1alpha1.AttachedDisk{
						{
							Name: "bar",
							AttachedDiskSource: corev1alpha1.AttachedDiskSource{
								EmptyDisk: &corev1alpha1.EmptyDiskSource{},
							},
						},
					},
				},
			}
			Expect(k8sClient.Patch(ctx, instance, client.Apply, fieldOwner)).To(Succeed())

			By("inspecting the instance's volumes")
			Expect(instance.Spec.Disks).To(Equal([]corev1alpha1.AttachedDisk{
				{
					Name:   "bar",
					Device: ptr.To("odb"),
					AttachedDiskSource: corev1alpha1.AttachedDiskSource{
						EmptyDisk: &corev1alpha1.EmptyDiskSource{},
					},
				},
			}))
		})

		It("should allow listing instances filtering by instance pool name", func() {
			const (
				fleet1 = "instance-pool-1"
				fleet2 = "instance-pool-2"
			)

			By("creating a instance on machine pool 1")
			instance1 := &corev1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:    ns.Name,
					GenerateName: "instance-",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
					FleetRef:        corev1alpha1.NewLocalObjRef(fleet1),
				},
			}
			Expect(k8sClient.Create(ctx, instance1)).To(Succeed())

			By("creating a instance on machine pool 2")
			instance2 := &corev1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:    ns.Name,
					GenerateName: "instance-",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
					FleetRef:        corev1alpha1.NewLocalObjRef(fleet2),
				},
			}
			Expect(k8sClient.Create(ctx, instance2)).To(Succeed())

			By("creating a instance on no machine pool")
			instance3 := &corev1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:    ns.Name,
					GenerateName: "instance-",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
				},
			}
			Expect(k8sClient.Create(ctx, instance3)).To(Succeed())

			By("listing all instances on machine pool 1")
			instancesOnFleet1List := &corev1alpha1.InstanceList{}
			Expect(k8sClient.List(ctx, instancesOnFleet1List,
				client.InNamespace(ns.Name),
				client.MatchingFields{corev1alpha1.InstanceFleetRefNameField: fleet1},
			)).To(Succeed())

			By("inspecting the items")
			Expect(instancesOnFleet1List.Items).To(ConsistOf(*instance1))

			By("listing all instances on machine pool 2")
			instancesOnFleet2List := &corev1alpha1.InstanceList{}
			Expect(k8sClient.List(ctx, instancesOnFleet2List,
				client.InNamespace(ns.Name),
				client.MatchingFields{corev1alpha1.InstanceFleetRefNameField: fleet2},
			)).To(Succeed())

			By("inspecting the items")
			Expect(instancesOnFleet2List.Items).To(ConsistOf(*instance2))

			By("listing all instances on no machine pool")
			instancesOnNoFleetList := &corev1alpha1.InstanceList{}
			Expect(k8sClient.List(ctx, instancesOnNoFleetList,
				client.InNamespace(ns.Name),
				client.MatchingFields{corev1alpha1.InstanceFleetRefNameField: ""},
			)).To(Succeed())

			By("inspecting the items")
			Expect(instancesOnNoFleetList.Items).To(ConsistOf(*instance3))
		})

		It("should allow listing instances by machine class name", func() {
			By("creating another instance class")
			instanceClass2 := &corev1alpha1.InstanceType{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "instance-class-",
				},
				Capabilities: corev1alpha1.ResourceList{
					corev1alpha1.ResourceCPU:    resource.MustParse("3"),
					corev1alpha1.ResourceMemory: resource.MustParse("10Gi"),
				},
			}
			Expect(k8sClient.Create(ctx, instanceClass2)).To(Succeed())
			DeferCleanup(k8sClient.Delete, ctx, instanceClass2)

			By("creating a instance")
			instance1 := &corev1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:    ns.Name,
					GenerateName: "instance-",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass.Name),
				},
			}
			Expect(k8sClient.Create(ctx, instance1)).To(Succeed())

			By("creating a instance with the other instance class")
			instance2 := &corev1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Namespace:    ns.Name,
					GenerateName: "instance-",
				},
				Spec: corev1alpha1.InstanceSpec{
					InstanceTypeRef: corev1alpha1.LocalObjRef(instanceClass2.Name),
				},
			}
			Expect(k8sClient.Create(ctx, instance2)).To(Succeed())

			By("listing instances with the first instance class name")
			instanceList := &corev1alpha1.InstanceList{}
			Expect(k8sClient.List(ctx, instanceList, client.MatchingFields{
				corev1alpha1.InstanceInstanceTypeRefNameField: instanceClass.Name,
			})).To(Succeed())

			By("inspecting the retrieved list to only have the instance with the correct instance class")
			Expect(instanceList.Items).To(ConsistOf(HaveField("UID", instance1.UID)))

			By("listing instances with the second instance class name")
			Expect(k8sClient.List(ctx, instanceList, client.MatchingFields{
				corev1alpha1.InstanceInstanceTypeRefNameField: instanceClass2.Name,
			})).To(Succeed())

			By("inspecting the retrieved list to only have the instance with the correct instance class")
			Expect(instanceList.Items).To(ConsistOf(HaveField("UID", instance2.UID)))
		})
	})
})
