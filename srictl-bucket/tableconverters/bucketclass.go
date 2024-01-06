// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	"spheric.cloud/spheric/srictl/api"
	"spheric.cloud/spheric/srictl/tableconverter"
)

var (
	bucketClassHeaders = []api.Header{
		{Name: "Name"},
		{Name: "TPS"},
		{Name: "IOPS"},
	}
)

var (
	BucketClass = tableconverter.Funcs[*sri.BucketClass]{
		Headers: tableconverter.Headers(bucketClassHeaders),
		Rows: tableconverter.SingleRowFrom(func(class *sri.BucketClass) (api.Row, error) {
			return api.Row{
				class.Name,
				class.Capabilities.Tps,
				class.Capabilities.Iops,
			}, nil
		}),
	}
	BucketClassSlice = tableconverter.SliceFuncs[*sri.BucketClass](BucketClass)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*sri.BucketClass](BucketClass),
		tableconverter.ToTagAndTypedAny[[]*sri.BucketClass](BucketClassSlice),
	)
}
