// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/storage"
	. "spheric.cloud/spheric/internal/apis/storage/validation"
	. "spheric.cloud/spheric/internal/testutils/validation"
)

var _ = Describe("VolumeTemplate", func() {
	DescribeTable("ValidateVolumeTemplateSpec",
		func(spec *storage.VolumeTemplateSpec, match types.GomegaMatcher) {
			errList := ValidateVolumeTemplateSpec(spec, field.NewPath("spec"))
			Expect(errList).To(match)
		},
		Entry("forbidden metadata name",
			&storage.VolumeTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
			ContainElement(ForbiddenField("spec.metadata.name")),
		),
		Entry("forbidden metadata namespace",
			&storage.VolumeTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo",
				},
			},
			ContainElement(ForbiddenField("spec.metadata.namespace")),
		),
		Entry("forbidden metadata generate name",
			&storage.VolumeTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "foo",
				},
			},
			ContainElement(ForbiddenField("spec.metadata.generateName")),
		),
	)
})
