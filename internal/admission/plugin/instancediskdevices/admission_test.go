// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instancediskdevices_test

import (
	"context"

	"k8s.io/utils/ptr"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"
	. "spheric.cloud/spheric/internal/admission/plugin/instancediskdevices"
	"spheric.cloud/spheric/internal/apis/core"
)

var _ = Describe("Admission", func() {
	var (
		plugin *InstanceDiskDevices
	)
	BeforeEach(func() {
		plugin = NewInstanceDiskDevices()
	})

	It("should ignore non-instance objects", func() {
		disk := &core.Disk{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "foo",
				Name:      "bar",
			},
		}
		origDisk := disk.DeepCopy()
		Expect(plugin.Admit(
			context.TODO(),
			admission.NewAttributesRecord(
				disk,
				nil,
				core.Kind("Disk").GroupKind().WithVersion("version"),
				disk.Namespace,
				disk.Name,
				core.Resource("disks").WithVersion("version"),
				"",
				admission.Create,
				&metav1.CreateOptions{},
				false,
				nil,
			),
			nil,
		)).NotTo(HaveOccurred())
		Expect(disk).To(Equal(origDisk))
	})

	It("should add disk device names when unset", func() {
		instance := &core.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "foo",
				Name:      "bar",
			},
			Spec: core.InstanceSpec{
				Disks: []core.AttachedDisk{
					{
						Device: ptr.To("odb"),
					},
					{},
					{},
				},
			},
		}
		Expect(plugin.Admit(
			context.TODO(),
			admission.NewAttributesRecord(
				instance,
				nil,
				core.Kind("Instance").GroupKind().WithVersion("version"),
				instance.Namespace,
				instance.Name,
				core.Resource("instances").WithVersion("version"),
				"",
				admission.Create,
				&metav1.CreateOptions{},
				false,
				nil,
			),
			nil,
		)).NotTo(HaveOccurred())

		Expect(instance.Spec.Disks).To(Equal([]core.AttachedDisk{
			{
				Device: ptr.To("odb"),
			},
			{
				Device: ptr.To("oda"),
			},
			{
				Device: ptr.To("odc"),
			},
		}))
	})
})
