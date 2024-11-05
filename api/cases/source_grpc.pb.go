// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: source.proto

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
	Sources_ListSources_FullMethodName  = "/webitel.cases.Sources/ListSources"
	Sources_CreateSource_FullMethodName = "/webitel.cases.Sources/CreateSource"
	Sources_UpdateSource_FullMethodName = "/webitel.cases.Sources/UpdateSource"
	Sources_DeleteSource_FullMethodName = "/webitel.cases.Sources/DeleteSource"
	Sources_LocateSource_FullMethodName = "/webitel.cases.Sources/LocateSource"
)

// SourcesClient is the client API for Sources service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SourcesClient interface {
	ListSources(ctx context.Context, in *ListSourceRequest, opts ...grpc.CallOption) (*SourceList, error)
	CreateSource(ctx context.Context, in *CreateSourceRequest, opts ...grpc.CallOption) (*Source, error)
	UpdateSource(ctx context.Context, in *UpdateSourceRequest, opts ...grpc.CallOption) (*Source, error)
	DeleteSource(ctx context.Context, in *DeleteSourceRequest, opts ...grpc.CallOption) (*Source, error)
	LocateSource(ctx context.Context, in *LocateSourceRequest, opts ...grpc.CallOption) (*LocateSourceResponse, error)
}

type sourcesClient struct {
	cc grpc.ClientConnInterface
}

func NewSourcesClient(cc grpc.ClientConnInterface) SourcesClient {
	return &sourcesClient{cc}
}

func (c *sourcesClient) ListSources(ctx context.Context, in *ListSourceRequest, opts ...grpc.CallOption) (*SourceList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SourceList)
	err := c.cc.Invoke(ctx, Sources_ListSources_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourcesClient) CreateSource(ctx context.Context, in *CreateSourceRequest, opts ...grpc.CallOption) (*Source, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Source)
	err := c.cc.Invoke(ctx, Sources_CreateSource_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourcesClient) UpdateSource(ctx context.Context, in *UpdateSourceRequest, opts ...grpc.CallOption) (*Source, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Source)
	err := c.cc.Invoke(ctx, Sources_UpdateSource_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourcesClient) DeleteSource(ctx context.Context, in *DeleteSourceRequest, opts ...grpc.CallOption) (*Source, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Source)
	err := c.cc.Invoke(ctx, Sources_DeleteSource_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sourcesClient) LocateSource(ctx context.Context, in *LocateSourceRequest, opts ...grpc.CallOption) (*LocateSourceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LocateSourceResponse)
	err := c.cc.Invoke(ctx, Sources_LocateSource_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SourcesServer is the server API for Sources service.
// All implementations must embed UnimplementedSourcesServer
// for forward compatibility.
type SourcesServer interface {
	ListSources(context.Context, *ListSourceRequest) (*SourceList, error)
	CreateSource(context.Context, *CreateSourceRequest) (*Source, error)
	UpdateSource(context.Context, *UpdateSourceRequest) (*Source, error)
	DeleteSource(context.Context, *DeleteSourceRequest) (*Source, error)
	LocateSource(context.Context, *LocateSourceRequest) (*LocateSourceResponse, error)
	mustEmbedUnimplementedSourcesServer()
}

// UnimplementedSourcesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSourcesServer struct{}

func (UnimplementedSourcesServer) ListSources(context.Context, *ListSourceRequest) (*SourceList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSources not implemented")
}
func (UnimplementedSourcesServer) CreateSource(context.Context, *CreateSourceRequest) (*Source, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSource not implemented")
}
func (UnimplementedSourcesServer) UpdateSource(context.Context, *UpdateSourceRequest) (*Source, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSource not implemented")
}
func (UnimplementedSourcesServer) DeleteSource(context.Context, *DeleteSourceRequest) (*Source, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSource not implemented")
}
func (UnimplementedSourcesServer) LocateSource(context.Context, *LocateSourceRequest) (*LocateSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateSource not implemented")
}
func (UnimplementedSourcesServer) mustEmbedUnimplementedSourcesServer() {}
func (UnimplementedSourcesServer) testEmbeddedByValue()                 {}

// UnsafeSourcesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SourcesServer will
// result in compilation errors.
type UnsafeSourcesServer interface {
	mustEmbedUnimplementedSourcesServer()
}

func RegisterSourcesServer(s grpc.ServiceRegistrar, srv SourcesServer) {
	// If the following call pancis, it indicates UnimplementedSourcesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Sources_ServiceDesc, srv)
}

func _Sources_ListSources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServer).ListSources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sources_ListSources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServer).ListSources(ctx, req.(*ListSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sources_CreateSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServer).CreateSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sources_CreateSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServer).CreateSource(ctx, req.(*CreateSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sources_UpdateSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServer).UpdateSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sources_UpdateSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServer).UpdateSource(ctx, req.(*UpdateSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sources_DeleteSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServer).DeleteSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sources_DeleteSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServer).DeleteSource(ctx, req.(*DeleteSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sources_LocateSource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateSourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SourcesServer).LocateSource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sources_LocateSource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SourcesServer).LocateSource(ctx, req.(*LocateSourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sources_ServiceDesc is the grpc.ServiceDesc for Sources service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sources_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.Sources",
	HandlerType: (*SourcesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSources",
			Handler:    _Sources_ListSources_Handler,
		},
		{
			MethodName: "CreateSource",
			Handler:    _Sources_CreateSource_Handler,
		},
		{
			MethodName: "UpdateSource",
			Handler:    _Sources_UpdateSource_Handler,
		},
		{
			MethodName: "DeleteSource",
			Handler:    _Sources_DeleteSource_Handler,
		},
		{
			MethodName: "LocateSource",
			Handler:    _Sources_LocateSource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "source.proto",
}
