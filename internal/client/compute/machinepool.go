// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package compute

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	computev1alpha1 "spheric.cloud/spheric/api/compute/v1alpha1"
)

const (
	MachinePoolAvailableMachineClassesField = "machinepool-available-machine-classes"
)

func SetupMachinePoolAvailableMachineClassesFieldIndexer(ctx context.Context, indexer client.FieldIndexer) error {
	return indexer.IndexField(ctx, &computev1alpha1.MachinePool{}, MachinePoolAvailableMachineClassesField, func(object client.Object) []string {
		machinePool := object.(*computev1alpha1.MachinePool)

		names := make([]string, 0, len(machinePool.Status.AvailableMachineClasses))
		for _, availableMachineClass := range machinePool.Status.AvailableMachineClasses {
			names = append(names, availableMachineClass.Name)
		}

		if len(names) == 0 {
			return []string{""}
		}
		return names
	})
}
