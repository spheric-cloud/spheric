// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package machine

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"spheric.cloud/spheric/sri/apis/machine"
	sri "spheric.cloud/spheric/sri/apis/machine/v1alpha1"
)

type remoteRuntime struct {
	client sri.MachineRuntimeClient
}

func NewRemoteRuntime(endpoint string) (machine.RuntimeService, error) {
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("error dialing: %w", err)
	}

	return &remoteRuntime{
		client: sri.NewMachineRuntimeClient(conn),
	}, nil
}

func (r *remoteRuntime) Version(ctx context.Context, req *sri.VersionRequest) (*sri.VersionResponse, error) {
	return r.client.Version(ctx, req)
}

func (r *remoteRuntime) ListMachines(ctx context.Context, req *sri.ListMachinesRequest) (*sri.ListMachinesResponse, error) {
	return r.client.ListMachines(ctx, req)
}

func (r *remoteRuntime) CreateMachine(ctx context.Context, req *sri.CreateMachineRequest) (*sri.CreateMachineResponse, error) {
	return r.client.CreateMachine(ctx, req)
}

func (r *remoteRuntime) DeleteMachine(ctx context.Context, req *sri.DeleteMachineRequest) (*sri.DeleteMachineResponse, error) {
	return r.client.DeleteMachine(ctx, req)
}

func (r *remoteRuntime) UpdateMachineAnnotations(ctx context.Context, req *sri.UpdateMachineAnnotationsRequest) (*sri.UpdateMachineAnnotationsResponse, error) {
	return r.client.UpdateMachineAnnotations(ctx, req)
}

func (r *remoteRuntime) UpdateMachinePower(ctx context.Context, req *sri.UpdateMachinePowerRequest) (*sri.UpdateMachinePowerResponse, error) {
	return r.client.UpdateMachinePower(ctx, req)
}

func (r *remoteRuntime) AttachVolume(ctx context.Context, req *sri.AttachVolumeRequest) (*sri.AttachVolumeResponse, error) {
	return r.client.AttachVolume(ctx, req)
}

func (r *remoteRuntime) DetachVolume(ctx context.Context, req *sri.DetachVolumeRequest) (*sri.DetachVolumeResponse, error) {
	return r.client.DetachVolume(ctx, req)
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
