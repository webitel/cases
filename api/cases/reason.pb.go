// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: cases/reason.proto

package cases

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/visibility"
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

// Reason message represents a reason entity with metadata
type Reason struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the reason
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the reason
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the reason
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Close Reason ID of the reason
	CloseReasonId int64 `protobuf:"varint,6,opt,name=close_reason_id,json=closeReasonId,proto3" json:"close_reason_id,omitempty"`
	// CreatedAt timestamp of the reason
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the reason
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the reason
	CreatedBy *Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the reason
	UpdatedBy *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
}

func (x *Reason) Reset() {
	*x = Reason{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reason) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reason) ProtoMessage() {}

func (x *Reason) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reason.ProtoReflect.Descriptor instead.
func (*Reason) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{0}
}

func (x *Reason) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Reason) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Reason) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Reason) GetCloseReasonId() int64 {
	if x != nil {
		return x.CloseReasonId
	}
	return 0
}

func (x *Reason) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Reason) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Reason) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Reason) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// InputReason message for inputting reason data
type InputReason struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the reason
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the reason
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *InputReason) Reset() {
	*x = InputReason{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputReason) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputReason) ProtoMessage() {}

func (x *InputReason) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputReason.ProtoReflect.Descriptor instead.
func (*InputReason) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{1}
}

func (x *InputReason) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputReason) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// ReasonList message contains a list of Reason items with pagination
type ReasonList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int32     `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next  bool      `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items []*Reason `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ReasonList) Reset() {
	*x = ReasonList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReasonList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReasonList) ProtoMessage() {}

func (x *ReasonList) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReasonList.ProtoReflect.Descriptor instead.
func (*ReasonList) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{2}
}

func (x *ReasonList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ReasonList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ReasonList) GetItems() []*Reason {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreateReasonRequest message for creating a new reason
type CreateReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CloseReasonId int64  `protobuf:"varint,1,opt,name=close_reason_id,json=closeReasonId,proto3" json:"close_reason_id,omitempty"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *CreateReasonRequest) Reset() {
	*x = CreateReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateReasonRequest) ProtoMessage() {}

func (x *CreateReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateReasonRequest.ProtoReflect.Descriptor instead.
func (*CreateReasonRequest) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{3}
}

func (x *CreateReasonRequest) GetCloseReasonId() int64 {
	if x != nil {
		return x.CloseReasonId
	}
	return 0
}

func (x *CreateReasonRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateReasonRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// UpdateReasonRequest message for updating an existing reason
type UpdateReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CloseReasonId int64        `protobuf:"varint,1,opt,name=close_reason_id,json=closeReasonId,proto3" json:"close_reason_id,omitempty"`
	Id            int64        `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Input         *InputReason `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	// JSON PATCH fields mask.
	// List of JPath fields specified in body(input).
	XJsonMask []string `protobuf:"bytes,4,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
}

func (x *UpdateReasonRequest) Reset() {
	*x = UpdateReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateReasonRequest) ProtoMessage() {}

func (x *UpdateReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateReasonRequest.ProtoReflect.Descriptor instead.
func (*UpdateReasonRequest) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateReasonRequest) GetCloseReasonId() int64 {
	if x != nil {
		return x.CloseReasonId
	}
	return 0
}

func (x *UpdateReasonRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateReasonRequest) GetInput() *InputReason {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateReasonRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// DeleteReasonRequest message for deleting an existing reason
type DeleteReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CloseReasonId int64 `protobuf:"varint,2,opt,name=close_reason_id,json=closeReasonId,proto3" json:"close_reason_id,omitempty"`
}

func (x *DeleteReasonRequest) Reset() {
	*x = DeleteReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReasonRequest) ProtoMessage() {}

func (x *DeleteReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReasonRequest.ProtoReflect.Descriptor instead.
func (*DeleteReasonRequest) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteReasonRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DeleteReasonRequest) GetCloseReasonId() int64 {
	if x != nil {
		return x.CloseReasonId
	}
	return 0
}

// ListReasonRequest message for listing or searching reasons
type ListReasonRequest struct {
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
	Id   []int64  `protobuf:"varint,5,rep,packed,name=id,proto3" json:"id,omitempty"`
	// Search query string for filtering by name. Supports:
	// - Wildcards (*) for substring matching
	// - Placeholder (?) for single character substitution
	// - Exact match for full names
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Upd close reason
	CloseReasonId int64 `protobuf:"varint,7,opt,name=close_reason_id,json=closeReasonId,proto3" json:"close_reason_id,omitempty"`
}

func (x *ListReasonRequest) Reset() {
	*x = ListReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListReasonRequest) ProtoMessage() {}

func (x *ListReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListReasonRequest.ProtoReflect.Descriptor instead.
func (*ListReasonRequest) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{6}
}

func (x *ListReasonRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListReasonRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListReasonRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListReasonRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListReasonRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListReasonRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListReasonRequest) GetCloseReasonId() int64 {
	if x != nil {
		return x.CloseReasonId
	}
	return 0
}

// LocateReasonRequest message for locating a specific reason by ID
type LocateReasonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CloseReasonId int64    `protobuf:"varint,2,opt,name=close_reason_id,json=closeReasonId,proto3" json:"close_reason_id,omitempty"`
	Fields        []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *LocateReasonRequest) Reset() {
	*x = LocateReasonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocateReasonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateReasonRequest) ProtoMessage() {}

func (x *LocateReasonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateReasonRequest.ProtoReflect.Descriptor instead.
func (*LocateReasonRequest) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{7}
}

func (x *LocateReasonRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateReasonRequest) GetCloseReasonId() int64 {
	if x != nil {
		return x.CloseReasonId
	}
	return 0
}

func (x *LocateReasonRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocateReasonResponse message contains a single reason entity
type LocateReasonResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reason *Reason `protobuf:"bytes,1,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (x *LocateReasonResponse) Reset() {
	*x = LocateReasonResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_reason_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocateReasonResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateReasonResponse) ProtoMessage() {}

func (x *LocateReasonResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cases_reason_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateReasonResponse.ProtoReflect.Descriptor instead.
func (*LocateReasonResponse) Descriptor() ([]byte, []int) {
	return file_cases_reason_proto_rawDescGZIP(), []int{8}
}

func (x *LocateReasonResponse) GetReason() *Reason {
	if x != nil {
		return x.Reason
	}
	return nil
}

var File_cases_reason_proto protoreflect.FileDescriptor

var file_cases_reason_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x1a, 0x12, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x6c, 0x6f, 0x6f, 0x6b, 0x75,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f,
	0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xa0, 0x02, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6c,
	0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x34, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12,
	0x34, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x17, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x42, 0x79, 0x22, 0x43, 0x0a, 0x0b, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x61, 0x0a, 0x0a, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74,
	0x12, 0x2b, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x8d, 0x01,
	0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x3a, 0x18, 0x92, 0x41, 0x15, 0x0a, 0x13, 0xd2, 0x01, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0xd2, 0x01, 0x09, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x22, 0xc6, 0x01,
	0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12,
	0x39, 0x0a, 0x0b, 0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02, 0x5e, 0x24,
	0xfa, 0xd2, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x52,
	0x09, 0x78, 0x4a, 0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a,
	0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x66, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a,
	0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x49, 0x64, 0x3a, 0x17, 0x92, 0x41, 0x14, 0x0a, 0x12, 0xd2, 0x01, 0x0f, 0x63,
	0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x22, 0xad,
	0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f,
	0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x65,
	0x0a, 0x13, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x45, 0x0a, 0x14, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a,
	0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x32, 0xf5, 0x06, 0x0a,
	0x07, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x12, 0xad, 0x01, 0x0a, 0x0b, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x12, 0x20, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x61, 0x92, 0x41, 0x2e, 0x12, 0x2c, 0x52, 0x65, 0x74, 0x72,
	0x69, 0x65, 0x76, 0x65, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x20, 0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2a, 0x12, 0x28,
	0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x7d,
	0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x12, 0x96, 0x01, 0x0a, 0x0c, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x22, 0x4b, 0x92, 0x41, 0x15, 0x12, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x2d, 0x3a, 0x01, 0x2a, 0x22, 0x28, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f,
	0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x73, 0x12, 0xde, 0x01, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x92, 0x01,
	0x92, 0x41, 0x1b, 0x12, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x65,
	0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x6e, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5a, 0x36, 0x3a, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x32, 0x2d, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x1a, 0x2d, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x12, 0x94, 0x01, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x49,
	0x92, 0x41, 0x11, 0x12, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61, 0x20, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2f, 0x2a, 0x2d, 0x2f, 0x63, 0x6c, 0x6f,
	0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x63, 0x6c, 0x6f, 0x73,
	0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0xa8, 0x01, 0x0a, 0x0c, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x4f, 0x92, 0x41, 0x17, 0x12, 0x15, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65,
	0x20, 0x61, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x2f, 0x12, 0x2d, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f, 0x7b, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x7d, 0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3b, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cases_reason_proto_rawDescOnce sync.Once
	file_cases_reason_proto_rawDescData = file_cases_reason_proto_rawDesc
)

func file_cases_reason_proto_rawDescGZIP() []byte {
	file_cases_reason_proto_rawDescOnce.Do(func() {
		file_cases_reason_proto_rawDescData = protoimpl.X.CompressGZIP(file_cases_reason_proto_rawDescData)
	})
	return file_cases_reason_proto_rawDescData
}

var file_cases_reason_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_cases_reason_proto_goTypes = []any{
	(*Reason)(nil),               // 0: webitel.cases.Reason
	(*InputReason)(nil),          // 1: webitel.cases.InputReason
	(*ReasonList)(nil),           // 2: webitel.cases.ReasonList
	(*CreateReasonRequest)(nil),  // 3: webitel.cases.CreateReasonRequest
	(*UpdateReasonRequest)(nil),  // 4: webitel.cases.UpdateReasonRequest
	(*DeleteReasonRequest)(nil),  // 5: webitel.cases.DeleteReasonRequest
	(*ListReasonRequest)(nil),    // 6: webitel.cases.ListReasonRequest
	(*LocateReasonRequest)(nil),  // 7: webitel.cases.LocateReasonRequest
	(*LocateReasonResponse)(nil), // 8: webitel.cases.LocateReasonResponse
	(*Lookup)(nil),               // 9: webitel.cases.Lookup
}
var file_cases_reason_proto_depIdxs = []int32{
	9,  // 0: webitel.cases.Reason.created_by:type_name -> webitel.cases.Lookup
	9,  // 1: webitel.cases.Reason.updated_by:type_name -> webitel.cases.Lookup
	0,  // 2: webitel.cases.ReasonList.items:type_name -> webitel.cases.Reason
	1,  // 3: webitel.cases.UpdateReasonRequest.input:type_name -> webitel.cases.InputReason
	0,  // 4: webitel.cases.LocateReasonResponse.reason:type_name -> webitel.cases.Reason
	6,  // 5: webitel.cases.Reasons.ListReasons:input_type -> webitel.cases.ListReasonRequest
	3,  // 6: webitel.cases.Reasons.CreateReason:input_type -> webitel.cases.CreateReasonRequest
	4,  // 7: webitel.cases.Reasons.UpdateReason:input_type -> webitel.cases.UpdateReasonRequest
	5,  // 8: webitel.cases.Reasons.DeleteReason:input_type -> webitel.cases.DeleteReasonRequest
	7,  // 9: webitel.cases.Reasons.LocateReason:input_type -> webitel.cases.LocateReasonRequest
	2,  // 10: webitel.cases.Reasons.ListReasons:output_type -> webitel.cases.ReasonList
	0,  // 11: webitel.cases.Reasons.CreateReason:output_type -> webitel.cases.Reason
	0,  // 12: webitel.cases.Reasons.UpdateReason:output_type -> webitel.cases.Reason
	0,  // 13: webitel.cases.Reasons.DeleteReason:output_type -> webitel.cases.Reason
	8,  // 14: webitel.cases.Reasons.LocateReason:output_type -> webitel.cases.LocateReasonResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_cases_reason_proto_init() }
func file_cases_reason_proto_init() {
	if File_cases_reason_proto != nil {
		return
	}
	file_cases_lookup_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_cases_reason_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Reason); i {
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
		file_cases_reason_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*InputReason); i {
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
		file_cases_reason_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ReasonList); i {
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
		file_cases_reason_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*CreateReasonRequest); i {
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
		file_cases_reason_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateReasonRequest); i {
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
		file_cases_reason_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteReasonRequest); i {
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
		file_cases_reason_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ListReasonRequest); i {
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
		file_cases_reason_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*LocateReasonRequest); i {
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
		file_cases_reason_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*LocateReasonResponse); i {
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
			RawDescriptor: file_cases_reason_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cases_reason_proto_goTypes,
		DependencyIndexes: file_cases_reason_proto_depIdxs,
		MessageInfos:      file_cases_reason_proto_msgTypes,
	}.Build()
	File_cases_reason_proto = out.File
	file_cases_reason_proto_rawDesc = nil
	file_cases_reason_proto_goTypes = nil
	file_cases_reason_proto_depIdxs = nil
}
