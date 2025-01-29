// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: close_reason_group.proto

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

// CloseReasonGroup message represents a close reason group entity with metadata
type CloseReasonGroup struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the close reason group
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the close reason group
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the close reason group
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// CreatedAt timestamp of the close reason group
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the close reason group
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the close reason group
	CreatedBy *Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the close reason group
	UpdatedBy     *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CloseReasonGroup) Reset() {
	*x = CloseReasonGroup{}
	mi := &file_close_reason_group_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CloseReasonGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseReasonGroup) ProtoMessage() {}

func (x *CloseReasonGroup) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseReasonGroup.ProtoReflect.Descriptor instead.
func (*CloseReasonGroup) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{0}
}

func (x *CloseReasonGroup) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CloseReasonGroup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CloseReasonGroup) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CloseReasonGroup) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *CloseReasonGroup) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *CloseReasonGroup) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *CloseReasonGroup) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// CloseReasonGroupList message contains a list of CloseReasonGroup items with pagination
type CloseReasonGroupList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next          bool                   `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*CloseReasonGroup    `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CloseReasonGroupList) Reset() {
	*x = CloseReasonGroupList{}
	mi := &file_close_reason_group_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CloseReasonGroupList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloseReasonGroupList) ProtoMessage() {}

func (x *CloseReasonGroupList) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloseReasonGroupList.ProtoReflect.Descriptor instead.
func (*CloseReasonGroupList) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{1}
}

func (x *CloseReasonGroupList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *CloseReasonGroupList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *CloseReasonGroupList) GetItems() []*CloseReasonGroup {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreateCloseReasonGroupRequest message for creating a new close reason group
type CreateCloseReasonGroupRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateCloseReasonGroupRequest) Reset() {
	*x = CreateCloseReasonGroupRequest{}
	mi := &file_close_reason_group_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateCloseReasonGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCloseReasonGroupRequest) ProtoMessage() {}

func (x *CreateCloseReasonGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCloseReasonGroupRequest.ProtoReflect.Descriptor instead.
func (*CreateCloseReasonGroupRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{2}
}

func (x *CreateCloseReasonGroupRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateCloseReasonGroupRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type InputCloseReasonGroup struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputCloseReasonGroup) Reset() {
	*x = InputCloseReasonGroup{}
	mi := &file_close_reason_group_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputCloseReasonGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputCloseReasonGroup) ProtoMessage() {}

func (x *InputCloseReasonGroup) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputCloseReasonGroup.ProtoReflect.Descriptor instead.
func (*InputCloseReasonGroup) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{3}
}

func (x *InputCloseReasonGroup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputCloseReasonGroup) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// UpdateCloseReasonGroupRequest message for updating an existing close reason group
type UpdateCloseReasonGroupRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Input *InputCloseReasonGroup `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,3,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCloseReasonGroupRequest) Reset() {
	*x = UpdateCloseReasonGroupRequest{}
	mi := &file_close_reason_group_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCloseReasonGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCloseReasonGroupRequest) ProtoMessage() {}

func (x *UpdateCloseReasonGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCloseReasonGroupRequest.ProtoReflect.Descriptor instead.
func (*UpdateCloseReasonGroupRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateCloseReasonGroupRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateCloseReasonGroupRequest) GetInput() *InputCloseReasonGroup {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateCloseReasonGroupRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// DeleteCloseReasonGroupRequest message for deleting an existing close reason group
type DeleteCloseReasonGroupRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCloseReasonGroupRequest) Reset() {
	*x = DeleteCloseReasonGroupRequest{}
	mi := &file_close_reason_group_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCloseReasonGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCloseReasonGroupRequest) ProtoMessage() {}

func (x *DeleteCloseReasonGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCloseReasonGroupRequest.ProtoReflect.Descriptor instead.
func (*DeleteCloseReasonGroupRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteCloseReasonGroupRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// ListCloseReasonGroupsRequest message for listing or searching close reason groups
type ListCloseReasonGroupsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Fields        []string               `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	Sort          []string               `protobuf:"bytes,4,rep,name=sort,proto3" json:"sort,omitempty"`
	Id            []int64                `protobuf:"varint,5,rep,packed,name=id,proto3" json:"id,omitempty"`
	Q             string                 `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCloseReasonGroupsRequest) Reset() {
	*x = ListCloseReasonGroupsRequest{}
	mi := &file_close_reason_group_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCloseReasonGroupsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCloseReasonGroupsRequest) ProtoMessage() {}

func (x *ListCloseReasonGroupsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCloseReasonGroupsRequest.ProtoReflect.Descriptor instead.
func (*ListCloseReasonGroupsRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{6}
}

func (x *ListCloseReasonGroupsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCloseReasonGroupsRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListCloseReasonGroupsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListCloseReasonGroupsRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListCloseReasonGroupsRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListCloseReasonGroupsRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

// LocateCloseReasonGroupRequest message for locating a specific close reason group by ID
type LocateCloseReasonGroupRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Fields        []string               `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateCloseReasonGroupRequest) Reset() {
	*x = LocateCloseReasonGroupRequest{}
	mi := &file_close_reason_group_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateCloseReasonGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateCloseReasonGroupRequest) ProtoMessage() {}

func (x *LocateCloseReasonGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateCloseReasonGroupRequest.ProtoReflect.Descriptor instead.
func (*LocateCloseReasonGroupRequest) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{7}
}

func (x *LocateCloseReasonGroupRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateCloseReasonGroupRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocateCloseReasonGroupResponse message contains a single close reason group entity
type LocateCloseReasonGroupResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	CloseReasonGroup *CloseReasonGroup      `protobuf:"bytes,1,opt,name=close_reason_group,json=closeReasonGroup,proto3" json:"close_reason_group,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *LocateCloseReasonGroupResponse) Reset() {
	*x = LocateCloseReasonGroupResponse{}
	mi := &file_close_reason_group_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateCloseReasonGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateCloseReasonGroupResponse) ProtoMessage() {}

func (x *LocateCloseReasonGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_close_reason_group_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateCloseReasonGroupResponse.ProtoReflect.Descriptor instead.
func (*LocateCloseReasonGroupResponse) Descriptor() ([]byte, []int) {
	return file_close_reason_group_proto_rawDescGZIP(), []int{8}
}

func (x *LocateCloseReasonGroupResponse) GetCloseReasonGroup() *CloseReasonGroup {
	if x != nil {
		return x.CloseReasonGroup
	}
	return nil
}

var File_close_reason_group_proto protoreflect.FileDescriptor

var file_close_reason_group_proto_rawDesc = string([]byte{
	0x0a, 0x18, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xf6, 0x01, 0x0a, 0x10, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x2e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x22, 0x75, 0x0a, 0x14, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x35, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22,
	0x63, 0x0a, 0x1d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x0c, 0x92, 0x41, 0x09, 0x0a, 0x07, 0xd2, 0x01, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x22, 0x4d, 0x0a, 0x15, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43, 0x6c, 0x6f,
	0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0xb2, 0x01, 0x0a, 0x1d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c,
	0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3a, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x12, 0x39, 0x0a, 0x0b, 0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02,
	0x5e, 0x24, 0xfa, 0xd2, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45,
	0x57, 0x52, 0x09, 0x78, 0x4a, 0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x3a, 0x0a, 0x92, 0x41,
	0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x3b, 0x0a, 0x1d, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05,
	0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x90, 0x01, 0x0a, 0x1c, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c,
	0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x22, 0x47, 0x0a, 0x1d, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x22, 0x6f, 0x0a, 0x1e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x12, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1f, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x10, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x32, 0xb0, 0x08, 0x0a, 0x11, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0xda, 0x01, 0x0a, 0x15, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x12, 0x2b, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x23, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x6f, 0x92, 0x41, 0x46, 0x12, 0x44, 0x52, 0x65, 0x74, 0x72, 0x69,
	0x65, 0x76, 0x65, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x63, 0x6c,
	0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x73, 0x20, 0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x63, 0x6c, 0x6f, 0x73,
	0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x90,
	0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x12, 0x1a, 0x2f, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0xb6, 0x01, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x12, 0x2c, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43,
	0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22,
	0x4d, 0x92, 0x41, 0x21, 0x12, 0x1f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x6e,
	0x65, 0x77, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x90, 0xb5, 0x18, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1f, 0x3a,
	0x01, 0x2a, 0x22, 0x1a, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65,
	0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0xf0,
	0x01, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x2c, 0x2e, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x86, 0x01, 0x92, 0x41, 0x27, 0x12, 0x25,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69,
	0x6e, 0x67, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x52, 0x3a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5a, 0x28, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x32,
	0x1f, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d,
	0x1a, 0x1f, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x7b, 0x69, 0x64,
	0x7d, 0x12, 0xb4, 0x01, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x2c, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x4b, 0x92, 0x41, 0x1d,
	0x12, 0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61, 0x20, 0x63, 0x6c, 0x6f, 0x73, 0x65,
	0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x90, 0xb5, 0x18,
	0x03, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x2a, 0x1f, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0xc8, 0x01, 0x0a, 0x16, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x12, 0x2c, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x52, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x51, 0x92, 0x41, 0x23, 0x12, 0x21, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x20, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x20, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x21, 0x12, 0x1f, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6c, 0x6f, 0x73, 0x65,
	0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x1a, 0x10, 0x8a, 0xb5, 0x18, 0x0c, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x61, 0x72, 0x69, 0x65, 0x73, 0x42, 0xa9, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0x15, 0x43, 0x6c, 0x6f,
	0x73, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3b, 0x63, 0x61, 0x73, 0x65, 0x73, 0xa2, 0x02,
	0x03, 0x57, 0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x43,
	0x61, 0x73, 0x65, 0x73, 0xca, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43,
	0x61, 0x73, 0x65, 0x73, 0xe2, 0x02, 0x19, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43,
	0x61, 0x73, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x0e, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a, 0x3a, 0x43, 0x61, 0x73, 0x65,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_close_reason_group_proto_rawDescOnce sync.Once
	file_close_reason_group_proto_rawDescData []byte
)

func file_close_reason_group_proto_rawDescGZIP() []byte {
	file_close_reason_group_proto_rawDescOnce.Do(func() {
		file_close_reason_group_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_close_reason_group_proto_rawDesc), len(file_close_reason_group_proto_rawDesc)))
	})
	return file_close_reason_group_proto_rawDescData
}

var file_close_reason_group_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_close_reason_group_proto_goTypes = []any{
	(*CloseReasonGroup)(nil),               // 0: webitel.cases.CloseReasonGroup
	(*CloseReasonGroupList)(nil),           // 1: webitel.cases.CloseReasonGroupList
	(*CreateCloseReasonGroupRequest)(nil),  // 2: webitel.cases.CreateCloseReasonGroupRequest
	(*InputCloseReasonGroup)(nil),          // 3: webitel.cases.InputCloseReasonGroup
	(*UpdateCloseReasonGroupRequest)(nil),  // 4: webitel.cases.UpdateCloseReasonGroupRequest
	(*DeleteCloseReasonGroupRequest)(nil),  // 5: webitel.cases.DeleteCloseReasonGroupRequest
	(*ListCloseReasonGroupsRequest)(nil),   // 6: webitel.cases.ListCloseReasonGroupsRequest
	(*LocateCloseReasonGroupRequest)(nil),  // 7: webitel.cases.LocateCloseReasonGroupRequest
	(*LocateCloseReasonGroupResponse)(nil), // 8: webitel.cases.LocateCloseReasonGroupResponse
	(*Lookup)(nil),                         // 9: general.Lookup
}
var file_close_reason_group_proto_depIdxs = []int32{
	9,  // 0: webitel.cases.CloseReasonGroup.created_by:type_name -> general.Lookup
	9,  // 1: webitel.cases.CloseReasonGroup.updated_by:type_name -> general.Lookup
	0,  // 2: webitel.cases.CloseReasonGroupList.items:type_name -> webitel.cases.CloseReasonGroup
	3,  // 3: webitel.cases.UpdateCloseReasonGroupRequest.input:type_name -> webitel.cases.InputCloseReasonGroup
	0,  // 4: webitel.cases.LocateCloseReasonGroupResponse.close_reason_group:type_name -> webitel.cases.CloseReasonGroup
	6,  // 5: webitel.cases.CloseReasonGroups.ListCloseReasonGroups:input_type -> webitel.cases.ListCloseReasonGroupsRequest
	2,  // 6: webitel.cases.CloseReasonGroups.CreateCloseReasonGroup:input_type -> webitel.cases.CreateCloseReasonGroupRequest
	4,  // 7: webitel.cases.CloseReasonGroups.UpdateCloseReasonGroup:input_type -> webitel.cases.UpdateCloseReasonGroupRequest
	5,  // 8: webitel.cases.CloseReasonGroups.DeleteCloseReasonGroup:input_type -> webitel.cases.DeleteCloseReasonGroupRequest
	7,  // 9: webitel.cases.CloseReasonGroups.LocateCloseReasonGroup:input_type -> webitel.cases.LocateCloseReasonGroupRequest
	1,  // 10: webitel.cases.CloseReasonGroups.ListCloseReasonGroups:output_type -> webitel.cases.CloseReasonGroupList
	0,  // 11: webitel.cases.CloseReasonGroups.CreateCloseReasonGroup:output_type -> webitel.cases.CloseReasonGroup
	0,  // 12: webitel.cases.CloseReasonGroups.UpdateCloseReasonGroup:output_type -> webitel.cases.CloseReasonGroup
	0,  // 13: webitel.cases.CloseReasonGroups.DeleteCloseReasonGroup:output_type -> webitel.cases.CloseReasonGroup
	8,  // 14: webitel.cases.CloseReasonGroups.LocateCloseReasonGroup:output_type -> webitel.cases.LocateCloseReasonGroupResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_close_reason_group_proto_init() }
func file_close_reason_group_proto_init() {
	if File_close_reason_group_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_close_reason_group_proto_rawDesc), len(file_close_reason_group_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_close_reason_group_proto_goTypes,
		DependencyIndexes: file_close_reason_group_proto_depIdxs,
		MessageInfos:      file_close_reason_group_proto_msgTypes,
	}.Build()
	File_close_reason_group_proto = out.File
	file_close_reason_group_proto_goTypes = nil
	file_close_reason_group_proto_depIdxs = nil
}
