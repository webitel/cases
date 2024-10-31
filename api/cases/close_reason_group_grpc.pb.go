// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: cases/close_reason_group.proto

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
	CloseReasonGroups_ListCloseReasonGroups_FullMethodName  = "/webitel.cases.CloseReasonGroups/ListCloseReasonGroups"
	CloseReasonGroups_CreateCloseReasonGroup_FullMethodName = "/webitel.cases.CloseReasonGroups/CreateCloseReasonGroup"
	CloseReasonGroups_UpdateCloseReasonGroup_FullMethodName = "/webitel.cases.CloseReasonGroups/UpdateCloseReasonGroup"
	CloseReasonGroups_DeleteCloseReasonGroup_FullMethodName = "/webitel.cases.CloseReasonGroups/DeleteCloseReasonGroup"
	CloseReasonGroups_LocateCloseReasonGroup_FullMethodName = "/webitel.cases.CloseReasonGroups/LocateCloseReasonGroup"
)

// CloseReasonGroupsClient is the client API for CloseReasonGroups service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// CloseReasonGroups service definition with RPC methods for managing close reason groups
type CloseReasonGroupsClient interface {
	ListCloseReasonGroups(ctx context.Context, in *ListCloseReasonGroupsRequest, opts ...grpc.CallOption) (*CloseReasonGroupList, error)
	CreateCloseReasonGroup(ctx context.Context, in *CreateCloseReasonGroupRequest, opts ...grpc.CallOption) (*CloseReasonGroup, error)
	UpdateCloseReasonGroup(ctx context.Context, in *UpdateCloseReasonGroupRequest, opts ...grpc.CallOption) (*CloseReasonGroup, error)
	DeleteCloseReasonGroup(ctx context.Context, in *DeleteCloseReasonGroupRequest, opts ...grpc.CallOption) (*CloseReasonGroup, error)
	LocateCloseReasonGroup(ctx context.Context, in *LocateCloseReasonGroupRequest, opts ...grpc.CallOption) (*LocateCloseReasonGroupResponse, error)
}

type closeReasonGroupsClient struct {
	cc grpc.ClientConnInterface
}

func NewCloseReasonGroupsClient(cc grpc.ClientConnInterface) CloseReasonGroupsClient {
	return &closeReasonGroupsClient{cc}
}

func (c *closeReasonGroupsClient) ListCloseReasonGroups(ctx context.Context, in *ListCloseReasonGroupsRequest, opts ...grpc.CallOption) (*CloseReasonGroupList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReasonGroupList)
	err := c.cc.Invoke(ctx, CloseReasonGroups_ListCloseReasonGroups_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonGroupsClient) CreateCloseReasonGroup(ctx context.Context, in *CreateCloseReasonGroupRequest, opts ...grpc.CallOption) (*CloseReasonGroup, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReasonGroup)
	err := c.cc.Invoke(ctx, CloseReasonGroups_CreateCloseReasonGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonGroupsClient) UpdateCloseReasonGroup(ctx context.Context, in *UpdateCloseReasonGroupRequest, opts ...grpc.CallOption) (*CloseReasonGroup, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReasonGroup)
	err := c.cc.Invoke(ctx, CloseReasonGroups_UpdateCloseReasonGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonGroupsClient) DeleteCloseReasonGroup(ctx context.Context, in *DeleteCloseReasonGroupRequest, opts ...grpc.CallOption) (*CloseReasonGroup, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CloseReasonGroup)
	err := c.cc.Invoke(ctx, CloseReasonGroups_DeleteCloseReasonGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *closeReasonGroupsClient) LocateCloseReasonGroup(ctx context.Context, in *LocateCloseReasonGroupRequest, opts ...grpc.CallOption) (*LocateCloseReasonGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LocateCloseReasonGroupResponse)
	err := c.cc.Invoke(ctx, CloseReasonGroups_LocateCloseReasonGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloseReasonGroupsServer is the server API for CloseReasonGroups service.
// All implementations should embed UnimplementedCloseReasonGroupsServer
// for forward compatibility.
//
// CloseReasonGroups service definition with RPC methods for managing close reason groups
type CloseReasonGroupsServer interface {
	ListCloseReasonGroups(context.Context, *ListCloseReasonGroupsRequest) (*CloseReasonGroupList, error)
	CreateCloseReasonGroup(context.Context, *CreateCloseReasonGroupRequest) (*CloseReasonGroup, error)
	UpdateCloseReasonGroup(context.Context, *UpdateCloseReasonGroupRequest) (*CloseReasonGroup, error)
	DeleteCloseReasonGroup(context.Context, *DeleteCloseReasonGroupRequest) (*CloseReasonGroup, error)
	LocateCloseReasonGroup(context.Context, *LocateCloseReasonGroupRequest) (*LocateCloseReasonGroupResponse, error)
}

// UnimplementedCloseReasonGroupsServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCloseReasonGroupsServer struct{}

func (UnimplementedCloseReasonGroupsServer) ListCloseReasonGroups(context.Context, *ListCloseReasonGroupsRequest) (*CloseReasonGroupList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCloseReasonGroups not implemented")
}
func (UnimplementedCloseReasonGroupsServer) CreateCloseReasonGroup(context.Context, *CreateCloseReasonGroupRequest) (*CloseReasonGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCloseReasonGroup not implemented")
}
func (UnimplementedCloseReasonGroupsServer) UpdateCloseReasonGroup(context.Context, *UpdateCloseReasonGroupRequest) (*CloseReasonGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCloseReasonGroup not implemented")
}
func (UnimplementedCloseReasonGroupsServer) DeleteCloseReasonGroup(context.Context, *DeleteCloseReasonGroupRequest) (*CloseReasonGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCloseReasonGroup not implemented")
}
func (UnimplementedCloseReasonGroupsServer) LocateCloseReasonGroup(context.Context, *LocateCloseReasonGroupRequest) (*LocateCloseReasonGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateCloseReasonGroup not implemented")
}
func (UnimplementedCloseReasonGroupsServer) testEmbeddedByValue() {}

// UnsafeCloseReasonGroupsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloseReasonGroupsServer will
// result in compilation errors.
type UnsafeCloseReasonGroupsServer interface {
	mustEmbedUnimplementedCloseReasonGroupsServer()
}

func RegisterCloseReasonGroupsServer(s grpc.ServiceRegistrar, srv CloseReasonGroupsServer) {
	// If the following call pancis, it indicates UnimplementedCloseReasonGroupsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CloseReasonGroups_ServiceDesc, srv)
}

func _CloseReasonGroups_ListCloseReasonGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCloseReasonGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonGroupsServer).ListCloseReasonGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasonGroups_ListCloseReasonGroups_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonGroupsServer).ListCloseReasonGroups(ctx, req.(*ListCloseReasonGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasonGroups_CreateCloseReasonGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCloseReasonGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonGroupsServer).CreateCloseReasonGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasonGroups_CreateCloseReasonGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonGroupsServer).CreateCloseReasonGroup(ctx, req.(*CreateCloseReasonGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasonGroups_UpdateCloseReasonGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCloseReasonGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonGroupsServer).UpdateCloseReasonGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasonGroups_UpdateCloseReasonGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonGroupsServer).UpdateCloseReasonGroup(ctx, req.(*UpdateCloseReasonGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasonGroups_DeleteCloseReasonGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCloseReasonGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonGroupsServer).DeleteCloseReasonGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasonGroups_DeleteCloseReasonGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonGroupsServer).DeleteCloseReasonGroup(ctx, req.(*DeleteCloseReasonGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloseReasonGroups_LocateCloseReasonGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateCloseReasonGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloseReasonGroupsServer).LocateCloseReasonGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloseReasonGroups_LocateCloseReasonGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloseReasonGroupsServer).LocateCloseReasonGroup(ctx, req.(*LocateCloseReasonGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CloseReasonGroups_ServiceDesc is the grpc.ServiceDesc for CloseReasonGroups service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloseReasonGroups_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.CloseReasonGroups",
	HandlerType: (*CloseReasonGroupsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCloseReasonGroups",
			Handler:    _CloseReasonGroups_ListCloseReasonGroups_Handler,
		},
		{
			MethodName: "CreateCloseReasonGroup",
			Handler:    _CloseReasonGroups_CreateCloseReasonGroup_Handler,
		},
		{
			MethodName: "UpdateCloseReasonGroup",
			Handler:    _CloseReasonGroups_UpdateCloseReasonGroup_Handler,
		},
		{
			MethodName: "DeleteCloseReasonGroup",
			Handler:    _CloseReasonGroups_DeleteCloseReasonGroup_Handler,
		},
		{
			MethodName: "LocateCloseReasonGroup",
			Handler:    _CloseReasonGroups_LocateCloseReasonGroup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cases/close_reason_group.proto",
}