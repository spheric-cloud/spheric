// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instance

import (
	"context"

	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
)

type RuntimeService interface {
	Version(context.Context, *iri.VersionRequest) (*iri.VersionResponse, error)
	ListInstances(context.Context, *iri.ListInstancesRequest) (*iri.ListInstancesResponse, error)
	CreateInstance(context.Context, *iri.CreateInstanceRequest) (*iri.CreateInstanceResponse, error)
	DeleteInstance(context.Context, *iri.DeleteInstanceRequest) (*iri.DeleteInstanceResponse, error)
	UpdateInstanceAnnotations(context.Context, *iri.UpdateInstanceAnnotationsRequest) (*iri.UpdateInstanceAnnotationsResponse, error)
	UpdateInstancePower(context.Context, *iri.UpdateInstancePowerRequest) (*iri.UpdateInstancePowerResponse, error)
	AttachDisk(context.Context, *iri.AttachDiskRequest) (*iri.AttachDiskResponse, error)
	DetachDisk(context.Context, *iri.DetachDiskRequest) (*iri.DetachDiskResponse, error)
	AttachNetworkInterface(context.Context, *iri.AttachNetworkInterfaceRequest) (*iri.AttachNetworkInterfaceResponse, error)
	DetachNetworkInterface(context.Context, *iri.DetachNetworkInterfaceRequest) (*iri.DetachNetworkInterfaceResponse, error)
	Status(context.Context, *iri.StatusRequest) (*iri.StatusResponse, error)
	Exec(context.Context, *iri.ExecRequest) (*iri.ExecResponse, error)
}
