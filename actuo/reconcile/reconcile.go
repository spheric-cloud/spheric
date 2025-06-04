// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package reconcile

import (
	"context"
	"time"
)

type Reconciler[Request any] interface {
	Reconcile(ctx context.Context, request Request) (Result, error)
}

type Result struct {
	Requeue      bool
	RequeueAfter time.Duration
}
