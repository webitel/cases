// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: case_timeline.proto

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
	CaseTimeline_GetTimeline_FullMethodName        = "/webitel.cases.CaseTimeline/GetTimeline"
	CaseTimeline_GetTimelineCounter_FullMethodName = "/webitel.cases.CaseTimeline/GetTimelineCounter"
)

// CaseTimelineClient is the client API for CaseTimeline service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// CloseReasons service definition with RPC methods for managing close reasons
type CaseTimelineClient interface {
	GetTimeline(ctx context.Context, in *GetTimelineRequest, opts ...grpc.CallOption) (*GetTimelineResponse, error)
	GetTimelineCounter(ctx context.Context, in *GetTimelineCounterRequest, opts ...grpc.CallOption) (*GetTimelineCounterResponse, error)
}

type caseTimelineClient struct {
	cc grpc.ClientConnInterface
}

func NewCaseTimelineClient(cc grpc.ClientConnInterface) CaseTimelineClient {
	return &caseTimelineClient{cc}
}

func (c *caseTimelineClient) GetTimeline(ctx context.Context, in *GetTimelineRequest, opts ...grpc.CallOption) (*GetTimelineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTimelineResponse)
	err := c.cc.Invoke(ctx, CaseTimeline_GetTimeline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseTimelineClient) GetTimelineCounter(ctx context.Context, in *GetTimelineCounterRequest, opts ...grpc.CallOption) (*GetTimelineCounterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTimelineCounterResponse)
	err := c.cc.Invoke(ctx, CaseTimeline_GetTimelineCounter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CaseTimelineServer is the server API for CaseTimeline service.
// All implementations must embed UnimplementedCaseTimelineServer
// for forward compatibility.
//
// CloseReasons service definition with RPC methods for managing close reasons
type CaseTimelineServer interface {
	GetTimeline(context.Context, *GetTimelineRequest) (*GetTimelineResponse, error)
	GetTimelineCounter(context.Context, *GetTimelineCounterRequest) (*GetTimelineCounterResponse, error)
	mustEmbedUnimplementedCaseTimelineServer()
}

// UnimplementedCaseTimelineServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCaseTimelineServer struct{}

func (UnimplementedCaseTimelineServer) GetTimeline(context.Context, *GetTimelineRequest) (*GetTimelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTimeline not implemented")
}
func (UnimplementedCaseTimelineServer) GetTimelineCounter(context.Context, *GetTimelineCounterRequest) (*GetTimelineCounterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTimelineCounter not implemented")
}
func (UnimplementedCaseTimelineServer) mustEmbedUnimplementedCaseTimelineServer() {}
func (UnimplementedCaseTimelineServer) testEmbeddedByValue()                      {}

// UnsafeCaseTimelineServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CaseTimelineServer will
// result in compilation errors.
type UnsafeCaseTimelineServer interface {
	mustEmbedUnimplementedCaseTimelineServer()
}

func RegisterCaseTimelineServer(s grpc.ServiceRegistrar, srv CaseTimelineServer) {
	// If the following call pancis, it indicates UnimplementedCaseTimelineServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CaseTimeline_ServiceDesc, srv)
}

func _CaseTimeline_GetTimeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTimelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseTimelineServer).GetTimeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseTimeline_GetTimeline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseTimelineServer).GetTimeline(ctx, req.(*GetTimelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseTimeline_GetTimelineCounter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTimelineCounterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseTimelineServer).GetTimelineCounter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseTimeline_GetTimelineCounter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseTimelineServer).GetTimelineCounter(ctx, req.(*GetTimelineCounterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CaseTimeline_ServiceDesc is the grpc.ServiceDesc for CaseTimeline service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CaseTimeline_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.CaseTimeline",
	HandlerType: (*CaseTimelineServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTimeline",
			Handler:    _CaseTimeline_GetTimeline_Handler,
		},
		{
			MethodName: "GetTimelineCounter",
			Handler:    _CaseTimeline_GetTimelineCounter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "case_timeline.proto",
}
