// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bcm

import (
	"context"
	"errors"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
)

var (
	ErrNoMatchingBucketClass        = errors.New("no matching bucket class")
	ErrAmbiguousMatchingBucketClass = errors.New("ambiguous matching bucket classes")
)

type BucketClassMapper interface {
	manager.Runnable
	GetBucketClassFor(ctx context.Context, name string, capabilities *sri.BucketClassCapabilities) (*sri.BucketClass, error)
	WaitForSync(ctx context.Context) error
}
