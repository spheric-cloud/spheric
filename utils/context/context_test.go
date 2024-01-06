// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package context_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "spheric.cloud/spheric/utils/context"
)

var _ = Describe("Context", func() {
	Describe("FromStopChannel", func() {
		It("should create a context from the given stop channel", func() {
			stopChan := make(chan struct{})
			ctx := FromStopChannel(stopChan)

			Expect(ctx.Done()).NotTo(BeClosed())
			Expect(ctx.Err()).To(Succeed())

			close(stopChan)

			Expect(ctx.Done()).To(BeClosed())
			Expect(ctx.Err()).To(Equal(context.Canceled))
		})

		It("should panic if a value is sent through the channel when determining whether it's closed", func() {
			stopChan := make(chan struct{}, 1)
			ctx := FromStopChannel(stopChan)

			stopChan <- struct{}{}
			Expect(func() { _ = ctx.Err() }).To(Panic())
		})
	})
})
