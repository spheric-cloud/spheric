// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"
	sri "spheric.cloud/spheric/sri/apis/bucket/v1alpha1"
	"spheric.cloud/spheric/srictl/api"
	"spheric.cloud/spheric/srictl/tableconverter"
)

var (
	bucketHeaders = []api.Header{
		{Name: "ID"},
		{Name: "Class"},
		{Name: "State"},
		{Name: "Age"},
	}
)

var (
	Bucket = tableconverter.Funcs[*sri.Bucket]{
		Headers: tableconverter.Headers(bucketHeaders),
		Rows: tableconverter.SingleRowFrom(func(bucket *sri.Bucket) (api.Row, error) {
			return api.Row{
				bucket.Metadata.Id,
				bucket.Spec.Class,
				bucket.Status.State.String(),
				duration.HumanDuration(time.Since(time.Unix(0, bucket.Metadata.CreatedAt))),
			}, nil
		}),
	}
	BucketSlice = tableconverter.SliceFuncs[*sri.Bucket](Bucket)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*sri.Bucket](Bucket),
		tableconverter.ToTagAndTypedAny[[]*sri.Bucket](BucketSlice),
	)
}
