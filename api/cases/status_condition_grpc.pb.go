// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: status_condition.proto

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
	StatusConditions_ListStatusConditions_FullMethodName  = "/webitel.cases.StatusConditions/ListStatusConditions"
	StatusConditions_CreateStatusCondition_FullMethodName = "/webitel.cases.StatusConditions/CreateStatusCondition"
	StatusConditions_UpdateStatusCondition_FullMethodName = "/webitel.cases.StatusConditions/UpdateStatusCondition"
	StatusConditions_DeleteStatusCondition_FullMethodName = "/webitel.cases.StatusConditions/DeleteStatusCondition"
	StatusConditions_LocateStatusCondition_FullMethodName = "/webitel.cases.StatusConditions/LocateStatusCondition"
)

// StatusConditionsClient is the client API for StatusConditions service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// StatusConditions service definition with RPC methods for managing statuses
type StatusConditionsClient interface {
	// RPC method to list or search statuses
	ListStatusConditions(ctx context.Context, in *ListStatusConditionRequest, opts ...grpc.CallOption) (*StatusConditionList, error)
	// RPC method to create a new status condition
	CreateStatusCondition(ctx context.Context, in *CreateStatusConditionRequest, opts ...grpc.CallOption) (*StatusCondition, error)
	// RPC method to update an existing status condition
	UpdateStatusCondition(ctx context.Context, in *UpdateStatusConditionRequest, opts ...grpc.CallOption) (*StatusCondition, error)
	// RPC method to delete an existing status condition
	DeleteStatusCondition(ctx context.Context, in *DeleteStatusConditionRequest, opts ...grpc.CallOption) (*StatusCondition, error)
	// RPC method to locate a specific status condition by ID
	LocateStatusCondition(ctx context.Context, in *LocateStatusConditionRequest, opts ...grpc.CallOption) (*LocateStatusConditionResponse, error)
}

type statusConditionsClient struct {
	cc grpc.ClientConnInterface
}

func NewStatusConditionsClient(cc grpc.ClientConnInterface) StatusConditionsClient {
	return &statusConditionsClient{cc}
}

func (c *statusConditionsClient) ListStatusConditions(ctx context.Context, in *ListStatusConditionRequest, opts ...grpc.CallOption) (*StatusConditionList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusConditionList)
	err := c.cc.Invoke(ctx, StatusConditions_ListStatusConditions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusConditionsClient) CreateStatusCondition(ctx context.Context, in *CreateStatusConditionRequest, opts ...grpc.CallOption) (*StatusCondition, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusCondition)
	err := c.cc.Invoke(ctx, StatusConditions_CreateStatusCondition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusConditionsClient) UpdateStatusCondition(ctx context.Context, in *UpdateStatusConditionRequest, opts ...grpc.CallOption) (*StatusCondition, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusCondition)
	err := c.cc.Invoke(ctx, StatusConditions_UpdateStatusCondition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusConditionsClient) DeleteStatusCondition(ctx context.Context, in *DeleteStatusConditionRequest, opts ...grpc.CallOption) (*StatusCondition, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatusCondition)
	err := c.cc.Invoke(ctx, StatusConditions_DeleteStatusCondition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statusConditionsClient) LocateStatusCondition(ctx context.Context, in *LocateStatusConditionRequest, opts ...grpc.CallOption) (*LocateStatusConditionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LocateStatusConditionResponse)
	err := c.cc.Invoke(ctx, StatusConditions_LocateStatusCondition_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatusConditionsServer is the server API for StatusConditions service.
// All implementations must embed UnimplementedStatusConditionsServer
// for forward compatibility.
//
// StatusConditions service definition with RPC methods for managing statuses
type StatusConditionsServer interface {
	// RPC method to list or search statuses
	ListStatusConditions(context.Context, *ListStatusConditionRequest) (*StatusConditionList, error)
	// RPC method to create a new status condition
	CreateStatusCondition(context.Context, *CreateStatusConditionRequest) (*StatusCondition, error)
	// RPC method to update an existing status condition
	UpdateStatusCondition(context.Context, *UpdateStatusConditionRequest) (*StatusCondition, error)
	// RPC method to delete an existing status condition
	DeleteStatusCondition(context.Context, *DeleteStatusConditionRequest) (*StatusCondition, error)
	// RPC method to locate a specific status condition by ID
	LocateStatusCondition(context.Context, *LocateStatusConditionRequest) (*LocateStatusConditionResponse, error)
	mustEmbedUnimplementedStatusConditionsServer()
}

// UnimplementedStatusConditionsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStatusConditionsServer struct{}

func (UnimplementedStatusConditionsServer) ListStatusConditions(context.Context, *ListStatusConditionRequest) (*StatusConditionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStatusConditions not implemented")
}
func (UnimplementedStatusConditionsServer) CreateStatusCondition(context.Context, *CreateStatusConditionRequest) (*StatusCondition, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStatusCondition not implemented")
}
func (UnimplementedStatusConditionsServer) UpdateStatusCondition(context.Context, *UpdateStatusConditionRequest) (*StatusCondition, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatusCondition not implemented")
}
func (UnimplementedStatusConditionsServer) DeleteStatusCondition(context.Context, *DeleteStatusConditionRequest) (*StatusCondition, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStatusCondition not implemented")
}
func (UnimplementedStatusConditionsServer) LocateStatusCondition(context.Context, *LocateStatusConditionRequest) (*LocateStatusConditionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateStatusCondition not implemented")
}
func (UnimplementedStatusConditionsServer) mustEmbedUnimplementedStatusConditionsServer() {}
func (UnimplementedStatusConditionsServer) testEmbeddedByValue()                          {}

// UnsafeStatusConditionsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatusConditionsServer will
// result in compilation errors.
type UnsafeStatusConditionsServer interface {
	mustEmbedUnimplementedStatusConditionsServer()
}

func RegisterStatusConditionsServer(s grpc.ServiceRegistrar, srv StatusConditionsServer) {
	// If the following call pancis, it indicates UnimplementedStatusConditionsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StatusConditions_ServiceDesc, srv)
}

func _StatusConditions_ListStatusConditions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatusConditionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusConditionsServer).ListStatusConditions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusConditions_ListStatusConditions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusConditionsServer).ListStatusConditions(ctx, req.(*ListStatusConditionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusConditions_CreateStatusCondition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStatusConditionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusConditionsServer).CreateStatusCondition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusConditions_CreateStatusCondition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusConditionsServer).CreateStatusCondition(ctx, req.(*CreateStatusConditionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusConditions_UpdateStatusCondition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusConditionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusConditionsServer).UpdateStatusCondition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusConditions_UpdateStatusCondition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusConditionsServer).UpdateStatusCondition(ctx, req.(*UpdateStatusConditionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusConditions_DeleteStatusCondition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteStatusConditionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusConditionsServer).DeleteStatusCondition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusConditions_DeleteStatusCondition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusConditionsServer).DeleteStatusCondition(ctx, req.(*DeleteStatusConditionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StatusConditions_LocateStatusCondition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateStatusConditionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatusConditionsServer).LocateStatusCondition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StatusConditions_LocateStatusCondition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatusConditionsServer).LocateStatusCondition(ctx, req.(*LocateStatusConditionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StatusConditions_ServiceDesc is the grpc.ServiceDesc for StatusConditions service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StatusConditions_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.StatusConditions",
	HandlerType: (*StatusConditionsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListStatusConditions",
			Handler:    _StatusConditions_ListStatusConditions_Handler,
		},
		{
			MethodName: "CreateStatusCondition",
			Handler:    _StatusConditions_CreateStatusCondition_Handler,
		},
		{
			MethodName: "UpdateStatusCondition",
			Handler:    _StatusConditions_UpdateStatusCondition_Handler,
		},
		{
			MethodName: "DeleteStatusCondition",
			Handler:    _StatusConditions_DeleteStatusCondition_Handler,
		},
		{
			MethodName: "LocateStatusCondition",
			Handler:    _StatusConditions_LocateStatusCondition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "status_condition.proto",
}
