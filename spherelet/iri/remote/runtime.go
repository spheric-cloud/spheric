// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package remote

import (
	"context"
	"fmt"

	"spheric.cloud/spheric/spherelet/instance"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	sri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
)

type remoteRuntime struct {
	client sri.RuntimeServiceClient
}

func NewRemoteRuntime(endpoint string) (instance.RuntimeService, error) {
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("error dialing: %w", err)
	}

	return &remoteRuntime{
		client: sri.NewRuntimeServiceClient(conn),
	}, nil
}

func (r *remoteRuntime) Version(ctx context.Context, req *sri.VersionRequest) (*sri.VersionResponse, error) {
	return r.client.Version(ctx, req)
}

func (r *remoteRuntime) ListInstances(ctx context.Context, req *sri.ListInstancesRequest) (*sri.ListInstancesResponse, error) {
	return r.client.ListInstances(ctx, req)
}

func (r *remoteRuntime) CreateInstance(ctx context.Context, req *sri.CreateInstanceRequest) (*sri.CreateInstanceResponse, error) {
	return r.client.CreateInstance(ctx, req)
}

func (r *remoteRuntime) DeleteInstance(ctx context.Context, req *sri.DeleteInstanceRequest) (*sri.DeleteInstanceResponse, error) {
	return r.client.DeleteInstance(ctx, req)
}

func (r *remoteRuntime) UpdateInstanceAnnotations(ctx context.Context, req *sri.UpdateInstanceAnnotationsRequest) (*sri.UpdateInstanceAnnotationsResponse, error) {
	return r.client.UpdateInstanceAnnotations(ctx, req)
}

func (r *remoteRuntime) UpdateInstancePower(ctx context.Context, req *sri.UpdateInstancePowerRequest) (*sri.UpdateInstancePowerResponse, error) {
	return r.client.UpdateInstancePower(ctx, req)
}

func (r *remoteRuntime) AttachDisk(ctx context.Context, req *sri.AttachDiskRequest) (*sri.AttachDiskResponse, error) {
	return r.client.AttachDisk(ctx, req)
}

func (r *remoteRuntime) DetachDisk(ctx context.Context, req *sri.DetachDiskRequest) (*sri.DetachDiskResponse, error) {
	return r.client.DetachDisk(ctx, req)
}

func (r *remoteRuntime) AttachNetworkInterface(ctx context.Context, req *sri.AttachNetworkInterfaceRequest) (*sri.AttachNetworkInterfaceResponse, error) {
	return r.client.AttachNetworkInterface(ctx, req)
}

func (r *remoteRuntime) DetachNetworkInterface(ctx context.Context, req *sri.DetachNetworkInterfaceRequest) (*sri.DetachNetworkInterfaceResponse, error) {
	return r.client.DetachNetworkInterface(ctx, req)
}

func (r *remoteRuntime) Status(ctx context.Context, req *sri.StatusRequest) (*sri.StatusResponse, error) {
	return r.client.Status(ctx, req)
}

func (r *remoteRuntime) Exec(ctx context.Context, req *sri.ExecRequest) (*sri.ExecResponse, error) {
	return r.client.Exec(ctx, req)
}
