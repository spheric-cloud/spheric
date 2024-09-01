// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package tableconverters

import (
	"time"

	"k8s.io/apimachinery/pkg/util/duration"

	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/irictl/api"
	"spheric.cloud/spheric/irictl/tableconverter"
)

var (
	instanceHeaders = []api.Header{
		{Name: "ID"},
		{Name: "Type"},
		{Name: "Image"},
		{Name: "State"},
		{Name: "Age"},
	}
)

var (
	Instance = tableconverter.Funcs[*iri.Instance]{
		Headers: tableconverter.Headers(instanceHeaders),
		Rows: tableconverter.SingleRowFrom(func(instance *iri.Instance) (api.Row, error) {
			return api.Row{
				instance.Metadata.Id,
				instance.Spec.Type,
				instance.Spec.GetImage().GetImage(),
				instance.Status.State.String(),
				duration.HumanDuration(time.Since(time.Unix(0, instance.Metadata.CreatedAt))),
			}, nil
		}),
	}
	InstanceSlice = tableconverter.SliceFuncs[*iri.Instance](Instance)
)

func init() {
	RegistryBuilder.Register(
		tableconverter.ToTagAndTypedAny[*iri.Instance](Instance),
		tableconverter.ToTagAndTypedAny[[]*iri.Instance](InstanceSlice),
	)
}
