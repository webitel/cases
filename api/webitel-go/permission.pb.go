// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: permission.proto

package api

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Permission struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`       // [required] e.g.: 'read', 'playback_record_file', ...
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`   // [optional] e.g.: 'Select Any'
	Usage         string                 `protobuf:"bytes,3,opt,name=usage,proto3" json:"usage,omitempty"` // [optional] e.g.: 'Grants permission to select any objects'
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Permission) Reset() {
	*x = Permission{}
	mi := &file_permission_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Permission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Permission) ProtoMessage() {}

func (x *Permission) ProtoReflect() protoreflect.Message {
	mi := &file_permission_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Permission.ProtoReflect.Descriptor instead.
func (*Permission) Descriptor() ([]byte, []int) {
	return file_permission_proto_rawDescGZIP(), []int{0}
}

func (x *Permission) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Permission) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Permission) GetUsage() string {
	if x != nil {
		return x.Usage
	}
	return ""
}

type SearchPermissionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Q     string                 `protobuf:"bytes,1,opt,name=q,proto3" json:"q,omitempty"`     // TermOfSearch
	Id    []string               `protobuf:"bytes,2,rep,name=id,proto3" json:"id,omitempty"`   // = ANY(id)
	Not   []string               `protobuf:"bytes,3,rep,name=not,proto3" json:"not,omitempty"` // <> ANY(not)
	// controls
	Fields        []string `protobuf:"bytes,10,rep,name=fields,proto3" json:"fields,omitempty"` // output
	Sort          []string `protobuf:"bytes,11,rep,name=sort,proto3" json:"sort,omitempty"`     // sort: "field" asc; "!field" desc
	Page          int32    `protobuf:"varint,12,opt,name=page,proto3" json:"page,omitempty"`    // page number
	Size          int32    `protobuf:"varint,13,opt,name=size,proto3" json:"size,omitempty"`    // page size
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchPermissionRequest) Reset() {
	*x = SearchPermissionRequest{}
	mi := &file_permission_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchPermissionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchPermissionRequest) ProtoMessage() {}

func (x *SearchPermissionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_permission_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchPermissionRequest.ProtoReflect.Descriptor instead.
func (*SearchPermissionRequest) Descriptor() ([]byte, []int) {
	return file_permission_proto_rawDescGZIP(), []int{1}
}

func (x *SearchPermissionRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchPermissionRequest) GetId() []string {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SearchPermissionRequest) GetNot() []string {
	if x != nil {
		return x.Not
	}
	return nil
}

func (x *SearchPermissionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchPermissionRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *SearchPermissionRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchPermissionRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

type SearchPermissionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,10,opt,name=page,proto3" json:"page,omitempty"` // select: offset {page}
	Size          int32                  `protobuf:"varint,11,opt,name=size,proto3" json:"size,omitempty"` // search: limit {size}
	Next          bool                   `protobuf:"varint,12,opt,name=next,proto3" json:"next,omitempty"` // search: has {next} page ?
	Items         []*Permission          `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"` // result entries
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchPermissionResponse) Reset() {
	*x = SearchPermissionResponse{}
	mi := &file_permission_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchPermissionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchPermissionResponse) ProtoMessage() {}

func (x *SearchPermissionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_permission_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchPermissionResponse.ProtoReflect.Descriptor instead.
func (*SearchPermissionResponse) Descriptor() ([]byte, []int) {
	return file_permission_proto_rawDescGZIP(), []int{2}
}

func (x *SearchPermissionResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchPermissionResponse) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchPermissionResponse) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *SearchPermissionResponse) GetItems() []*Permission {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_permission_proto protoreflect.FileDescriptor

var file_permission_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x0a, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x22, 0x9d, 0x01,
	0x0a, 0x17, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x6f, 0x74, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x6e, 0x6f, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x22, 0x7d, 0x0a,
	0x18, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0x72, 0x0a, 0x0b,
	0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x63, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1c, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x42, 0x5a, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x42, 0x0f, 0x50, 0x65, 0x72,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x12,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x61,
	0x70, 0x69, 0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x03, 0x41, 0x70, 0x69, 0xca, 0x02,
	0x03, 0x41, 0x70, 0x69, 0xe2, 0x02, 0x0f, 0x41, 0x70, 0x69, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x03, 0x41, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_permission_proto_rawDescOnce sync.Once
	file_permission_proto_rawDescData []byte
)

func file_permission_proto_rawDescGZIP() []byte {
	file_permission_proto_rawDescOnce.Do(func() {
		file_permission_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_permission_proto_rawDesc), len(file_permission_proto_rawDesc)))
	})
	return file_permission_proto_rawDescData
}

var file_permission_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_permission_proto_goTypes = []any{
	(*Permission)(nil),               // 0: api.Permission
	(*SearchPermissionRequest)(nil),  // 1: api.SearchPermissionRequest
	(*SearchPermissionResponse)(nil), // 2: api.SearchPermissionResponse
}
var file_permission_proto_depIdxs = []int32{
	0, // 0: api.SearchPermissionResponse.items:type_name -> api.Permission
	1, // 1: api.Permissions.GetPermissions:input_type -> api.SearchPermissionRequest
	2, // 2: api.Permissions.GetPermissions:output_type -> api.SearchPermissionResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_permission_proto_init() }
func file_permission_proto_init() {
	if File_permission_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_permission_proto_rawDesc), len(file_permission_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_permission_proto_goTypes,
		DependencyIndexes: file_permission_proto_depIdxs,
		MessageInfos:      file_permission_proto_msgTypes,
	}.Build()
	File_permission_proto = out.File
	file_permission_proto_goTypes = nil
	file_permission_proto_depIdxs = nil
}
