// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: case_link.proto

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

type CaseLink struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Ver       int32   `protobuf:"varint,2,opt,name=ver,proto3" json:"ver,omitempty"`
	Etag      string  `protobuf:"bytes,3,opt,name=etag,proto3" json:"etag,omitempty"` // main field required for read, update and delete
	CreatedBy *Lookup `protobuf:"bytes,4,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	CreatedAt int64   `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"` // unixmilli
	UpdatedBy *Lookup `protobuf:"bytes,6,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	UpdatedAt int64   `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Author    *Lookup `protobuf:"bytes,8,opt,name=author,proto3" json:"author,omitempty"` // contact-author calculated on-flight (optional)
	Name      string  `protobuf:"bytes,9,opt,name=name,proto3" json:"name,omitempty"`     // link name (optional)
	Url       string  `protobuf:"bytes,11,opt,name=url,proto3" json:"url,omitempty"`      // URL
}

func (x *CaseLink) Reset() {
	*x = CaseLink{}
	mi := &file_case_link_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CaseLink) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseLink) ProtoMessage() {}

func (x *CaseLink) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseLink.ProtoReflect.Descriptor instead.
func (*CaseLink) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{0}
}

func (x *CaseLink) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CaseLink) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *CaseLink) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *CaseLink) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *CaseLink) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *CaseLink) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *CaseLink) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *CaseLink) GetAuthor() *Lookup {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *CaseLink) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CaseLink) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type InputCaseLink struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	Url  string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Name string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *InputCaseLink) Reset() {
	*x = InputCaseLink{}
	mi := &file_case_link_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputCaseLink) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputCaseLink) ProtoMessage() {}

func (x *InputCaseLink) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputCaseLink.ProtoReflect.Descriptor instead.
func (*InputCaseLink) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{1}
}

func (x *InputCaseLink) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *InputCaseLink) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *InputCaseLink) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type CaseLinkList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int64       `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next  bool        `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items []*CaseLink `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *CaseLinkList) Reset() {
	*x = CaseLinkList{}
	mi := &file_case_link_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CaseLinkList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseLinkList) ProtoMessage() {}

func (x *CaseLinkList) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseLinkList.ProtoReflect.Descriptor instead.
func (*CaseLinkList) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{2}
}

func (x *CaseLinkList) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *CaseLinkList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *CaseLinkList) GetItems() []*CaseLink {
	if x != nil {
		return x.Items
	}
	return nil
}

type LocateLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Etag   string   `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"` // (id allowed)
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *LocateLinkRequest) Reset() {
	*x = LocateLinkRequest{}
	mi := &file_case_link_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateLinkRequest) ProtoMessage() {}

func (x *LocateLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateLinkRequest.ProtoReflect.Descriptor instead.
func (*LocateLinkRequest) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{3}
}

func (x *LocateLinkRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *LocateLinkRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type UpdateLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Etag      string         `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	Fields    []string       `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"` // on return
	XJsonMask []string       `protobuf:"bytes,3,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	Input     *InputCaseLink `protobuf:"bytes,4,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *UpdateLinkRequest) Reset() {
	*x = UpdateLinkRequest{}
	mi := &file_case_link_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateLinkRequest) ProtoMessage() {}

func (x *UpdateLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateLinkRequest.ProtoReflect.Descriptor instead.
func (*UpdateLinkRequest) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateLinkRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *UpdateLinkRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *UpdateLinkRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

func (x *UpdateLinkRequest) GetInput() *InputCaseLink {
	if x != nil {
		return x.Input
	}
	return nil
}

type DeleteLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"` //
}

func (x *DeleteLinkRequest) Reset() {
	*x = DeleteLinkRequest{}
	mi := &file_case_link_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteLinkRequest) ProtoMessage() {}

func (x *DeleteLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteLinkRequest.ProtoReflect.Descriptor instead.
func (*DeleteLinkRequest) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteLinkRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

type ListLinksRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page     int32    `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size     int32    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q        string   `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	Ids      []string `protobuf:"bytes,4,rep,name=ids,proto3" json:"ids,omitempty"`
	Sort     string   `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields   []string `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	CaseEtag string   `protobuf:"bytes,9,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"`
}

func (x *ListLinksRequest) Reset() {
	*x = ListLinksRequest{}
	mi := &file_case_link_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListLinksRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLinksRequest) ProtoMessage() {}

func (x *ListLinksRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLinksRequest.ProtoReflect.Descriptor instead.
func (*ListLinksRequest) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{6}
}

func (x *ListLinksRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListLinksRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListLinksRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListLinksRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ListLinksRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *ListLinksRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListLinksRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

type CreateLinkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields   []string       `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`                     // result fields needed on the front-end for each comment
	CaseEtag string         `protobuf:"bytes,3,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"` // new comment link
	Input    *InputCaseLink `protobuf:"bytes,4,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *CreateLinkRequest) Reset() {
	*x = CreateLinkRequest{}
	mi := &file_case_link_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateLinkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLinkRequest) ProtoMessage() {}

func (x *CreateLinkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_link_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLinkRequest.ProtoReflect.Descriptor instead.
func (*CreateLinkRequest) Descriptor() ([]byte, []int) {
	return file_case_link_proto_rawDescGZIP(), []int{7}
}

func (x *CreateLinkRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *CreateLinkRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

func (x *CreateLinkRequest) GetInput() *InputCaseLink {
	if x != nil {
		return x.Input
	}
	return nil
}

var File_case_link_proto protoreflect.FileDescriptor

var file_case_link_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x1a, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad, 0x02, 0x0a, 0x08, 0x43, 0x61, 0x73, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61, 0x67, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2e, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x27, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x49, 0x0a, 0x0d, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43,
	0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x65, 0x0a, 0x0c, 0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x2d, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e,
	0x6b, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x3f, 0x0a, 0x11, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x65, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61,
	0x67, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0xae, 0x01, 0x0a, 0x11, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65,
	0x74, 0x61, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x39, 0x0a, 0x0b, 0x78,
	0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02, 0x5e, 0x24, 0xfa, 0xd2, 0xe4, 0x93,
	0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x52, 0x09, 0x78, 0x4a, 0x73,
	0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x12, 0x32, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43, 0x61, 0x73, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x22, 0x27, 0x0a, 0x11, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65,
	0x74, 0x61, 0x67, 0x22, 0xa3, 0x01, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x10,
	0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x73, 0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x06,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09,
	0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x63, 0x61, 0x73, 0x65, 0x45, 0x74, 0x61, 0x67, 0x22, 0x7c, 0x0a, 0x11, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65,
	0x74, 0x61, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x73, 0x65, 0x45,
	0x74, 0x61, 0x67, 0x12, 0x32, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b,
	0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x32, 0xe2, 0x04, 0x0a, 0x09, 0x43, 0x61, 0x73, 0x65,
	0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x68, 0x0a, 0x0a, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x12, 0x20, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0x1f,
	0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12, 0x13, 0x2f, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x2f, 0x7b, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12,
	0x6d, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x20, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0x24, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x1a, 0x22, 0x18, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x6e,
	0x6b, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12, 0x99,
	0x01, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x20, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0x50, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x46, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5a, 0x22, 0x3a, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x32, 0x19, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x6e,
	0x6b, 0x73, 0x2f, 0x7b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x1a,
	0x19, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x2f, 0x7b, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x2e, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12, 0x68, 0x0a, 0x0a, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x20, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x4c,
	0x69, 0x6e, 0x6b, 0x22, 0x1f, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x2a,
	0x13, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x2f, 0x7b, 0x65,
	0x74, 0x61, 0x67, 0x7d, 0x12, 0x6b, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b,
	0x73, 0x12, 0x1f, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x12, 0x18, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f,
	0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x2f, 0x6c, 0x69, 0x6e, 0x6b,
	0x73, 0x1a, 0x09, 0x8a, 0xb5, 0x18, 0x05, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0xa1, 0x01, 0x0a,
	0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x42, 0x0d, 0x43, 0x61, 0x73, 0x65, 0x4c, 0x69, 0x6e, 0x6b, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3b, 0x63, 0x61, 0x73, 0x65, 0x73, 0xa2, 0x02, 0x03,
	0x57, 0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x43, 0x61,
	0x73, 0x65, 0x73, 0xca, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61,
	0x73, 0x65, 0x73, 0xe2, 0x02, 0x19, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61,
	0x73, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0e, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a, 0x3a, 0x43, 0x61, 0x73, 0x65, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_case_link_proto_rawDescOnce sync.Once
	file_case_link_proto_rawDescData = file_case_link_proto_rawDesc
)

func file_case_link_proto_rawDescGZIP() []byte {
	file_case_link_proto_rawDescOnce.Do(func() {
		file_case_link_proto_rawDescData = protoimpl.X.CompressGZIP(file_case_link_proto_rawDescData)
	})
	return file_case_link_proto_rawDescData
}

var file_case_link_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_case_link_proto_goTypes = []any{
	(*CaseLink)(nil),          // 0: webitel.cases.CaseLink
	(*InputCaseLink)(nil),     // 1: webitel.cases.InputCaseLink
	(*CaseLinkList)(nil),      // 2: webitel.cases.CaseLinkList
	(*LocateLinkRequest)(nil), // 3: webitel.cases.LocateLinkRequest
	(*UpdateLinkRequest)(nil), // 4: webitel.cases.UpdateLinkRequest
	(*DeleteLinkRequest)(nil), // 5: webitel.cases.DeleteLinkRequest
	(*ListLinksRequest)(nil),  // 6: webitel.cases.ListLinksRequest
	(*CreateLinkRequest)(nil), // 7: webitel.cases.CreateLinkRequest
	(*Lookup)(nil),            // 8: general.Lookup
}
var file_case_link_proto_depIdxs = []int32{
	8,  // 0: webitel.cases.CaseLink.created_by:type_name -> general.Lookup
	8,  // 1: webitel.cases.CaseLink.updated_by:type_name -> general.Lookup
	8,  // 2: webitel.cases.CaseLink.author:type_name -> general.Lookup
	0,  // 3: webitel.cases.CaseLinkList.items:type_name -> webitel.cases.CaseLink
	1,  // 4: webitel.cases.UpdateLinkRequest.input:type_name -> webitel.cases.InputCaseLink
	1,  // 5: webitel.cases.CreateLinkRequest.input:type_name -> webitel.cases.InputCaseLink
	3,  // 6: webitel.cases.CaseLinks.LocateLink:input_type -> webitel.cases.LocateLinkRequest
	7,  // 7: webitel.cases.CaseLinks.CreateLink:input_type -> webitel.cases.CreateLinkRequest
	4,  // 8: webitel.cases.CaseLinks.UpdateLink:input_type -> webitel.cases.UpdateLinkRequest
	5,  // 9: webitel.cases.CaseLinks.DeleteLink:input_type -> webitel.cases.DeleteLinkRequest
	6,  // 10: webitel.cases.CaseLinks.ListLinks:input_type -> webitel.cases.ListLinksRequest
	0,  // 11: webitel.cases.CaseLinks.LocateLink:output_type -> webitel.cases.CaseLink
	0,  // 12: webitel.cases.CaseLinks.CreateLink:output_type -> webitel.cases.CaseLink
	0,  // 13: webitel.cases.CaseLinks.UpdateLink:output_type -> webitel.cases.CaseLink
	0,  // 14: webitel.cases.CaseLinks.DeleteLink:output_type -> webitel.cases.CaseLink
	2,  // 15: webitel.cases.CaseLinks.ListLinks:output_type -> webitel.cases.CaseLinkList
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_case_link_proto_init() }
func file_case_link_proto_init() {
	if File_case_link_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_case_link_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_case_link_proto_goTypes,
		DependencyIndexes: file_case_link_proto_depIdxs,
		MessageInfos:      file_case_link_proto_msgTypes,
	}.Build()
	File_case_link_proto = out.File
	file_case_link_proto_rawDesc = nil
	file_case_link_proto_goTypes = nil
	file_case_link_proto_depIdxs = nil
}
