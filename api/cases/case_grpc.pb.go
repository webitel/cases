// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: cases/case.proto

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
type CasesClient interface {
	SearchCases(ctx context.Context, in *SearchCasesRequest, opts ...grpc.CallOption) (*CaseList, error)
	LocateCase(ctx context.Context, in *LocateCaseRequest, opts ...grpc.CallOption) (*Case, error)
	// on create, we should accept service and all parameters that correspond to it,
	// priority have the fields that were directly set from the front-end and if they are empty we should
	// fill them from service (we can't change the SLA and SLA conditions)
	// etag doesn't play any role on this API
	CreateCase(ctx context.Context, in *CreateCaseRequest, opts ...grpc.CallOption) (*Case, error)
	// on update, we should be able to accept service and all parameters that correspond to it,
	// if service and corresponding to it fields were changed simultaneously then priority have
	// service and dependent fields set from the service automatically (we can't change the SLA, SLA conditions )
	// etag is required to update the true version of the case
	UpdateCase(ctx context.Context, in *UpdateCaseRequest, opts ...grpc.CallOption) (*Case, error)
	// on delete, we should require etag, to understand if user has right version of the case
	// also will be deleted all objects connected to the case, such as comments, related cases, links and files
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
// All implementations should embed UnimplementedCasesServer
// for forward compatibility.
type CasesServer interface {
	SearchCases(context.Context, *SearchCasesRequest) (*CaseList, error)
	LocateCase(context.Context, *LocateCaseRequest) (*Case, error)
	// on create, we should accept service and all parameters that correspond to it,
	// priority have the fields that were directly set from the front-end and if they are empty we should
	// fill them from service (we can't change the SLA and SLA conditions)
	// etag doesn't play any role on this API
	CreateCase(context.Context, *CreateCaseRequest) (*Case, error)
	// on update, we should be able to accept service and all parameters that correspond to it,
	// if service and corresponding to it fields were changed simultaneously then priority have
	// service and dependent fields set from the service automatically (we can't change the SLA, SLA conditions )
	// etag is required to update the true version of the case
	UpdateCase(context.Context, *UpdateCaseRequest) (*Case, error)
	// on delete, we should require etag, to understand if user has right version of the case
	// also will be deleted all objects connected to the case, such as comments, related cases, links and files
	DeleteCase(context.Context, *DeleteCaseRequest) (*Case, error)
}

// UnimplementedCasesServer should be embedded to have
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
func (UnimplementedCasesServer) testEmbeddedByValue() {}

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
	Metadata: "cases/case.proto",
}

const (
	RelatedCases_LocateRelatedCase_FullMethodName = "/webitel.cases.RelatedCases/LocateRelatedCase"
	RelatedCases_UpdateRelatedCase_FullMethodName = "/webitel.cases.RelatedCases/UpdateRelatedCase"
	RelatedCases_DeleteRelatedCase_FullMethodName = "/webitel.cases.RelatedCases/DeleteRelatedCase"
	RelatedCases_ListRelatedCases_FullMethodName  = "/webitel.cases.RelatedCases/ListRelatedCases"
	RelatedCases_MergeRelatedCases_FullMethodName = "/webitel.cases.RelatedCases/MergeRelatedCases"
	RelatedCases_ResetRelatedCases_FullMethodName = "/webitel.cases.RelatedCases/ResetRelatedCases"
)

// RelatedCasesClient is the client API for RelatedCases service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelatedCasesClient interface {
	LocateRelatedCase(ctx context.Context, in *LocateRelatedCaseRequest, opts ...grpc.CallOption) (*RelatedCase, error)
	UpdateRelatedCase(ctx context.Context, in *UpdateRelatedCaseRequest, opts ...grpc.CallOption) (*RelatedCase, error)
	DeleteRelatedCase(ctx context.Context, in *DeleteRelatedCaseRequest, opts ...grpc.CallOption) (*RelatedCase, error)
	// The related cases can be obtained bidirectionally as child or parent, but we should consider them from the perspective of the requested case, by inverting their connection type
	// Requested case always a parent and related cases a children
	ListRelatedCases(ctx context.Context, in *ListRelatedCasesRequest, opts ...grpc.CallOption) (*RelatedCaseList, error)
	MergeRelatedCases(ctx context.Context, in *MergeRelatedCasesRequest, opts ...grpc.CallOption) (*RelatedCaseList, error)
	ResetRelatedCases(ctx context.Context, in *ResetRelatedCasesRequest, opts ...grpc.CallOption) (*RelatedCaseList, error)
}

type relatedCasesClient struct {
	cc grpc.ClientConnInterface
}

func NewRelatedCasesClient(cc grpc.ClientConnInterface) RelatedCasesClient {
	return &relatedCasesClient{cc}
}

func (c *relatedCasesClient) LocateRelatedCase(ctx context.Context, in *LocateRelatedCaseRequest, opts ...grpc.CallOption) (*RelatedCase, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RelatedCase)
	err := c.cc.Invoke(ctx, RelatedCases_LocateRelatedCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relatedCasesClient) UpdateRelatedCase(ctx context.Context, in *UpdateRelatedCaseRequest, opts ...grpc.CallOption) (*RelatedCase, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RelatedCase)
	err := c.cc.Invoke(ctx, RelatedCases_UpdateRelatedCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relatedCasesClient) DeleteRelatedCase(ctx context.Context, in *DeleteRelatedCaseRequest, opts ...grpc.CallOption) (*RelatedCase, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RelatedCase)
	err := c.cc.Invoke(ctx, RelatedCases_DeleteRelatedCase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relatedCasesClient) ListRelatedCases(ctx context.Context, in *ListRelatedCasesRequest, opts ...grpc.CallOption) (*RelatedCaseList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RelatedCaseList)
	err := c.cc.Invoke(ctx, RelatedCases_ListRelatedCases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relatedCasesClient) MergeRelatedCases(ctx context.Context, in *MergeRelatedCasesRequest, opts ...grpc.CallOption) (*RelatedCaseList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RelatedCaseList)
	err := c.cc.Invoke(ctx, RelatedCases_MergeRelatedCases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relatedCasesClient) ResetRelatedCases(ctx context.Context, in *ResetRelatedCasesRequest, opts ...grpc.CallOption) (*RelatedCaseList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RelatedCaseList)
	err := c.cc.Invoke(ctx, RelatedCases_ResetRelatedCases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelatedCasesServer is the server API for RelatedCases service.
// All implementations should embed UnimplementedRelatedCasesServer
// for forward compatibility.
type RelatedCasesServer interface {
	LocateRelatedCase(context.Context, *LocateRelatedCaseRequest) (*RelatedCase, error)
	UpdateRelatedCase(context.Context, *UpdateRelatedCaseRequest) (*RelatedCase, error)
	DeleteRelatedCase(context.Context, *DeleteRelatedCaseRequest) (*RelatedCase, error)
	// The related cases can be obtained bidirectionally as child or parent, but we should consider them from the perspective of the requested case, by inverting their connection type
	// Requested case always a parent and related cases a children
	ListRelatedCases(context.Context, *ListRelatedCasesRequest) (*RelatedCaseList, error)
	MergeRelatedCases(context.Context, *MergeRelatedCasesRequest) (*RelatedCaseList, error)
	ResetRelatedCases(context.Context, *ResetRelatedCasesRequest) (*RelatedCaseList, error)
}

// UnimplementedRelatedCasesServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRelatedCasesServer struct{}

func (UnimplementedRelatedCasesServer) LocateRelatedCase(context.Context, *LocateRelatedCaseRequest) (*RelatedCase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocateRelatedCase not implemented")
}
func (UnimplementedRelatedCasesServer) UpdateRelatedCase(context.Context, *UpdateRelatedCaseRequest) (*RelatedCase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRelatedCase not implemented")
}
func (UnimplementedRelatedCasesServer) DeleteRelatedCase(context.Context, *DeleteRelatedCaseRequest) (*RelatedCase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRelatedCase not implemented")
}
func (UnimplementedRelatedCasesServer) ListRelatedCases(context.Context, *ListRelatedCasesRequest) (*RelatedCaseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRelatedCases not implemented")
}
func (UnimplementedRelatedCasesServer) MergeRelatedCases(context.Context, *MergeRelatedCasesRequest) (*RelatedCaseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeRelatedCases not implemented")
}
func (UnimplementedRelatedCasesServer) ResetRelatedCases(context.Context, *ResetRelatedCasesRequest) (*RelatedCaseList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetRelatedCases not implemented")
}
func (UnimplementedRelatedCasesServer) testEmbeddedByValue() {}

// UnsafeRelatedCasesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelatedCasesServer will
// result in compilation errors.
type UnsafeRelatedCasesServer interface {
	mustEmbedUnimplementedRelatedCasesServer()
}

func RegisterRelatedCasesServer(s grpc.ServiceRegistrar, srv RelatedCasesServer) {
	// If the following call pancis, it indicates UnimplementedRelatedCasesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RelatedCases_ServiceDesc, srv)
}

func _RelatedCases_LocateRelatedCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocateRelatedCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelatedCasesServer).LocateRelatedCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelatedCases_LocateRelatedCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelatedCasesServer).LocateRelatedCase(ctx, req.(*LocateRelatedCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelatedCases_UpdateRelatedCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRelatedCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelatedCasesServer).UpdateRelatedCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelatedCases_UpdateRelatedCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelatedCasesServer).UpdateRelatedCase(ctx, req.(*UpdateRelatedCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelatedCases_DeleteRelatedCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRelatedCaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelatedCasesServer).DeleteRelatedCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelatedCases_DeleteRelatedCase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelatedCasesServer).DeleteRelatedCase(ctx, req.(*DeleteRelatedCaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelatedCases_ListRelatedCases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRelatedCasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelatedCasesServer).ListRelatedCases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelatedCases_ListRelatedCases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelatedCasesServer).ListRelatedCases(ctx, req.(*ListRelatedCasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelatedCases_MergeRelatedCases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeRelatedCasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelatedCasesServer).MergeRelatedCases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelatedCases_MergeRelatedCases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelatedCasesServer).MergeRelatedCases(ctx, req.(*MergeRelatedCasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelatedCases_ResetRelatedCases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetRelatedCasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelatedCasesServer).ResetRelatedCases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelatedCases_ResetRelatedCases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelatedCasesServer).ResetRelatedCases(ctx, req.(*ResetRelatedCasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RelatedCases_ServiceDesc is the grpc.ServiceDesc for RelatedCases service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelatedCases_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.cases.RelatedCases",
	HandlerType: (*RelatedCasesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LocateRelatedCase",
			Handler:    _RelatedCases_LocateRelatedCase_Handler,
		},
		{
			MethodName: "UpdateRelatedCase",
			Handler:    _RelatedCases_UpdateRelatedCase_Handler,
		},
		{
			MethodName: "DeleteRelatedCase",
			Handler:    _RelatedCases_DeleteRelatedCase_Handler,
		},
		{
			MethodName: "ListRelatedCases",
			Handler:    _RelatedCases_ListRelatedCases_Handler,
		},
		{
			MethodName: "MergeRelatedCases",
			Handler:    _RelatedCases_MergeRelatedCases_Handler,
		},
		{
			MethodName: "ResetRelatedCases",
			Handler:    _RelatedCases_ResetRelatedCases_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cases/case.proto",
}
