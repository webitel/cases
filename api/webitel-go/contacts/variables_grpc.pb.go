// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: contacts/variables.proto

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
	Variables_ListVariables_FullMethodName   = "/webitel.contacts.Variables/ListVariables"
	Variables_MergeVariables_FullMethodName  = "/webitel.contacts.Variables/MergeVariables"
	Variables_ResetVariables_FullMethodName  = "/webitel.contacts.Variables/ResetVariables"
	Variables_DeleteVariables_FullMethodName = "/webitel.contacts.Variables/DeleteVariables"
	Variables_UpdateVariable_FullMethodName  = "/webitel.contacts.Variables/UpdateVariable"
	Variables_DeleteVariable_FullMethodName  = "/webitel.contacts.Variables/DeleteVariable"
)

// VariablesClient is the client API for Variables service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Variables service catalog.
type VariablesClient interface {
	// List variables of the contact
	ListVariables(ctx context.Context, in *SearchVariablesRequest, opts ...grpc.CallOption) (*VariableList, error)
	// Update or append variables to the contact
	MergeVariables(ctx context.Context, in *MergeVariablesRequest, opts ...grpc.CallOption) (*VariableList, error)
	// Reset all variables of the contact
	ResetVariables(ctx context.Context, in *ResetVariablesRequest, opts ...grpc.CallOption) (*VariableList, error)
	// Remove variable(s) of the contact
	DeleteVariables(ctx context.Context, in *DeleteVariablesRequest, opts ...grpc.CallOption) (*VariableList, error)
	// Update contact variable
	UpdateVariable(ctx context.Context, in *UpdateVariableRequest, opts ...grpc.CallOption) (*VariableList, error)
	// Remove the contact's variable by etag
	DeleteVariable(ctx context.Context, in *DeleteVariableRequest, opts ...grpc.CallOption) (*Variable, error)
}

type variablesClient struct {
	cc grpc.ClientConnInterface
}

func NewVariablesClient(cc grpc.ClientConnInterface) VariablesClient {
	return &variablesClient{cc}
}

func (c *variablesClient) ListVariables(ctx context.Context, in *SearchVariablesRequest, opts ...grpc.CallOption) (*VariableList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VariableList)
	err := c.cc.Invoke(ctx, Variables_ListVariables_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *variablesClient) MergeVariables(ctx context.Context, in *MergeVariablesRequest, opts ...grpc.CallOption) (*VariableList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VariableList)
	err := c.cc.Invoke(ctx, Variables_MergeVariables_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *variablesClient) ResetVariables(ctx context.Context, in *ResetVariablesRequest, opts ...grpc.CallOption) (*VariableList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VariableList)
	err := c.cc.Invoke(ctx, Variables_ResetVariables_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *variablesClient) DeleteVariables(ctx context.Context, in *DeleteVariablesRequest, opts ...grpc.CallOption) (*VariableList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VariableList)
	err := c.cc.Invoke(ctx, Variables_DeleteVariables_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *variablesClient) UpdateVariable(ctx context.Context, in *UpdateVariableRequest, opts ...grpc.CallOption) (*VariableList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VariableList)
	err := c.cc.Invoke(ctx, Variables_UpdateVariable_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *variablesClient) DeleteVariable(ctx context.Context, in *DeleteVariableRequest, opts ...grpc.CallOption) (*Variable, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Variable)
	err := c.cc.Invoke(ctx, Variables_DeleteVariable_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VariablesServer is the server API for Variables service.
// All implementations must embed UnimplementedVariablesServer
// for forward compatibility.
//
// Variables service catalog.
type VariablesServer interface {
	// List variables of the contact
	ListVariables(context.Context, *SearchVariablesRequest) (*VariableList, error)
	// Update or append variables to the contact
	MergeVariables(context.Context, *MergeVariablesRequest) (*VariableList, error)
	// Reset all variables of the contact
	ResetVariables(context.Context, *ResetVariablesRequest) (*VariableList, error)
	// Remove variable(s) of the contact
	DeleteVariables(context.Context, *DeleteVariablesRequest) (*VariableList, error)
	// Update contact variable
	UpdateVariable(context.Context, *UpdateVariableRequest) (*VariableList, error)
	// Remove the contact's variable by etag
	DeleteVariable(context.Context, *DeleteVariableRequest) (*Variable, error)
	mustEmbedUnimplementedVariablesServer()
}

// UnimplementedVariablesServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVariablesServer struct{}

func (UnimplementedVariablesServer) ListVariables(context.Context, *SearchVariablesRequest) (*VariableList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVariables not implemented")
}
func (UnimplementedVariablesServer) MergeVariables(context.Context, *MergeVariablesRequest) (*VariableList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MergeVariables not implemented")
}
func (UnimplementedVariablesServer) ResetVariables(context.Context, *ResetVariablesRequest) (*VariableList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetVariables not implemented")
}
func (UnimplementedVariablesServer) DeleteVariables(context.Context, *DeleteVariablesRequest) (*VariableList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVariables not implemented")
}
func (UnimplementedVariablesServer) UpdateVariable(context.Context, *UpdateVariableRequest) (*VariableList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVariable not implemented")
}
func (UnimplementedVariablesServer) DeleteVariable(context.Context, *DeleteVariableRequest) (*Variable, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVariable not implemented")
}
func (UnimplementedVariablesServer) mustEmbedUnimplementedVariablesServer() {}
func (UnimplementedVariablesServer) testEmbeddedByValue()                   {}

// UnsafeVariablesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VariablesServer will
// result in compilation errors.
type UnsafeVariablesServer interface {
	mustEmbedUnimplementedVariablesServer()
}

func RegisterVariablesServer(s grpc.ServiceRegistrar, srv VariablesServer) {
	// If the following call pancis, it indicates UnimplementedVariablesServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Variables_ServiceDesc, srv)
}

func _Variables_ListVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchVariablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VariablesServer).ListVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Variables_ListVariables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VariablesServer).ListVariables(ctx, req.(*SearchVariablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Variables_MergeVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MergeVariablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VariablesServer).MergeVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Variables_MergeVariables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VariablesServer).MergeVariables(ctx, req.(*MergeVariablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Variables_ResetVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetVariablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VariablesServer).ResetVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Variables_ResetVariables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VariablesServer).ResetVariables(ctx, req.(*ResetVariablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Variables_DeleteVariables_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVariablesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VariablesServer).DeleteVariables(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Variables_DeleteVariables_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VariablesServer).DeleteVariables(ctx, req.(*DeleteVariablesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Variables_UpdateVariable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVariableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VariablesServer).UpdateVariable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Variables_UpdateVariable_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VariablesServer).UpdateVariable(ctx, req.(*UpdateVariableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Variables_DeleteVariable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVariableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VariablesServer).DeleteVariable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Variables_DeleteVariable_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VariablesServer).DeleteVariable(ctx, req.(*DeleteVariableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Variables_ServiceDesc is the grpc.ServiceDesc for Variables service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Variables_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.contacts.Variables",
	HandlerType: (*VariablesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListVariables",
			Handler:    _Variables_ListVariables_Handler,
		},
		{
			MethodName: "MergeVariables",
			Handler:    _Variables_MergeVariables_Handler,
		},
		{
			MethodName: "ResetVariables",
			Handler:    _Variables_ResetVariables_Handler,
		},
		{
			MethodName: "DeleteVariables",
			Handler:    _Variables_DeleteVariables_Handler,
		},
		{
			MethodName: "UpdateVariable",
			Handler:    _Variables_UpdateVariable_Handler,
		},
		{
			MethodName: "DeleteVariable",
			Handler:    _Variables_DeleteVariable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "contacts/variables.proto",
}
