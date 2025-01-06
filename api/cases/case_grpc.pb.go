// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: case.proto

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
	Cases_SearchCases_FullMethodName = "/webitel.cases.Cases/SearchCases"
	Cases_LocateCase_FullMethodName  = "/webitel.cases.Cases/LocateCase"
	Cases_CreateCase_FullMethodName  = "/webitel.cases.Cases/CreateCase"
	Cases_UpdateCase_FullMethodName  = "/webitel.cases.Cases/UpdateCase"
	Cases_DeleteCase_FullMethodName  = "/webitel.cases.Cases/DeleteCase"
)

// CasesClient is the client API for Cases service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Service definition for managing cases.
type CasesClient interface {
	// RPC method for searching cases.
	SearchCases(ctx context.Context, in *SearchCasesRequest, opts ...grpc.CallOption) (*CaseList, error)
	// RPC method to retrieve a specific case by its etag identifier.
	LocateCase(ctx context.Context, in *LocateCaseRequest, opts ...grpc.CallOption) (*Case, error)
	// RPC method for creating a new case.
	CreateCase(ctx context.Context, in *CreateCaseRequest, opts ...grpc.CallOption) (*Case, error)
	// RPC method for updating an existing case.
	UpdateCase(ctx context.Context, in *UpdateCaseRequest, opts ...grpc.CallOption) (*Case, error)
	// RPC method for deleting an existing case by its etag.
	DeleteCase(ctx context.Context, in *DeleteCaseRequest, opts ...grpc.CallOption) (*Case, error)
}

type casesClient struct {
	cc grpc.ClientConnInterface
}

func NewCasesClient(cc grpc.ClientConnInterface) CasesClient {
	return &casesClient{cc}
}

func (c *casesClient) SearchCases(ctx context.Context, in *SearchCasesRequest, opts ...grpc.CallOption) (*CaseList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseList)
	err := c.cc.Invoke(ctx, Cases_SearchCases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *casesClient) LocateCase(ctx context.Context, in *LocateCaseRequest, opts ...grpc.CallOption) (*Case, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Case)
	err := c.cc.Invoke(ctx, Cases_LocateCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *casesClient) CreateCase(ctx context.Context, in *CreateCaseRequest, opts ...grpc.CallOption) (*Case, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Case)
	err := c.cc.Invoke(ctx, Cases_CreateCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *casesClient) UpdateCase(ctx context.Context, in *UpdateCaseRequest, opts ...grpc.CallOption) (*Case, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Case)
	err := c.cc.Invoke(ctx, Cases_UpdateCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *casesClient) DeleteCase(ctx context.Context, in *DeleteCaseRequest, opts ...grpc.CallOption) (*Case, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Case)
	err := c.cc.Invoke(ctx, Cases_DeleteCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CasesServer is the server API for Cases service.
// All implementations must embed UnimplementedCasesServer
// for forward compatibility.
//
// Service definition for managing cases.
type CasesServer interface {
	// RPC method for searching cases.
	SearchCases(context.Context, *SearchCasesRequest) (*CaseList, error)
	// RPC method to retrieve a specific case by its etag identifier.
	LocateCase(context.Context, *LocateCaseRequest) (*Case, error)
	// RPC method for creating a new case.
	CreateCase(context.Context, *CreateCaseRequest) (*Case, error)
	// RPC method for updating an existing case.
	UpdateCase(context.Context, *UpdateCaseRequest) (*Case, error)
	// RPC method for deleting an existing case by its etag.
	DeleteCase(context.Context, *DeleteCaseRequest) (*Case, error)
	mustEmbedUnimplementedCasesServer()
}

// UnimplementedCasesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCasesServer struct{}

func (UnimplementedCasesServer) SearchCases(context.Context, *SearchCasesRequest) (*CaseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchCases not implemented")
}
func (UnimplementedCasesServer) LocateCase(context.Context, *LocateCaseRequest) (*Case, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateCase not implemented")
}
func (UnimplementedCasesServer) CreateCase(context.Context, *CreateCaseRequest) (*Case, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCase not implemented")
}
func (UnimplementedCasesServer) UpdateCase(context.Context, *UpdateCaseRequest) (*Case, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCase not implemented")
}
func (UnimplementedCasesServer) DeleteCase(context.Context, *DeleteCaseRequest) (*Case, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCase not implemented")
}
func (UnimplementedCasesServer) mustEmbedUnimplementedCasesServer() {}
func (UnimplementedCasesServer) testEmbeddedByValue()               {}

// UnsafeCasesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CasesServer will
// result in compilation errors.
type UnsafeCasesServer interface {
	mustEmbedUnimplementedCasesServer()
}

func RegisterCasesServer(s grpc.ServiceRegistrar, srv CasesServer) {
	// If the following call pancis, it indicates UnimplementedCasesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Cases_ServiceDesc, srv)
}

func _Cases_SearchCases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchCasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CasesServer).SearchCases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cases_SearchCases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CasesServer).SearchCases(ctx, req.(*SearchCasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cases_LocateCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CasesServer).LocateCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cases_LocateCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CasesServer).LocateCase(ctx, req.(*LocateCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cases_CreateCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CasesServer).CreateCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cases_CreateCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CasesServer).CreateCase(ctx, req.(*CreateCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cases_UpdateCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CasesServer).UpdateCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cases_UpdateCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CasesServer).UpdateCase(ctx, req.(*UpdateCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cases_DeleteCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CasesServer).DeleteCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cases_DeleteCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CasesServer).DeleteCase(ctx, req.(*DeleteCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cases_ServiceDesc is the grpc.ServiceDesc for Cases service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cases_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.Cases",
	HandlerType: (*CasesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchCases",
			Handler:    _Cases_SearchCases_Handler,
		},
		{
			MethodName: "LocateCase",
			Handler:    _Cases_LocateCase_Handler,
		},
		{
			MethodName: "CreateCase",
			Handler:    _Cases_CreateCase_Handler,
		},
		{
			MethodName: "UpdateCase",
			Handler:    _Cases_UpdateCase_Handler,
		},
		{
			MethodName: "DeleteCase",
			Handler:    _Cases_DeleteCase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "case.proto",
}
