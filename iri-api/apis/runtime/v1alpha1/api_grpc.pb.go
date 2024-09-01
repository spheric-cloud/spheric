// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.26.1
// source: iri-api/apis/runtime/v1alpha1/api.proto

package v1alpha1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RuntimeService_Version_FullMethodName                   = "/runtime.v1alpha1.RuntimeService/Version"
	RuntimeService_ListInstances_FullMethodName             = "/runtime.v1alpha1.RuntimeService/ListInstances"
	RuntimeService_CreateInstance_FullMethodName            = "/runtime.v1alpha1.RuntimeService/CreateInstance"
	RuntimeService_DeleteInstance_FullMethodName            = "/runtime.v1alpha1.RuntimeService/DeleteInstance"
	RuntimeService_UpdateInstanceAnnotations_FullMethodName = "/runtime.v1alpha1.RuntimeService/UpdateInstanceAnnotations"
	RuntimeService_UpdateInstancePower_FullMethodName       = "/runtime.v1alpha1.RuntimeService/UpdateInstancePower"
	RuntimeService_AttachDisk_FullMethodName                = "/runtime.v1alpha1.RuntimeService/AttachDisk"
	RuntimeService_DetachDisk_FullMethodName                = "/runtime.v1alpha1.RuntimeService/DetachDisk"
	RuntimeService_AttachNetworkInterface_FullMethodName    = "/runtime.v1alpha1.RuntimeService/AttachNetworkInterface"
	RuntimeService_DetachNetworkInterface_FullMethodName    = "/runtime.v1alpha1.RuntimeService/DetachNetworkInterface"
	RuntimeService_Status_FullMethodName                    = "/runtime.v1alpha1.RuntimeService/Status"
	RuntimeService_Exec_FullMethodName                      = "/runtime.v1alpha1.RuntimeService/Exec"
)

// RuntimeServiceClient is the client API for RuntimeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RuntimeServiceClient interface {
	Version(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error)
	ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error)
	CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (*CreateInstanceResponse, error)
	DeleteInstance(ctx context.Context, in *DeleteInstanceRequest, opts ...grpc.CallOption) (*DeleteInstanceResponse, error)
	UpdateInstanceAnnotations(ctx context.Context, in *UpdateInstanceAnnotationsRequest, opts ...grpc.CallOption) (*UpdateInstanceAnnotationsResponse, error)
	UpdateInstancePower(ctx context.Context, in *UpdateInstancePowerRequest, opts ...grpc.CallOption) (*UpdateInstancePowerResponse, error)
	AttachDisk(ctx context.Context, in *AttachDiskRequest, opts ...grpc.CallOption) (*AttachDiskResponse, error)
	DetachDisk(ctx context.Context, in *DetachDiskRequest, opts ...grpc.CallOption) (*DetachDiskResponse, error)
	AttachNetworkInterface(ctx context.Context, in *AttachNetworkInterfaceRequest, opts ...grpc.CallOption) (*AttachNetworkInterfaceResponse, error)
	DetachNetworkInterface(ctx context.Context, in *DetachNetworkInterfaceRequest, opts ...grpc.CallOption) (*DetachNetworkInterfaceResponse, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResponse, error)
}

type runtimeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRuntimeServiceClient(cc grpc.ClientConnInterface) RuntimeServiceClient {
	return &runtimeServiceClient{cc}
}

func (c *runtimeServiceClient) Version(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, RuntimeService_Version_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListInstancesResponse)
	err := c.cc.Invoke(ctx, RuntimeService_ListInstances_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (*CreateInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateInstanceResponse)
	err := c.cc.Invoke(ctx, RuntimeService_CreateInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) DeleteInstance(ctx context.Context, in *DeleteInstanceRequest, opts ...grpc.CallOption) (*DeleteInstanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteInstanceResponse)
	err := c.cc.Invoke(ctx, RuntimeService_DeleteInstance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) UpdateInstanceAnnotations(ctx context.Context, in *UpdateInstanceAnnotationsRequest, opts ...grpc.CallOption) (*UpdateInstanceAnnotationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateInstanceAnnotationsResponse)
	err := c.cc.Invoke(ctx, RuntimeService_UpdateInstanceAnnotations_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) UpdateInstancePower(ctx context.Context, in *UpdateInstancePowerRequest, opts ...grpc.CallOption) (*UpdateInstancePowerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateInstancePowerResponse)
	err := c.cc.Invoke(ctx, RuntimeService_UpdateInstancePower_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) AttachDisk(ctx context.Context, in *AttachDiskRequest, opts ...grpc.CallOption) (*AttachDiskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AttachDiskResponse)
	err := c.cc.Invoke(ctx, RuntimeService_AttachDisk_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) DetachDisk(ctx context.Context, in *DetachDiskRequest, opts ...grpc.CallOption) (*DetachDiskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetachDiskResponse)
	err := c.cc.Invoke(ctx, RuntimeService_DetachDisk_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) AttachNetworkInterface(ctx context.Context, in *AttachNetworkInterfaceRequest, opts ...grpc.CallOption) (*AttachNetworkInterfaceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AttachNetworkInterfaceResponse)
	err := c.cc.Invoke(ctx, RuntimeService_AttachNetworkInterface_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) DetachNetworkInterface(ctx context.Context, in *DetachNetworkInterfaceRequest, opts ...grpc.CallOption) (*DetachNetworkInterfaceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetachNetworkInterfaceResponse)
	err := c.cc.Invoke(ctx, RuntimeService_DetachNetworkInterface_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, RuntimeService_Status_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExecResponse)
	err := c.cc.Invoke(ctx, RuntimeService_Exec_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuntimeServiceServer is the server API for RuntimeService service.
// All implementations must embed UnimplementedRuntimeServiceServer
// for forward compatibility.
type RuntimeServiceServer interface {
	Version(context.Context, *VersionRequest) (*VersionResponse, error)
	ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error)
	CreateInstance(context.Context, *CreateInstanceRequest) (*CreateInstanceResponse, error)
	DeleteInstance(context.Context, *DeleteInstanceRequest) (*DeleteInstanceResponse, error)
	UpdateInstanceAnnotations(context.Context, *UpdateInstanceAnnotationsRequest) (*UpdateInstanceAnnotationsResponse, error)
	UpdateInstancePower(context.Context, *UpdateInstancePowerRequest) (*UpdateInstancePowerResponse, error)
	AttachDisk(context.Context, *AttachDiskRequest) (*AttachDiskResponse, error)
	DetachDisk(context.Context, *DetachDiskRequest) (*DetachDiskResponse, error)
	AttachNetworkInterface(context.Context, *AttachNetworkInterfaceRequest) (*AttachNetworkInterfaceResponse, error)
	DetachNetworkInterface(context.Context, *DetachNetworkInterfaceRequest) (*DetachNetworkInterfaceResponse, error)
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	Exec(context.Context, *ExecRequest) (*ExecResponse, error)
	mustEmbedUnimplementedRuntimeServiceServer()
}

// UnimplementedRuntimeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRuntimeServiceServer struct{}

func (UnimplementedRuntimeServiceServer) Version(context.Context, *VersionRequest) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Version not implemented")
}
func (UnimplementedRuntimeServiceServer) ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInstances not implemented")
}
func (UnimplementedRuntimeServiceServer) CreateInstance(context.Context, *CreateInstanceRequest) (*CreateInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInstance not implemented")
}
func (UnimplementedRuntimeServiceServer) DeleteInstance(context.Context, *DeleteInstanceRequest) (*DeleteInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteInstance not implemented")
}
func (UnimplementedRuntimeServiceServer) UpdateInstanceAnnotations(context.Context, *UpdateInstanceAnnotationsRequest) (*UpdateInstanceAnnotationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateInstanceAnnotations not implemented")
}
func (UnimplementedRuntimeServiceServer) UpdateInstancePower(context.Context, *UpdateInstancePowerRequest) (*UpdateInstancePowerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateInstancePower not implemented")
}
func (UnimplementedRuntimeServiceServer) AttachDisk(context.Context, *AttachDiskRequest) (*AttachDiskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttachDisk not implemented")
}
func (UnimplementedRuntimeServiceServer) DetachDisk(context.Context, *DetachDiskRequest) (*DetachDiskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetachDisk not implemented")
}
func (UnimplementedRuntimeServiceServer) AttachNetworkInterface(context.Context, *AttachNetworkInterfaceRequest) (*AttachNetworkInterfaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AttachNetworkInterface not implemented")
}
func (UnimplementedRuntimeServiceServer) DetachNetworkInterface(context.Context, *DetachNetworkInterfaceRequest) (*DetachNetworkInterfaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetachNetworkInterface not implemented")
}
func (UnimplementedRuntimeServiceServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedRuntimeServiceServer) Exec(context.Context, *ExecRequest) (*ExecResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exec not implemented")
}
func (UnimplementedRuntimeServiceServer) mustEmbedUnimplementedRuntimeServiceServer() {}
func (UnimplementedRuntimeServiceServer) testEmbeddedByValue()                        {}

// UnsafeRuntimeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RuntimeServiceServer will
// result in compilation errors.
type UnsafeRuntimeServiceServer interface {
	mustEmbedUnimplementedRuntimeServiceServer()
}

func RegisterRuntimeServiceServer(s grpc.ServiceRegistrar, srv RuntimeServiceServer) {
	// If the following call pancis, it indicates UnimplementedRuntimeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RuntimeService_ServiceDesc, srv)
}

func _RuntimeService_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_Version_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).Version(ctx, req.(*VersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_ListInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).ListInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_ListInstances_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).ListInstances(ctx, req.(*ListInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_CreateInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).CreateInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_CreateInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).CreateInstance(ctx, req.(*CreateInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_DeleteInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).DeleteInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_DeleteInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).DeleteInstance(ctx, req.(*DeleteInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_UpdateInstanceAnnotations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateInstanceAnnotationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).UpdateInstanceAnnotations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_UpdateInstanceAnnotations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).UpdateInstanceAnnotations(ctx, req.(*UpdateInstanceAnnotationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_UpdateInstancePower_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateInstancePowerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).UpdateInstancePower(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_UpdateInstancePower_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).UpdateInstancePower(ctx, req.(*UpdateInstancePowerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_AttachDisk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttachDiskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).AttachDisk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_AttachDisk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).AttachDisk(ctx, req.(*AttachDiskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_DetachDisk_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetachDiskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).DetachDisk(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_DetachDisk_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).DetachDisk(ctx, req.(*DetachDiskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_AttachNetworkInterface_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AttachNetworkInterfaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).AttachNetworkInterface(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_AttachNetworkInterface_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).AttachNetworkInterface(ctx, req.(*AttachNetworkInterfaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_DetachNetworkInterface_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetachNetworkInterfaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).DetachNetworkInterface(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_DetachNetworkInterface_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).DetachNetworkInterface(ctx, req.(*DetachNetworkInterfaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_Status_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_Exec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).Exec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RuntimeService_Exec_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).Exec(ctx, req.(*ExecRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RuntimeService_ServiceDesc is the grpc.ServiceDesc for RuntimeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RuntimeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "runtime.v1alpha1.RuntimeService",
	HandlerType: (*RuntimeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _RuntimeService_Version_Handler,
		},
		{
			MethodName: "ListInstances",
			Handler:    _RuntimeService_ListInstances_Handler,
		},
		{
			MethodName: "CreateInstance",
			Handler:    _RuntimeService_CreateInstance_Handler,
		},
		{
			MethodName: "DeleteInstance",
			Handler:    _RuntimeService_DeleteInstance_Handler,
		},
		{
			MethodName: "UpdateInstanceAnnotations",
			Handler:    _RuntimeService_UpdateInstanceAnnotations_Handler,
		},
		{
			MethodName: "UpdateInstancePower",
			Handler:    _RuntimeService_UpdateInstancePower_Handler,
		},
		{
			MethodName: "AttachDisk",
			Handler:    _RuntimeService_AttachDisk_Handler,
		},
		{
			MethodName: "DetachDisk",
			Handler:    _RuntimeService_DetachDisk_Handler,
		},
		{
			MethodName: "AttachNetworkInterface",
			Handler:    _RuntimeService_AttachNetworkInterface_Handler,
		},
		{
			MethodName: "DetachNetworkInterface",
			Handler:    _RuntimeService_DetachNetworkInterface_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _RuntimeService_Status_Handler,
		},
		{
			MethodName: "Exec",
			Handler:    _RuntimeService_Exec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iri-api/apis/runtime/v1alpha1/api.proto",
}