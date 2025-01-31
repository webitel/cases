// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: user_settings.proto

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
	UserSettings_GetUserSettings_FullMethodName = "/api.UserSettings/GetUserSettings"
	UserSettings_SetUserSettings_FullMethodName = "/api.UserSettings/SetUserSettings"
)

// UserSettingsClient is the client API for UserSettings service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserSettingsClient interface {
	// GET /user/settings/{key=*}
	GetUserSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*UserSetting, error)
	// Create -or- Update
	//
	// PUT /user/settings/{key=*}
	// Content-Type: application/json; charset=utf-8
	//
	// ${value=.body}
	SetUserSettings(ctx context.Context, in *SetSettingsRequest, opts ...grpc.CallOption) (*UserSetting, error)
}

type userSettingsClient struct {
	cc grpc.ClientConnInterface
}

func NewUserSettingsClient(cc grpc.ClientConnInterface) UserSettingsClient {
	return &userSettingsClient{cc}
}

func (c *userSettingsClient) GetUserSettings(ctx context.Context, in *GetSettingsRequest, opts ...grpc.CallOption) (*UserSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserSetting)
	err := c.cc.Invoke(ctx, UserSettings_GetUserSettings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userSettingsClient) SetUserSettings(ctx context.Context, in *SetSettingsRequest, opts ...grpc.CallOption) (*UserSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserSetting)
	err := c.cc.Invoke(ctx, UserSettings_SetUserSettings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserSettingsServer is the server API for UserSettings service.
// All implementations must embed UnimplementedUserSettingsServer
// for forward compatibility.
type UserSettingsServer interface {
	// GET /user/settings/{key=*}
	GetUserSettings(context.Context, *GetSettingsRequest) (*UserSetting, error)
	// Create -or- Update
	//
	// PUT /user/settings/{key=*}
	// Content-Type: application/json; charset=utf-8
	//
	// ${value=.body}
	SetUserSettings(context.Context, *SetSettingsRequest) (*UserSetting, error)
	mustEmbedUnimplementedUserSettingsServer()
}

// UnimplementedUserSettingsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserSettingsServer struct{}

func (UnimplementedUserSettingsServer) GetUserSettings(context.Context, *GetSettingsRequest) (*UserSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSettings not implemented")
}
func (UnimplementedUserSettingsServer) SetUserSettings(context.Context, *SetSettingsRequest) (*UserSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserSettings not implemented")
}
func (UnimplementedUserSettingsServer) mustEmbedUnimplementedUserSettingsServer() {}
func (UnimplementedUserSettingsServer) testEmbeddedByValue()                      {}

// UnsafeUserSettingsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserSettingsServer will
// result in compilation errors.
type UnsafeUserSettingsServer interface {
	mustEmbedUnimplementedUserSettingsServer()
}

func RegisterUserSettingsServer(s grpc.ServiceRegistrar, srv UserSettingsServer) {
	// If the following call pancis, it indicates UnimplementedUserSettingsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserSettings_ServiceDesc, srv)
}

func _UserSettings_GetUserSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSettingsServer).GetUserSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserSettings_GetUserSettings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSettingsServer).GetUserSettings(ctx, req.(*GetSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserSettings_SetUserSettings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetSettingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserSettingsServer).SetUserSettings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserSettings_SetUserSettings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserSettingsServer).SetUserSettings(ctx, req.(*SetSettingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserSettings_ServiceDesc is the grpc.ServiceDesc for UserSettings service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserSettings_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.UserSettings",
	HandlerType: (*UserSettingsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserSettings",
			Handler:    _UserSettings_GetUserSettings_Handler,
		},
		{
			MethodName: "SetUserSettings",
			Handler:    _UserSettings_SetUserSettings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user_settings.proto",
}
