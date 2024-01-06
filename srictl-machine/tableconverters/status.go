// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"k8s.io/apimachinery/pkg/api/resource"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/srictl/api"
	"spheric.cloud/spheric/srictl/tableconverter"
)

var (
	machineClassHeaders = []api.Header{
		{Name: "Name"},
		{Name: "CPU"},
		{Name: "Memory"},
		{Name: "Quantity"},
	}

	MachineClassStatus = tableconverter.Funcs[*sri.MachineClassStatus]{
		Headers: tableconverter.Headers(machineClassHeaders),
		Rows: tableconverter.SingleRowFrom(func(status *sri.MachineClassStatus) (api.Row, error) {
			return api.Row{
				status.MachineClass.Name,
				resource.NewMilliQuantity(status.MachineClass.Capabilities.CpuMillis, resource.DecimalSI).String(),
				resource.NewQuantity(int64(status.MachineClass.Capabilities.MemoryBytes), resource.DecimalSI).String(),
				resource.NewQuantity(status.Quantity, resource.DecimalSI).String(),
			}, nil
		}),
	}

	MachineClassStatusSlice = tableconverter.SliceFuncs[*sri.MachineClassStatus](MachineClassStatus)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*sri.MachineClassStatus](MachineClassStatus),
		tableconverter.ToTagAndTypedAny[[]*sri.MachineClassStatus](MachineClassStatusSlice),
	)
}
