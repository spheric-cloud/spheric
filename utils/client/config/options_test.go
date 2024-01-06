// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"spheric.cloud/spheric/utils/client/config"
)

var _ = Describe("Options", func() {
	Describe("WithNamePrefix", func() {
		It("should set a prefixing function as name function", func() {
			o := &config.BindFlagOptions{}
			config.WithNamePrefix("foo")(o)

			Expect(o.NameFunc).NotTo(BeNil())
			Expect(o.NameFunc("bar")).To(Equal("foobar"))
		})
	})

	Describe("WithNameSuffix", func() {
		It("should set suffixing function as name function", func() {
			o := &config.BindFlagOptions{}
			config.WithNameSuffix("bar")(o)

			Expect(o.NameFunc).NotTo(BeNil())
			Expect(o.NameFunc("foo")).To(Equal("foobar"))
		})
	})
})
