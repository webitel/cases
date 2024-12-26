// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: source.proto

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

// Represents a source type for the source entity.
type Type int32

const (
	// Unspecified source type.
	Type_TYPE_UNSPECIFIED Type = 0
	// Phone call source type.
	Type_CALL Type = 1
	// Chat source type.
	Type_CHAT Type = 2
	// Social media source type.
	Type_SOCIAL_MEDIA Type = 3
	// Email source type.
	Type_EMAIL Type = 4
	// API source type.
	Type_API Type = 5
	// Manual source type.
	Type_MANUAL Type = 6
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "CALL",
		2: "CHAT",
		3: "SOCIAL_MEDIA",
		4: "EMAIL",
		5: "API",
		6: "MANUAL",
	}
	Type_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"CALL":             1,
		"CHAT":             2,
		"SOCIAL_MEDIA":     3,
		"EMAIL":            4,
		"API":              5,
		"MANUAL":           6,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_source_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_source_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{0}
}

// Represents a source entity in the contact system.
type Source struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique identifier of the source
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the source
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the source
	Description string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	// Source type of the source
	Type Type `protobuf:"varint,5,opt,name=type,proto3,enum=webitel.cases.Type" json:"type,omitempty"`
	// CreatedAt timestamp of the source
	CreatedAt int64 `protobuf:"varint,20,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// UpdatedAt timestamp of the source
	UpdatedAt int64 `protobuf:"varint,21,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// CreatedBy user of the source
	CreatedBy *Lookup `protobuf:"bytes,22,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// UpdatedBy user of the source
	UpdatedBy     *Lookup `protobuf:"bytes,23,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Source) Reset() {
	*x = Source{}
	mi := &file_source_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{0}
}

func (x *Source) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Source) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Source) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Source) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_TYPE_UNSPECIFIED
}

func (x *Source) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Source) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Source) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *Source) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

// A list of sources.
type SourceList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number of the partial result.
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Have more records.
	Next bool `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	// List of sources.
	Items         []*Source `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SourceList) Reset() {
	*x = SourceList{}
	mi := &file_source_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SourceList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceList) ProtoMessage() {}

func (x *SourceList) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceList.ProtoReflect.Descriptor instead.
func (*SourceList) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{1}
}

func (x *SourceList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SourceList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *SourceList) GetItems() []*Source {
	if x != nil {
		return x.Items
	}
	return nil
}

// Request message for creating a new source.
type CreateSourceRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the source.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the source.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// The source type of the source.
	Type          Type `protobuf:"varint,4,opt,name=type,proto3,enum=webitel.cases.Type" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateSourceRequest) Reset() {
	*x = CreateSourceRequest{}
	mi := &file_source_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateSourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSourceRequest) ProtoMessage() {}

func (x *CreateSourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSourceRequest.ProtoReflect.Descriptor instead.
func (*CreateSourceRequest) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{2}
}

func (x *CreateSourceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateSourceRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateSourceRequest) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_TYPE_UNSPECIFIED
}

type InputSource struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The name of the source.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The description of the source.
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// The source type of the source.
	Type          Type `protobuf:"varint,3,opt,name=type,proto3,enum=webitel.cases.Type" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputSource) Reset() {
	*x = InputSource{}
	mi := &file_source_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputSource) ProtoMessage() {}

func (x *InputSource) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputSource.ProtoReflect.Descriptor instead.
func (*InputSource) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{3}
}

func (x *InputSource) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InputSource) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *InputSource) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_TYPE_UNSPECIFIED
}

// Request message for updating an existing source.
type UpdateSourceRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	Id    int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Input *InputSource           `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	// ---- JSON PATCH fields mask ----
	// List of JPath fields specified in body(input).
	XJsonMask     []string `protobuf:"bytes,3,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateSourceRequest) Reset() {
	*x = UpdateSourceRequest{}
	mi := &file_source_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateSourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSourceRequest) ProtoMessage() {}

func (x *UpdateSourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSourceRequest.ProtoReflect.Descriptor instead.
func (*UpdateSourceRequest) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateSourceRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateSourceRequest) GetInput() *InputSource {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *UpdateSourceRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

// Request message for deleting a source.
type DeleteSourceRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the source to delete.
	Id            int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteSourceRequest) Reset() {
	*x = DeleteSourceRequest{}
	mi := &file_source_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteSourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSourceRequest) ProtoMessage() {}

func (x *DeleteSourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSourceRequest.ProtoReflect.Descriptor instead.
func (*DeleteSourceRequest) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteSourceRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

// Request message for listing sources.
type ListSourceRequest struct {
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
	Q string `protobuf:"bytes,6,opt,name=q,proto3" json:"q,omitempty"`
	// Filter by source type.
	Type          []Type `protobuf:"varint,7,rep,packed,name=type,proto3,enum=webitel.cases.Type" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSourceRequest) Reset() {
	*x = ListSourceRequest{}
	mi := &file_source_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSourceRequest) ProtoMessage() {}

func (x *ListSourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSourceRequest.ProtoReflect.Descriptor instead.
func (*ListSourceRequest) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{6}
}

func (x *ListSourceRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListSourceRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListSourceRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListSourceRequest) GetSort() []string {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *ListSourceRequest) GetId() []int64 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ListSourceRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListSourceRequest) GetType() []Type {
	if x != nil {
		return x.Type
	}
	return nil
}

// Request message for locating a source by ID.
type LocateSourceRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The unique ID of the source to locate.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Fields to be retrieved into result.
	Fields        []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateSourceRequest) Reset() {
	*x = LocateSourceRequest{}
	mi := &file_source_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateSourceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateSourceRequest) ProtoMessage() {}

func (x *LocateSourceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateSourceRequest.ProtoReflect.Descriptor instead.
func (*LocateSourceRequest) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{7}
}

func (x *LocateSourceRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LocateSourceRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Response message for locating a source.
type LocateSourceResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The located source.
	Source        *Source `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateSourceResponse) Reset() {
	*x = LocateSourceResponse{}
	mi := &file_source_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateSourceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateSourceResponse) ProtoMessage() {}

func (x *LocateSourceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_source_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateSourceResponse.ProtoReflect.Descriptor instead.
func (*LocateSourceResponse) Descriptor() ([]byte, []int) {
	return file_source_proto_rawDescGZIP(), []int{8}
}

func (x *LocateSourceResponse) GetSource() *Source {
	if x != nil {
		return x.Source
	}
	return nil
}

var File_source_proto protoreflect.FileDescriptor

var file_source_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x0d, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c,
	0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x95, 0x02, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70,
	0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x2e, 0x0a, 0x0a, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70,
	0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x22, 0x61, 0x0a, 0x0a, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78,
	0x74, 0x12, 0x2b, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x89,
	0x01, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x13, 0x92, 0x41, 0x10, 0x0a, 0x0e, 0xd2, 0x01, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0xd2, 0x01, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x6c, 0x0a, 0x0b, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x9e, 0x01, 0x0a, 0x13, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x30, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x12, 0x39, 0x0a, 0x0b, 0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73,
	0x6b, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01,
	0x02, 0x5e, 0x24, 0xfa, 0xd2, 0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49,
	0x45, 0x57, 0x52, 0x09, 0x78, 0x4a, 0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x3a, 0x0a, 0x92,
	0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0x31, 0x0a, 0x13, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x3a, 0x0a, 0x92, 0x41, 0x07, 0x0a, 0x05, 0xd2, 0x01, 0x02, 0x69, 0x64, 0x22, 0xae, 0x01, 0x0a,
	0x11, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x01, 0x71, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x03,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3d, 0x0a,
	0x13, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x45, 0x0a, 0x14,
	0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x2a, 0x62, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x41, 0x4c, 0x4c, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x43,
	0x48, 0x41, 0x54, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x4f, 0x43, 0x49, 0x41, 0x4c, 0x5f,
	0x4d, 0x45, 0x44, 0x49, 0x41, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x4d, 0x41, 0x49, 0x4c,
	0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x4d,
	0x41, 0x4e, 0x55, 0x41, 0x4c, 0x10, 0x06, 0x32, 0xfd, 0x05, 0x0a, 0x07, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x12, 0x97, 0x01, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x69, 0x73, 0x74,
	0x22, 0x4b, 0x92, 0x41, 0x2e, 0x12, 0x2c, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20,
	0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x73, 0x20, 0x6f, 0x72, 0x20, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x20, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x73, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x80, 0x01,
	0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x22,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x35, 0x92, 0x41, 0x15, 0x12, 0x13,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x90, 0xb5, 0x18, 0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x3a, 0x01, 0x2a,
	0x22, 0x0e, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x12, 0xad, 0x01, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x62, 0x92, 0x41,
	0x1b, 0x12, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x6e, 0x20, 0x65, 0x78, 0x69,
	0x73, 0x74, 0x69, 0x6e, 0x67, 0x20, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x90, 0xb5, 0x18, 0x02,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x3a, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5a, 0x1c, 0x3a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x32, 0x13, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x13, 0x2f, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d,
	0x12, 0x7e, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x33, 0x92, 0x41, 0x11,
	0x12, 0x0f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61, 0x20, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x90, 0xb5, 0x18, 0x03, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x2a, 0x13, 0x2f, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d,
	0x12, 0x92, 0x01, 0x0a, 0x0c, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x39, 0x92, 0x41, 0x17, 0x12,
	0x15, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x20, 0x62, 0x79, 0x20, 0x49, 0x44, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15,
	0x12, 0x13, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x1a, 0x10, 0x8a, 0xb5, 0x18, 0x0c, 0x64, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x61, 0x72, 0x69, 0x65, 0x73, 0x42, 0x9f, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0x0b, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69,
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
	file_source_proto_rawDescOnce sync.Once
	file_source_proto_rawDescData = file_source_proto_rawDesc
)

func file_source_proto_rawDescGZIP() []byte {
	file_source_proto_rawDescOnce.Do(func() {
		file_source_proto_rawDescData = protoimpl.X.CompressGZIP(file_source_proto_rawDescData)
	})
	return file_source_proto_rawDescData
}

var file_source_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_source_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_source_proto_goTypes = []any{
	(Type)(0),                    // 0: webitel.cases.Type
	(*Source)(nil),               // 1: webitel.cases.Source
	(*SourceList)(nil),           // 2: webitel.cases.SourceList
	(*CreateSourceRequest)(nil),  // 3: webitel.cases.CreateSourceRequest
	(*InputSource)(nil),          // 4: webitel.cases.InputSource
	(*UpdateSourceRequest)(nil),  // 5: webitel.cases.UpdateSourceRequest
	(*DeleteSourceRequest)(nil),  // 6: webitel.cases.DeleteSourceRequest
	(*ListSourceRequest)(nil),    // 7: webitel.cases.ListSourceRequest
	(*LocateSourceRequest)(nil),  // 8: webitel.cases.LocateSourceRequest
	(*LocateSourceResponse)(nil), // 9: webitel.cases.LocateSourceResponse
	(*Lookup)(nil),               // 10: general.Lookup
}
var file_source_proto_depIdxs = []int32{
	0,  // 0: webitel.cases.Source.type:type_name -> webitel.cases.Type
	10, // 1: webitel.cases.Source.created_by:type_name -> general.Lookup
	10, // 2: webitel.cases.Source.updated_by:type_name -> general.Lookup
	1,  // 3: webitel.cases.SourceList.items:type_name -> webitel.cases.Source
	0,  // 4: webitel.cases.CreateSourceRequest.type:type_name -> webitel.cases.Type
	0,  // 5: webitel.cases.InputSource.type:type_name -> webitel.cases.Type
	4,  // 6: webitel.cases.UpdateSourceRequest.input:type_name -> webitel.cases.InputSource
	0,  // 7: webitel.cases.ListSourceRequest.type:type_name -> webitel.cases.Type
	1,  // 8: webitel.cases.LocateSourceResponse.source:type_name -> webitel.cases.Source
	7,  // 9: webitel.cases.Sources.ListSources:input_type -> webitel.cases.ListSourceRequest
	3,  // 10: webitel.cases.Sources.CreateSource:input_type -> webitel.cases.CreateSourceRequest
	5,  // 11: webitel.cases.Sources.UpdateSource:input_type -> webitel.cases.UpdateSourceRequest
	6,  // 12: webitel.cases.Sources.DeleteSource:input_type -> webitel.cases.DeleteSourceRequest
	8,  // 13: webitel.cases.Sources.LocateSource:input_type -> webitel.cases.LocateSourceRequest
	2,  // 14: webitel.cases.Sources.ListSources:output_type -> webitel.cases.SourceList
	1,  // 15: webitel.cases.Sources.CreateSource:output_type -> webitel.cases.Source
	1,  // 16: webitel.cases.Sources.UpdateSource:output_type -> webitel.cases.Source
	1,  // 17: webitel.cases.Sources.DeleteSource:output_type -> webitel.cases.Source
	9,  // 18: webitel.cases.Sources.LocateSource:output_type -> webitel.cases.LocateSourceResponse
	14, // [14:19] is the sub-list for method output_type
	9,  // [9:14] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_source_proto_init() }
func file_source_proto_init() {
	if File_source_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_source_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_source_proto_goTypes,
		DependencyIndexes: file_source_proto_depIdxs,
		EnumInfos:         file_source_proto_enumTypes,
		MessageInfos:      file_source_proto_msgTypes,
	}.Build()
	File_source_proto = out.File
	file_source_proto_rawDesc = nil
	file_source_proto_goTypes = nil
	file_source_proto_depIdxs = nil
}
