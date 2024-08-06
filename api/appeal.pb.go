// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: appeal.proto

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

// Represents a appeal in the contact system.
type Appeal struct {
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

func (x *Appeal) Reset() {
	*x = Appeal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Appeal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Appeal) ProtoMessage() {}

func (x *Appeal) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Appeal.ProtoReflect.Descriptor instead.
func (*Appeal) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{0}
}

func (x *Appeal) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Appeal) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Appeal) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Appeal) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Appeal) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Appeal) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Appeal) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// A list of appeals.
type AppealList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Page number of the partial result.
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Have more records.
	Next bool `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	// List of appeals.
	Items []*Appeal `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *AppealList) Reset() {
	*x = AppealList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppealList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppealList) ProtoMessage() {}

func (x *AppealList) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppealList.ProtoReflect.Descriptor instead.
func (*AppealList) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{1}
}

func (x *AppealList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *AppealList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *AppealList) GetItems() []*Appeal {
	if x != nil {
		return x.Items
	}
	return nil
}

// Request message for creating a new appeal.
type CreateAppealRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the appeal.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the appeal.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *CreateAppealRequest) Reset() {
	*x = CreateAppealRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAppealRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAppealRequest) ProtoMessage() {}

func (x *CreateAppealRequest) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAppealRequest.ProtoReflect.Descriptor instead.
func (*CreateAppealRequest) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{2}
}

func (x *CreateAppealRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAppealRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// Request message for updating an existing appeal.
type UpdateAppealRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID of the appeal to update.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// The new name of the appeal.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The new description of the appeal.
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *UpdateAppealRequest) Reset() {
	*x = UpdateAppealRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateAppealRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAppealRequest) ProtoMessage() {}

func (x *UpdateAppealRequest) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAppealRequest.ProtoReflect.Descriptor instead.
func (*UpdateAppealRequest) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateAppealRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateAppealRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateAppealRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// Request message for deleting a appeal.
type DeleteAppealRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID of the appeal to delete.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteAppealRequest) Reset() {
	*x = DeleteAppealRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAppealRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAppealRequest) ProtoMessage() {}

func (x *DeleteAppealRequest) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAppealRequest.ProtoReflect.Descriptor instead.
func (*DeleteAppealRequest) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteAppealRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// Request message for listing appeals.
type ListAppealRequest struct {
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
	// Search term: appeal name;
	// `?` - matches any one character
	// `*` - matches 0 or more characters
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Filter by appeal name.
	Name string `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ListAppealRequest) Reset() {
	*x = ListAppealRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAppealRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAppealRequest) ProtoMessage() {}

func (x *ListAppealRequest) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAppealRequest.ProtoReflect.Descriptor instead.
func (*ListAppealRequest) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{5}
}

func (x *ListAppealRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListAppealRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListAppealRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListAppealRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListAppealRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListAppealRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListAppealRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Request message for locating a appeal by ID.
type LocateAppealRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID of the appeal to locate.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Fields to be retrieved into result.
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *LocateAppealRequest) Reset() {
	*x = LocateAppealRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocateAppealRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateAppealRequest) ProtoMessage() {}

func (x *LocateAppealRequest) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateAppealRequest.ProtoReflect.Descriptor instead.
func (*LocateAppealRequest) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{6}
}

func (x *LocateAppealRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateAppealRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Response message for locating a appeal.
type LocateAppealResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The located appeal.
	Appeal *Appeal `protobuf:"bytes,1,opt,name=appeal,proto3" json:"appeal,omitempty"`
}

func (x *LocateAppealResponse) Reset() {
	*x = LocateAppealResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_appeal_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocateAppealResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateAppealResponse) ProtoMessage() {}

func (x *LocateAppealResponse) ProtoReflect() protoreflect.Message {
	mi := &file_appeal_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateAppealResponse.ProtoReflect.Descriptor instead.
func (*LocateAppealResponse) Descriptor() ([]byte, []int) {
	return file_appeal_proto_rawDescGZIP(), []int{7}
}

func (x *LocateAppealResponse) GetAppeal() *Appeal {
	if x != nil {
		return x.Appeal
	}
	return nil
}

var File_appeal_proto protoreflect.FileDescriptor

var file_appeal_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x0c, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xe8, 0x01, 0x0a, 0x06, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x2c, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x16,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f,
	0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x2c,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x17, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75,
	0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x22, 0x59, 0x0a, 0x0a,
	0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65,
	0x78, 0x74, 0x12, 0x23, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x59, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x0c, 0x92, 0x41, 0x09, 0x0a, 0x07, 0xd2, 0x01, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0x67, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65,
	0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a,
	0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x31, 0x0a, 0x13, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x99,
	0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3d, 0x0a, 0x13, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x3d, 0x0a, 0x14, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x25, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c,
	0x52, 0x06, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x32, 0xd9, 0x04, 0x0a, 0x07, 0x41, 0x70, 0x70,
	0x65, 0x61, 0x6c, 0x73, 0x12, 0x7d, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x65,
	0x61, 0x6c, 0x73, 0x12, 0x18, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x4c, 0x69, 0x73, 0x74,
	0x22, 0x41, 0x92, 0x41, 0x2e, 0x12, 0x2c, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20,
	0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c,
	0x73, 0x20, 0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x61, 0x70, 0x70, 0x65,
	0x61, 0x6c, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x61, 0x70, 0x70, 0x65,
	0x61, 0x6c, 0x73, 0x12, 0x66, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70,
	0x65, 0x61, 0x6c, 0x12, 0x1a, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x22, 0x2b,
	0x92, 0x41, 0x15, 0x12, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x65,
	0x77, 0x20, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x3a, 0x01,
	0x2a, 0x22, 0x08, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x73, 0x12, 0x85, 0x01, 0x0a, 0x0c,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x12, 0x1a, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x22, 0x4a, 0x92, 0x41, 0x1b, 0x12, 0x19, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x20, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x26, 0x3a, 0x01, 0x2a,
	0x5a, 0x12, 0x3a, 0x01, 0x2a, 0x32, 0x0d, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x0d, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x12, 0x65, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x70, 0x70,
	0x65, 0x61, 0x6c, 0x12, 0x1a, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0d, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x22, 0x2a,
	0x92, 0x41, 0x12, 0x12, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x61,
	0x70, 0x70, 0x65, 0x61, 0x6c, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x2a, 0x0d, 0x2f, 0x61, 0x70,
	0x70, 0x65, 0x61, 0x6c, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x78, 0x0a, 0x0c, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x12, 0x1a, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x65, 0x41, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x2f, 0x92, 0x41, 0x17, 0x12, 0x15, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65,
	0x20, 0x61, 0x20, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x61, 0x6c, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x42, 0x0b, 0x5a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_appeal_proto_rawDescOnce sync.Once
	file_appeal_proto_rawDescData = file_appeal_proto_rawDesc
)

func file_appeal_proto_rawDescGZIP() []byte {
	file_appeal_proto_rawDescOnce.Do(func() {
		file_appeal_proto_rawDescData = protoimpl.X.CompressGZIP(file_appeal_proto_rawDescData)
	})
	return file_appeal_proto_rawDescData
}

var file_appeal_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_appeal_proto_goTypes = []any{
	(*Appeal)(nil),               // 0: cases.Appeal
	(*AppealList)(nil),           // 1: cases.AppealList
	(*CreateAppealRequest)(nil),  // 2: cases.CreateAppealRequest
	(*UpdateAppealRequest)(nil),  // 3: cases.UpdateAppealRequest
	(*DeleteAppealRequest)(nil),  // 4: cases.DeleteAppealRequest
	(*ListAppealRequest)(nil),    // 5: cases.ListAppealRequest
	(*LocateAppealRequest)(nil),  // 6: cases.LocateAppealRequest
	(*LocateAppealResponse)(nil), // 7: cases.LocateAppealResponse
	(*Lookup)(nil),               // 8: cases.Lookup
}
var file_appeal_proto_depIdxs = []int32{
	8, // 0: cases.Appeal.created_by:type_name -> cases.Lookup
	8, // 1: cases.Appeal.updated_by:type_name -> cases.Lookup
	0, // 2: cases.AppealList.items:type_name -> cases.Appeal
	0, // 3: cases.LocateAppealResponse.appeal:type_name -> cases.Appeal
	5, // 4: cases.Appeals.ListAppeals:input_type -> cases.ListAppealRequest
	2, // 5: cases.Appeals.CreateAppeal:input_type -> cases.CreateAppealRequest
	3, // 6: cases.Appeals.UpdateAppeal:input_type -> cases.UpdateAppealRequest
	4, // 7: cases.Appeals.DeleteAppeal:input_type -> cases.DeleteAppealRequest
	6, // 8: cases.Appeals.LocateAppeal:input_type -> cases.LocateAppealRequest
	1, // 9: cases.Appeals.ListAppeals:output_type -> cases.AppealList
	0, // 10: cases.Appeals.CreateAppeal:output_type -> cases.Appeal
	0, // 11: cases.Appeals.UpdateAppeal:output_type -> cases.Appeal
	0, // 12: cases.Appeals.DeleteAppeal:output_type -> cases.Appeal
	7, // 13: cases.Appeals.LocateAppeal:output_type -> cases.LocateAppealResponse
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_appeal_proto_init() }
func file_appeal_proto_init() {
	if File_appeal_proto != nil {
		return
	}
	file_lookup_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_appeal_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Appeal); i {
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
		file_appeal_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AppealList); i {
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
		file_appeal_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CreateAppealRequest); i {
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
		file_appeal_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateAppealRequest); i {
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
		file_appeal_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteAppealRequest); i {
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
		file_appeal_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ListAppealRequest); i {
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
		file_appeal_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*LocateAppealRequest); i {
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
		file_appeal_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*LocateAppealResponse); i {
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
			RawDescriptor: file_appeal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_appeal_proto_goTypes,
		DependencyIndexes: file_appeal_proto_depIdxs,
		MessageInfos:      file_appeal_proto_msgTypes,
	}.Build()
	File_appeal_proto = out.File
	file_appeal_proto_rawDesc = nil
	file_appeal_proto_goTypes = nil
	file_appeal_proto_depIdxs = nil
}
