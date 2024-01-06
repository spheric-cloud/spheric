// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package mcm

import (
	"context"
	"errors"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	"spheric.cloud/spheric/poollet/srievent"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
)

var (
	ErrNoMatchingMachineClass        = errors.New("no matching machine class")
	ErrAmbiguousMatchingMachineClass = errors.New("ambiguous matching machine classes")
)

type MachineClassMapper interface {
	manager.Runnable
	GetMachineClassFor(ctx context.Context, name string, capabilities *sri.MachineClassCapabilities) (*sri.MachineClass, int64, error)
	WaitForSync(ctx context.Context) error
	AddListener(listener srievent.Listener) (srievent.ListenerRegistration, error)
	RemoveListener(reg srievent.ListenerRegistration) error
}
