// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controllers_test

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	commonv1alpha1 "spheric.cloud/spheric/api/common/v1alpha1"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
	networkingv1alpha1 "spheric.cloud/spheric/api/networking/v1alpha1"
	storagev1alpha1 "spheric.cloud/spheric/api/storage/v1alpha1"
	machinepoolletv1alpha1 "spheric.cloud/spheric/poollet/machinepoollet/api/v1alpha1"
	machinepoolletmachine "spheric.cloud/spheric/poollet/machinepoollet/machine"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	testingmachine "spheric.cloud/spheric/sri/testing/machine"
)

var _ = Describe("MachineController", func() {
	ns, mp, mc, srv := SetupTest()

	It("should create a machine", func(ctx SpecContext) {
		By("creating a network")
		network := &networkingv1alpha1.Network{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "network-",
			},
			Spec: networkingv1alpha1.NetworkSpec{
				ProviderID: "foo",
			},
		}
		Expect(k8sClient.Create(ctx, network)).To(Succeed())

		By("patching the network to be available")
		Eventually(UpdateStatus(network, func() {
			network.Status.State = networkingv1alpha1.NetworkStateAvailable
		})).Should(Succeed())

		By("creating a network interface")
		nic := &networkingv1alpha1.NetworkInterface{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "nic-",
			},
			Spec: networkingv1alpha1.NetworkInterfaceSpec{
				NetworkRef: corev1.LocalObjectReference{Name: network.Name},
				IPs: []networkingv1alpha1.IPSource{
					{Value: commonv1alpha1.MustParseNewIP("10.0.0.1")},
				},
			},
		}
		Expect(k8sClient.Create(ctx, nic)).To(Succeed())

		By("creating a volume")
		volume := &storagev1alpha1.Volume{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "volume-",
			},
			Spec: storagev1alpha1.VolumeSpec{},
		}
		Expect(k8sClient.Create(ctx, volume)).To(Succeed())

		By("patching the volume to be available")
		Eventually(UpdateStatus(volume, func() {
			volume.Status.State = storagev1alpha1.VolumeStateAvailable
			volume.Status.Access = &storagev1alpha1.VolumeAccess{
				Driver: "test",
				Handle: "testhandle",
			}
		})).Should(Succeed())

		By("creating a machine")
		const fooAnnotationValue = "bar"
		machine := &computev1alpha1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "machine-",
				Annotations: map[string]string{
					fooAnnotation: fooAnnotationValue,
				},
			},
			Spec: computev1alpha1.MachineSpec{
				MachineClassRef: corev1.LocalObjectReference{Name: mc.Name},
				MachinePoolRef:  &corev1.LocalObjectReference{Name: mp.Name},
				Volumes: []computev1alpha1.Volume{
					{
						Name: "primary",
						VolumeSource: computev1alpha1.VolumeSource{
							VolumeRef: &corev1.LocalObjectReference{Name: volume.Name},
						},
					},
				},
				NetworkInterfaces: []computev1alpha1.NetworkInterface{
					{
						Name: "primary",
						NetworkInterfaceSource: computev1alpha1.NetworkInterfaceSource{
							NetworkInterfaceRef: &corev1.LocalObjectReference{Name: nic.Name},
						},
					},
				},
			},
		}
		Expect(k8sClient.Create(ctx, machine)).To(Succeed())

		By("waiting for the runtime to report the machine, volume and network interface")
		Eventually(srv).Should(SatisfyAll(
			HaveField("Machines", HaveLen(1)),
		))
		_, sriMachine := GetSingleMapEntry(srv.Machines)

		By("inspecting the sri machine")
		Expect(sriMachine.Metadata.Labels).To(HaveKeyWithValue(machinepoolletv1alpha1.DownwardAPILabel(fooDownwardAPILabel), fooAnnotationValue))
		Expect(sriMachine.Spec.Class).To(Equal(mc.Name))
		Expect(sriMachine.Spec.Power).To(Equal(sri.Power_POWER_ON))
		Expect(sriMachine.Spec.Volumes).To(ConsistOf(&sri.Volume{
			Name:   "primary",
			Device: "oda",
			Connection: &sri.VolumeConnection{
				Driver: "test",
				Handle: "testhandle",
			},
		}))
		Expect(sriMachine.Spec.NetworkInterfaces).To(ConsistOf(&sri.NetworkInterface{
			Name:      "primary",
			NetworkId: "foo",
			Ips:       []string{"10.0.0.1"},
		}))

		By("waiting for the spheric machine status to be up-to-date")
		expectedMachineID := machinepoolletmachine.MakeID(testingmachine.FakeRuntimeName, sriMachine.Metadata.Id)
		Eventually(Object(machine)).Should(SatisfyAll(
			HaveField("Status.MachineID", expectedMachineID.String()),
			HaveField("Status.ObservedGeneration", machine.Generation),
		))

		By("setting the network interface id in the machine status")
		sriMachine = &testingmachine.FakeMachine{Machine: *proto.Clone(&sriMachine.Machine).(*sri.Machine)}
		sriMachine.Metadata.Generation = 1
		sriMachine.Status.ObservedGeneration = 1
		sriMachine.Status.NetworkInterfaces = []*sri.NetworkInterfaceStatus{
			{
				Name:   "primary",
				Handle: "primary-handle",
				State:  sri.NetworkInterfaceState_NETWORK_INTERFACE_ATTACHED,
			},
		}
		srv.SetMachines([]*testingmachine.FakeMachine{sriMachine})

		By("waiting for the spheric network interface to have a provider id set")
		Eventually(Object(nic)).Should(HaveField("Spec.ProviderID", "primary-handle"))
		Eventually(Object(machine)).Should(HaveField("Status.NetworkInterfaces", ConsistOf(MatchFields(IgnoreExtras, Fields{
			"Name":   Equal("primary"),
			"Handle": Equal("primary-handle"),
			"State":  Equal(computev1alpha1.NetworkInterfaceStateAttached),
		}))))
	})

	It("should correctly manage the power state of a machine", func(ctx SpecContext) {
		By("creating a machine")
		machine := &computev1alpha1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "machine-",
			},
			Spec: computev1alpha1.MachineSpec{
				MachineClassRef: corev1.LocalObjectReference{Name: mc.Name},
				MachinePoolRef:  &corev1.LocalObjectReference{Name: mp.Name},
			},
		}
		Expect(k8sClient.Create(ctx, machine)).To(Succeed())

		By("waiting for the machine to be created")
		Eventually(srv).Should(HaveField("Machines", HaveLen(1)))

		By("inspecting the machine")
		_, sriMachine := GetSingleMapEntry(srv.Machines)
		Expect(sriMachine.Spec.Power).To(Equal(sri.Power_POWER_ON))

		By("updating the machine power")
		base := machine.DeepCopy()
		machine.Spec.Power = computev1alpha1.PowerOff
		Expect(k8sClient.Patch(ctx, machine, client.MergeFrom(base))).To(Succeed())

		By("waiting for the sri machine to be updated")
		Eventually(sriMachine).Should(HaveField("Spec.Power", Equal(sri.Power_POWER_OFF)))
	})
})

func GetSingleMapEntry[K comparable, V any](m map[K]V) (K, V) {
	if n := len(m); n != 1 {
		Fail(fmt.Sprintf("Expected for map to have a single entry but got %d", n), 1)
	}
	for k, v := range m {
		return k, v
	}
	panic("unreachable")
}
