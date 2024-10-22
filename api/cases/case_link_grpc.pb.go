// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: cases/case_link.proto

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
	CaseLinks_LocateLink_FullMethodName = "/webitel.cases.CaseLinks/LocateLink"
	CaseLinks_UpdateLink_FullMethodName = "/webitel.cases.CaseLinks/UpdateLink"
	CaseLinks_DeleteLink_FullMethodName = "/webitel.cases.CaseLinks/DeleteLink"
	CaseLinks_ListLinks_FullMethodName  = "/webitel.cases.CaseLinks/ListLinks"
	CaseLinks_MergeLinks_FullMethodName = "/webitel.cases.CaseLinks/MergeLinks"
	CaseLinks_ResetLinks_FullMethodName = "/webitel.cases.CaseLinks/ResetLinks"
)

// CaseLinksClient is the client API for CaseLinks service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CaseLinksClient interface {
	LocateLink(ctx context.Context, in *LocateLinkRequest, opts ...grpc.CallOption) (*CaseLink, error)
	UpdateLink(ctx context.Context, in *UpdateLinkRequest, opts ...grpc.CallOption) (*CaseLink, error)
	DeleteLink(ctx context.Context, in *DeleteLinkRequest, opts ...grpc.CallOption) (*CaseLink, error)
	ListLinks(ctx context.Context, in *ListLinksRequest, opts ...grpc.CallOption) (*CaseLinkList, error)
	MergeLinks(ctx context.Context, in *MergeLinksRequest, opts ...grpc.CallOption) (*CaseLinkList, error)
	ResetLinks(ctx context.Context, in *ResetLinksRequest, opts ...grpc.CallOption) (*CaseLinkList, error)
}

type caseLinksClient struct {
	cc grpc.ClientConnInterface
}

func NewCaseLinksClient(cc grpc.ClientConnInterface) CaseLinksClient {
	return &caseLinksClient{cc}
}

func (c *caseLinksClient) LocateLink(ctx context.Context, in *LocateLinkRequest, opts ...grpc.CallOption) (*CaseLink, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseLink)
	err := c.cc.Invoke(ctx, CaseLinks_LocateLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseLinksClient) UpdateLink(ctx context.Context, in *UpdateLinkRequest, opts ...grpc.CallOption) (*CaseLink, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseLink)
	err := c.cc.Invoke(ctx, CaseLinks_UpdateLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseLinksClient) DeleteLink(ctx context.Context, in *DeleteLinkRequest, opts ...grpc.CallOption) (*CaseLink, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseLink)
	err := c.cc.Invoke(ctx, CaseLinks_DeleteLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseLinksClient) ListLinks(ctx context.Context, in *ListLinksRequest, opts ...grpc.CallOption) (*CaseLinkList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseLinkList)
	err := c.cc.Invoke(ctx, CaseLinks_ListLinks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseLinksClient) MergeLinks(ctx context.Context, in *MergeLinksRequest, opts ...grpc.CallOption) (*CaseLinkList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseLinkList)
	err := c.cc.Invoke(ctx, CaseLinks_MergeLinks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *caseLinksClient) ResetLinks(ctx context.Context, in *ResetLinksRequest, opts ...grpc.CallOption) (*CaseLinkList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CaseLinkList)
	err := c.cc.Invoke(ctx, CaseLinks_ResetLinks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CaseLinksServer is the server API for CaseLinks service.
// All implementations should embed UnimplementedCaseLinksServer
// for forward compatibility.
type CaseLinksServer interface {
	LocateLink(context.Context, *LocateLinkRequest) (*CaseLink, error)
	UpdateLink(context.Context, *UpdateLinkRequest) (*CaseLink, error)
	DeleteLink(context.Context, *DeleteLinkRequest) (*CaseLink, error)
	ListLinks(context.Context, *ListLinksRequest) (*CaseLinkList, error)
	MergeLinks(context.Context, *MergeLinksRequest) (*CaseLinkList, error)
	ResetLinks(context.Context, *ResetLinksRequest) (*CaseLinkList, error)
}

// UnimplementedCaseLinksServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCaseLinksServer struct{}

func (UnimplementedCaseLinksServer) LocateLink(context.Context, *LocateLinkRequest) (*CaseLink, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateLink not implemented")
}
func (UnimplementedCaseLinksServer) UpdateLink(context.Context, *UpdateLinkRequest) (*CaseLink, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLink not implemented")
}
func (UnimplementedCaseLinksServer) DeleteLink(context.Context, *DeleteLinkRequest) (*CaseLink, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLink not implemented")
}
func (UnimplementedCaseLinksServer) ListLinks(context.Context, *ListLinksRequest) (*CaseLinkList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLinks not implemented")
}
func (UnimplementedCaseLinksServer) MergeLinks(context.Context, *MergeLinksRequest) (*CaseLinkList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeLinks not implemented")
}
func (UnimplementedCaseLinksServer) ResetLinks(context.Context, *ResetLinksRequest) (*CaseLinkList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetLinks not implemented")
}
func (UnimplementedCaseLinksServer) testEmbeddedByValue() {}

// UnsafeCaseLinksServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CaseLinksServer will
// result in compilation errors.
type UnsafeCaseLinksServer interface {
	mustEmbedUnimplementedCaseLinksServer()
}

func RegisterCaseLinksServer(s grpc.ServiceRegistrar, srv CaseLinksServer) {
	// If the following call pancis, it indicates UnimplementedCaseLinksServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CaseLinks_ServiceDesc, srv)
}

func _CaseLinks_LocateLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseLinksServer).LocateLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseLinks_LocateLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseLinksServer).LocateLink(ctx, req.(*LocateLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseLinks_UpdateLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseLinksServer).UpdateLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseLinks_UpdateLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseLinksServer).UpdateLink(ctx, req.(*UpdateLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseLinks_DeleteLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseLinksServer).DeleteLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseLinks_DeleteLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseLinksServer).DeleteLink(ctx, req.(*DeleteLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseLinks_ListLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseLinksServer).ListLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseLinks_ListLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseLinksServer).ListLinks(ctx, req.(*ListLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseLinks_MergeLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseLinksServer).MergeLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseLinks_MergeLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseLinksServer).MergeLinks(ctx, req.(*MergeLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CaseLinks_ResetLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CaseLinksServer).ResetLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CaseLinks_ResetLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CaseLinksServer).ResetLinks(ctx, req.(*ResetLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CaseLinks_ServiceDesc is the grpc.ServiceDesc for CaseLinks service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CaseLinks_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.CaseLinks",
	HandlerType: (*CaseLinksServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LocateLink",
			Handler:    _CaseLinks_LocateLink_Handler,
		},
		{
			MethodName: "UpdateLink",
			Handler:    _CaseLinks_UpdateLink_Handler,
		},
		{
			MethodName: "DeleteLink",
			Handler:    _CaseLinks_DeleteLink_Handler,
		},
		{
			MethodName: "ListLinks",
			Handler:    _CaseLinks_ListLinks_Handler,
		},
		{
			MethodName: "MergeLinks",
			Handler:    _CaseLinks_MergeLinks_Handler,
		},
		{
			MethodName: "ResetLinks",
			Handler:    _CaseLinks_ResetLinks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cases/case_link.proto",
}
