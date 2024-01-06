// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"spheric.cloud/spheric/internal/apis/ipam"
	. "spheric.cloud/spheric/internal/apis/ipam/validation"
	. "spheric.cloud/spheric/internal/testutils/validation"
)

var _ = Describe("PrefixTemplate", func() {
	DescribeTable("ValidatePrefixTemplateSpec",
		func(spec *ipam.PrefixTemplateSpec, match types.GomegaMatcher) {
			errList := ValidatePrefixTemplateSpec(spec, field.NewPath("spec"))
			Expect(errList).To(match)
		},
		Entry("forbidden metadata name",
			&ipam.PrefixTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
			},
			ContainElement(ForbiddenField("spec.metadata.name")),
		),
		Entry("forbidden metadata namespace",
			&ipam.PrefixTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "foo",
				},
			},
			ContainElement(ForbiddenField("spec.metadata.namespace")),
		),
		Entry("forbidden metadata generate name",
			&ipam.PrefixTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "foo",
				},
			},
			ContainElement(ForbiddenField("spec.metadata.generateName")),
		),
		Entry("valid prefix template spec",
			&ipam.PrefixTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{"bar": "baz"},
				},
				Spec: ipam.PrefixSpec{
					IPFamily:     corev1.IPv4Protocol,
					PrefixLength: 28,
					ParentRef:    &corev1.LocalObjectReference{Name: "parent"},
				},
			},
			BeEmpty(),
		),
	)
})
