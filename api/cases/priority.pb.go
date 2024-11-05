// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: priority.proto

package cases

import (
	_go "buf.build/gen/go/webitel/general/protocolbuffers/go"
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

// Priority message represents a priority entity with metadata
type Priority struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique identifier of the priority
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the priority
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the priority
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// CreatedAt timestamp of the priority
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the priority
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the priority
	CreatedBy *_go.Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the priority
	UpdatedBy *_go.Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// Color of the priority
	Color string `protobuf:"bytes,24,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *Priority) Reset() {
	*x = Priority{}
	mi := &file_priority_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Priority) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Priority) ProtoMessage() {}

func (x *Priority) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Priority.ProtoReflect.Descriptor instead.
func (*Priority) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{0}
}

func (x *Priority) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Priority) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Priority) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Priority) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Priority) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Priority) GetCreatedBy() *_go.Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Priority) GetUpdatedBy() *_go.Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *Priority) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

// PriorityList message contains a list of Priority items with pagination
type PriorityList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int32       `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next  bool        `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items []*Priority `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *PriorityList) Reset() {
	*x = PriorityList{}
	mi := &file_priority_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PriorityList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PriorityList) ProtoMessage() {}

func (x *PriorityList) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PriorityList.ProtoReflect.Descriptor instead.
func (*PriorityList) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{1}
}

func (x *PriorityList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *PriorityList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *PriorityList) GetItems() []*Priority {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreatePriorityRequest message for creating a new priority
type CreatePriorityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Color       string `protobuf:"bytes,4,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *CreatePriorityRequest) Reset() {
	*x = CreatePriorityRequest{}
	mi := &file_priority_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePriorityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePriorityRequest) ProtoMessage() {}

func (x *CreatePriorityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePriorityRequest.ProtoReflect.Descriptor instead.
func (*CreatePriorityRequest) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{2}
}

func (x *CreatePriorityRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreatePriorityRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreatePriorityRequest) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

// InputPriority message for creating a new priority
type InputPriority struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the priority
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the priority
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// Color of the priority
	Color string `protobuf:"bytes,3,opt,name=color,proto3" json:"color,omitempty"`
}

func (x *InputPriority) Reset() {
	*x = InputPriority{}
	mi := &file_priority_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputPriority) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputPriority) ProtoMessage() {}

func (x *InputPriority) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputPriority.ProtoReflect.Descriptor instead.
func (*InputPriority) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{3}
}

func (x *InputPriority) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputPriority) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *InputPriority) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

// UpdatePriorityRequest message for updating an existing priority
type UpdatePriorityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64          `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Input *InputPriority `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask []string `protobuf:"bytes,3,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
}

func (x *UpdatePriorityRequest) Reset() {
	*x = UpdatePriorityRequest{}
	mi := &file_priority_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePriorityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePriorityRequest) ProtoMessage() {}

func (x *UpdatePriorityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePriorityRequest.ProtoReflect.Descriptor instead.
func (*UpdatePriorityRequest) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{4}
}

func (x *UpdatePriorityRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdatePriorityRequest) GetInput() *InputPriority {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdatePriorityRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// DeletePriorityRequest message for deleting an existing priority
type DeletePriorityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeletePriorityRequest) Reset() {
	*x = DeletePriorityRequest{}
	mi := &file_priority_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeletePriorityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePriorityRequest) ProtoMessage() {}

func (x *DeletePriorityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePriorityRequest.ProtoReflect.Descriptor instead.
func (*DeletePriorityRequest) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{5}
}

func (x *DeletePriorityRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// ListPriorityRequest message for listing or searching priority
type ListPriorityRequest struct {
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
	// Search query string for filtering by name. Supports:
	// - Wildcards (*) for substring matching
	// - Placeholder (?) for single character substitution
	// - Exact match for full names
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
}

func (x *ListPriorityRequest) Reset() {
	*x = ListPriorityRequest{}
	mi := &file_priority_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPriorityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPriorityRequest) ProtoMessage() {}

func (x *ListPriorityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPriorityRequest.ProtoReflect.Descriptor instead.
func (*ListPriorityRequest) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{6}
}

func (x *ListPriorityRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListPriorityRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListPriorityRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListPriorityRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListPriorityRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListPriorityRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

// LocatePriorityRequest message for locating a specific priority by ID
type LocatePriorityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the priority to be located
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *LocatePriorityRequest) Reset() {
	*x = LocatePriorityRequest{}
	mi := &file_priority_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocatePriorityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocatePriorityRequest) ProtoMessage() {}

func (x *LocatePriorityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocatePriorityRequest.ProtoReflect.Descriptor instead.
func (*LocatePriorityRequest) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{7}
}

func (x *LocatePriorityRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocatePriorityRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocatePriorityResponse message contains a single priority entity
type LocatePriorityResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Priority *Priority `protobuf:"bytes,1,opt,name=priority,proto3" json:"priority,omitempty"`
}

func (x *LocatePriorityResponse) Reset() {
	*x = LocatePriorityResponse{}
	mi := &file_priority_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocatePriorityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocatePriorityResponse) ProtoMessage() {}

func (x *LocatePriorityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_priority_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocatePriorityResponse.ProtoReflect.Descriptor instead.
func (*LocatePriorityResponse) Descriptor() ([]byte, []int) {
	return file_priority_proto_rawDescGZIP(), []int{8}
}

func (x *LocatePriorityResponse) GetPriority() *Priority {
	if x != nil {
		return x.Priority
	}
	return nil
}

var File_priority_proto protoreflect.FileDescriptor

var file_priority_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a,
	0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x02, 0x0a, 0x08, 0x50, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x2e, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x6c, 0x6f, 0x72, 0x18, 0x18, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72,
	0x22, 0x65, 0x0a, 0x0c, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x2d, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x71, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x3a, 0x0c, 0x92, 0x41,
	0x09, 0x0a, 0x07, 0xd2, 0x01, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5b, 0x0a, 0x0d, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x22, 0xa2, 0x01, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x32, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x39, 0x0a, 0x0b, 0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f,
	0x6d, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40,
	0x01, 0x8a, 0x01, 0x02, 0x5e, 0x24, 0xfa, 0xd2, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52,
	0x45, 0x56, 0x49, 0x45, 0x57, 0x52, 0x09, 0x78, 0x4a, 0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b,
	0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x33, 0x0a, 0x15,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69,
	0x64, 0x22, 0x87, 0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72,
	0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a,
	0x01, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x22, 0x3f, 0x0a, 0x15, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x4d, 0x0a, 0x16,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x32, 0x9a, 0x06, 0x0a, 0x0a,
	0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0xa3, 0x01, 0x0a, 0x0e, 0x4c,
	0x69, 0x73, 0x74, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x22, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x50,
	0x92, 0x41, 0x34, 0x12, 0x32, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20, 0x61, 0x20,
	0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x20, 0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x12, 0x11, 0x2f,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73,
	0x12, 0x87, 0x01, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72,
	0x69, 0x74, 0x79, 0x12, 0x24, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x22, 0x36, 0x92, 0x41, 0x17, 0x12, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20,
	0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x16, 0x3a, 0x01, 0x2a, 0x22, 0x11, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f,
	0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0xb7, 0x01, 0x0a, 0x0e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x24, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x22, 0x66, 0x92, 0x41,
	0x1d, 0x12, 0x1b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x78, 0x69,
	0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x40, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5a, 0x1f, 0x3a, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x32, 0x16, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x16, 0x2f, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x12, 0x85, 0x01, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50,
	0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x24, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x50, 0x72,
	0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x22, 0x34, 0x92, 0x41, 0x13, 0x12, 0x11, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x20, 0x61, 0x20, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x18, 0x2a, 0x16, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x99, 0x01, 0x0a,
	0x0e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x24, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x50, 0x72, 0x69, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3a, 0x92, 0x41,
	0x19, 0x12, 0x17, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x70, 0x72, 0x69, 0x6f,
	0x72, 0x69, 0x74, 0x79, 0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18,
	0x12, 0x16, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x9b, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0x0d,
	0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61,
	0x73, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x57, 0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x73, 0xca, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0xe2, 0x02, 0x19, 0x57, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a,
	0x3a, 0x43, 0x61, 0x73, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_priority_proto_rawDescOnce sync.Once
	file_priority_proto_rawDescData = file_priority_proto_rawDesc
)

func file_priority_proto_rawDescGZIP() []byte {
	file_priority_proto_rawDescOnce.Do(func() {
		file_priority_proto_rawDescData = protoimpl.X.CompressGZIP(file_priority_proto_rawDescData)
	})
	return file_priority_proto_rawDescData
}

var file_priority_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_priority_proto_goTypes = []any{
	(*Priority)(nil),               // 0: webitel.cases.Priority
	(*PriorityList)(nil),           // 1: webitel.cases.PriorityList
	(*CreatePriorityRequest)(nil),  // 2: webitel.cases.CreatePriorityRequest
	(*InputPriority)(nil),          // 3: webitel.cases.InputPriority
	(*UpdatePriorityRequest)(nil),  // 4: webitel.cases.UpdatePriorityRequest
	(*DeletePriorityRequest)(nil),  // 5: webitel.cases.DeletePriorityRequest
	(*ListPriorityRequest)(nil),    // 6: webitel.cases.ListPriorityRequest
	(*LocatePriorityRequest)(nil),  // 7: webitel.cases.LocatePriorityRequest
	(*LocatePriorityResponse)(nil), // 8: webitel.cases.LocatePriorityResponse
	(*_go.Lookup)(nil),             // 9: general.Lookup
}
var file_priority_proto_depIdxs = []int32{
	9,  // 0: webitel.cases.Priority.created_by:type_name -> general.Lookup
	9,  // 1: webitel.cases.Priority.updated_by:type_name -> general.Lookup
	0,  // 2: webitel.cases.PriorityList.items:type_name -> webitel.cases.Priority
	3,  // 3: webitel.cases.UpdatePriorityRequest.input:type_name -> webitel.cases.InputPriority
	0,  // 4: webitel.cases.LocatePriorityResponse.priority:type_name -> webitel.cases.Priority
	6,  // 5: webitel.cases.Priorities.ListPriorities:input_type -> webitel.cases.ListPriorityRequest
	2,  // 6: webitel.cases.Priorities.CreatePriority:input_type -> webitel.cases.CreatePriorityRequest
	4,  // 7: webitel.cases.Priorities.UpdatePriority:input_type -> webitel.cases.UpdatePriorityRequest
	5,  // 8: webitel.cases.Priorities.DeletePriority:input_type -> webitel.cases.DeletePriorityRequest
	7,  // 9: webitel.cases.Priorities.LocatePriority:input_type -> webitel.cases.LocatePriorityRequest
	1,  // 10: webitel.cases.Priorities.ListPriorities:output_type -> webitel.cases.PriorityList
	0,  // 11: webitel.cases.Priorities.CreatePriority:output_type -> webitel.cases.Priority
	0,  // 12: webitel.cases.Priorities.UpdatePriority:output_type -> webitel.cases.Priority
	0,  // 13: webitel.cases.Priorities.DeletePriority:output_type -> webitel.cases.Priority
	8,  // 14: webitel.cases.Priorities.LocatePriority:output_type -> webitel.cases.LocatePriorityResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_priority_proto_init() }
func file_priority_proto_init() {
	if File_priority_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_priority_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_priority_proto_goTypes,
		DependencyIndexes: file_priority_proto_depIdxs,
		MessageInfos:      file_priority_proto_msgTypes,
	}.Build()
	File_priority_proto = out.File
	file_priority_proto_rawDesc = nil
	file_priority_proto_goTypes = nil
	file_priority_proto_depIdxs = nil
}
