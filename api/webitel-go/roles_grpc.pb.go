// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: roles.proto

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
	Roles_ReadRole_FullMethodName                   = "/api.Roles/ReadRole"
	Roles_UpdateRole_FullMethodName                 = "/api.Roles/UpdateRole"
	Roles_DeleteRole_FullMethodName                 = "/api.Roles/DeleteRole"
	Roles_CreateRole_FullMethodName                 = "/api.Roles/CreateRole"
	Roles_SearchRoles_FullMethodName                = "/api.Roles/SearchRoles"
	Roles_SearchRecordAvailableRoles_FullMethodName = "/api.Roles/SearchRecordAvailableRoles"
	Roles_GrantRole_FullMethodName                  = "/api.Roles/GrantRole"
	Roles_RevokeRole_FullMethodName                 = "/api.Roles/RevokeRole"
	Roles_RoleMembers_FullMethodName                = "/api.Roles/RoleMembers"
	Roles_RoleMetadata_FullMethodName               = "/api.Roles/RoleMetadata"
)

// RolesClient is the client API for Roles service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RolesClient interface {
	ReadRole(ctx context.Context, in *ReadRoleRequest, opts ...grpc.CallOption) (*ReadRoleResponse, error)
	UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...grpc.CallOption) (*UpdateRoleResponse, error)
	DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error)
	CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error)
	SearchRoles(ctx context.Context, in *SearchRolesRequest, opts ...grpc.CallOption) (*SearchRolesResponse, error)
	SearchRecordAvailableRoles(ctx context.Context, in *SearchRecordAvailableRolesRequest, opts ...grpc.CallOption) (*SearchRolesResponse, error)
	GrantRole(ctx context.Context, in *GrantRoleRequest, opts ...grpc.CallOption) (*GrantRoleResponse, error)
	RevokeRole(ctx context.Context, in *RevokeRoleRequest, opts ...grpc.CallOption) (*RevokeRoleResponse, error)
	RoleMembers(ctx context.Context, in *RoleMembersRequest, opts ...grpc.CallOption) (*RoleMembersResponse, error)
	RoleMetadata(ctx context.Context, in *RoleMetadataRequest, opts ...grpc.CallOption) (*RoleMetadataResponse, error)
}

type rolesClient struct {
	cc grpc.ClientConnInterface
}

func NewRolesClient(cc grpc.ClientConnInterface) RolesClient {
	return &rolesClient{cc}
}

func (c *rolesClient) ReadRole(ctx context.Context, in *ReadRoleRequest, opts ...grpc.CallOption) (*ReadRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadRoleResponse)
	err := c.cc.Invoke(ctx, Roles_ReadRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...grpc.CallOption) (*UpdateRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateRoleResponse)
	err := c.cc.Invoke(ctx, Roles_UpdateRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteRoleResponse)
	err := c.cc.Invoke(ctx, Roles_DeleteRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...grpc.CallOption) (*CreateRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateRoleResponse)
	err := c.cc.Invoke(ctx, Roles_CreateRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) SearchRoles(ctx context.Context, in *SearchRolesRequest, opts ...grpc.CallOption) (*SearchRolesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchRolesResponse)
	err := c.cc.Invoke(ctx, Roles_SearchRoles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) SearchRecordAvailableRoles(ctx context.Context, in *SearchRecordAvailableRolesRequest, opts ...grpc.CallOption) (*SearchRolesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchRolesResponse)
	err := c.cc.Invoke(ctx, Roles_SearchRecordAvailableRoles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) GrantRole(ctx context.Context, in *GrantRoleRequest, opts ...grpc.CallOption) (*GrantRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GrantRoleResponse)
	err := c.cc.Invoke(ctx, Roles_GrantRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) RevokeRole(ctx context.Context, in *RevokeRoleRequest, opts ...grpc.CallOption) (*RevokeRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RevokeRoleResponse)
	err := c.cc.Invoke(ctx, Roles_RevokeRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) RoleMembers(ctx context.Context, in *RoleMembersRequest, opts ...grpc.CallOption) (*RoleMembersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RoleMembersResponse)
	err := c.cc.Invoke(ctx, Roles_RoleMembers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rolesClient) RoleMetadata(ctx context.Context, in *RoleMetadataRequest, opts ...grpc.CallOption) (*RoleMetadataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RoleMetadataResponse)
	err := c.cc.Invoke(ctx, Roles_RoleMetadata_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RolesServer is the server API for Roles service.
// All implementations must embed UnimplementedRolesServer
// for forward compatibility.
type RolesServer interface {
	ReadRole(context.Context, *ReadRoleRequest) (*ReadRoleResponse, error)
	UpdateRole(context.Context, *UpdateRoleRequest) (*UpdateRoleResponse, error)
	DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error)
	CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error)
	SearchRoles(context.Context, *SearchRolesRequest) (*SearchRolesResponse, error)
	SearchRecordAvailableRoles(context.Context, *SearchRecordAvailableRolesRequest) (*SearchRolesResponse, error)
	GrantRole(context.Context, *GrantRoleRequest) (*GrantRoleResponse, error)
	RevokeRole(context.Context, *RevokeRoleRequest) (*RevokeRoleResponse, error)
	RoleMembers(context.Context, *RoleMembersRequest) (*RoleMembersResponse, error)
	RoleMetadata(context.Context, *RoleMetadataRequest) (*RoleMetadataResponse, error)
	mustEmbedUnimplementedRolesServer()
}

// UnimplementedRolesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRolesServer struct{}

func (UnimplementedRolesServer) ReadRole(context.Context, *ReadRoleRequest) (*ReadRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadRole not implemented")
}
func (UnimplementedRolesServer) UpdateRole(context.Context, *UpdateRoleRequest) (*UpdateRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRole not implemented")
}
func (UnimplementedRolesServer) DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}
func (UnimplementedRolesServer) CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRole not implemented")
}
func (UnimplementedRolesServer) SearchRoles(context.Context, *SearchRolesRequest) (*SearchRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchRoles not implemented")
}
func (UnimplementedRolesServer) SearchRecordAvailableRoles(context.Context, *SearchRecordAvailableRolesRequest) (*SearchRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchRecordAvailableRoles not implemented")
}
func (UnimplementedRolesServer) GrantRole(context.Context, *GrantRoleRequest) (*GrantRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GrantRole not implemented")
}
func (UnimplementedRolesServer) RevokeRole(context.Context, *RevokeRoleRequest) (*RevokeRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeRole not implemented")
}
func (UnimplementedRolesServer) RoleMembers(context.Context, *RoleMembersRequest) (*RoleMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RoleMembers not implemented")
}
func (UnimplementedRolesServer) RoleMetadata(context.Context, *RoleMetadataRequest) (*RoleMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RoleMetadata not implemented")
}
func (UnimplementedRolesServer) mustEmbedUnimplementedRolesServer() {}
func (UnimplementedRolesServer) testEmbeddedByValue()               {}

// UnsafeRolesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RolesServer will
// result in compilation errors.
type UnsafeRolesServer interface {
	mustEmbedUnimplementedRolesServer()
}

func RegisterRolesServer(s grpc.ServiceRegistrar, srv RolesServer) {
	// If the following call pancis, it indicates UnimplementedRolesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Roles_ServiceDesc, srv)
}

func _Roles_ReadRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).ReadRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_ReadRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).ReadRole(ctx, req.(*ReadRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_UpdateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).UpdateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_UpdateRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).UpdateRole(ctx, req.(*UpdateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_DeleteRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).DeleteRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_DeleteRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).DeleteRole(ctx, req.(*DeleteRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_CreateRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).CreateRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_CreateRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).CreateRole(ctx, req.(*CreateRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_SearchRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).SearchRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_SearchRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).SearchRoles(ctx, req.(*SearchRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_SearchRecordAvailableRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchRecordAvailableRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).SearchRecordAvailableRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_SearchRecordAvailableRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).SearchRecordAvailableRoles(ctx, req.(*SearchRecordAvailableRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_GrantRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrantRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).GrantRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_GrantRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).GrantRole(ctx, req.(*GrantRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_RevokeRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RevokeRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).RevokeRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_RevokeRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).RevokeRole(ctx, req.(*RevokeRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_RoleMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).RoleMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_RoleMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).RoleMembers(ctx, req.(*RoleMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Roles_RoleMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RolesServer).RoleMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Roles_RoleMetadata_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RolesServer).RoleMetadata(ctx, req.(*RoleMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Roles_ServiceDesc is the grpc.ServiceDesc for Roles service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Roles_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Roles",
	HandlerType: (*RolesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReadRole",
			Handler:    _Roles_ReadRole_Handler,
		},
		{
			MethodName: "UpdateRole",
			Handler:    _Roles_UpdateRole_Handler,
		},
		{
			MethodName: "DeleteRole",
			Handler:    _Roles_DeleteRole_Handler,
		},
		{
			MethodName: "CreateRole",
			Handler:    _Roles_CreateRole_Handler,
		},
		{
			MethodName: "SearchRoles",
			Handler:    _Roles_SearchRoles_Handler,
		},
		{
			MethodName: "SearchRecordAvailableRoles",
			Handler:    _Roles_SearchRecordAvailableRoles_Handler,
		},
		{
			MethodName: "GrantRole",
			Handler:    _Roles_GrantRole_Handler,
		},
		{
			MethodName: "RevokeRole",
			Handler:    _Roles_RevokeRole_Handler,
		},
		{
			MethodName: "RoleMembers",
			Handler:    _Roles_RoleMembers_Handler,
		},
		{
			MethodName: "RoleMetadata",
			Handler:    _Roles_RoleMetadata_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "roles.proto",
}
