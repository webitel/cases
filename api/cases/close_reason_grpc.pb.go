// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: cases/close_reason.proto

package api

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
	CloseReasons_ListCloseReasons_FullMethodName  = "/cases.CloseReasons/ListCloseReasons"
	CloseReasons_CreateCloseReason_FullMethodName = "/cases.CloseReasons/CreateCloseReason"
	CloseReasons_UpdateCloseReason_FullMethodName = "/cases.CloseReasons/UpdateCloseReason"
	CloseReasons_DeleteCloseReason_FullMethodName = "/cases.CloseReasons/DeleteCloseReason"
	CloseReasons_LocateCloseReason_FullMethodName = "/cases.CloseReasons/LocateCloseReason"
)

// CloseReasonsClient is the client API for CloseReasons service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// CloseReasons service definition with RPC methods for managing close reasons
type CloseReasonsClient interface {
	// RPC method to list or search close reasons
	ListCloseReasons(ctx context.Context, in *ListCloseReasonRequest, opts ...grpc.CallOption) (*CloseReasonList, error)
	// RPC method to create a new close reason
	CreateCloseReason(ctx context.Context, in *CreateCloseReasonRequest, opts ...grpc.CallOption) (*CloseReason, error)
	// RPC method to update an existing close reason
	UpdateCloseReason(ctx context.Context, in *UpdateCloseReasonRequest, opts ...grpc.CallOption) (*CloseReason, error)
	// RPC method to delete an existing close reason
	DeleteCloseReason(ctx context.Context, in *DeleteCloseReasonRequest, opts ...grpc.CallOption) (*CloseReason, error)
	// RPC method to locate a specific close reason by ID
	LocateCloseReason(ctx context.Context, in *LocateCloseReasonRequest, opts ...grpc.CallOption) (*LocateCloseReasonResponse, error)
}

type closeReasonsClient struct {
	cc grpc.ClientConnInterface
}

func NewCloseReasonsClient(cc grpc.ClientConnInterface) CloseReasonsClient {
	return &closeReasonsClient{cc}
}

func (c *closeReasonsClient) ListCloseReasons(ctx context.Context, in *ListCloseReasonRequest, opts ...grpc.CallOption) (*CloseReasonList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReasonList)
	err := c.cc.Invoke(ctx, CloseReasons_ListCloseReasons_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonsClient) CreateCloseReason(ctx context.Context, in *CreateCloseReasonRequest, opts ...grpc.CallOption) (*CloseReason, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReason)
	err := c.cc.Invoke(ctx, CloseReasons_CreateCloseReason_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonsClient) UpdateCloseReason(ctx context.Context, in *UpdateCloseReasonRequest, opts ...grpc.CallOption) (*CloseReason, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReason)
	err := c.cc.Invoke(ctx, CloseReasons_UpdateCloseReason_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonsClient) DeleteCloseReason(ctx context.Context, in *DeleteCloseReasonRequest, opts ...grpc.CallOption) (*CloseReason, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReason)
	err := c.cc.Invoke(ctx, CloseReasons_DeleteCloseReason_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonsClient) LocateCloseReason(ctx context.Context, in *LocateCloseReasonRequest, opts ...grpc.CallOption) (*LocateCloseReasonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LocateCloseReasonResponse)
	err := c.cc.Invoke(ctx, CloseReasons_LocateCloseReason_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloseReasonsServer is the server API for CloseReasons service.
// All implementations should embed UnimplementedCloseReasonsServer
// for forward compatibility.
//
// CloseReasons service definition with RPC methods for managing close reasons
type CloseReasonsServer interface {
	// RPC method to list or search close reasons
	ListCloseReasons(context.Context, *ListCloseReasonRequest) (*CloseReasonList, error)
	// RPC method to create a new close reason
	CreateCloseReason(context.Context, *CreateCloseReasonRequest) (*CloseReason, error)
	// RPC method to update an existing close reason
	UpdateCloseReason(context.Context, *UpdateCloseReasonRequest) (*CloseReason, error)
	// RPC method to delete an existing close reason
	DeleteCloseReason(context.Context, *DeleteCloseReasonRequest) (*CloseReason, error)
	// RPC method to locate a specific close reason by ID
	LocateCloseReason(context.Context, *LocateCloseReasonRequest) (*LocateCloseReasonResponse, error)
}

// UnimplementedCloseReasonsServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCloseReasonsServer struct{}

func (UnimplementedCloseReasonsServer) ListCloseReasons(context.Context, *ListCloseReasonRequest) (*CloseReasonList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCloseReasons not implemented")
}
func (UnimplementedCloseReasonsServer) CreateCloseReason(context.Context, *CreateCloseReasonRequest) (*CloseReason, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCloseReason not implemented")
}
func (UnimplementedCloseReasonsServer) UpdateCloseReason(context.Context, *UpdateCloseReasonRequest) (*CloseReason, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCloseReason not implemented")
}
func (UnimplementedCloseReasonsServer) DeleteCloseReason(context.Context, *DeleteCloseReasonRequest) (*CloseReason, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCloseReason not implemented")
}
func (UnimplementedCloseReasonsServer) LocateCloseReason(context.Context, *LocateCloseReasonRequest) (*LocateCloseReasonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateCloseReason not implemented")
}
func (UnimplementedCloseReasonsServer) testEmbeddedByValue() {}

// UnsafeCloseReasonsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloseReasonsServer will
// result in compilation errors.
type UnsafeCloseReasonsServer interface {
	mustEmbedUnimplementedCloseReasonsServer()
}

func RegisterCloseReasonsServer(s grpc.ServiceRegistrar, srv CloseReasonsServer) {
	// If the following call pancis, it indicates UnimplementedCloseReasonsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CloseReasons_ServiceDesc, srv)
}

func _CloseReasons_ListCloseReasons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCloseReasonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonsServer).ListCloseReasons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasons_ListCloseReasons_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonsServer).ListCloseReasons(ctx, req.(*ListCloseReasonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasons_CreateCloseReason_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCloseReasonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonsServer).CreateCloseReason(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasons_CreateCloseReason_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonsServer).CreateCloseReason(ctx, req.(*CreateCloseReasonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasons_UpdateCloseReason_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCloseReasonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonsServer).UpdateCloseReason(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasons_UpdateCloseReason_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonsServer).UpdateCloseReason(ctx, req.(*UpdateCloseReasonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasons_DeleteCloseReason_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCloseReasonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonsServer).DeleteCloseReason(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasons_DeleteCloseReason_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonsServer).DeleteCloseReason(ctx, req.(*DeleteCloseReasonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasons_LocateCloseReason_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateCloseReasonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonsServer).LocateCloseReason(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasons_LocateCloseReason_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonsServer).LocateCloseReason(ctx, req.(*LocateCloseReasonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CloseReasons_ServiceDesc is the grpc.ServiceDesc for CloseReasons service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloseReasons_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cases.CloseReasons",
	HandlerType: (*CloseReasonsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCloseReasons",
			Handler:    _CloseReasons_ListCloseReasons_Handler,
		},
		{
			MethodName: "CreateCloseReason",
			Handler:    _CloseReasons_CreateCloseReason_Handler,
		},
		{
			MethodName: "UpdateCloseReason",
			Handler:    _CloseReasons_UpdateCloseReason_Handler,
		},
		{
			MethodName: "DeleteCloseReason",
			Handler:    _CloseReasons_DeleteCloseReason_Handler,
		},
		{
			MethodName: "LocateCloseReason",
			Handler:    _CloseReasons_LocateCloseReason_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cases/close_reason.proto",
}
