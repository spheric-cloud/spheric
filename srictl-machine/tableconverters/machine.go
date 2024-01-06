// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
	"spheric.cloud/spheric/srictl/api"
	"spheric.cloud/spheric/srictl/tableconverter"
)

var (
	machineHeaders = []api.Header{
		{Name: "ID"},
		{Name: "Class"},
		{Name: "Image"},
		{Name: "State"},
		{Name: "Age"},
	}
)

var (
	Machine = tableconverter.Funcs[*sri.Machine]{
		Headers: tableconverter.Headers(machineHeaders),
		Rows: tableconverter.SingleRowFrom(func(machine *sri.Machine) (api.Row, error) {
			return api.Row{
				machine.Metadata.Id,
				machine.Spec.Class,
				machine.Spec.GetImage().GetImage(),
				machine.Status.State.String(),
				duration.HumanDuration(time.Since(time.Unix(0, machine.Metadata.CreatedAt))),
			}, nil
		}),
	}
	MachineSlice = tableconverter.SliceFuncs[*sri.Machine](Machine)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*sri.Machine](Machine),
		tableconverter.ToTagAndTypedAny[[]*sri.Machine](MachineSlice),
	)
}
