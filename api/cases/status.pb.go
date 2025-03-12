// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: status.proto

package cases

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "github.com/webitel/webitel-go-kit/cmd/protoc-gen-go-webitel/gen/go/proto/webitel"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	_ "google.golang.org/genproto/googleapis/api/visibility"
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

// Status message represents a status entity with metadata
type Status struct {
	state protoimpl.MessageState `protogen:"open.v1"`
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
	UpdatedBy     *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Status) Reset() {
	*x = Status{}
	mi := &file_status_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{0}
}

func (x *Status) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Status) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Status) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Status) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Status) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Status) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Status) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// StatusList message contains a list of Status items with pagination
type StatusList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next          bool                   `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*Status              `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatusList) Reset() {
	*x = StatusList{}
	mi := &file_status_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatusList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusList) ProtoMessage() {}

func (x *StatusList) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusList.ProtoReflect.Descriptor instead.
func (*StatusList) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{1}
}

func (x *StatusList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *StatusList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *StatusList) GetItems() []*Status {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreateStatusRequest message for creating a new status
type CreateStatusRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Input *InputStatus           `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	// Fields to be retrieved as a result.
	Fields        []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateStatusRequest) Reset() {
	*x = CreateStatusRequest{}
	mi := &file_status_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateStatusRequest) ProtoMessage() {}

func (x *CreateStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateStatusRequest.ProtoReflect.Descriptor instead.
func (*CreateStatusRequest) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{2}
}

func (x *CreateStatusRequest) GetInput() *InputStatus {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *CreateStatusRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type InputStatus struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputStatus) Reset() {
	*x = InputStatus{}
	mi := &file_status_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputStatus) ProtoMessage() {}

func (x *InputStatus) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputStatus.ProtoReflect.Descriptor instead.
func (*InputStatus) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{3}
}

func (x *InputStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputStatus) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// UpdateStatusRequest message for updating an existing status
type UpdateStatusRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Fields to be retrieved as a result.
	Fields []string     `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	Input  *InputStatus `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,4,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateStatusRequest) Reset() {
	*x = UpdateStatusRequest{}
	mi := &file_status_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateStatusRequest) ProtoMessage() {}

func (x *UpdateStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateStatusRequest) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateStatusRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateStatusRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *UpdateStatusRequest) GetInput() *InputStatus {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateStatusRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// DeleteStatusRequest message for deleting an existing status
type DeleteStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteStatusRequest) Reset() {
	*x = DeleteStatusRequest{}
	mi := &file_status_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteStatusRequest) ProtoMessage() {}

func (x *DeleteStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteStatusRequest.ProtoReflect.Descriptor instead.
func (*DeleteStatusRequest) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteStatusRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// ListStatusesRequest message for listing or searching statuses
type ListStatusRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of result dataset records. offset = (page*size)
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Size count of records on result page. limit = (size++)
	Size int32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// Fields to be retrieved as a result.
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Sort the result according to fields.
	Sort string `protobuf:"bytes,4,opt,name=sort,proto3" json:"sort,omitempty"`
	// Filter by unique IDs.
	Id []int64 `protobuf:"varint,5,rep,packed,name=id,proto3" json:"id,omitempty"`
	// Search query string for filtering by name. Supports:
	// - Wildcards (*) for substring matching
	// - Placeholder (?) for single character substitution
	// - Exact match for full names
	Q             string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListStatusRequest) Reset() {
	*x = ListStatusRequest{}
	mi := &file_status_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStatusRequest) ProtoMessage() {}

func (x *ListStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStatusRequest.ProtoReflect.Descriptor instead.
func (*ListStatusRequest) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{6}
}

func (x *ListStatusRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListStatusRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListStatusRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListStatusRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *ListStatusRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListStatusRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

// LocateStatusRequest message for locating a specific status by ID
type LocateStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Fields        []string               `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateStatusRequest) Reset() {
	*x = LocateStatusRequest{}
	mi := &file_status_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateStatusRequest) ProtoMessage() {}

func (x *LocateStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateStatusRequest.ProtoReflect.Descriptor instead.
func (*LocateStatusRequest) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{7}
}

func (x *LocateStatusRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateStatusRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocateStatusResponse message contains a single status entity
type LocateStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        *Status                `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateStatusResponse) Reset() {
	*x = LocateStatusResponse{}
	mi := &file_status_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateStatusResponse) ProtoMessage() {}

func (x *LocateStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_status_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateStatusResponse.ProtoReflect.Descriptor instead.
func (*LocateStatusResponse) Descriptor() ([]byte, []int) {
	return file_status_proto_rawDescGZIP(), []int{8}
}

func (x *LocateStatusResponse) GetStatus() *Status {
	if x != nil {
		return x.Status
	}
	return nil
}

var File_status_proto protoreflect.FileDescriptor

var file_status_proto_rawDesc = string([]byte{
	0x0a, 0x0c, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x0d, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d,
	0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xec, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c,
	0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x42, 0x79, 0x12, 0x2e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79,
	0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c,
	0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x42, 0x79, 0x22, 0x61, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x2b, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x6d, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x3a, 0x0c, 0x92, 0x41, 0x09, 0x0a, 0x07, 0xd2, 0x01, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x43, 0x0a, 0x0b, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xb6, 0x01, 0x0a, 0x13, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x30, 0x0a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x39, 0x0a, 0x0b, 0x78,
	0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02, 0x5e, 0x24, 0xfa, 0xd2, 0xe4, 0x93,
	0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x52, 0x09, 0x78, 0x4a, 0x73,
	0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02,
	0x69, 0x64, 0x22, 0x31, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05,
	0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x85, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x0c, 0x0a, 0x01, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x22, 0x3d, 0x0a,
	0x13, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x45, 0x0a, 0x14,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x32, 0x8b, 0x06, 0x0a, 0x08, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73,
	0x12, 0x9b, 0x01, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65,
	0x73, 0x12, 0x20, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x4e,
	0x92, 0x41, 0x30, 0x12, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20, 0x61, 0x20,
	0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73,
	0x20, 0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x65, 0x73, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0x85,
	0x01, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3a, 0x92, 0x41, 0x15, 0x12,
	0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x90, 0xb5, 0x18, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x05,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x22, 0x0f, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x12, 0xaf, 0x01, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x64, 0x92, 0x41, 0x1b, 0x12, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20,
	0x61, 0x6e, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x3c, 0x3a, 0x05, 0x69, 0x6e,
	0x70, 0x75, 0x74, 0x5a, 0x1d, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x32, 0x14, 0x2f, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x1a, 0x14, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x7f, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x34, 0x92, 0x41, 0x11, 0x12, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x20, 0x61, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x90, 0xb5, 0x18, 0x03, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x16, 0x2a, 0x14, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x93, 0x01, 0x0a, 0x0c, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c,
	0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x3a, 0x92, 0x41, 0x17, 0x12, 0x15, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65,
	0x20, 0x61, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x90,
	0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x16, 0x12, 0x14, 0x2f, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a,
	0x10, 0x8a, 0xb5, 0x18, 0x0c, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x69, 0x65,
	0x73, 0x42, 0x9f, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0x0b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3b, 0x63, 0x61, 0x73, 0x65, 0x73,
	0xa2, 0x02, 0x03, 0x57, 0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2e, 0x43, 0x61, 0x73, 0x65, 0x73, 0xca, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0xe2, 0x02, 0x19, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x0e, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a, 0x3a, 0x43, 0x61,
	0x73, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_status_proto_rawDescOnce sync.Once
	file_status_proto_rawDescData []byte
)

func file_status_proto_rawDescGZIP() []byte {
	file_status_proto_rawDescOnce.Do(func() {
		file_status_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_status_proto_rawDesc), len(file_status_proto_rawDesc)))
	})
	return file_status_proto_rawDescData
}

var file_status_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_status_proto_goTypes = []any{
	(*Status)(nil),               // 0: webitel.cases.Status
	(*StatusList)(nil),           // 1: webitel.cases.StatusList
	(*CreateStatusRequest)(nil),  // 2: webitel.cases.CreateStatusRequest
	(*InputStatus)(nil),          // 3: webitel.cases.InputStatus
	(*UpdateStatusRequest)(nil),  // 4: webitel.cases.UpdateStatusRequest
	(*DeleteStatusRequest)(nil),  // 5: webitel.cases.DeleteStatusRequest
	(*ListStatusRequest)(nil),    // 6: webitel.cases.ListStatusRequest
	(*LocateStatusRequest)(nil),  // 7: webitel.cases.LocateStatusRequest
	(*LocateStatusResponse)(nil), // 8: webitel.cases.LocateStatusResponse
	(*Lookup)(nil),               // 9: general.Lookup
}
var file_status_proto_depIdxs = []int32{
	9,  // 0: webitel.cases.Status.created_by:type_name -> general.Lookup
	9,  // 1: webitel.cases.Status.updated_by:type_name -> general.Lookup
	0,  // 2: webitel.cases.StatusList.items:type_name -> webitel.cases.Status
	3,  // 3: webitel.cases.CreateStatusRequest.input:type_name -> webitel.cases.InputStatus
	3,  // 4: webitel.cases.UpdateStatusRequest.input:type_name -> webitel.cases.InputStatus
	0,  // 5: webitel.cases.LocateStatusResponse.status:type_name -> webitel.cases.Status
	6,  // 6: webitel.cases.Statuses.ListStatuses:input_type -> webitel.cases.ListStatusRequest
	2,  // 7: webitel.cases.Statuses.CreateStatus:input_type -> webitel.cases.CreateStatusRequest
	4,  // 8: webitel.cases.Statuses.UpdateStatus:input_type -> webitel.cases.UpdateStatusRequest
	5,  // 9: webitel.cases.Statuses.DeleteStatus:input_type -> webitel.cases.DeleteStatusRequest
	7,  // 10: webitel.cases.Statuses.LocateStatus:input_type -> webitel.cases.LocateStatusRequest
	1,  // 11: webitel.cases.Statuses.ListStatuses:output_type -> webitel.cases.StatusList
	0,  // 12: webitel.cases.Statuses.CreateStatus:output_type -> webitel.cases.Status
	0,  // 13: webitel.cases.Statuses.UpdateStatus:output_type -> webitel.cases.Status
	0,  // 14: webitel.cases.Statuses.DeleteStatus:output_type -> webitel.cases.Status
	8,  // 15: webitel.cases.Statuses.LocateStatus:output_type -> webitel.cases.LocateStatusResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_status_proto_init() }
func file_status_proto_init() {
	if File_status_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_status_proto_rawDesc), len(file_status_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_status_proto_goTypes,
		DependencyIndexes: file_status_proto_depIdxs,
		MessageInfos:      file_status_proto_msgTypes,
	}.Build()
	File_status_proto = out.File
	file_status_proto_goTypes = nil
	file_status_proto_depIdxs = nil
}
