// SPDX-FileCopyrightText: 2024 Axel Christ and Spheric contributors
// SPDX-License-Identifier: Apache-2.0

package iriserver

import (
	"context"
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	iri "spheric.cloud/spheric/iri-api/apis/runtime/v1alpha1"
	"spheric.cloud/spheric/vee/version"
)

type Server struct {
	iri.UnimplementedRuntimeServiceServer
	dir string
}

func New(dir string) (*Server, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("error ensuring directory at %q: %w", dir, err)
	}

	return &Server{
		dir: dir,
	}, nil
}

func (s *Server) Version(ctx context.Context, request *iri.VersionRequest) (*iri.VersionResponse, error) {
	var runtimeVersion string
	switch {
	case version.Version != "":
		runtimeVersion = version.Version
	case version.Commit != "":
		v, err := semver.NewBuildVersion(version.Commit)
		if err != nil {
			runtimeVersion = "0.0.0"
		} else {
			runtimeVersion = v
		}
	default:
		runtimeVersion = "0.0.0"
	}

	return &iri.VersionResponse{
		RuntimeName:    version.RuntimeName,
		RuntimeVersion: runtimeVersion,
	}, nil
}

func (s *Server) ListInstances(ctx context.Context, request *iri.ListInstancesRequest) (*iri.ListInstancesResponse, error) {
	panic("implement me")
}

func (s *Server) CreateInstance(ctx context.Context, request *iri.CreateInstanceRequest) (*iri.CreateInstanceResponse, error) {
	panic("implement me")
}

func (s *Server) DeleteInstance(ctx context.Context, request *iri.DeleteInstanceRequest) (*iri.DeleteInstanceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) UpdateInstanceAnnotations(ctx context.Context, request *iri.UpdateInstanceAnnotationsRequest) (*iri.UpdateInstanceAnnotationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) UpdateInstancePower(ctx context.Context, request *iri.UpdateInstancePowerRequest) (*iri.UpdateInstancePowerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) AttachDisk(ctx context.Context, request *iri.AttachDiskRequest) (*iri.AttachDiskResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DetachDisk(ctx context.Context, request *iri.DetachDiskRequest) (*iri.DetachDiskResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) AttachNetworkInterface(ctx context.Context, request *iri.AttachNetworkInterfaceRequest) (*iri.AttachNetworkInterfaceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DetachNetworkInterface(ctx context.Context, request *iri.DetachNetworkInterfaceRequest) (*iri.DetachNetworkInterfaceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Status(ctx context.Context, request *iri.StatusRequest) (*iri.StatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Exec(ctx context.Context, request *iri.ExecRequest) (*iri.ExecResponse, error) {
	//TODO implement me
	panic("implement me")
}
