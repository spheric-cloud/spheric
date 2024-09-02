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
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
)

type remoteRuntime struct {
	client iri.RuntimeServiceClient
}

func NewRemoteRuntime(endpoint string) (instance.RuntimeService, error) {
	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("error dialing: %w", err)
	}

	return &remoteRuntime{
		client: iri.NewRuntimeServiceClient(conn),
	}, nil
}

func (r *remoteRuntime) Version(ctx context.Context, req *iri.VersionRequest) (*iri.VersionResponse, error) {
	return r.client.Version(ctx, req)
}

func (r *remoteRuntime) ListInstances(ctx context.Context, req *iri.ListInstancesRequest) (*iri.ListInstancesResponse, error) {
	return r.client.ListInstances(ctx, req)
}

func (r *remoteRuntime) CreateInstance(ctx context.Context, req *iri.CreateInstanceRequest) (*iri.CreateInstanceResponse, error) {
	return r.client.CreateInstance(ctx, req)
}

func (r *remoteRuntime) DeleteInstance(ctx context.Context, req *iri.DeleteInstanceRequest) (*iri.DeleteInstanceResponse, error) {
	return r.client.DeleteInstance(ctx, req)
}

func (r *remoteRuntime) UpdateInstanceAnnotations(ctx context.Context, req *iri.UpdateInstanceAnnotationsRequest) (*iri.UpdateInstanceAnnotationsResponse, error) {
	return r.client.UpdateInstanceAnnotations(ctx, req)
}

func (r *remoteRuntime) UpdateInstancePower(ctx context.Context, req *iri.UpdateInstancePowerRequest) (*iri.UpdateInstancePowerResponse, error) {
	return r.client.UpdateInstancePower(ctx, req)
}

func (r *remoteRuntime) AttachDisk(ctx context.Context, req *iri.AttachDiskRequest) (*iri.AttachDiskResponse, error) {
	return r.client.AttachDisk(ctx, req)
}

func (r *remoteRuntime) DetachDisk(ctx context.Context, req *iri.DetachDiskRequest) (*iri.DetachDiskResponse, error) {
	return r.client.DetachDisk(ctx, req)
}

func (r *remoteRuntime) AttachNetworkInterface(ctx context.Context, req *iri.AttachNetworkInterfaceRequest) (*iri.AttachNetworkInterfaceResponse, error) {
	return r.client.AttachNetworkInterface(ctx, req)
}

func (r *remoteRuntime) DetachNetworkInterface(ctx context.Context, req *iri.DetachNetworkInterfaceRequest) (*iri.DetachNetworkInterfaceResponse, error) {
	return r.client.DetachNetworkInterface(ctx, req)
}

func (r *remoteRuntime) Status(ctx context.Context, req *iri.StatusRequest) (*iri.StatusResponse, error) {
	return r.client.Status(ctx, req)
}

func (r *remoteRuntime) Exec(ctx context.Context, req *iri.ExecRequest) (*iri.ExecResponse, error) {
	return r.client.Exec(ctx, req)
}
