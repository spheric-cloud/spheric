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
	"spheric.cloud/spheric/internal/apis/storage"
	. "spheric.cloud/spheric/internal/apis/storage/validation"
	. "spheric.cloud/spheric/internal/testutils/validation"
)

var _ = Describe("BucketPool", func() {
	DescribeTable("ValidateBucketPool",
		func(bucketPool *storage.BucketPool, match types.GomegaMatcher) {
			errList := ValidateBucketPool(bucketPool)
			Expect(errList).To(match)
		},
		Entry("missing name",
			&storage.BucketPool{},
			ContainElement(RequiredField("metadata.name")),
		),
		Entry("bad name",
			&storage.BucketPool{ObjectMeta: metav1.ObjectMeta{Name: "foo*"}},
			ContainElement(InvalidField("metadata.name")),
		),
		Entry("dns subdomain name",
			&storage.BucketPool{ObjectMeta: metav1.ObjectMeta{Name: "foo.bar.baz"}},
			Not(ContainElement(InvalidField("metadata.name"))),
		),
	)

	DescribeTable("ValidateBucketUpdate",
		func(newBucketPool, oldBucketPool *storage.BucketPool, match types.GomegaMatcher) {
			errList := ValidateBucketPoolUpdate(newBucketPool, oldBucketPool)
			Expect(errList).To(match)
		},
		Entry("immutable providerID if set",
			&storage.BucketPool{
				Spec: storage.BucketPoolSpec{
					ProviderID: "foo",
				},
			},
			&storage.BucketPool{
				Spec: storage.BucketPoolSpec{
					ProviderID: "bar",
				},
			},
			ContainElement(ImmutableField("spec.providerID")),
		),
		Entry("mutable providerID if not set",
			&storage.BucketPool{
				Spec: storage.BucketPoolSpec{
					ProviderID: "foo",
				},
			},
			&storage.BucketPool{},
			Not(ContainElement(ImmutableField("spec.providerID"))),
		),
	)
})
