// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"k8s.io/apimachinery/pkg/api/resource"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
	"spheric.cloud/spheric/srictl/api"
	"spheric.cloud/spheric/srictl/tableconverter"
)

var (
	volumeClassHeaders = []api.Header{
		{Name: "Name"},
		{Name: "TPS"},
		{Name: "IOPS"},
		{Name: "Quantity"},
	}
)

var (
	VolumeClassStatus = tableconverter.Funcs[*sri.VolumeClassStatus]{
		Headers: tableconverter.Headers(volumeClassHeaders),
		Rows: tableconverter.SingleRowFrom(func(status *sri.VolumeClassStatus) (api.Row, error) {
			return api.Row{
				status.VolumeClass.Name,
				resource.NewQuantity(status.VolumeClass.Capabilities.Tps, resource.BinarySI).String(),
				resource.NewQuantity(status.VolumeClass.Capabilities.Iops, resource.DecimalSI).String(),
				resource.NewQuantity(status.Quantity, resource.BinarySI).String(),
			}, nil
		}),
	}
	VolumeClassStatusSlice = tableconverter.SliceFuncs[*sri.VolumeClassStatus](VolumeClassStatus)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*sri.VolumeClassStatus](VolumeClassStatus),
		tableconverter.ToTagAndTypedAny[[]*sri.VolumeClassStatus](VolumeClassStatusSlice),
	)
}
