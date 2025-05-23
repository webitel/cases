// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: sla_condition.proto

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

// SLACondition message represents an SLACondition entity with metadata
type SLACondition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the SLACondition
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the SLACondition - required
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Priorities associated with the SLACondition - fetched as Lookup entities [ Priority name + ID ]
	Priorities []*Lookup `protobuf:"bytes,3,rep,name=priorities,proto3" json:"priorities,omitempty"`
	// Reaction time - required
	ReactionTime int64 `protobuf:"varint,4,opt,name=reaction_time,json=reactionTime,proto3" json:"reaction_time,omitempty"`
	// Resolution time - required
	ResolutionTime int64 `protobuf:"varint,5,opt,name=resolution_time,json=resolutionTime,proto3" json:"resolution_time,omitempty"`
	// SLA ID associated with the SLACondition
	SlaId int64 `protobuf:"varint,6,opt,name=sla_id,json=slaId,proto3" json:"sla_id,omitempty"`
	// CreatedAt timestamp of the SLACondition
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the SLACondition
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the SLACondition
	CreatedBy *Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the SLACondition
	UpdatedBy     *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SLACondition) Reset() {
	*x = SLACondition{}
	mi := &file_sla_condition_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SLACondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SLACondition) ProtoMessage() {}

func (x *SLACondition) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SLACondition.ProtoReflect.Descriptor instead.
func (*SLACondition) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{0}
}

func (x *SLACondition) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SLACondition) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SLACondition) GetPriorities() []*Lookup {
	if x != nil {
		return x.Priorities
	}
	return nil
}

func (x *SLACondition) GetReactionTime() int64 {
	if x != nil {
		return x.ReactionTime
	}
	return 0
}

func (x *SLACondition) GetResolutionTime() int64 {
	if x != nil {
		return x.ResolutionTime
	}
	return 0
}

func (x *SLACondition) GetSlaId() int64 {
	if x != nil {
		return x.SlaId
	}
	return 0
}

func (x *SLACondition) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *SLACondition) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *SLACondition) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *SLACondition) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// InputSLACondition message for inputting SLACondition data
type InputSLACondition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Name  string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// List of priority IDs for creation
	Priorities     []*Lookup `protobuf:"bytes,3,rep,name=priorities,proto3" json:"priorities,omitempty"`
	ReactionTime   int64     `protobuf:"varint,4,opt,name=reaction_time,json=reactionTime,proto3" json:"reaction_time,omitempty"`
	ResolutionTime int64     `protobuf:"varint,5,opt,name=resolution_time,json=resolutionTime,proto3" json:"resolution_time,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *InputSLACondition) Reset() {
	*x = InputSLACondition{}
	mi := &file_sla_condition_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputSLACondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputSLACondition) ProtoMessage() {}

func (x *InputSLACondition) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputSLACondition.ProtoReflect.Descriptor instead.
func (*InputSLACondition) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{1}
}

func (x *InputSLACondition) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputSLACondition) GetPriorities() []*Lookup {
	if x != nil {
		return x.Priorities
	}
	return nil
}

func (x *InputSLACondition) GetReactionTime() int64 {
	if x != nil {
		return x.ReactionTime
	}
	return 0
}

func (x *InputSLACondition) GetResolutionTime() int64 {
	if x != nil {
		return x.ResolutionTime
	}
	return 0
}

// SLAConditionList message contains a list of SLACondition items with pagination
type SLAConditionList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next          bool                   `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*SLACondition        `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SLAConditionList) Reset() {
	*x = SLAConditionList{}
	mi := &file_sla_condition_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SLAConditionList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SLAConditionList) ProtoMessage() {}

func (x *SLAConditionList) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SLAConditionList.ProtoReflect.Descriptor instead.
func (*SLAConditionList) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{2}
}

func (x *SLAConditionList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SLAConditionList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *SLAConditionList) GetItems() []*SLACondition {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreateSLAConditionRequest message for creating a new SLACondition
type CreateSLAConditionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Input         *InputSLACondition     `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	SlaId         int64                  `protobuf:"varint,2,opt,name=sla_id,json=slaId,proto3" json:"sla_id,omitempty"`
	Fields        []string               `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateSLAConditionRequest) Reset() {
	*x = CreateSLAConditionRequest{}
	mi := &file_sla_condition_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSLAConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSLAConditionRequest) ProtoMessage() {}

func (x *CreateSLAConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSLAConditionRequest.ProtoReflect.Descriptor instead.
func (*CreateSLAConditionRequest) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{3}
}

func (x *CreateSLAConditionRequest) GetInput() *InputSLACondition {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *CreateSLAConditionRequest) GetSlaId() int64 {
	if x != nil {
		return x.SlaId
	}
	return 0
}

func (x *CreateSLAConditionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// UpdateSLAConditionRequest message for updating an existing SLACondition
type UpdateSLAConditionRequest struct {
	state  protoimpl.MessageState `protogen:"open.v1"`
	Id     int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	SlaId  int64                  `protobuf:"varint,2,opt,name=sla_id,json=slaId,proto3" json:"sla_id,omitempty"`
	Input  *InputSLACondition     `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`
	Fields []string               `protobuf:"bytes,4,rep,name=fields,proto3" json:"fields,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,5,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateSLAConditionRequest) Reset() {
	*x = UpdateSLAConditionRequest{}
	mi := &file_sla_condition_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateSLAConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSLAConditionRequest) ProtoMessage() {}

func (x *UpdateSLAConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSLAConditionRequest.ProtoReflect.Descriptor instead.
func (*UpdateSLAConditionRequest) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateSLAConditionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateSLAConditionRequest) GetSlaId() int64 {
	if x != nil {
		return x.SlaId
	}
	return 0
}

func (x *UpdateSLAConditionRequest) GetInput() *InputSLACondition {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateSLAConditionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *UpdateSLAConditionRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// DeleteSLAConditionRequest message for deleting an existing SLACondition
type DeleteSLAConditionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SlaId         int64                  `protobuf:"varint,1,opt,name=sla_id,json=slaId,proto3" json:"sla_id,omitempty"`
	Id            int64                  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteSLAConditionRequest) Reset() {
	*x = DeleteSLAConditionRequest{}
	mi := &file_sla_condition_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteSLAConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSLAConditionRequest) ProtoMessage() {}

func (x *DeleteSLAConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSLAConditionRequest.ProtoReflect.Descriptor instead.
func (*DeleteSLAConditionRequest) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteSLAConditionRequest) GetSlaId() int64 {
	if x != nil {
		return x.SlaId
	}
	return 0
}

func (x *DeleteSLAConditionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// ListSLAConditionRequest message for listing or searching SLAConditions
type ListSLAConditionRequest struct {
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
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Filter by SLACondition Id.
	SlaConditionId int64 `protobuf:"varint,7,opt,name=sla_condition_id,json=slaConditionId,proto3" json:"sla_condition_id,omitempty"`
	// Filter by SLA Id
	SlaId int64 `protobuf:"varint,8,opt,name=sla_id,json=slaId,proto3" json:"sla_id,omitempty"`
	// filter by priority id
	PriorityId    int64 `protobuf:"varint,9,opt,name=priority_id,json=priorityId,proto3" json:"priority_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSLAConditionRequest) Reset() {
	*x = ListSLAConditionRequest{}
	mi := &file_sla_condition_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSLAConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSLAConditionRequest) ProtoMessage() {}

func (x *ListSLAConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSLAConditionRequest.ProtoReflect.Descriptor instead.
func (*ListSLAConditionRequest) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{6}
}

func (x *ListSLAConditionRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListSLAConditionRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListSLAConditionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListSLAConditionRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *ListSLAConditionRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListSLAConditionRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListSLAConditionRequest) GetSlaConditionId() int64 {
	if x != nil {
		return x.SlaConditionId
	}
	return 0
}

func (x *ListSLAConditionRequest) GetSlaId() int64 {
	if x != nil {
		return x.SlaId
	}
	return 0
}

func (x *ListSLAConditionRequest) GetPriorityId() int64 {
	if x != nil {
		return x.PriorityId
	}
	return 0
}

// LocateSLAConditionRequest message for locating a specific SLACondition by ID
type LocateSLAConditionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the SLACondition
	SlaId int64 `protobuf:"varint,1,opt,name=sla_id,json=slaId,proto3" json:"sla_id,omitempty"`
	// Unique identifier of the SLACondition
	Id int64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	// Fields to be retrieved as a result.
	Fields        []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateSLAConditionRequest) Reset() {
	*x = LocateSLAConditionRequest{}
	mi := &file_sla_condition_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateSLAConditionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateSLAConditionRequest) ProtoMessage() {}

func (x *LocateSLAConditionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateSLAConditionRequest.ProtoReflect.Descriptor instead.
func (*LocateSLAConditionRequest) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{7}
}

func (x *LocateSLAConditionRequest) GetSlaId() int64 {
	if x != nil {
		return x.SlaId
	}
	return 0
}

func (x *LocateSLAConditionRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateSLAConditionRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocateSLAConditionResponse message contains a single SLACondition entity
type LocateSLAConditionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SlaCondition  *SLACondition          `protobuf:"bytes,1,opt,name=sla_condition,json=slaCondition,proto3" json:"sla_condition,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateSLAConditionResponse) Reset() {
	*x = LocateSLAConditionResponse{}
	mi := &file_sla_condition_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateSLAConditionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateSLAConditionResponse) ProtoMessage() {}

func (x *LocateSLAConditionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sla_condition_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateSLAConditionResponse.ProtoReflect.Descriptor instead.
func (*LocateSLAConditionResponse) Descriptor() ([]byte, []int) {
	return file_sla_condition_proto_rawDescGZIP(), []int{8}
}

func (x *LocateSLAConditionResponse) GetSlaCondition() *SLACondition {
	if x != nil {
		return x.SlaCondition
	}
	return nil
}

var File_sla_condition_proto protoreflect.FileDescriptor

const file_sla_condition_proto_rawDesc = "" +
	"\n" +
	"\x13sla_condition.proto\x12\rwebitel.cases\x1a\rgeneral.proto\x1a\x1bgoogle/api/visibility.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\x1a\x1aproto/webitel/option.proto\"\xe6\x02\n" +
	"\fSLACondition\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12/\n" +
	"\n" +
	"priorities\x18\x03 \x03(\v2\x0f.general.LookupR\n" +
	"priorities\x12#\n" +
	"\rreaction_time\x18\x04 \x01(\x03R\freactionTime\x12'\n" +
	"\x0fresolution_time\x18\x05 \x01(\x03R\x0eresolutionTime\x12\x15\n" +
	"\x06sla_id\x18\x06 \x01(\x03R\x05slaId\x12\x1d\n" +
	"\n" +
	"created_at\x18\x14 \x01(\x03R\tcreatedAt\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x15 \x01(\x03R\tupdatedAt\x12.\n" +
	"\n" +
	"created_by\x18\x16 \x01(\v2\x0f.general.LookupR\tcreatedBy\x12.\n" +
	"\n" +
	"updated_by\x18\x17 \x01(\v2\x0f.general.LookupR\tupdatedBy\"\xa6\x01\n" +
	"\x11InputSLACondition\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12/\n" +
	"\n" +
	"priorities\x18\x03 \x03(\v2\x0f.general.LookupR\n" +
	"priorities\x12#\n" +
	"\rreaction_time\x18\x04 \x01(\x03R\freactionTime\x12'\n" +
	"\x0fresolution_time\x18\x05 \x01(\x03R\x0eresolutionTime\"m\n" +
	"\x10SLAConditionList\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x02 \x01(\bR\x04next\x121\n" +
	"\x05items\x18\x03 \x03(\v2\x1b.webitel.cases.SLAConditionR\x05items\"\x82\x01\n" +
	"\x19CreateSLAConditionRequest\x126\n" +
	"\x05input\x18\x01 \x01(\v2 .webitel.cases.InputSLAConditionR\x05input\x12\x15\n" +
	"\x06sla_id\x18\x02 \x01(\x03R\x05slaId\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\"\xd9\x01\n" +
	"\x19UpdateSLAConditionRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x15\n" +
	"\x06sla_id\x18\x02 \x01(\x03R\x05slaId\x126\n" +
	"\x05input\x18\x03 \x01(\v2 .webitel.cases.InputSLAConditionR\x05input\x12\x16\n" +
	"\x06fields\x18\x04 \x03(\tR\x06fields\x129\n" +
	"\vx_json_mask\x18\x05 \x03(\tB\x19\x92A\a@\x01\x8a\x01\x02^$\xfa\xd2\xe4\x93\x02\t\x12\aPREVIEWR\txJsonMask:\n" +
	"\x92A\a\n" +
	"\x05\xd2\x01\x02id\"\\\n" +
	"\x19DeleteSLAConditionRequest\x12\x15\n" +
	"\x06sla_id\x18\x01 \x01(\x03R\x05slaId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\x03R\x02id:\x18\x92A\x15\n" +
	"\x13\xd2\x01\x10sla_condition_id\"\xed\x01\n" +
	"\x17ListSLAConditionRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x12\n" +
	"\x04size\x18\x02 \x01(\x05R\x04size\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\x12\x12\n" +
	"\x04sort\x18\x04 \x01(\tR\x04sort\x12\x0e\n" +
	"\x02id\x18\x05 \x03(\x03R\x02id\x12\f\n" +
	"\x01q\x18\x06 \x01(\tR\x01q\x12(\n" +
	"\x10sla_condition_id\x18\a \x01(\x03R\x0eslaConditionId\x12\x15\n" +
	"\x06sla_id\x18\b \x01(\x03R\x05slaId\x12\x1f\n" +
	"\vpriority_id\x18\t \x01(\x03R\n" +
	"priorityId\"Z\n" +
	"\x19LocateSLAConditionRequest\x12\x15\n" +
	"\x06sla_id\x18\x01 \x01(\x03R\x05slaId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\x03R\x02id\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\"^\n" +
	"\x1aLocateSLAConditionResponse\x12@\n" +
	"\rsla_condition\x18\x01 \x01(\v2\x1b.webitel.cases.SLAConditionR\fslaCondition2\xe8\a\n" +
	"\rSLAConditions\x12\xce\x01\n" +
	"\x11ListSLAConditions\x12&.webitel.cases.ListSLAConditionRequest\x1a\x1f.webitel.cases.SLAConditionList\"p\x92AD\x12BRetrieve a list of SLAConditions or search SLACondition conditions\x90\xb5\x18\x01\x82\xd3\xe4\x93\x02\x1f\x12\x1d/slas/{sla_id}/sla_conditions\x12\xaa\x01\n" +
	"\x12CreateSLACondition\x12(.webitel.cases.CreateSLAConditionRequest\x1a\x1b.webitel.cases.SLACondition\"M\x92A\x1b\x12\x19Create a new SLACondition\x90\xb5\x18\x02\x82\xd3\xe4\x93\x02%:\x05input\"\x1c/slas/{sla_id}/sla_condition\x12\xe2\x01\n" +
	"\x12UpdateSLACondition\x12(.webitel.cases.UpdateSLAConditionRequest\x1a\x1b.webitel.cases.SLACondition\"\x84\x01\x92A!\x12\x1fUpdate an existing SLACondition\x90\xb5\x18\x02\x82\xd3\xe4\x93\x02V:\x05inputZ*:\x05input2!/slas/{sla_id}/sla_condition/{id}\x1a!/slas/{sla_id}/sla_condition/{id}\x12\xa5\x01\n" +
	"\x12DeleteSLACondition\x12(.webitel.cases.DeleteSLAConditionRequest\x1a\x1b.webitel.cases.SLACondition\"H\x92A\x18\x12\x16Delete an SLACondition\x90\xb5\x18\x02\x82\xd3\xe4\x93\x02#*!/slas/{sla_id}/sla_condition/{id}\x12\xb9\x01\n" +
	"\x12LocateSLACondition\x12(.webitel.cases.LocateSLAConditionRequest\x1a).webitel.cases.LocateSLAConditionResponse\"N\x92A\x1e\x12\x1cLocate an SLACondition by ID\x90\xb5\x18\x01\x82\xd3\xe4\x93\x02#\x12!/slas/{sla_id}/sla_condition/{id}\x1a\x10\x8a\xb5\x18\fdictionariesB\xa5\x01\n" +
	"\x11com.webitel.casesB\x11SlaConditionProtoP\x01Z(github.com/webitel/cases/api/cases;cases\xa2\x02\x03WCX\xaa\x02\rWebitel.Cases\xca\x02\rWebitel\\Cases\xe2\x02\x19Webitel\\Cases\\GPBMetadata\xea\x02\x0eWebitel::Casesb\x06proto3"

var (
	file_sla_condition_proto_rawDescOnce sync.Once
	file_sla_condition_proto_rawDescData []byte
)

func file_sla_condition_proto_rawDescGZIP() []byte {
	file_sla_condition_proto_rawDescOnce.Do(func() {
		file_sla_condition_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_sla_condition_proto_rawDesc), len(file_sla_condition_proto_rawDesc)))
	})
	return file_sla_condition_proto_rawDescData
}

var file_sla_condition_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_sla_condition_proto_goTypes = []any{
	(*SLACondition)(nil),               // 0: webitel.cases.SLACondition
	(*InputSLACondition)(nil),          // 1: webitel.cases.InputSLACondition
	(*SLAConditionList)(nil),           // 2: webitel.cases.SLAConditionList
	(*CreateSLAConditionRequest)(nil),  // 3: webitel.cases.CreateSLAConditionRequest
	(*UpdateSLAConditionRequest)(nil),  // 4: webitel.cases.UpdateSLAConditionRequest
	(*DeleteSLAConditionRequest)(nil),  // 5: webitel.cases.DeleteSLAConditionRequest
	(*ListSLAConditionRequest)(nil),    // 6: webitel.cases.ListSLAConditionRequest
	(*LocateSLAConditionRequest)(nil),  // 7: webitel.cases.LocateSLAConditionRequest
	(*LocateSLAConditionResponse)(nil), // 8: webitel.cases.LocateSLAConditionResponse
	(*Lookup)(nil),                     // 9: general.Lookup
}
var file_sla_condition_proto_depIdxs = []int32{
	9,  // 0: webitel.cases.SLACondition.priorities:type_name -> general.Lookup
	9,  // 1: webitel.cases.SLACondition.created_by:type_name -> general.Lookup
	9,  // 2: webitel.cases.SLACondition.updated_by:type_name -> general.Lookup
	9,  // 3: webitel.cases.InputSLACondition.priorities:type_name -> general.Lookup
	0,  // 4: webitel.cases.SLAConditionList.items:type_name -> webitel.cases.SLACondition
	1,  // 5: webitel.cases.CreateSLAConditionRequest.input:type_name -> webitel.cases.InputSLACondition
	1,  // 6: webitel.cases.UpdateSLAConditionRequest.input:type_name -> webitel.cases.InputSLACondition
	0,  // 7: webitel.cases.LocateSLAConditionResponse.sla_condition:type_name -> webitel.cases.SLACondition
	6,  // 8: webitel.cases.SLAConditions.ListSLAConditions:input_type -> webitel.cases.ListSLAConditionRequest
	3,  // 9: webitel.cases.SLAConditions.CreateSLACondition:input_type -> webitel.cases.CreateSLAConditionRequest
	4,  // 10: webitel.cases.SLAConditions.UpdateSLACondition:input_type -> webitel.cases.UpdateSLAConditionRequest
	5,  // 11: webitel.cases.SLAConditions.DeleteSLACondition:input_type -> webitel.cases.DeleteSLAConditionRequest
	7,  // 12: webitel.cases.SLAConditions.LocateSLACondition:input_type -> webitel.cases.LocateSLAConditionRequest
	2,  // 13: webitel.cases.SLAConditions.ListSLAConditions:output_type -> webitel.cases.SLAConditionList
	0,  // 14: webitel.cases.SLAConditions.CreateSLACondition:output_type -> webitel.cases.SLACondition
	0,  // 15: webitel.cases.SLAConditions.UpdateSLACondition:output_type -> webitel.cases.SLACondition
	0,  // 16: webitel.cases.SLAConditions.DeleteSLACondition:output_type -> webitel.cases.SLACondition
	8,  // 17: webitel.cases.SLAConditions.LocateSLACondition:output_type -> webitel.cases.LocateSLAConditionResponse
	13, // [13:18] is the sub-list for method output_type
	8,  // [8:13] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_sla_condition_proto_init() }
func file_sla_condition_proto_init() {
	if File_sla_condition_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_sla_condition_proto_rawDesc), len(file_sla_condition_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sla_condition_proto_goTypes,
		DependencyIndexes: file_sla_condition_proto_depIdxs,
		MessageInfos:      file_sla_condition_proto_msgTypes,
	}.Build()
	File_sla_condition_proto = out.File
	file_sla_condition_proto_goTypes = nil
	file_sla_condition_proto_depIdxs = nil
}
