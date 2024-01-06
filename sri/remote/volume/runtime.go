// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package volume

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"spheric.cloud/spheric/sri/apis/volume"
	sri "spheric.cloud/spheric/sri/apis/volume/v1alpha1"
)

type remoteRuntime struct {
	client sri.VolumeRuntimeClient
}

func NewRemoteRuntime(endpoint string) (volume.RuntimeService, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("error dialing: %w", err)
	}

	return &remoteRuntime{
		client: sri.NewVolumeRuntimeClient(conn),
	}, nil
}

func (r *remoteRuntime) ListVolumes(ctx context.Context, request *sri.ListVolumesRequest) (*sri.ListVolumesResponse, error) {
	return r.client.ListVolumes(ctx, request)
}

func (r *remoteRuntime) CreateVolume(ctx context.Context, request *sri.CreateVolumeRequest) (*sri.CreateVolumeResponse, error) {
	return r.client.CreateVolume(ctx, request)
}

func (r *remoteRuntime) ExpandVolume(ctx context.Context, request *sri.ExpandVolumeRequest) (*sri.ExpandVolumeResponse, error) {
	return r.client.ExpandVolume(ctx, request)
}

func (r *remoteRuntime) DeleteVolume(ctx context.Context, request *sri.DeleteVolumeRequest) (*sri.DeleteVolumeResponse, error) {
	return r.client.DeleteVolume(ctx, request)
}

func (r *remoteRuntime) Status(ctx context.Context, request *sri.StatusRequest) (*sri.StatusResponse, error) {
	return r.client.Status(ctx, request)
}
