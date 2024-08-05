// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: close_reason.proto

package api

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CloseReason message represents a close reason entity with metadata
type CloseReason struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the close reason
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the close reason
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the close reason
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// CreatedAt timestamp of the close reason
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the close reason
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the close reason
	CreatedBy *Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the close reason
	UpdatedBy *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
}

func (x *CloseReason) Reset() {
	*x = CloseReason{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseReason) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseReason) ProtoMessage() {}

func (x *CloseReason) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseReason.ProtoReflect.Descriptor instead.
func (*CloseReason) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{0}
}

func (x *CloseReason) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CloseReason) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CloseReason) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CloseReason) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *CloseReason) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *CloseReason) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *CloseReason) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// CloseReasonList message contains a list of CloseReason items with pagination
type CloseReasonList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int32          `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next  bool           `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items []*CloseReason `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *CloseReasonList) Reset() {
	*x = CloseReasonList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloseReasonList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseReasonList) ProtoMessage() {}

func (x *CloseReasonList) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseReasonList.ProtoReflect.Descriptor instead.
func (*CloseReasonList) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{1}
}

func (x *CloseReasonList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *CloseReasonList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *CloseReasonList) GetItems() []*CloseReason {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreateCloseReasonRequest message for creating a new close reason
type CreateCloseReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *CreateCloseReasonRequest) Reset() {
	*x = CreateCloseReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCloseReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCloseReasonRequest) ProtoMessage() {}

func (x *CreateCloseReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCloseReasonRequest.ProtoReflect.Descriptor instead.
func (*CreateCloseReasonRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{2}
}

func (x *CreateCloseReasonRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCloseReasonRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// UpdateCloseReasonRequest message for updating an existing close reason
type UpdateCloseReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *UpdateCloseReasonRequest) Reset() {
	*x = UpdateCloseReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCloseReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCloseReasonRequest) ProtoMessage() {}

func (x *UpdateCloseReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCloseReasonRequest.ProtoReflect.Descriptor instead.
func (*UpdateCloseReasonRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateCloseReasonRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateCloseReasonRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateCloseReasonRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// DeleteCloseReasonRequest message for deleting an existing close reason
type DeleteCloseReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteCloseReasonRequest) Reset() {
	*x = DeleteCloseReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteCloseReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCloseReasonRequest) ProtoMessage() {}

func (x *DeleteCloseReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCloseReasonRequest.ProtoReflect.Descriptor instead.
func (*DeleteCloseReasonRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteCloseReasonRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// ListCloseReasonsRequest message for listing or searching close reasons
type ListCloseReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Page number of result dataset records. offset = (page*size)
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Size count of records on result page. limit = (size++)
	Size int32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Sort the result according to fields.
	Sort []string `protobuf:"bytes,4,rep,name=sort,proto3" json:"sort,omitempty"`
	// Filter by unique IDs.
	Id []int64 `protobuf:"varint,5,rep,packed,name=id,proto3" json:"id,omitempty"`
	// Search term: close reason name;
	// `?` - matches any one character
	// `*` - matches 0 or more characters
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Filter by close reason name.
	Name string `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ListCloseReasonRequest) Reset() {
	*x = ListCloseReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListCloseReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCloseReasonRequest) ProtoMessage() {}

func (x *ListCloseReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCloseReasonRequest.ProtoReflect.Descriptor instead.
func (*ListCloseReasonRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{5}
}

func (x *ListCloseReasonRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCloseReasonRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListCloseReasonRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListCloseReasonRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListCloseReasonRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListCloseReasonRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListCloseReasonRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// LocateCloseReasonRequest message for locating a specific close reason by ID
type LocateCloseReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *LocateCloseReasonRequest) Reset() {
	*x = LocateCloseReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocateCloseReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateCloseReasonRequest) ProtoMessage() {}

func (x *LocateCloseReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateCloseReasonRequest.ProtoReflect.Descriptor instead.
func (*LocateCloseReasonRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{6}
}

func (x *LocateCloseReasonRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateCloseReasonRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocateCloseReasonResponse message contains a single close reason entity
type LocateCloseReasonResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CloseReason *CloseReason `protobuf:"bytes,1,opt,name=close_reason,json=closeReason,proto3" json:"close_reason,omitempty"`
}

func (x *LocateCloseReasonResponse) Reset() {
	*x = LocateCloseReasonResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_close_reason_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocateCloseReasonResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateCloseReasonResponse) ProtoMessage() {}

func (x *LocateCloseReasonResponse) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateCloseReasonResponse.ProtoReflect.Descriptor instead.
func (*LocateCloseReasonResponse) Descriptor() ([]byte, []int) {
	return file_close_reason_proto_rawDescGZIP(), []int{7}
}

func (x *LocateCloseReasonResponse) GetCloseReason() *CloseReason {
	if x != nil {
		return x.CloseReason
	}
	return nil
}

var File_close_reason_proto protoreflect.FileDescriptor

var file_close_reason_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x0c, 0x6c, 0x6f, 0x6f,
	0x6b, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d,
	0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xed, 0x01, 0x0a, 0x0b, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2c, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x2c, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x22, 0x63, 0x0a, 0x0f, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65,
	0x78, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x5e, 0x0a, 0x18,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x0c,
	0x92, 0x41, 0x09, 0x0a, 0x07, 0xd2, 0x01, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x6c, 0x0a, 0x18,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x0a,
	0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x36, 0x0a, 0x18, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02,
	0x69, 0x64, 0x22, 0x9e, 0x01, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0x42, 0x0a, 0x18, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f,
	0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x52, 0x0a, 0x19, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x0c, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x0b,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x32, 0xed, 0x05, 0x0a, 0x0c,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x12, 0x9d, 0x01, 0x0a,
	0x10, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x73, 0x12, 0x1d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c,
	0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x52, 0x92, 0x41, 0x3a, 0x12, 0x38, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f,
	0x66, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x20,
	0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20,
	0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x12, 0x80, 0x01, 0x0a,
	0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x36, 0x92, 0x41, 0x1b, 0x12, 0x19, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65,
	0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x3a, 0x01, 0x2a,
	0x22, 0x0d, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x12,
	0xa4, 0x01, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x5a, 0x92, 0x41, 0x21, 0x12,
	0x1f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74,
	0x69, 0x6e, 0x67, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x3a, 0x01, 0x2a, 0x5a, 0x17, 0x3a, 0x01, 0x2a, 0x32, 0x12,
	0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x1a, 0x12, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x7e, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x22, 0x34, 0x92, 0x41, 0x17, 0x12, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61, 0x20,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x14, 0x2a, 0x12, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x92, 0x01, 0x0a, 0x11, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x3a, 0x92, 0x41, 0x1d, 0x12, 0x1b, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x63,
	0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20, 0x62, 0x79, 0x20, 0x49,
	0x44, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x0b, 0x5a, 0x09, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_close_reason_proto_rawDescOnce sync.Once
	file_close_reason_proto_rawDescData = file_close_reason_proto_rawDesc
)

func file_close_reason_proto_rawDescGZIP() []byte {
	file_close_reason_proto_rawDescOnce.Do(func() {
		file_close_reason_proto_rawDescData = protoimpl.X.CompressGZIP(file_close_reason_proto_rawDescData)
	})
	return file_close_reason_proto_rawDescData
}

var file_close_reason_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_close_reason_proto_goTypes = []any{
	(*CloseReason)(nil),               // 0: cases.CloseReason
	(*CloseReasonList)(nil),           // 1: cases.CloseReasonList
	(*CreateCloseReasonRequest)(nil),  // 2: cases.CreateCloseReasonRequest
	(*UpdateCloseReasonRequest)(nil),  // 3: cases.UpdateCloseReasonRequest
	(*DeleteCloseReasonRequest)(nil),  // 4: cases.DeleteCloseReasonRequest
	(*ListCloseReasonRequest)(nil),    // 5: cases.ListCloseReasonRequest
	(*LocateCloseReasonRequest)(nil),  // 6: cases.LocateCloseReasonRequest
	(*LocateCloseReasonResponse)(nil), // 7: cases.LocateCloseReasonResponse
	(*Lookup)(nil),                    // 8: cases.Lookup
}
var file_close_reason_proto_depIdxs = []int32{
	8, // 0: cases.CloseReason.created_by:type_name -> cases.Lookup
	8, // 1: cases.CloseReason.updated_by:type_name -> cases.Lookup
	0, // 2: cases.CloseReasonList.items:type_name -> cases.CloseReason
	0, // 3: cases.LocateCloseReasonResponse.close_reason:type_name -> cases.CloseReason
	5, // 4: cases.CloseReasons.ListCloseReasons:input_type -> cases.ListCloseReasonRequest
	2, // 5: cases.CloseReasons.CreateCloseReason:input_type -> cases.CreateCloseReasonRequest
	3, // 6: cases.CloseReasons.UpdateCloseReason:input_type -> cases.UpdateCloseReasonRequest
	4, // 7: cases.CloseReasons.DeleteCloseReason:input_type -> cases.DeleteCloseReasonRequest
	6, // 8: cases.CloseReasons.LocateCloseReason:input_type -> cases.LocateCloseReasonRequest
	1, // 9: cases.CloseReasons.ListCloseReasons:output_type -> cases.CloseReasonList
	0, // 10: cases.CloseReasons.CreateCloseReason:output_type -> cases.CloseReason
	0, // 11: cases.CloseReasons.UpdateCloseReason:output_type -> cases.CloseReason
	0, // 12: cases.CloseReasons.DeleteCloseReason:output_type -> cases.CloseReason
	7, // 13: cases.CloseReasons.LocateCloseReason:output_type -> cases.LocateCloseReasonResponse
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_close_reason_proto_init() }
func file_close_reason_proto_init() {
	if File_close_reason_proto != nil {
		return
	}
	file_lookup_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_close_reason_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CloseReason); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CloseReasonList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CreateCloseReasonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateCloseReasonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteCloseReasonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ListCloseReasonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*LocateCloseReasonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_close_reason_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*LocateCloseReasonResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_close_reason_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_close_reason_proto_goTypes,
		DependencyIndexes: file_close_reason_proto_depIdxs,
		MessageInfos:      file_close_reason_proto_msgTypes,
	}.Build()
	File_close_reason_proto = out.File
	file_close_reason_proto_rawDesc = nil
	file_close_reason_proto_goTypes = nil
	file_close_reason_proto_depIdxs = nil
}
