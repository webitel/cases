// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: case_comment.proto

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
	CaseComments_LocateComment_FullMethodName = "/webitel.cases.CaseComments/LocateComment"
	CaseComments_UpdateComment_FullMethodName = "/webitel.cases.CaseComments/UpdateComment"
	CaseComments_DeleteComment_FullMethodName = "/webitel.cases.CaseComments/DeleteComment"
	CaseComments_ListComments_FullMethodName  = "/webitel.cases.CaseComments/ListComments"
	CaseComments_MergeComments_FullMethodName = "/webitel.cases.CaseComments/MergeComments"
)

// CaseCommentsClient is the client API for CaseComments service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CaseCommentsClient interface {
	// Itself
	LocateComment(ctx context.Context, in *LocateCommentRequest, opts ...grpc.CallOption) (*CaseComment, error)
	UpdateComment(ctx context.Context, in *UpdateCommentRequest, opts ...grpc.CallOption) (*CaseComment, error)
	DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*CaseComment, error)
	ListComments(ctx context.Context, in *ListCommentsRequest, opts ...grpc.CallOption) (*CaseCommentList, error)
	MergeComments(ctx context.Context, in *MergeCommentsRequest, opts ...grpc.CallOption) (*CaseCommentList, error)
}

type caseCommentsClient struct {
	cc grpc.ClientConnInterface
}

func NewCaseCommentsClient(cc grpc.ClientConnInterface) CaseCommentsClient {
	return &caseCommentsClient{cc}
}

func (c *caseCommentsClient) LocateComment(ctx context.Context, in *LocateCommentRequest, opts ...grpc.CallOption) (*CaseComment, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseComment)
	err := c.cc.Invoke(ctx, CaseComments_LocateComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseCommentsClient) UpdateComment(ctx context.Context, in *UpdateCommentRequest, opts ...grpc.CallOption) (*CaseComment, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseComment)
	err := c.cc.Invoke(ctx, CaseComments_UpdateComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseCommentsClient) DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*CaseComment, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseComment)
	err := c.cc.Invoke(ctx, CaseComments_DeleteComment_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseCommentsClient) ListComments(ctx context.Context, in *ListCommentsRequest, opts ...grpc.CallOption) (*CaseCommentList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseCommentList)
	err := c.cc.Invoke(ctx, CaseComments_ListComments_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseCommentsClient) MergeComments(ctx context.Context, in *MergeCommentsRequest, opts ...grpc.CallOption) (*CaseCommentList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseCommentList)
	err := c.cc.Invoke(ctx, CaseComments_MergeComments_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CaseCommentsServer is the server API for CaseComments service.
// All implementations must embed UnimplementedCaseCommentsServer
// for forward compatibility.
type CaseCommentsServer interface {
	// Itself
	LocateComment(context.Context, *LocateCommentRequest) (*CaseComment, error)
	UpdateComment(context.Context, *UpdateCommentRequest) (*CaseComment, error)
	DeleteComment(context.Context, *DeleteCommentRequest) (*CaseComment, error)
	ListComments(context.Context, *ListCommentsRequest) (*CaseCommentList, error)
	MergeComments(context.Context, *MergeCommentsRequest) (*CaseCommentList, error)
	mustEmbedUnimplementedCaseCommentsServer()
}

// UnimplementedCaseCommentsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCaseCommentsServer struct{}

func (UnimplementedCaseCommentsServer) LocateComment(context.Context, *LocateCommentRequest) (*CaseComment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateComment not implemented")
}
func (UnimplementedCaseCommentsServer) UpdateComment(context.Context, *UpdateCommentRequest) (*CaseComment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateComment not implemented")
}
func (UnimplementedCaseCommentsServer) DeleteComment(context.Context, *DeleteCommentRequest) (*CaseComment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (UnimplementedCaseCommentsServer) ListComments(context.Context, *ListCommentsRequest) (*CaseCommentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListComments not implemented")
}
func (UnimplementedCaseCommentsServer) MergeComments(context.Context, *MergeCommentsRequest) (*CaseCommentList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeComments not implemented")
}
func (UnimplementedCaseCommentsServer) mustEmbedUnimplementedCaseCommentsServer() {}
func (UnimplementedCaseCommentsServer) testEmbeddedByValue()                      {}

// UnsafeCaseCommentsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CaseCommentsServer will
// result in compilation errors.
type UnsafeCaseCommentsServer interface {
	mustEmbedUnimplementedCaseCommentsServer()
}

func RegisterCaseCommentsServer(s grpc.ServiceRegistrar, srv CaseCommentsServer) {
	// If the following call pancis, it indicates UnimplementedCaseCommentsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CaseComments_ServiceDesc, srv)
}

func _CaseComments_LocateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseCommentsServer).LocateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseComments_LocateComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseCommentsServer).LocateComment(ctx, req.(*LocateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseComments_UpdateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseCommentsServer).UpdateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseComments_UpdateComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseCommentsServer).UpdateComment(ctx, req.(*UpdateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseComments_DeleteComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseCommentsServer).DeleteComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseComments_DeleteComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseCommentsServer).DeleteComment(ctx, req.(*DeleteCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseComments_ListComments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCommentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseCommentsServer).ListComments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseComments_ListComments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseCommentsServer).ListComments(ctx, req.(*ListCommentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseComments_MergeComments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeCommentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseCommentsServer).MergeComments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseComments_MergeComments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseCommentsServer).MergeComments(ctx, req.(*MergeCommentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CaseComments_ServiceDesc is the grpc.ServiceDesc for CaseComments service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CaseComments_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.CaseComments",
	HandlerType: (*CaseCommentsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LocateComment",
			Handler:    _CaseComments_LocateComment_Handler,
		},
		{
			MethodName: "UpdateComment",
			Handler:    _CaseComments_UpdateComment_Handler,
		},
		{
			MethodName: "DeleteComment",
			Handler:    _CaseComments_DeleteComment_Handler,
		},
		{
			MethodName: "ListComments",
			Handler:    _CaseComments_ListComments_Handler,
		},
		{
			MethodName: "MergeComments",
			Handler:    _CaseComments_MergeComments_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "case_comment.proto",
}
