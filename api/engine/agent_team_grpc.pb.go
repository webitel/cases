// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: agent_team.proto

package engine

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
	AgentTeamService_CreateAgentTeam_FullMethodName = "/engine.AgentTeamService/CreateAgentTeam"
	AgentTeamService_SearchAgentTeam_FullMethodName = "/engine.AgentTeamService/SearchAgentTeam"
	AgentTeamService_ReadAgentTeam_FullMethodName   = "/engine.AgentTeamService/ReadAgentTeam"
	AgentTeamService_UpdateAgentTeam_FullMethodName = "/engine.AgentTeamService/UpdateAgentTeam"
	AgentTeamService_DeleteAgentTeam_FullMethodName = "/engine.AgentTeamService/DeleteAgentTeam"
)

// AgentTeamServiceClient is the client API for AgentTeamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentTeamServiceClient interface {
	// Create AgentTeam
	CreateAgentTeam(ctx context.Context, in *CreateAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error)
	// List of AgentTeam
	SearchAgentTeam(ctx context.Context, in *SearchAgentTeamRequest, opts ...grpc.CallOption) (*ListAgentTeam, error)
	// AgentTeam item
	ReadAgentTeam(ctx context.Context, in *ReadAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error)
	// Update AgentTeam
	UpdateAgentTeam(ctx context.Context, in *UpdateAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error)
	// Remove AgentTeam
	DeleteAgentTeam(ctx context.Context, in *DeleteAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error)
}

type agentTeamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentTeamServiceClient(cc grpc.ClientConnInterface) AgentTeamServiceClient {
	return &agentTeamServiceClient{cc}
}

func (c *agentTeamServiceClient) CreateAgentTeam(ctx context.Context, in *CreateAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AgentTeam)
	err := c.cc.Invoke(ctx, AgentTeamService_CreateAgentTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentTeamServiceClient) SearchAgentTeam(ctx context.Context, in *SearchAgentTeamRequest, opts ...grpc.CallOption) (*ListAgentTeam, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAgentTeam)
	err := c.cc.Invoke(ctx, AgentTeamService_SearchAgentTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentTeamServiceClient) ReadAgentTeam(ctx context.Context, in *ReadAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AgentTeam)
	err := c.cc.Invoke(ctx, AgentTeamService_ReadAgentTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentTeamServiceClient) UpdateAgentTeam(ctx context.Context, in *UpdateAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AgentTeam)
	err := c.cc.Invoke(ctx, AgentTeamService_UpdateAgentTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentTeamServiceClient) DeleteAgentTeam(ctx context.Context, in *DeleteAgentTeamRequest, opts ...grpc.CallOption) (*AgentTeam, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AgentTeam)
	err := c.cc.Invoke(ctx, AgentTeamService_DeleteAgentTeam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentTeamServiceServer is the server API for AgentTeamService service.
// All implementations must embed UnimplementedAgentTeamServiceServer
// for forward compatibility.
type AgentTeamServiceServer interface {
	// Create AgentTeam
	CreateAgentTeam(context.Context, *CreateAgentTeamRequest) (*AgentTeam, error)
	// List of AgentTeam
	SearchAgentTeam(context.Context, *SearchAgentTeamRequest) (*ListAgentTeam, error)
	// AgentTeam item
	ReadAgentTeam(context.Context, *ReadAgentTeamRequest) (*AgentTeam, error)
	// Update AgentTeam
	UpdateAgentTeam(context.Context, *UpdateAgentTeamRequest) (*AgentTeam, error)
	// Remove AgentTeam
	DeleteAgentTeam(context.Context, *DeleteAgentTeamRequest) (*AgentTeam, error)
	mustEmbedUnimplementedAgentTeamServiceServer()
}

// UnimplementedAgentTeamServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAgentTeamServiceServer struct{}

func (UnimplementedAgentTeamServiceServer) CreateAgentTeam(context.Context, *CreateAgentTeamRequest) (*AgentTeam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAgentTeam not implemented")
}
func (UnimplementedAgentTeamServiceServer) SearchAgentTeam(context.Context, *SearchAgentTeamRequest) (*ListAgentTeam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAgentTeam not implemented")
}
func (UnimplementedAgentTeamServiceServer) ReadAgentTeam(context.Context, *ReadAgentTeamRequest) (*AgentTeam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAgentTeam not implemented")
}
func (UnimplementedAgentTeamServiceServer) UpdateAgentTeam(context.Context, *UpdateAgentTeamRequest) (*AgentTeam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAgentTeam not implemented")
}
func (UnimplementedAgentTeamServiceServer) DeleteAgentTeam(context.Context, *DeleteAgentTeamRequest) (*AgentTeam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAgentTeam not implemented")
}
func (UnimplementedAgentTeamServiceServer) mustEmbedUnimplementedAgentTeamServiceServer() {}
func (UnimplementedAgentTeamServiceServer) testEmbeddedByValue()                          {}

// UnsafeAgentTeamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentTeamServiceServer will
// result in compilation errors.
type UnsafeAgentTeamServiceServer interface {
	mustEmbedUnimplementedAgentTeamServiceServer()
}

func RegisterAgentTeamServiceServer(s grpc.ServiceRegistrar, srv AgentTeamServiceServer) {
	// If the following call pancis, it indicates UnimplementedAgentTeamServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AgentTeamService_ServiceDesc, srv)
}

func _AgentTeamService_CreateAgentTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAgentTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentTeamServiceServer).CreateAgentTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AgentTeamService_CreateAgentTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentTeamServiceServer).CreateAgentTeam(ctx, req.(*CreateAgentTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentTeamService_SearchAgentTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchAgentTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentTeamServiceServer).SearchAgentTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AgentTeamService_SearchAgentTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentTeamServiceServer).SearchAgentTeam(ctx, req.(*SearchAgentTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentTeamService_ReadAgentTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadAgentTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentTeamServiceServer).ReadAgentTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AgentTeamService_ReadAgentTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentTeamServiceServer).ReadAgentTeam(ctx, req.(*ReadAgentTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentTeamService_UpdateAgentTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAgentTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentTeamServiceServer).UpdateAgentTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AgentTeamService_UpdateAgentTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentTeamServiceServer).UpdateAgentTeam(ctx, req.(*UpdateAgentTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AgentTeamService_DeleteAgentTeam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAgentTeamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentTeamServiceServer).DeleteAgentTeam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AgentTeamService_DeleteAgentTeam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentTeamServiceServer).DeleteAgentTeam(ctx, req.(*DeleteAgentTeamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AgentTeamService_ServiceDesc is the grpc.ServiceDesc for AgentTeamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AgentTeamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "engine.AgentTeamService",
	HandlerType: (*AgentTeamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAgentTeam",
			Handler:    _AgentTeamService_CreateAgentTeam_Handler,
		},
		{
			MethodName: "SearchAgentTeam",
			Handler:    _AgentTeamService_SearchAgentTeam_Handler,
		},
		{
			MethodName: "ReadAgentTeam",
			Handler:    _AgentTeamService_ReadAgentTeam_Handler,
		},
		{
			MethodName: "UpdateAgentTeam",
			Handler:    _AgentTeamService_UpdateAgentTeam_Handler,
		},
		{
			MethodName: "DeleteAgentTeam",
			Handler:    _AgentTeamService_DeleteAgentTeam_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent_team.proto",
}
