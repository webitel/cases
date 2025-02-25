// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: regiong.proto

package engine

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

type Region struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Timezone      *Lookup                `protobuf:"bytes,4,opt,name=timezone,proto3" json:"timezone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Region) Reset() {
	*x = Region{}
	mi := &file_regiong_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Region) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Region) ProtoMessage() {}

func (x *Region) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Region.ProtoReflect.Descriptor instead.
func (*Region) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{0}
}

func (x *Region) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Region) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Region) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Region) GetTimezone() *Lookup {
	if x != nil {
		return x.Timezone
	}
	return nil
}

type CreateRegionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Timezone      *Lookup                `protobuf:"bytes,3,opt,name=timezone,proto3" json:"timezone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateRegionRequest) Reset() {
	*x = CreateRegionRequest{}
	mi := &file_regiong_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRegionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRegionRequest) ProtoMessage() {}

func (x *CreateRegionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRegionRequest.ProtoReflect.Descriptor instead.
func (*CreateRegionRequest) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRegionRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRegionRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateRegionRequest) GetTimezone() *Lookup {
	if x != nil {
		return x.Timezone
	}
	return nil
}

type SearchRegionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q             string                 `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	Sort          string                 `protobuf:"bytes,4,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields        []string               `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            []int64                `protobuf:"varint,6,rep,packed,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,8,opt,name=description,proto3" json:"description,omitempty"`
	TimezoneId    []uint32               `protobuf:"varint,9,rep,packed,name=timezone_id,json=timezoneId,proto3" json:"timezone_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchRegionRequest) Reset() {
	*x = SearchRegionRequest{}
	mi := &file_regiong_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchRegionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRegionRequest) ProtoMessage() {}

func (x *SearchRegionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRegionRequest.ProtoReflect.Descriptor instead.
func (*SearchRegionRequest) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{2}
}

func (x *SearchRegionRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchRegionRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchRegionRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchRegionRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchRegionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchRegionRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SearchRegionRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SearchRegionRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SearchRegionRequest) GetTimezoneId() []uint32 {
	if x != nil {
		return x.TimezoneId
	}
	return nil
}

type ListRegion struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*Region              `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListRegion) Reset() {
	*x = ListRegion{}
	mi := &file_regiong_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRegion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRegion) ProtoMessage() {}

func (x *ListRegion) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRegion.ProtoReflect.Descriptor instead.
func (*ListRegion) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{3}
}

func (x *ListRegion) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListRegion) GetItems() []*Region {
	if x != nil {
		return x.Items
	}
	return nil
}

type ReadRegionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadRegionRequest) Reset() {
	*x = ReadRegionRequest{}
	mi := &file_regiong_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadRegionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadRegionRequest) ProtoMessage() {}

func (x *ReadRegionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadRegionRequest.ProtoReflect.Descriptor instead.
func (*ReadRegionRequest) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{4}
}

func (x *ReadRegionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type PatchRegionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Fields        []string               `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            int64                  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Timezone      *Lookup                `protobuf:"bytes,5,opt,name=timezone,proto3" json:"timezone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchRegionRequest) Reset() {
	*x = PatchRegionRequest{}
	mi := &file_regiong_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchRegionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchRegionRequest) ProtoMessage() {}

func (x *PatchRegionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchRegionRequest.ProtoReflect.Descriptor instead.
func (*PatchRegionRequest) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{5}
}

func (x *PatchRegionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *PatchRegionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PatchRegionRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PatchRegionRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PatchRegionRequest) GetTimezone() *Lookup {
	if x != nil {
		return x.Timezone
	}
	return nil
}

type UpdateRegionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Timezone      *Lookup                `protobuf:"bytes,4,opt,name=timezone,proto3" json:"timezone,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateRegionRequest) Reset() {
	*x = UpdateRegionRequest{}
	mi := &file_regiong_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateRegionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRegionRequest) ProtoMessage() {}

func (x *UpdateRegionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRegionRequest.ProtoReflect.Descriptor instead.
func (*UpdateRegionRequest) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateRegionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateRegionRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateRegionRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateRegionRequest) GetTimezone() *Lookup {
	if x != nil {
		return x.Timezone
	}
	return nil
}

type DeleteRegionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteRegionRequest) Reset() {
	*x = DeleteRegionRequest{}
	mi := &file_regiong_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRegionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRegionRequest) ProtoMessage() {}

func (x *DeleteRegionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_regiong_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRegionRequest.ProtoReflect.Descriptor instead.
func (*DeleteRegionRequest) Descriptor() ([]byte, []int) {
	return file_regiong_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteRegionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_regiong_proto protoreflect.FileDescriptor

var file_regiong_proto_rawDesc = string([]byte{
	0x0a, 0x0d, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x7a, 0x0a, 0x06, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x22, 0x77,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x08, 0x74,
	0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x08, 0x74,
	0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x22, 0xde, 0x01, 0x0a, 0x13, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70,
	0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x7a,
	0x6f, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x0a, 0x74, 0x69,
	0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x49, 0x64, 0x22, 0x46, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x24, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x22, 0x23, 0x0a, 0x11, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x9e, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x08, 0x74, 0x69,
	0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x08, 0x74, 0x69,
	0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65, 0x22, 0x87, 0x01, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e,
	0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x7a, 0x6f, 0x6e, 0x65,
	0x22, 0x25, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x32, 0x84, 0x04, 0x0a, 0x0d, 0x52, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x0c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e,
	0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x3a, 0x01,
	0x2a, 0x22, 0x08, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x51, 0x0a, 0x0c, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x67, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x10, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x4e,
	0x0a, 0x0a, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x19, 0x2e, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12,
	0x0d, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x53,
	0x0a, 0x0b, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x2e,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x50, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x12, 0x3a, 0x01, 0x2a, 0x32, 0x0d, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x12, 0x55, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x67,
	0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x3a, 0x01, 0x2a, 0x1a, 0x0d, 0x2f, 0x72, 0x65,
	0x67, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x52, 0x0a, 0x0c, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x15, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x2a,
	0x0d, 0x2f, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x74,
	0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x42, 0x0c, 0x52, 0x65,
	0x67, 0x69, 0x6f, 0x6e, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x20, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0xa2, 0x02,
	0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0xca, 0x02, 0x06,
	0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0xe2, 0x02, 0x12, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x06, 0x45, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_regiong_proto_rawDescOnce sync.Once
	file_regiong_proto_rawDescData []byte
)

func file_regiong_proto_rawDescGZIP() []byte {
	file_regiong_proto_rawDescOnce.Do(func() {
		file_regiong_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_regiong_proto_rawDesc), len(file_regiong_proto_rawDesc)))
	})
	return file_regiong_proto_rawDescData
}

var file_regiong_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_regiong_proto_goTypes = []any{
	(*Region)(nil),              // 0: engine.Region
	(*CreateRegionRequest)(nil), // 1: engine.CreateRegionRequest
	(*SearchRegionRequest)(nil), // 2: engine.SearchRegionRequest
	(*ListRegion)(nil),          // 3: engine.ListRegion
	(*ReadRegionRequest)(nil),   // 4: engine.ReadRegionRequest
	(*PatchRegionRequest)(nil),  // 5: engine.PatchRegionRequest
	(*UpdateRegionRequest)(nil), // 6: engine.UpdateRegionRequest
	(*DeleteRegionRequest)(nil), // 7: engine.DeleteRegionRequest
	(*Lookup)(nil),              // 8: engine.Lookup
}
var file_regiong_proto_depIdxs = []int32{
	8,  // 0: engine.Region.timezone:type_name -> engine.Lookup
	8,  // 1: engine.CreateRegionRequest.timezone:type_name -> engine.Lookup
	0,  // 2: engine.ListRegion.items:type_name -> engine.Region
	8,  // 3: engine.PatchRegionRequest.timezone:type_name -> engine.Lookup
	8,  // 4: engine.UpdateRegionRequest.timezone:type_name -> engine.Lookup
	1,  // 5: engine.RegionService.CreateRegion:input_type -> engine.CreateRegionRequest
	2,  // 6: engine.RegionService.SearchRegion:input_type -> engine.SearchRegionRequest
	4,  // 7: engine.RegionService.ReadRegion:input_type -> engine.ReadRegionRequest
	5,  // 8: engine.RegionService.PatchRegion:input_type -> engine.PatchRegionRequest
	6,  // 9: engine.RegionService.UpdateRegion:input_type -> engine.UpdateRegionRequest
	7,  // 10: engine.RegionService.DeleteRegion:input_type -> engine.DeleteRegionRequest
	0,  // 11: engine.RegionService.CreateRegion:output_type -> engine.Region
	3,  // 12: engine.RegionService.SearchRegion:output_type -> engine.ListRegion
	0,  // 13: engine.RegionService.ReadRegion:output_type -> engine.Region
	0,  // 14: engine.RegionService.PatchRegion:output_type -> engine.Region
	0,  // 15: engine.RegionService.UpdateRegion:output_type -> engine.Region
	0,  // 16: engine.RegionService.DeleteRegion:output_type -> engine.Region
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_regiong_proto_init() }
func file_regiong_proto_init() {
	if File_regiong_proto != nil {
		return
	}
	file_const_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_regiong_proto_rawDesc), len(file_regiong_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_regiong_proto_goTypes,
		DependencyIndexes: file_regiong_proto_depIdxs,
		MessageInfos:      file_regiong_proto_msgTypes,
	}.Build()
	File_regiong_proto = out.File
	file_regiong_proto_goTypes = nil
	file_regiong_proto_depIdxs = nil
}
