// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package vcm

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"spheric.cloud/spheric/poollet/srievent"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
)

var (
	ErrNoMatchingVolumeClass        = errors.New("no matching volume class")
	ErrAmbiguousMatchingVolumeClass = errors.New("ambiguous matching volume classes")
)

type VolumeClassMapper interface {
	manager.Runnable
	GetVolumeClassFor(ctx context.Context, name string, capabilities *sri.VolumeClassCapabilities) (*sri.VolumeClass, *resource.Quantity, error)
	WaitForSync(ctx context.Context) error
	AddListener(listener srievent.Listener) (srievent.ListenerRegistration, error)
	RemoveListener(reg srievent.ListenerRegistration) error
}
