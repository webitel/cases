// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: sla.proto

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
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SLA message represents an SLA entity with metadata
type SLA struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the SLA
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the SLA - required
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the SLA - optional
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Validity period for the SLA - optional
	ValidFrom int64 `protobuf:"varint,4,opt,name=valid_from,json=validFrom,proto3" json:"valid_from,omitempty"`
	ValidTo   int64 `protobuf:"varint,5,opt,name=valid_to,json=validTo,proto3" json:"valid_to,omitempty"`
	// Calendar ID from the "Calendars" dictionary - required
	Calendar *Lookup `protobuf:"bytes,6,opt,name=calendar,proto3" json:"calendar,omitempty"`
	// Reaction time - required
	ReactionTime int64 `protobuf:"varint,7,opt,name=reaction_time,json=reactionTime,proto3" json:"reaction_time,omitempty"`
	// Resolution time - required
	ResolutionTime int64 `protobuf:"varint,8,opt,name=resolution_time,json=resolutionTime,proto3" json:"resolution_time,omitempty"`
	// CreatedAt timestamp of the SLA
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the SLA
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the SLA
	CreatedBy *Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the SLA
	UpdatedBy     *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SLA) Reset() {
	*x = SLA{}
	mi := &file_sla_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SLA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SLA) ProtoMessage() {}

func (x *SLA) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SLA.ProtoReflect.Descriptor instead.
func (*SLA) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{0}
}

func (x *SLA) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SLA) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SLA) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SLA) GetValidFrom() int64 {
	if x != nil {
		return x.ValidFrom
	}
	return 0
}

func (x *SLA) GetValidTo() int64 {
	if x != nil {
		return x.ValidTo
	}
	return 0
}

func (x *SLA) GetCalendar() *Lookup {
	if x != nil {
		return x.Calendar
	}
	return nil
}

func (x *SLA) GetReactionTime() int64 {
	if x != nil {
		return x.ReactionTime
	}
	return 0
}

func (x *SLA) GetResolutionTime() int64 {
	if x != nil {
		return x.ResolutionTime
	}
	return 0
}

func (x *SLA) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *SLA) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *SLA) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *SLA) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// InputSLA message for inputting SLA data
type InputSLA struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Name of the SLA
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the SLA
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Validity period for the SLA
	ValidFrom int64 `protobuf:"varint,4,opt,name=valid_from,json=validFrom,proto3" json:"valid_from,omitempty"`
	ValidTo   int64 `protobuf:"varint,5,opt,name=valid_to,json=validTo,proto3" json:"valid_to,omitempty"`
	// Calendar ID from the "Calendars" dictionary
	CalendarId int64 `protobuf:"varint,6,opt,name=calendar_id,json=calendarId,proto3" json:"calendar_id,omitempty"`
	// Reaction time
	ReactionTime int64 `protobuf:"varint,7,opt,name=reaction_time,json=reactionTime,proto3" json:"reaction_time,omitempty"`
	// Resolution time
	ResolutionTime int64 `protobuf:"varint,8,opt,name=resolution_time,json=resolutionTime,proto3" json:"resolution_time,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *InputSLA) Reset() {
	*x = InputSLA{}
	mi := &file_sla_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputSLA) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputSLA) ProtoMessage() {}

func (x *InputSLA) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputSLA.ProtoReflect.Descriptor instead.
func (*InputSLA) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{1}
}

func (x *InputSLA) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputSLA) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *InputSLA) GetValidFrom() int64 {
	if x != nil {
		return x.ValidFrom
	}
	return 0
}

func (x *InputSLA) GetValidTo() int64 {
	if x != nil {
		return x.ValidTo
	}
	return 0
}

func (x *InputSLA) GetCalendarId() int64 {
	if x != nil {
		return x.CalendarId
	}
	return 0
}

func (x *InputSLA) GetReactionTime() int64 {
	if x != nil {
		return x.ReactionTime
	}
	return 0
}

func (x *InputSLA) GetResolutionTime() int64 {
	if x != nil {
		return x.ResolutionTime
	}
	return 0
}

// SLAList message contains a list of SLA items with pagination
type SLAList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next          bool                   `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*SLA                 `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SLAList) Reset() {
	*x = SLAList{}
	mi := &file_sla_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SLAList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SLAList) ProtoMessage() {}

func (x *SLAList) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SLAList.ProtoReflect.Descriptor instead.
func (*SLAList) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{2}
}

func (x *SLAList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SLAList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *SLAList) GetItems() []*SLA {
	if x != nil {
		return x.Items
	}
	return nil
}

// CreateSLARequest message for creating a new SLA
type CreateSLARequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// InputSLA message for creating a new SLA
	// Name of the SLA
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the SLA
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Validity period for the SLA
	ValidFrom int64 `protobuf:"varint,4,opt,name=valid_from,json=validFrom,proto3" json:"valid_from,omitempty"`
	ValidTo   int64 `protobuf:"varint,5,opt,name=valid_to,json=validTo,proto3" json:"valid_to,omitempty"`
	// Calendar ID from the "Calendars" dictionary
	CalendarId int64 `protobuf:"varint,6,opt,name=calendar_id,json=calendarId,proto3" json:"calendar_id,omitempty"`
	// Reaction time
	ReactionTime int64 `protobuf:"varint,7,opt,name=reaction_time,json=reactionTime,proto3" json:"reaction_time,omitempty"`
	// Resolution time
	ResolutionTime int64 `protobuf:"varint,8,opt,name=resolution_time,json=resolutionTime,proto3" json:"resolution_time,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *CreateSLARequest) Reset() {
	*x = CreateSLARequest{}
	mi := &file_sla_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSLARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSLARequest) ProtoMessage() {}

func (x *CreateSLARequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSLARequest.ProtoReflect.Descriptor instead.
func (*CreateSLARequest) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{3}
}

func (x *CreateSLARequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateSLARequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateSLARequest) GetValidFrom() int64 {
	if x != nil {
		return x.ValidFrom
	}
	return 0
}

func (x *CreateSLARequest) GetValidTo() int64 {
	if x != nil {
		return x.ValidTo
	}
	return 0
}

func (x *CreateSLARequest) GetCalendarId() int64 {
	if x != nil {
		return x.CalendarId
	}
	return 0
}

func (x *CreateSLARequest) GetReactionTime() int64 {
	if x != nil {
		return x.ReactionTime
	}
	return 0
}

func (x *CreateSLARequest) GetResolutionTime() int64 {
	if x != nil {
		return x.ResolutionTime
	}
	return 0
}

// UpdateSLARequest message for updating an existing SLA
type UpdateSLARequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Input *InputSLA              `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,4,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateSLARequest) Reset() {
	*x = UpdateSLARequest{}
	mi := &file_sla_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateSLARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSLARequest) ProtoMessage() {}

func (x *UpdateSLARequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSLARequest.ProtoReflect.Descriptor instead.
func (*UpdateSLARequest) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateSLARequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateSLARequest) GetInput() *InputSLA {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateSLARequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// DeleteSLARequest message for deleting an existing SLA
type DeleteSLARequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteSLARequest) Reset() {
	*x = DeleteSLARequest{}
	mi := &file_sla_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteSLARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSLARequest) ProtoMessage() {}

func (x *DeleteSLARequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSLARequest.ProtoReflect.Descriptor instead.
func (*DeleteSLARequest) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteSLARequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// ListSLARequest message for listing or searching SLAs
type ListSLARequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
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
	Q             string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSLARequest) Reset() {
	*x = ListSLARequest{}
	mi := &file_sla_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSLARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSLARequest) ProtoMessage() {}

func (x *ListSLARequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSLARequest.ProtoReflect.Descriptor instead.
func (*ListSLARequest) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{6}
}

func (x *ListSLARequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListSLARequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListSLARequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListSLARequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListSLARequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListSLARequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

// LocateSLARequest message for locating a specific SLA by ID
type LocateSLARequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Fields        []string               `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateSLARequest) Reset() {
	*x = LocateSLARequest{}
	mi := &file_sla_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateSLARequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateSLARequest) ProtoMessage() {}

func (x *LocateSLARequest) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateSLARequest.ProtoReflect.Descriptor instead.
func (*LocateSLARequest) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{7}
}

func (x *LocateSLARequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateSLARequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// LocateSLAResponse message contains a single SLA entity
type LocateSLAResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sla           *SLA                   `protobuf:"bytes,1,opt,name=sla,proto3" json:"sla,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateSLAResponse) Reset() {
	*x = LocateSLAResponse{}
	mi := &file_sla_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateSLAResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateSLAResponse) ProtoMessage() {}

func (x *LocateSLAResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sla_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateSLAResponse.ProtoReflect.Descriptor instead.
func (*LocateSLAResponse) Descriptor() ([]byte, []int) {
	return file_sla_proto_rawDescGZIP(), []int{8}
}

func (x *LocateSLAResponse) GetSla() *SLA {
	if x != nil {
		return x.Sla
	}
	return nil
}

var File_sla_proto protoreflect.FileDescriptor

var file_sla_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x6c, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x0d, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x9e, 0x03, 0x0a, 0x03, 0x53, 0x4c, 0x41, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d,
	0x0a, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x12, 0x19, 0x0a,
	0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x74, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x12, 0x2b, 0x0a, 0x08, 0x63, 0x61, 0x6c, 0x65,
	0x6e, 0x64, 0x61, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e,
	0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x08, 0x63, 0x61, 0x6c,
	0x65, 0x6e, 0x64, 0x61, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65,
	0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0e, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18,
	0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e,
	0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42,
	0x79, 0x12, 0x2e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18,
	0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e,
	0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42,
	0x79, 0x22, 0xe9, 0x01, 0x0a, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x4c, 0x41, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x66, 0x72,
	0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x46,
	0x72, 0x6f, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x74, 0x6f, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x12, 0x1f,
	0x0a, 0x0b, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x72,
	0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x5b, 0x0a,
	0x07, 0x53, 0x4c, 0x41, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74,
	0x12, 0x28, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x12, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x53, 0x4c, 0x41, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xaf, 0x02, 0x0a, 0x10, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x66,
	0x72, 0x6f, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x46, 0x72, 0x6f, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x74, 0x6f,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x12,
	0x1f, 0x0a, 0x0b, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x49, 0x64,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e,
	0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x3a, 0x3c,
	0x92, 0x41, 0x39, 0x0a, 0x37, 0xd2, 0x01, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0xd2, 0x01, 0x0b, 0x63,
	0x61, 0x6c, 0x65, 0x6e, 0x64, 0x61, 0x72, 0x5f, 0x69, 0x64, 0xd2, 0x01, 0x0d, 0x72, 0x65, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0xd2, 0x01, 0x0f, 0x72, 0x65, 0x73,
	0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x98, 0x01, 0x0a,
	0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x2d, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x4c, 0x41, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x12, 0x39, 0x0a, 0x0b, 0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02, 0x5e,
	0x24, 0xfa, 0xd2, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57,
	0x52, 0x09, 0x78, 0x4a, 0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x3a, 0x0a, 0x92, 0x41, 0x07,
	0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x2e, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x3a, 0x0a, 0x92, 0x41, 0x07,
	0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x82, 0x01, 0x0a, 0x0e, 0x4c, 0x69, 0x73, 0x74,
	0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f,
	0x72, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c,
	0x0a, 0x01, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x22, 0x3a, 0x0a, 0x10,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x39, 0x0a, 0x11, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a,
	0x03, 0x73, 0x6c, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x4c, 0x41, 0x52, 0x03,
	0x73, 0x6c, 0x61, 0x32, 0xb4, 0x05, 0x0a, 0x04, 0x53, 0x4c, 0x41, 0x73, 0x12, 0x8f, 0x01, 0x0a,
	0x08, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x4c, 0x41, 0x73, 0x12, 0x1d, 0x2e, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x4c,
	0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x4c, 0x41, 0x4c, 0x69, 0x73, 0x74,
	0x22, 0x4c, 0x92, 0x41, 0x32, 0x12, 0x30, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20,
	0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x53, 0x4c, 0x41, 0x73, 0x20, 0x6f,
	0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x53, 0x4c, 0x41, 0x20, 0x63, 0x6f, 0x6e,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0d, 0x12, 0x0b, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6c, 0x61, 0x73, 0x12, 0x71,
	0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x12, 0x1f, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x4c, 0x41,
	0x22, 0x2f, 0x92, 0x41, 0x12, 0x12, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20,
	0x6e, 0x65, 0x77, 0x20, 0x53, 0x4c, 0x41, 0x90, 0xb5, 0x18, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x10, 0x3a, 0x01, 0x2a, 0x22, 0x0b, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6c, 0x61,
	0x73, 0x12, 0x9b, 0x01, 0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x12,
	0x1f, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x53, 0x4c, 0x41, 0x22, 0x59, 0x92, 0x41, 0x18, 0x12, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x78, 0x69, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x53, 0x4c,
	0x41, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x34, 0x3a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x5a, 0x19, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x32, 0x10, 0x2f, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2f, 0x73, 0x6c, 0x61, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x10, 0x2f,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6c, 0x61, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12,
	0x70, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x12, 0x1f, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x4c,
	0x41, 0x22, 0x2e, 0x92, 0x41, 0x0f, 0x12, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61,
	0x6e, 0x20, 0x53, 0x4c, 0x41, 0x90, 0xb5, 0x18, 0x03, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x2a,
	0x10, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6c, 0x61, 0x73, 0x2f, 0x7b, 0x69, 0x64,
	0x7d, 0x12, 0x84, 0x01, 0x0a, 0x09, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x12,
	0x1f, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x20, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x4c, 0x41, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x34, 0x92, 0x41, 0x15, 0x12, 0x13, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x20,
	0x61, 0x6e, 0x20, 0x53, 0x4c, 0x41, 0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x90, 0xb5, 0x18, 0x01,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73,
	0x6c, 0x61, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x10, 0x8a, 0xb5, 0x18, 0x0c, 0x64, 0x69,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x72, 0x69, 0x65, 0x73, 0x42, 0x9c, 0x01, 0x0a, 0x11, 0x63,
	0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x42, 0x08, 0x53, 0x6c, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x3b, 0x63, 0x61, 0x73, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x57, 0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x73, 0xca, 0x02, 0x0d, 0x57,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0xe2, 0x02, 0x19, 0x57,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x57, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x3a, 0x3a, 0x43, 0x61, 0x73, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_sla_proto_rawDescOnce sync.Once
	file_sla_proto_rawDescData = file_sla_proto_rawDesc
)

func file_sla_proto_rawDescGZIP() []byte {
	file_sla_proto_rawDescOnce.Do(func() {
		file_sla_proto_rawDescData = protoimpl.X.CompressGZIP(file_sla_proto_rawDescData)
	})
	return file_sla_proto_rawDescData
}

var file_sla_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_sla_proto_goTypes = []any{
	(*SLA)(nil),               // 0: webitel.cases.SLA
	(*InputSLA)(nil),          // 1: webitel.cases.InputSLA
	(*SLAList)(nil),           // 2: webitel.cases.SLAList
	(*CreateSLARequest)(nil),  // 3: webitel.cases.CreateSLARequest
	(*UpdateSLARequest)(nil),  // 4: webitel.cases.UpdateSLARequest
	(*DeleteSLARequest)(nil),  // 5: webitel.cases.DeleteSLARequest
	(*ListSLARequest)(nil),    // 6: webitel.cases.ListSLARequest
	(*LocateSLARequest)(nil),  // 7: webitel.cases.LocateSLARequest
	(*LocateSLAResponse)(nil), // 8: webitel.cases.LocateSLAResponse
	(*Lookup)(nil),            // 9: general.Lookup
}
var file_sla_proto_depIdxs = []int32{
	9,  // 0: webitel.cases.SLA.calendar:type_name -> general.Lookup
	9,  // 1: webitel.cases.SLA.created_by:type_name -> general.Lookup
	9,  // 2: webitel.cases.SLA.updated_by:type_name -> general.Lookup
	0,  // 3: webitel.cases.SLAList.items:type_name -> webitel.cases.SLA
	1,  // 4: webitel.cases.UpdateSLARequest.input:type_name -> webitel.cases.InputSLA
	0,  // 5: webitel.cases.LocateSLAResponse.sla:type_name -> webitel.cases.SLA
	6,  // 6: webitel.cases.SLAs.ListSLAs:input_type -> webitel.cases.ListSLARequest
	3,  // 7: webitel.cases.SLAs.CreateSLA:input_type -> webitel.cases.CreateSLARequest
	4,  // 8: webitel.cases.SLAs.UpdateSLA:input_type -> webitel.cases.UpdateSLARequest
	5,  // 9: webitel.cases.SLAs.DeleteSLA:input_type -> webitel.cases.DeleteSLARequest
	7,  // 10: webitel.cases.SLAs.LocateSLA:input_type -> webitel.cases.LocateSLARequest
	2,  // 11: webitel.cases.SLAs.ListSLAs:output_type -> webitel.cases.SLAList
	0,  // 12: webitel.cases.SLAs.CreateSLA:output_type -> webitel.cases.SLA
	0,  // 13: webitel.cases.SLAs.UpdateSLA:output_type -> webitel.cases.SLA
	0,  // 14: webitel.cases.SLAs.DeleteSLA:output_type -> webitel.cases.SLA
	8,  // 15: webitel.cases.SLAs.LocateSLA:output_type -> webitel.cases.LocateSLAResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_sla_proto_init() }
func file_sla_proto_init() {
	if File_sla_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sla_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sla_proto_goTypes,
		DependencyIndexes: file_sla_proto_depIdxs,
		MessageInfos:      file_sla_proto_msgTypes,
	}.Build()
	File_sla_proto = out.File
	file_sla_proto_rawDesc = nil
	file_sla_proto_goTypes = nil
	file_sla_proto_depIdxs = nil
}
