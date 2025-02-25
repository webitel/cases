// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: queue_communications.proto

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

type CommunicationChannels int32

const (
	CommunicationChannels_Undefined CommunicationChannels = 0
	CommunicationChannels_Phone     CommunicationChannels = 1
	CommunicationChannels_Email     CommunicationChannels = 2
	CommunicationChannels_Messaging CommunicationChannels = 3
)

// Enum value maps for CommunicationChannels.
var (
	CommunicationChannels_name = map[int32]string{
		0: "Undefined",
		1: "Phone",
		2: "Email",
		3: "Messaging",
	}
	CommunicationChannels_value = map[string]int32{
		"Undefined": 0,
		"Phone":     1,
		"Email":     2,
		"Messaging": 3,
	}
)

func (x CommunicationChannels) Enum() *CommunicationChannels {
	p := new(CommunicationChannels)
	*p = x
	return p
}

func (x CommunicationChannels) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommunicationChannels) Descriptor() protoreflect.EnumDescriptor {
	return file_queue_communications_proto_enumTypes[0].Descriptor()
}

func (CommunicationChannels) Type() protoreflect.EnumType {
	return &file_queue_communications_proto_enumTypes[0]
}

func (x CommunicationChannels) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommunicationChannels.Descriptor instead.
func (CommunicationChannels) EnumDescriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{0}
}

type DeleteCommunicationTypeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId      int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCommunicationTypeRequest) Reset() {
	*x = DeleteCommunicationTypeRequest{}
	mi := &file_queue_communications_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCommunicationTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCommunicationTypeRequest) ProtoMessage() {}

func (x *DeleteCommunicationTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCommunicationTypeRequest.ProtoReflect.Descriptor instead.
func (*DeleteCommunicationTypeRequest) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteCommunicationTypeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DeleteCommunicationTypeRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

type UpdateCommunicationTypeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Code          string                 `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	Channel       CommunicationChannels  `protobuf:"varint,4,opt,name=channel,proto3,enum=engine.CommunicationChannels" json:"channel,omitempty"`
	Description   string                 `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Default       bool                   `protobuf:"varint,7,opt,name=default,proto3" json:"default,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCommunicationTypeRequest) Reset() {
	*x = UpdateCommunicationTypeRequest{}
	mi := &file_queue_communications_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCommunicationTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCommunicationTypeRequest) ProtoMessage() {}

func (x *UpdateCommunicationTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCommunicationTypeRequest.ProtoReflect.Descriptor instead.
func (*UpdateCommunicationTypeRequest) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateCommunicationTypeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateCommunicationTypeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateCommunicationTypeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *UpdateCommunicationTypeRequest) GetChannel() CommunicationChannels {
	if x != nil {
		return x.Channel
	}
	return CommunicationChannels_Undefined
}

func (x *UpdateCommunicationTypeRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateCommunicationTypeRequest) GetDefault() bool {
	if x != nil {
		return x.Default
	}
	return false
}

type PatchCommunicationTypeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Fields        []string               `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            int64                  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Code          string                 `protobuf:"bytes,5,opt,name=code,proto3" json:"code,omitempty"`
	Channel       CommunicationChannels  `protobuf:"varint,6,opt,name=channel,proto3,enum=engine.CommunicationChannels" json:"channel,omitempty"`
	Default       bool                   `protobuf:"varint,7,opt,name=default,proto3" json:"default,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PatchCommunicationTypeRequest) Reset() {
	*x = PatchCommunicationTypeRequest{}
	mi := &file_queue_communications_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PatchCommunicationTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PatchCommunicationTypeRequest) ProtoMessage() {}

func (x *PatchCommunicationTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PatchCommunicationTypeRequest.ProtoReflect.Descriptor instead.
func (*PatchCommunicationTypeRequest) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{2}
}

func (x *PatchCommunicationTypeRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *PatchCommunicationTypeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PatchCommunicationTypeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PatchCommunicationTypeRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PatchCommunicationTypeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *PatchCommunicationTypeRequest) GetChannel() CommunicationChannels {
	if x != nil {
		return x.Channel
	}
	return CommunicationChannels_Undefined
}

func (x *PatchCommunicationTypeRequest) GetDefault() bool {
	if x != nil {
		return x.Default
	}
	return false
}

type ListCommunicationType struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*CommunicationType   `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCommunicationType) Reset() {
	*x = ListCommunicationType{}
	mi := &file_queue_communications_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCommunicationType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommunicationType) ProtoMessage() {}

func (x *ListCommunicationType) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommunicationType.ProtoReflect.Descriptor instead.
func (*ListCommunicationType) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{3}
}

func (x *ListCommunicationType) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListCommunicationType) GetItems() []*CommunicationType {
	if x != nil {
		return x.Items
	}
	return nil
}

type SearchCommunicationTypeRequest struct {
	state         protoimpl.MessageState  `protogen:"open.v1"`
	Page          int32                   `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                   `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q             string                  `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	Sort          string                  `protobuf:"bytes,4,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields        []string                `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            []uint32                `protobuf:"varint,6,rep,packed,name=id,proto3" json:"id,omitempty"`
	Channel       []CommunicationChannels `protobuf:"varint,7,rep,packed,name=channel,proto3,enum=engine.CommunicationChannels" json:"channel,omitempty"`
	Default       bool                    `protobuf:"varint,8,opt,name=default,proto3" json:"default,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchCommunicationTypeRequest) Reset() {
	*x = SearchCommunicationTypeRequest{}
	mi := &file_queue_communications_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchCommunicationTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchCommunicationTypeRequest) ProtoMessage() {}

func (x *SearchCommunicationTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchCommunicationTypeRequest.ProtoReflect.Descriptor instead.
func (*SearchCommunicationTypeRequest) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{4}
}

func (x *SearchCommunicationTypeRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchCommunicationTypeRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchCommunicationTypeRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchCommunicationTypeRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchCommunicationTypeRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchCommunicationTypeRequest) GetId() []uint32 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SearchCommunicationTypeRequest) GetChannel() []CommunicationChannels {
	if x != nil {
		return x.Channel
	}
	return nil
}

func (x *SearchCommunicationTypeRequest) GetDefault() bool {
	if x != nil {
		return x.Default
	}
	return false
}

type ReadCommunicationTypeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId      int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadCommunicationTypeRequest) Reset() {
	*x = ReadCommunicationTypeRequest{}
	mi := &file_queue_communications_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadCommunicationTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadCommunicationTypeRequest) ProtoMessage() {}

func (x *ReadCommunicationTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadCommunicationTypeRequest.ProtoReflect.Descriptor instead.
func (*ReadCommunicationTypeRequest) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{5}
}

func (x *ReadCommunicationTypeRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ReadCommunicationTypeRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

type CommunicationTypeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Code          string                 `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Channel       CommunicationChannels  `protobuf:"varint,3,opt,name=channel,proto3,enum=engine.CommunicationChannels" json:"channel,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Default       bool                   `protobuf:"varint,6,opt,name=default,proto3" json:"default,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommunicationTypeRequest) Reset() {
	*x = CommunicationTypeRequest{}
	mi := &file_queue_communications_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommunicationTypeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommunicationTypeRequest) ProtoMessage() {}

func (x *CommunicationTypeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommunicationTypeRequest.ProtoReflect.Descriptor instead.
func (*CommunicationTypeRequest) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{6}
}

func (x *CommunicationTypeRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CommunicationTypeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *CommunicationTypeRequest) GetChannel() CommunicationChannels {
	if x != nil {
		return x.Channel
	}
	return CommunicationChannels_Undefined
}

func (x *CommunicationTypeRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CommunicationTypeRequest) GetDefault() bool {
	if x != nil {
		return x.Default
	}
	return false
}

type CommunicationType struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId      int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Code          string                 `protobuf:"bytes,4,opt,name=code,proto3" json:"code,omitempty"`
	Channel       CommunicationChannels  `protobuf:"varint,5,opt,name=channel,proto3,enum=engine.CommunicationChannels" json:"channel,omitempty"`
	Description   string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Default       bool                   `protobuf:"varint,7,opt,name=default,proto3" json:"default,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommunicationType) Reset() {
	*x = CommunicationType{}
	mi := &file_queue_communications_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommunicationType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommunicationType) ProtoMessage() {}

func (x *CommunicationType) ProtoReflect() protoreflect.Message {
	mi := &file_queue_communications_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommunicationType.ProtoReflect.Descriptor instead.
func (*CommunicationType) Descriptor() ([]byte, []int) {
	return file_queue_communications_proto_rawDescGZIP(), []int{7}
}

func (x *CommunicationType) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CommunicationType) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *CommunicationType) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CommunicationType) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *CommunicationType) GetChannel() CommunicationChannels {
	if x != nil {
		return x.Channel
	}
	return CommunicationChannels_Undefined
}

func (x *CommunicationType) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CommunicationType) GetDefault() bool {
	if x != nil {
		return x.Default
	}
	return false
}

var File_queue_communications_proto protoreflect.FileDescriptor

var file_queue_communications_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x1e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49,
	0x64, 0x22, 0xcd, 0x01, 0x0a, 0x1e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d,
	0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x37, 0x0a, 0x07,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x22, 0xe4, 0x01, 0x0a, 0x1d, 0x50, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x75,
	0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x37, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x18,
	0x0a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x22, 0x5c, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74,
	0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xe5, 0x01, 0x0a, 0x1e, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73,
	0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x37, 0x0a, 0x07, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x07, 0x63, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x22, 0x4b,
	0x0a, 0x1c, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x22, 0xb7, 0x01, 0x0a, 0x18,
	0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x37, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75,
	0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x22, 0xdd, 0x01, 0x0a, 0x11, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x37, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75,
	0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73,
	0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x2a, 0x4b, 0x0a, 0x15, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x12, 0x0d,
	0x0a, 0x09, 0x55, 0x6e, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x10, 0x00, 0x12, 0x09, 0x0a,
	0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69,
	0x6c, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67,
	0x10, 0x03, 0x32, 0xdf, 0x06, 0x0a, 0x18, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x82, 0x01, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x20, 0x2e, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24,
	0x3a, 0x01, 0x2a, 0x22, 0x1f, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65,
	0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x89, 0x01, 0x0a, 0x17, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43,
	0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x26, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x12,
	0x1f, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x86, 0x01, 0x0a, 0x15, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x24, 0x2e, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x19, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2c, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x26, 0x12, 0x24, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65,
	0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x8d, 0x01, 0x0a, 0x17, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x29,
	0x3a, 0x01, 0x2a, 0x1a, 0x24, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65,
	0x72, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x8b, 0x01, 0x0a, 0x16, 0x50, 0x61,
	0x74, 0x63, 0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x25, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x50, 0x61,
	0x74, 0x63, 0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x29, 0x3a, 0x01,
	0x2a, 0x32, 0x24, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x8a, 0x01, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x26, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x22, 0x2c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26, 0x2a, 0x24,
	0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x42, 0x80, 0x01, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x42, 0x18, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0xa2, 0x02, 0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0xca, 0x02, 0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0xe2, 0x02, 0x12, 0x45, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_queue_communications_proto_rawDescOnce sync.Once
	file_queue_communications_proto_rawDescData []byte
)

func file_queue_communications_proto_rawDescGZIP() []byte {
	file_queue_communications_proto_rawDescOnce.Do(func() {
		file_queue_communications_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_queue_communications_proto_rawDesc), len(file_queue_communications_proto_rawDesc)))
	})
	return file_queue_communications_proto_rawDescData
}

var file_queue_communications_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_queue_communications_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_queue_communications_proto_goTypes = []any{
	(CommunicationChannels)(0),             // 0: engine.CommunicationChannels
	(*DeleteCommunicationTypeRequest)(nil), // 1: engine.DeleteCommunicationTypeRequest
	(*UpdateCommunicationTypeRequest)(nil), // 2: engine.UpdateCommunicationTypeRequest
	(*PatchCommunicationTypeRequest)(nil),  // 3: engine.PatchCommunicationTypeRequest
	(*ListCommunicationType)(nil),          // 4: engine.ListCommunicationType
	(*SearchCommunicationTypeRequest)(nil), // 5: engine.SearchCommunicationTypeRequest
	(*ReadCommunicationTypeRequest)(nil),   // 6: engine.ReadCommunicationTypeRequest
	(*CommunicationTypeRequest)(nil),       // 7: engine.CommunicationTypeRequest
	(*CommunicationType)(nil),              // 8: engine.CommunicationType
}
var file_queue_communications_proto_depIdxs = []int32{
	0,  // 0: engine.UpdateCommunicationTypeRequest.channel:type_name -> engine.CommunicationChannels
	0,  // 1: engine.PatchCommunicationTypeRequest.channel:type_name -> engine.CommunicationChannels
	8,  // 2: engine.ListCommunicationType.items:type_name -> engine.CommunicationType
	0,  // 3: engine.SearchCommunicationTypeRequest.channel:type_name -> engine.CommunicationChannels
	0,  // 4: engine.CommunicationTypeRequest.channel:type_name -> engine.CommunicationChannels
	0,  // 5: engine.CommunicationType.channel:type_name -> engine.CommunicationChannels
	7,  // 6: engine.CommunicationTypeService.CreateCommunicationType:input_type -> engine.CommunicationTypeRequest
	5,  // 7: engine.CommunicationTypeService.SearchCommunicationType:input_type -> engine.SearchCommunicationTypeRequest
	6,  // 8: engine.CommunicationTypeService.ReadCommunicationType:input_type -> engine.ReadCommunicationTypeRequest
	2,  // 9: engine.CommunicationTypeService.UpdateCommunicationType:input_type -> engine.UpdateCommunicationTypeRequest
	3,  // 10: engine.CommunicationTypeService.PatchCommunicationType:input_type -> engine.PatchCommunicationTypeRequest
	1,  // 11: engine.CommunicationTypeService.DeleteCommunicationType:input_type -> engine.DeleteCommunicationTypeRequest
	8,  // 12: engine.CommunicationTypeService.CreateCommunicationType:output_type -> engine.CommunicationType
	4,  // 13: engine.CommunicationTypeService.SearchCommunicationType:output_type -> engine.ListCommunicationType
	8,  // 14: engine.CommunicationTypeService.ReadCommunicationType:output_type -> engine.CommunicationType
	8,  // 15: engine.CommunicationTypeService.UpdateCommunicationType:output_type -> engine.CommunicationType
	8,  // 16: engine.CommunicationTypeService.PatchCommunicationType:output_type -> engine.CommunicationType
	8,  // 17: engine.CommunicationTypeService.DeleteCommunicationType:output_type -> engine.CommunicationType
	12, // [12:18] is the sub-list for method output_type
	6,  // [6:12] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_queue_communications_proto_init() }
func file_queue_communications_proto_init() {
	if File_queue_communications_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_queue_communications_proto_rawDesc), len(file_queue_communications_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_queue_communications_proto_goTypes,
		DependencyIndexes: file_queue_communications_proto_depIdxs,
		EnumInfos:         file_queue_communications_proto_enumTypes,
		MessageInfos:      file_queue_communications_proto_msgTypes,
	}.Build()
	File_queue_communications_proto = out.File
	file_queue_communications_proto_goTypes = nil
	file_queue_communications_proto_depIdxs = nil
}
