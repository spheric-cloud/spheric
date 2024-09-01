// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"k8s.io/apimachinery/pkg/api/resource"
	. "spheric.cloud/spheric/utils/testing"

	testinginstance "spheric.cloud/spheric/spherelet/iri/remote/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	corev1alpha1 "spheric.cloud/spheric/api/core/v1alpha1"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	sphereletv1alpha1 "spheric.cloud/spheric/spherelet/api/v1alpha1"
	sphereletinstance "spheric.cloud/spheric/spherelet/instance"
)

var _ = Describe("InstanceController", func() {
	ns := SetupNamespace(k8sClient)
	instanceType := SetupInstanceType()

	It("should create a instance", func(ctx SpecContext) {
		By("creating a network")
		network := &corev1alpha1.Network{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "network-",
			},
			Spec: corev1alpha1.NetworkSpec{},
		}
		Expect(k8sClient.Create(ctx, network)).To(Succeed())

		By("patching the network to be available")
		Eventually(UpdateStatus(network, func() {
			network.Status.State = corev1alpha1.NetworkStateAvailable
		})).Should(Succeed())

		By("creating a subnet")
		subnet := &corev1alpha1.Subnet{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "subnet-",
			},
			Spec: corev1alpha1.SubnetSpec{
				NetworkRef: corev1alpha1.LocalObjectReference{
					Name: network.Name,
				},
				CIDRs: []string{"10.0.0.0/24"},
			},
		}
		Expect(k8sClient.Create(ctx, subnet)).To(Succeed())

		By("creating a disk")
		disk := &corev1alpha1.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "disk-",
			},
			Spec: corev1alpha1.DiskSpec{
				Resources: corev1alpha1.ResourceList{
					corev1alpha1.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
		}
		Expect(k8sClient.Create(ctx, disk)).To(Succeed())

		By("patching the disk to be available")
		Eventually(UpdateStatus(disk, func() {
			disk.Status.State = corev1alpha1.DiskStateAvailable
			disk.Status.Access = &corev1alpha1.DiskAccess{
				Driver: "test-driver",
				Handle: "test-handle",
			}
		})).Should(Succeed())

		By("creating a instance")
		const fooAnnotationValue = "bar"
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "instance-",
				Annotations: map[string]string{
					fooAnnotation: fooAnnotationValue,
				},
			},
			Spec: corev1alpha1.InstanceSpec{
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
				FleetRef:        corev1alpha1.NewLocalObjRef(fleetName),
				Disks: []corev1alpha1.AttachedDisk{
					{
						Name: "primary",
						AttachedDiskSource: corev1alpha1.AttachedDiskSource{
							DiskRef: corev1alpha1.NewLocalObjRef(disk.Name),
						},
					},
				},
				NetworkInterfaces: []corev1alpha1.NetworkInterface{
					{
						Name: "primary",
						SubnetRef: corev1alpha1.SubnetReference{
							NetworkName: network.Name,
							Name:        subnet.Name,
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("waiting for the runtime to report the instance, disk and network interface")
		iriInstance := NewFakeInstanceWithUID(instance.UID)
		Eventually(GetInstanceByUID(srv, iriInstance)).Should(Succeed())

		By("inspecting the iri instance")
		Expect(iriInstance.Metadata.Labels).To(HaveKeyWithValue(sphereletv1alpha1.DownwardAPILabel(fooDownwardAPILabel), fooAnnotationValue))
		Expect(iriInstance.Spec.Type).To(Equal(instanceType.Name))
		Expect(iriInstance.Spec.Power).To(Equal(iri.Power_POWER_ON))
		Expect(iriInstance.Spec.Disks).To(ConsistOf(ProtoEqual(&iri.Disk{
			Name:   "primary",
			Device: "oda",
			Connection: &iri.DiskConnection{
				Driver: "test-driver",
				Handle: "test-handle",
			},
		})))
		Expect(iriInstance.Spec.NetworkInterfaces).To(ConsistOf(ProtoEqual(&iri.NetworkInterface{
			Name: "primary",
			SubnetMetadata: &iri.NetworkInterfaceSubnetMetadata{
				NetworkName: network.Name,
				NetworkUid:  string(network.UID),
				SubnetName:  subnet.Name,
				SubnetUid:   string(subnet.UID),
			},
			SubnetCidrs: []string{
				"10.0.0.0/24",
			},
		})))

		By("waiting for the spheric instance status to be up-to-date")
		expectedInstanceID := sphereletinstance.MakeID(testinginstance.RuntimeName, iriInstance.Metadata.Id)
		Eventually(Object(instance)).Should(SatisfyAll(
			HaveField("Status.InstanceID", expectedInstanceID.String()),
			HaveField("Status.ObservedGeneration", instance.Generation),
		))

		By("setting the network interface id in the instance status")
		Eventually(UpdateInstance(srv, iriInstance, func() {
			iriInstance.Metadata.Generation = 1
			iriInstance.Status.ObservedGeneration = 1
			iriInstance.Status.NetworkInterfaces = []*iri.NetworkInterfaceStatus{
				{
					Name:   "primary",
					Handle: "primary-handle",
					State:  iri.NetworkInterfaceState_NETWORK_INTERFACE_ATTACHED,
				},
			}
		})).Should(Succeed())

		By("waiting for the spheric network interface to have a provider id set")
		Eventually(Object(instance)).Should(HaveField("Status.NetworkInterfaces", ConsistOf(MatchFields(IgnoreExtras, Fields{
			"Name":  Equal("primary"),
			"State": Equal(corev1alpha1.NetworkInterfaceStateAttached),
		}))))
	})

	It("should correctly manage the power state of a instance", func(ctx SpecContext) {
		By("creating a instance")
		instance := &corev1alpha1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "instance-",
			},
			Spec: corev1alpha1.InstanceSpec{
				InstanceTypeRef: corev1alpha1.LocalObjRef(instanceType.Name),
				FleetRef:        corev1alpha1.NewLocalObjRef(fleetName),
			},
		}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("waiting for the instance to be created")
		iriInstance := NewFakeInstanceWithUID(instance.UID)
		Eventually(GetInstanceByUID(srv, iriInstance)).Should(Succeed())

		By("inspecting the instance")
		Expect(iriInstance.Spec.Power).To(Equal(iri.Power_POWER_ON))

		By("updating the instance power")
		base := instance.DeepCopy()
		instance.Spec.Power = corev1alpha1.PowerOff
		Expect(k8sClient.Patch(ctx, instance, client.MergeFrom(base))).To(Succeed())

		By("waiting for the iri instance to be updated")
		Eventually(Instance(srv, iriInstance)).Should(HaveField("Spec.Power", Equal(iri.Power_POWER_OFF)))
	})
})
