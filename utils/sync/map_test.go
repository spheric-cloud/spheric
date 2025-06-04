// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package sync

import (
	"iter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("Map", func() {
	Describe("Get/Set", func() {
		It("should correctly get and set the values", func() {
			m := NewMap[string, int]()

			m.Set("a", 1)

			a, ok := m.GetOK("a")
			Expect(ok).To(BeTrue())
			Expect(a).To(Equal(1))

			b, ok := m.GetOK("b")
			Expect(ok).To(BeFalse())
			Expect(b).To(BeZero())
		})
	})

	Describe("Values", func() {
		It("should be possible to do simultaneous reads but not writes / global writes", func() {
			m := NewMap[string, int]()
			m.Set("a", 1)
			m.Set("b", 2)

			next, stop := iter.Pull2(m.All())
			defer stop()

			k, v, ok := next()
			Expect(ok).To(BeTrue())
			Expect(k).To(Or(Equal("a"), Equal("b")))
			Expect(v).To(Or(Equal(1), Equal(2)))

			setCDone := make(chan struct{})
			go func() {
				defer close(setCDone)
				m.Set("c", 3)
			}()

			By("checking if the set is still pending")
			Consistently(setCDone).ShouldNot(BeClosed())

			k, v, ok = next()
			Expect(ok).To(BeTrue())
			Expect(k).To(Or(Equal("a"), Equal("b")))
			Expect(v).To(Or(Equal(1), Equal(2)))

			By("checking if the set is still pending")
			Consistently(setCDone).ShouldNot(BeClosed())

			_, _, ok = next()
			Expect(ok).To(BeFalse())

			By("checking if the set is done")
			Consistently(setCDone).Should(BeClosed())
		})
	})
})
