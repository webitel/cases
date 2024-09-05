// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: cases/sla.proto

package cases

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
	SLAs_ListSLAs_FullMethodName  = "/webitel.cases.SLAs/ListSLAs"
	SLAs_CreateSLA_FullMethodName = "/webitel.cases.SLAs/CreateSLA"
	SLAs_UpdateSLA_FullMethodName = "/webitel.cases.SLAs/UpdateSLA"
	SLAs_DeleteSLA_FullMethodName = "/webitel.cases.SLAs/DeleteSLA"
	SLAs_LocateSLA_FullMethodName = "/webitel.cases.SLAs/LocateSLA"
)

// SLAsClient is the client API for SLAs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// SLAs service definition with RPC methods for managing SLAs
type SLAsClient interface {
	// RPC method to list or search SLAs
	ListSLAs(ctx context.Context, in *ListSLARequest, opts ...grpc.CallOption) (*SLAList, error)
	// RPC method to create a new SLA
	CreateSLA(ctx context.Context, in *CreateSLARequest, opts ...grpc.CallOption) (*SLA, error)
	// RPC method to update an existing SLA
	UpdateSLA(ctx context.Context, in *UpdateSLARequest, opts ...grpc.CallOption) (*SLA, error)
	// RPC method to delete an existing SLA
	DeleteSLA(ctx context.Context, in *DeleteSLARequest, opts ...grpc.CallOption) (*SLA, error)
	// RPC method to locate a specific SLA by ID
	LocateSLA(ctx context.Context, in *LocateSLARequest, opts ...grpc.CallOption) (*LocateSLAResponse, error)
}

type sLAsClient struct {
	cc grpc.ClientConnInterface
}

func NewSLAsClient(cc grpc.ClientConnInterface) SLAsClient {
	return &sLAsClient{cc}
}

func (c *sLAsClient) ListSLAs(ctx context.Context, in *ListSLARequest, opts ...grpc.CallOption) (*SLAList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SLAList)
	err := c.cc.Invoke(ctx, SLAs_ListSLAs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLAsClient) CreateSLA(ctx context.Context, in *CreateSLARequest, opts ...grpc.CallOption) (*SLA, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SLA)
	err := c.cc.Invoke(ctx, SLAs_CreateSLA_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLAsClient) UpdateSLA(ctx context.Context, in *UpdateSLARequest, opts ...grpc.CallOption) (*SLA, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SLA)
	err := c.cc.Invoke(ctx, SLAs_UpdateSLA_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLAsClient) DeleteSLA(ctx context.Context, in *DeleteSLARequest, opts ...grpc.CallOption) (*SLA, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SLA)
	err := c.cc.Invoke(ctx, SLAs_DeleteSLA_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sLAsClient) LocateSLA(ctx context.Context, in *LocateSLARequest, opts ...grpc.CallOption) (*LocateSLAResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LocateSLAResponse)
	err := c.cc.Invoke(ctx, SLAs_LocateSLA_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SLAsServer is the server API for SLAs service.
// All implementations should embed UnimplementedSLAsServer
// for forward compatibility.
//
// SLAs service definition with RPC methods for managing SLAs
type SLAsServer interface {
	// RPC method to list or search SLAs
	ListSLAs(context.Context, *ListSLARequest) (*SLAList, error)
	// RPC method to create a new SLA
	CreateSLA(context.Context, *CreateSLARequest) (*SLA, error)
	// RPC method to update an existing SLA
	UpdateSLA(context.Context, *UpdateSLARequest) (*SLA, error)
	// RPC method to delete an existing SLA
	DeleteSLA(context.Context, *DeleteSLARequest) (*SLA, error)
	// RPC method to locate a specific SLA by ID
	LocateSLA(context.Context, *LocateSLARequest) (*LocateSLAResponse, error)
}

// UnimplementedSLAsServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSLAsServer struct{}

func (UnimplementedSLAsServer) ListSLAs(context.Context, *ListSLARequest) (*SLAList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSLAs not implemented")
}
func (UnimplementedSLAsServer) CreateSLA(context.Context, *CreateSLARequest) (*SLA, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSLA not implemented")
}
func (UnimplementedSLAsServer) UpdateSLA(context.Context, *UpdateSLARequest) (*SLA, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSLA not implemented")
}
func (UnimplementedSLAsServer) DeleteSLA(context.Context, *DeleteSLARequest) (*SLA, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSLA not implemented")
}
func (UnimplementedSLAsServer) LocateSLA(context.Context, *LocateSLARequest) (*LocateSLAResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateSLA not implemented")
}
func (UnimplementedSLAsServer) testEmbeddedByValue() {}

// UnsafeSLAsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SLAsServer will
// result in compilation errors.
type UnsafeSLAsServer interface {
	mustEmbedUnimplementedSLAsServer()
}

func RegisterSLAsServer(s grpc.ServiceRegistrar, srv SLAsServer) {
	// If the following call pancis, it indicates UnimplementedSLAsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SLAs_ServiceDesc, srv)
}

func _SLAs_ListSLAs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSLARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLAsServer).ListSLAs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLAs_ListSLAs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLAsServer).ListSLAs(ctx, req.(*ListSLARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLAs_CreateSLA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSLARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLAsServer).CreateSLA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLAs_CreateSLA_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLAsServer).CreateSLA(ctx, req.(*CreateSLARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLAs_UpdateSLA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSLARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLAsServer).UpdateSLA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLAs_UpdateSLA_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLAsServer).UpdateSLA(ctx, req.(*UpdateSLARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLAs_DeleteSLA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSLARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLAsServer).DeleteSLA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLAs_DeleteSLA_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLAsServer).DeleteSLA(ctx, req.(*DeleteSLARequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SLAs_LocateSLA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateSLARequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SLAsServer).LocateSLA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SLAs_LocateSLA_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SLAsServer).LocateSLA(ctx, req.(*LocateSLARequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SLAs_ServiceDesc is the grpc.ServiceDesc for SLAs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SLAs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.SLAs",
	HandlerType: (*SLAsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSLAs",
			Handler:    _SLAs_ListSLAs_Handler,
		},
		{
			MethodName: "CreateSLA",
			Handler:    _SLAs_CreateSLA_Handler,
		},
		{
			MethodName: "UpdateSLA",
			Handler:    _SLAs_UpdateSLA_Handler,
		},
		{
			MethodName: "DeleteSLA",
			Handler:    _SLAs_DeleteSLA_Handler,
		},
		{
			MethodName: "LocateSLA",
			Handler:    _SLAs_LocateSLA_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cases/sla.proto",
}