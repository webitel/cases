// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: contacts/photos.proto

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
	Photos_UploadPhoto_FullMethodName  = "/webitel.contacts.Photos/UploadPhoto"
	Photos_LocatePhoto_FullMethodName  = "/webitel.contacts.Photos/LocatePhoto"
	Photos_SearchPhotos_FullMethodName = "/webitel.contacts.Photos/SearchPhotos"
	Photos_CreatePhotos_FullMethodName = "/webitel.contacts.Photos/CreatePhotos"
	Photos_UpdatePhotos_FullMethodName = "/webitel.contacts.Photos/UpdatePhotos"
	Photos_UpdatePhoto_FullMethodName  = "/webitel.contacts.Photos/UpdatePhoto"
	Photos_DeletePhotos_FullMethodName = "/webitel.contacts.Photos/DeletePhotos"
	Photos_DeletePhoto_FullMethodName  = "/webitel.contacts.Photos/DeletePhoto"
)

// PhotosClient is the client API for Photos service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Photos service catalog.
type PhotosClient interface {
	// Upload an image or photo
	UploadPhoto(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[UploadMediaRequest, UploadMediaResponse], error)
	// Locate the contact's photo link.
	LocatePhoto(ctx context.Context, in *LocatePhotoRequest, opts ...grpc.CallOption) (*Photo, error)
	// Search the contact's photo(s)
	SearchPhotos(ctx context.Context, in *SearchPhotosRequest, opts ...grpc.CallOption) (*PhotoList, error)
	// Link photo(s) with the contact
	CreatePhotos(ctx context.Context, in *CreatePhotosRequest, opts ...grpc.CallOption) (*PhotoList, error)
	// Reset the contact's photos to fit given data set.
	UpdatePhotos(ctx context.Context, in *UpdatePhotosRequest, opts ...grpc.CallOption) (*PhotoList, error)
	// Update the contact's photo link details
	UpdatePhoto(ctx context.Context, in *UpdatePhotoRequest, opts ...grpc.CallOption) (*Photo, error)
	// Remove the contact's photo link(s)
	DeletePhotos(ctx context.Context, in *DeletePhotosRequest, opts ...grpc.CallOption) (*PhotoList, error)
	// Remove the contact's photo
	DeletePhoto(ctx context.Context, in *DeletePhotoRequest, opts ...grpc.CallOption) (*Photo, error)
}

type photosClient struct {
	cc grpc.ClientConnInterface
}

func NewPhotosClient(cc grpc.ClientConnInterface) PhotosClient {
	return &photosClient{cc}
}

func (c *photosClient) UploadPhoto(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[UploadMediaRequest, UploadMediaResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Photos_ServiceDesc.Streams[0], Photos_UploadPhoto_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadMediaRequest, UploadMediaResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Photos_UploadPhotoClient = grpc.BidiStreamingClient[UploadMediaRequest, UploadMediaResponse]

func (c *photosClient) LocatePhoto(ctx context.Context, in *LocatePhotoRequest, opts ...grpc.CallOption) (*Photo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Photo)
	err := c.cc.Invoke(ctx, Photos_LocatePhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photosClient) SearchPhotos(ctx context.Context, in *SearchPhotosRequest, opts ...grpc.CallOption) (*PhotoList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PhotoList)
	err := c.cc.Invoke(ctx, Photos_SearchPhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photosClient) CreatePhotos(ctx context.Context, in *CreatePhotosRequest, opts ...grpc.CallOption) (*PhotoList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PhotoList)
	err := c.cc.Invoke(ctx, Photos_CreatePhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photosClient) UpdatePhotos(ctx context.Context, in *UpdatePhotosRequest, opts ...grpc.CallOption) (*PhotoList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PhotoList)
	err := c.cc.Invoke(ctx, Photos_UpdatePhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photosClient) UpdatePhoto(ctx context.Context, in *UpdatePhotoRequest, opts ...grpc.CallOption) (*Photo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Photo)
	err := c.cc.Invoke(ctx, Photos_UpdatePhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photosClient) DeletePhotos(ctx context.Context, in *DeletePhotosRequest, opts ...grpc.CallOption) (*PhotoList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PhotoList)
	err := c.cc.Invoke(ctx, Photos_DeletePhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photosClient) DeletePhoto(ctx context.Context, in *DeletePhotoRequest, opts ...grpc.CallOption) (*Photo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Photo)
	err := c.cc.Invoke(ctx, Photos_DeletePhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PhotosServer is the server API for Photos service.
// All implementations must embed UnimplementedPhotosServer
// for forward compatibility.
//
// Photos service catalog.
type PhotosServer interface {
	// Upload an image or photo
	UploadPhoto(grpc.BidiStreamingServer[UploadMediaRequest, UploadMediaResponse]) error
	// Locate the contact's photo link.
	LocatePhoto(context.Context, *LocatePhotoRequest) (*Photo, error)
	// Search the contact's photo(s)
	SearchPhotos(context.Context, *SearchPhotosRequest) (*PhotoList, error)
	// Link photo(s) with the contact
	CreatePhotos(context.Context, *CreatePhotosRequest) (*PhotoList, error)
	// Reset the contact's photos to fit given data set.
	UpdatePhotos(context.Context, *UpdatePhotosRequest) (*PhotoList, error)
	// Update the contact's photo link details
	UpdatePhoto(context.Context, *UpdatePhotoRequest) (*Photo, error)
	// Remove the contact's photo link(s)
	DeletePhotos(context.Context, *DeletePhotosRequest) (*PhotoList, error)
	// Remove the contact's photo
	DeletePhoto(context.Context, *DeletePhotoRequest) (*Photo, error)
	mustEmbedUnimplementedPhotosServer()
}

// UnimplementedPhotosServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPhotosServer struct{}

func (UnimplementedPhotosServer) UploadPhoto(grpc.BidiStreamingServer[UploadMediaRequest, UploadMediaResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadPhoto not implemented")
}
func (UnimplementedPhotosServer) LocatePhoto(context.Context, *LocatePhotoRequest) (*Photo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LocatePhoto not implemented")
}
func (UnimplementedPhotosServer) SearchPhotos(context.Context, *SearchPhotosRequest) (*PhotoList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchPhotos not implemented")
}
func (UnimplementedPhotosServer) CreatePhotos(context.Context, *CreatePhotosRequest) (*PhotoList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePhotos not implemented")
}
func (UnimplementedPhotosServer) UpdatePhotos(context.Context, *UpdatePhotosRequest) (*PhotoList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePhotos not implemented")
}
func (UnimplementedPhotosServer) UpdatePhoto(context.Context, *UpdatePhotoRequest) (*Photo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePhoto not implemented")
}
func (UnimplementedPhotosServer) DeletePhotos(context.Context, *DeletePhotosRequest) (*PhotoList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePhotos not implemented")
}
func (UnimplementedPhotosServer) DeletePhoto(context.Context, *DeletePhotoRequest) (*Photo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePhoto not implemented")
}
func (UnimplementedPhotosServer) mustEmbedUnimplementedPhotosServer() {}
func (UnimplementedPhotosServer) testEmbeddedByValue()                {}

// UnsafePhotosServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PhotosServer will
// result in compilation errors.
type UnsafePhotosServer interface {
	mustEmbedUnimplementedPhotosServer()
}

func RegisterPhotosServer(s grpc.ServiceRegistrar, srv PhotosServer) {
	// If the following call pancis, it indicates UnimplementedPhotosServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Photos_ServiceDesc, srv)
}

func _Photos_UploadPhoto_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PhotosServer).UploadPhoto(&grpc.GenericServerStream[UploadMediaRequest, UploadMediaResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Photos_UploadPhotoServer = grpc.BidiStreamingServer[UploadMediaRequest, UploadMediaResponse]

func _Photos_LocatePhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocatePhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).LocatePhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_LocatePhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).LocatePhoto(ctx, req.(*LocatePhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Photos_SearchPhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchPhotosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).SearchPhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_SearchPhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).SearchPhotos(ctx, req.(*SearchPhotosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Photos_CreatePhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePhotosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).CreatePhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_CreatePhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).CreatePhotos(ctx, req.(*CreatePhotosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Photos_UpdatePhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePhotosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).UpdatePhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_UpdatePhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).UpdatePhotos(ctx, req.(*UpdatePhotosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Photos_UpdatePhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).UpdatePhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_UpdatePhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).UpdatePhoto(ctx, req.(*UpdatePhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Photos_DeletePhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePhotosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).DeletePhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_DeletePhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).DeletePhotos(ctx, req.(*DeletePhotosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Photos_DeletePhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotosServer).DeletePhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Photos_DeletePhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotosServer).DeletePhoto(ctx, req.(*DeletePhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Photos_ServiceDesc is the grpc.ServiceDesc for Photos service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Photos_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel.contacts.Photos",
	HandlerType: (*PhotosServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LocatePhoto",
			Handler:    _Photos_LocatePhoto_Handler,
		},
		{
			MethodName: "SearchPhotos",
			Handler:    _Photos_SearchPhotos_Handler,
		},
		{
			MethodName: "CreatePhotos",
			Handler:    _Photos_CreatePhotos_Handler,
		},
		{
			MethodName: "UpdatePhotos",
			Handler:    _Photos_UpdatePhotos_Handler,
		},
		{
			MethodName: "UpdatePhoto",
			Handler:    _Photos_UpdatePhoto_Handler,
		},
		{
			MethodName: "DeletePhotos",
			Handler:    _Photos_DeletePhotos_Handler,
		},
		{
			MethodName: "DeletePhoto",
			Handler:    _Photos_DeletePhoto_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadPhoto",
			Handler:       _Photos_UploadPhoto_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "contacts/photos.proto",
}
