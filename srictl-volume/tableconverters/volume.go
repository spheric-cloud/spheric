// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
	"spheric.cloud/spheric/srictl/api"
	"spheric.cloud/spheric/srictl/tableconverter"
)

var (
	volumeHeaders = []api.Header{
		{Name: "ID"},
		{Name: "Class"},
		{Name: "Image"},
		{Name: "State"},
		{Name: "Age"},
	}
)

var (
	Volume = tableconverter.Funcs[*sri.Volume]{
		Headers: tableconverter.Headers(volumeHeaders),
		Rows: tableconverter.SingleRowFrom(func(volume *sri.Volume) (api.Row, error) {
			return api.Row{
				volume.Metadata.Id,
				volume.Spec.Class,
				volume.Spec.Image,
				volume.Status.State.String(),
				duration.HumanDuration(time.Since(time.Unix(0, volume.Metadata.CreatedAt))),
			}, nil
		}),
	}
	VolumeSlice = tableconverter.SliceFuncs[*sri.Volume](Volume)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*sri.Volume](Volume),
		tableconverter.ToTagAndTypedAny[[]*sri.Volume](VolumeSlice),
	)
}
