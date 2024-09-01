// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package instance

import (
	"context"

	api "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
)

type RuntimeService interface {
	Version(context.Context, *api.VersionRequest) (*api.VersionResponse, error)
	ListInstances(context.Context, *api.ListInstancesRequest) (*api.ListInstancesResponse, error)
	CreateInstance(context.Context, *api.CreateInstanceRequest) (*api.CreateInstanceResponse, error)
	DeleteInstance(context.Context, *api.DeleteInstanceRequest) (*api.DeleteInstanceResponse, error)
	UpdateInstanceAnnotations(context.Context, *api.UpdateInstanceAnnotationsRequest) (*api.UpdateInstanceAnnotationsResponse, error)
	UpdateInstancePower(context.Context, *api.UpdateInstancePowerRequest) (*api.UpdateInstancePowerResponse, error)
	AttachDisk(context.Context, *api.AttachDiskRequest) (*api.AttachDiskResponse, error)
	DetachDisk(context.Context, *api.DetachDiskRequest) (*api.DetachDiskResponse, error)
	AttachNetworkInterface(context.Context, *api.AttachNetworkInterfaceRequest) (*api.AttachNetworkInterfaceResponse, error)
	DetachNetworkInterface(context.Context, *api.DetachNetworkInterfaceRequest) (*api.DetachNetworkInterfaceResponse, error)
	Status(context.Context, *api.StatusRequest) (*api.StatusResponse, error)
	Exec(context.Context, *api.ExecRequest) (*api.ExecResponse, error)
}
