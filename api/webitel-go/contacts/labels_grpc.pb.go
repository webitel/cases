// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: contacts/labels.proto

package contacts

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
	Labels_GetLabels_FullMethodName    = "/webitel.contacts.Labels/GetLabels"
	Labels_ListLabels_FullMethodName   = "/webitel.contacts.Labels/ListLabels"
	Labels_MergeLabels_FullMethodName  = "/webitel.contacts.Labels/MergeLabels"
	Labels_ResetLabels_FullMethodName  = "/webitel.contacts.Labels/ResetLabels"
	Labels_DeleteLabels_FullMethodName = "/webitel.contacts.Labels/DeleteLabels"
)

// LabelsClient is the client API for Labels service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Labels service catalog.
type LabelsClient interface {
	// Search for Contacts engaged Label(s).
	GetLabels(ctx context.Context, in *GetLabelsRequest, opts ...grpc.CallOption) (*LabelTags, error)
	// Locate the Contact's associated Label(s).
	ListLabels(ctx context.Context, in *ListLabelsRequest, opts ...grpc.CallOption) (*LabelList, error)
	// Associate NEW Labels to the Contact.
	MergeLabels(ctx context.Context, in *MergeLabelsRequest, opts ...grpc.CallOption) (*LabelList, error)
	// Reset Labels to fit the specified final set.
	ResetLabels(ctx context.Context, in *ResetLabelsRequest, opts ...grpc.CallOption) (*LabelList, error)
	// Remove Contact Labels associations.
	DeleteLabels(ctx context.Context, in *DeleteLabelsRequest, opts ...grpc.CallOption) (*LabelList, error)
}

type labelsClient struct {
	cc grpc.ClientConnInterface
}

func NewLabelsClient(cc grpc.ClientConnInterface) LabelsClient {
	return &labelsClient{cc}
}

func (c *labelsClient) GetLabels(ctx context.Context, in *GetLabelsRequest, opts ...grpc.CallOption) (*LabelTags, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LabelTags)
	err := c.cc.Invoke(ctx, Labels_GetLabels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *labelsClient) ListLabels(ctx context.Context, in *ListLabelsRequest, opts ...grpc.CallOption) (*LabelList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LabelList)
	err := c.cc.Invoke(ctx, Labels_ListLabels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *labelsClient) MergeLabels(ctx context.Context, in *MergeLabelsRequest, opts ...grpc.CallOption) (*LabelList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LabelList)
	err := c.cc.Invoke(ctx, Labels_MergeLabels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *labelsClient) ResetLabels(ctx context.Context, in *ResetLabelsRequest, opts ...grpc.CallOption) (*LabelList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LabelList)
	err := c.cc.Invoke(ctx, Labels_ResetLabels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *labelsClient) DeleteLabels(ctx context.Context, in *DeleteLabelsRequest, opts ...grpc.CallOption) (*LabelList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LabelList)
	err := c.cc.Invoke(ctx, Labels_DeleteLabels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LabelsServer is the server API for Labels service.
// All implementations must embed UnimplementedLabelsServer
// for forward compatibility.
//
// Labels service catalog.
type LabelsServer interface {
	// Search for Contacts engaged Label(s).
	GetLabels(context.Context, *GetLabelsRequest) (*LabelTags, error)
	// Locate the Contact's associated Label(s).
	ListLabels(context.Context, *ListLabelsRequest) (*LabelList, error)
	// Associate NEW Labels to the Contact.
	MergeLabels(context.Context, *MergeLabelsRequest) (*LabelList, error)
	// Reset Labels to fit the specified final set.
	ResetLabels(context.Context, *ResetLabelsRequest) (*LabelList, error)
	// Remove Contact Labels associations.
	DeleteLabels(context.Context, *DeleteLabelsRequest) (*LabelList, error)
	mustEmbedUnimplementedLabelsServer()
}

// UnimplementedLabelsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLabelsServer struct{}

func (UnimplementedLabelsServer) GetLabels(context.Context, *GetLabelsRequest) (*LabelTags, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLabels not implemented")
}
func (UnimplementedLabelsServer) ListLabels(context.Context, *ListLabelsRequest) (*LabelList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLabels not implemented")
}
func (UnimplementedLabelsServer) MergeLabels(context.Context, *MergeLabelsRequest) (*LabelList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeLabels not implemented")
}
func (UnimplementedLabelsServer) ResetLabels(context.Context, *ResetLabelsRequest) (*LabelList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetLabels not implemented")
}
func (UnimplementedLabelsServer) DeleteLabels(context.Context, *DeleteLabelsRequest) (*LabelList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLabels not implemented")
}
func (UnimplementedLabelsServer) mustEmbedUnimplementedLabelsServer() {}
func (UnimplementedLabelsServer) testEmbeddedByValue()                {}

// UnsafeLabelsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LabelsServer will
// result in compilation errors.
type UnsafeLabelsServer interface {
	mustEmbedUnimplementedLabelsServer()
}

func RegisterLabelsServer(s grpc.ServiceRegistrar, srv LabelsServer) {
	// If the following call pancis, it indicates UnimplementedLabelsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Labels_ServiceDesc, srv)
}

func _Labels_GetLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabelsServer).GetLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Labels_GetLabels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabelsServer).GetLabels(ctx, req.(*GetLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Labels_ListLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabelsServer).ListLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Labels_ListLabels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabelsServer).ListLabels(ctx, req.(*ListLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Labels_MergeLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabelsServer).MergeLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Labels_MergeLabels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabelsServer).MergeLabels(ctx, req.(*MergeLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Labels_ResetLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabelsServer).ResetLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Labels_ResetLabels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabelsServer).ResetLabels(ctx, req.(*ResetLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Labels_DeleteLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteLabelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LabelsServer).DeleteLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Labels_DeleteLabels_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LabelsServer).DeleteLabels(ctx, req.(*DeleteLabelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Labels_ServiceDesc is the grpc.ServiceDesc for Labels service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Labels_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.contacts.Labels",
	HandlerType: (*LabelsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLabels",
			Handler:    _Labels_GetLabels_Handler,
		},
		{
			MethodName: "ListLabels",
			Handler:    _Labels_ListLabels_Handler,
		},
		{
			MethodName: "MergeLabels",
			Handler:    _Labels_MergeLabels_Handler,
		},
		{
			MethodName: "ResetLabels",
			Handler:    _Labels_ResetLabels_Handler,
		},
		{
			MethodName: "DeleteLabels",
			Handler:    _Labels_DeleteLabels_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contacts/labels.proto",
}
