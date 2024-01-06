// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package volume

import (
	"context"

	api "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
)

type RuntimeService interface {
	ListVolumes(context.Context, *api.ListVolumesRequest) (*api.ListVolumesResponse, error)
	CreateVolume(context.Context, *api.CreateVolumeRequest) (*api.CreateVolumeResponse, error)
	ExpandVolume(ctx context.Context, request *api.ExpandVolumeRequest) (*api.ExpandVolumeResponse, error)
	DeleteVolume(context.Context, *api.DeleteVolumeRequest) (*api.DeleteVolumeResponse, error)

	Status(context.Context, *api.StatusRequest) (*api.StatusResponse, error)
}
