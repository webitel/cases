// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: case_comment.proto

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

// Represents a comment associated with a case.
type CaseComment struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Main identifier for read, update, and delete operations.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// entity tag representing id + ver
	Etag string `protobuf:"bytes,2,opt,name=etag,proto3" json:"etag,omitempty"`
	// Version number of the comment, used for concurrency control.
	Ver int32 `protobuf:"varint,3,opt,name=ver,proto3" json:"ver,omitempty"`
	// User who created the comment.
	CreatedBy *Lookup `protobuf:"bytes,4,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// Timestamp (in milliseconds) of when the comment was created.
	CreatedAt int64 `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// User who last updated the comment.
	UpdatedBy *Lookup `protobuf:"bytes,6,opt,name=updated_by,json=updatedBy,proto3" json:"updated_by,omitempty"`
	// Timestamp (in milliseconds) of the last update.
	UpdatedAt int64 `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// The content of the comment.
	Text string `protobuf:"bytes,8,opt,name=text,proto3" json:"text,omitempty"`
	// Indicates if the comment was edited; true if created_at < updated_at.
	Edited bool `protobuf:"varint,9,opt,name=edited,proto3" json:"edited,omitempty"`
	// Indicates if the comment can be edited by current user.
	CanEdit bool `protobuf:"varint,10,opt,name=can_edit,json=canEdit,proto3" json:"can_edit,omitempty"`
	// Contact-author of the comment.
	Author *Lookup `protobuf:"bytes,11,opt,name=author,proto3" json:"author,omitempty"`
	// Optional relation to the associated case.
	CaseId        int64 `protobuf:"varint,12,opt,name=case_id,json=caseId,proto3" json:"case_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CaseComment) Reset() {
	*x = CaseComment{}
	mi := &file_case_comment_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CaseComment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseComment) ProtoMessage() {}

func (x *CaseComment) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseComment.ProtoReflect.Descriptor instead.
func (*CaseComment) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{0}
}

func (x *CaseComment) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CaseComment) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *CaseComment) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *CaseComment) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *CaseComment) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *CaseComment) GetUpdatedBy() *Lookup {
	if x != nil {
		return x.UpdatedBy
	}
	return nil
}

func (x *CaseComment) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *CaseComment) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *CaseComment) GetEdited() bool {
	if x != nil {
		return x.Edited
	}
	return false
}

func (x *CaseComment) GetCanEdit() bool {
	if x != nil {
		return x.CanEdit
	}
	return false
}

func (x *CaseComment) GetAuthor() *Lookup {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *CaseComment) GetCaseId() int64 {
	if x != nil {
		return x.CaseId
	}
	return 0
}

// Contains a paginated list of comments.
type CaseCommentList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Current page number.
	Page int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Flag to indicate if more pages are available.
	Next bool `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	// List of comments on the current page.
	Items         []*CaseComment `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CaseCommentList) Reset() {
	*x = CaseCommentList{}
	mi := &file_case_comment_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CaseCommentList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseCommentList) ProtoMessage() {}

func (x *CaseCommentList) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseCommentList.ProtoReflect.Descriptor instead.
func (*CaseCommentList) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{1}
}

func (x *CaseCommentList) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *CaseCommentList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *CaseCommentList) GetItems() []*CaseComment {
	if x != nil {
		return x.Items
	}
	return nil
}

// Input structure for creating or updating a case comment.
type InputCaseComment struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Identifier for the comment.
	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	// Content of the comment.
	Text          string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InputCaseComment) Reset() {
	*x = InputCaseComment{}
	mi := &file_case_comment_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputCaseComment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputCaseComment) ProtoMessage() {}

func (x *InputCaseComment) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputCaseComment.ProtoReflect.Descriptor instead.
func (*InputCaseComment) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{2}
}

func (x *InputCaseComment) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *InputCaseComment) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

// Request to locate a comment based on its etag.
type LocateCommentRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Identifier of the comment to retrieve.
	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	// Specific fields to return for the comment.
	Fields        []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocateCommentRequest) Reset() {
	*x = LocateCommentRequest{}
	mi := &file_case_comment_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocateCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocateCommentRequest) ProtoMessage() {}

func (x *LocateCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocateCommentRequest.ProtoReflect.Descriptor instead.
func (*LocateCommentRequest) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{3}
}

func (x *LocateCommentRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *LocateCommentRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Request to update a comment with specified data.
type UpdateCommentRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// JSON fields specified in front-end request (patch).
	XJsonMask []string `protobuf:"bytes,1,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	// Fields to include in the response.
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Data for the comment to update.
	Input         *InputCaseComment `protobuf:"bytes,4,opt,name=input,proto3" json:"input,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCommentRequest) Reset() {
	*x = UpdateCommentRequest{}
	mi := &file_case_comment_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCommentRequest) ProtoMessage() {}

func (x *UpdateCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCommentRequest.ProtoReflect.Descriptor instead.
func (*UpdateCommentRequest) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateCommentRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

func (x *UpdateCommentRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *UpdateCommentRequest) GetInput() *InputCaseComment {
	if x != nil {
		return x.Input
	}
	return nil
}

// Request to delete a comment based on its etag.
type DeleteCommentRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Identifier of the comment to delete.
	Etag string `protobuf:"bytes,1,opt,name=etag,proto3" json:"etag,omitempty"`
	// Fields to return after deletion.
	Fields        []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteCommentRequest) Reset() {
	*x = DeleteCommentRequest{}
	mi := &file_case_comment_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteCommentRequest) ProtoMessage() {}

func (x *DeleteCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteCommentRequest.ProtoReflect.Descriptor instead.
func (*DeleteCommentRequest) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteCommentRequest) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *DeleteCommentRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Request to list comments for a specific case.
type ListCommentsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Page number for pagination.
	Page int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Number of comments per page.
	Size int32 `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	// Query string for search.
	Q string `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	// Array of requested id.
	Ids []string `protobuf:"bytes,4,rep,name=ids,proto3" json:"ids,omitempty"`
	// Sorting order.
	Sort string `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	// Fields to return for each comment.
	Fields []string `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	// Etag or ID of the case for which comments are requested.
	CaseEtag      string `protobuf:"bytes,9,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCommentsRequest) Reset() {
	*x = ListCommentsRequest{}
	mi := &file_case_comment_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCommentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommentsRequest) ProtoMessage() {}

func (x *ListCommentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommentsRequest.ProtoReflect.Descriptor instead.
func (*ListCommentsRequest) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{6}
}

func (x *ListCommentsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCommentsRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListCommentsRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListCommentsRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ListCommentsRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *ListCommentsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListCommentsRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

// Request to publish comment into a case.
type PublishCommentRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// JSON fields specified in the front-end request.
	XJsonMask []string `protobuf:"bytes,1,rep,name=x_json_mask,json=xJsonMask,proto3" json:"x_json_mask,omitempty"`
	// Result fields to include in the response.
	Fields []string `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
	// Etag or ID of the case to which comments belong.
	CaseEtag string `protobuf:"bytes,3,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"`
	// Comment to publish.
	Input         *InputCaseComment `protobuf:"bytes,4,opt,name=input,proto3" json:"input,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PublishCommentRequest) Reset() {
	*x = PublishCommentRequest{}
	mi := &file_case_comment_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PublishCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PublishCommentRequest) ProtoMessage() {}

func (x *PublishCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_comment_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PublishCommentRequest.ProtoReflect.Descriptor instead.
func (*PublishCommentRequest) Descriptor() ([]byte, []int) {
	return file_case_comment_proto_rawDescGZIP(), []int{7}
}

func (x *PublishCommentRequest) GetXJsonMask() []string {
	if x != nil {
		return x.XJsonMask
	}
	return nil
}

func (x *PublishCommentRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *PublishCommentRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

func (x *PublishCommentRequest) GetInput() *InputCaseComment {
	if x != nil {
		return x.Input
	}
	return nil
}

var File_case_comment_proto protoreflect.FileDescriptor

var file_case_comment_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x1a, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70,
	0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xea, 0x02, 0x0a, 0x0b, 0x43, 0x61,
	0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x74, 0x61,
	0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61, 0x67, 0x12, 0x10, 0x0a,
	0x03, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x76, 0x65, 0x72, 0x12,
	0x2e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2e,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f,
	0x6b, 0x75, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1d,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x61, 0x6e,
	0x5f, 0x65, 0x64, 0x69, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x63, 0x61, 0x6e,
	0x45, 0x64, 0x69, 0x74, 0x12, 0x27, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x0b,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c,
	0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x17, 0x0a,
	0x07, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x63, 0x61, 0x73, 0x65, 0x49, 0x64, 0x22, 0x6b, 0x0a, 0x0f, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78,
	0x74, 0x12, 0x30, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x22, 0x5a, 0x0a, 0x10, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43, 0x61, 0x73, 0x65,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x3a,
	0x1e, 0x92, 0x41, 0x1b, 0x32, 0x19, 0x7b, 0x22, 0x74, 0x65, 0x78, 0x74, 0x22, 0x3a, 0x22, 0x4d,
	0x79, 0x20, 0x6e, 0x65, 0x77, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x7d, 0x22,
	0x42, 0x0a, 0x14, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x22, 0xa0, 0x01, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0b,
	0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02, 0x5e, 0x24, 0xfa, 0xd2, 0xe4,
	0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x52, 0x09, 0x78, 0x4a,
	0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12,
	0x35, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x22, 0x42, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74,
	0x61, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0xa6, 0x01, 0x0a, 0x13, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f,
	0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65,
	0x74, 0x61, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x73, 0x65, 0x45,
	0x74, 0x61, 0x67, 0x22, 0xbe, 0x01, 0x0a, 0x15, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x43,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a,
	0x0b, 0x78, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x42, 0x19, 0x92, 0x41, 0x07, 0x40, 0x01, 0x8a, 0x01, 0x02, 0x5e, 0x24, 0xfa, 0xd2,
	0xe4, 0x93, 0x02, 0x09, 0x12, 0x07, 0x50, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x52, 0x09, 0x78,
	0x4a, 0x73, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
	0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x73, 0x65, 0x45, 0x74, 0x61, 0x67, 0x12, 0x35, 0x0a,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x32, 0xa7, 0x07, 0x0a, 0x0c, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0xa0, 0x01, 0x0a, 0x0d, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65,
	0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73,
	0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x4e, 0x92, 0x41, 0x29, 0x12, 0x27, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20, 0x61, 0x20, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66,
	0x69, 0x63, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x62, 0x79, 0x20, 0x69, 0x74,
	0x73, 0x20, 0x65, 0x74, 0x61, 0x67, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18,
	0x12, 0x16, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x2f, 0x7b, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12, 0xd3, 0x01, 0x0a, 0x0d, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1a, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e,
	0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x80, 0x01, 0x92, 0x41,
	0x27, 0x12, 0x25, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x62, 0x79, 0x20,
	0x69, 0x74, 0x73, 0x20, 0x65, 0x74, 0x61, 0x67, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x4c, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5a, 0x25, 0x3a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x32, 0x1c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x65, 0x74, 0x61, 0x67, 0x7d,
	0x1a, 0x1c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x2f, 0x7b, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12, 0x9e,
	0x01, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x23, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x22, 0x4c, 0x92, 0x41, 0x27, 0x12, 0x25, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x20, 0x61,
	0x20, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x20, 0x62, 0x79, 0x20, 0x69, 0x74, 0x73, 0x20, 0x65, 0x74, 0x61, 0x67, 0x90, 0xb5, 0x18,
	0x03, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x2a, 0x16, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x7b, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12,
	0xbb, 0x01, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73,
	0x12, 0x22, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x67, 0x92, 0x41, 0x3d, 0x12, 0x3b, 0x52, 0x65, 0x74, 0x72, 0x69,
	0x65, 0x76, 0x65, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x63, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x20, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x65,
	0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x61, 0x20, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x63, 0x20, 0x63, 0x61, 0x73, 0x65, 0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1d,
	0x12, 0x1b, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65,
	0x74, 0x61, 0x67, 0x7d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0xab, 0x01,
	0x0a, 0x0e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x24, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x22, 0x57, 0x92, 0x41, 0x26, 0x12, 0x24, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x69, 0x6e, 0x74, 0x6f, 0x20, 0x61, 0x20,
	0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x20, 0x63, 0x61, 0x73, 0x65, 0x90, 0xb5, 0x18,
	0x00, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x3a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x22, 0x1b,
	0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61,
	0x67, 0x7d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x11, 0x8a, 0xb5, 0x18,
	0x0d, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x42, 0xa4,
	0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x42, 0x10, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3b, 0x63, 0x61, 0x73,
	0x65, 0x73, 0xa2, 0x02, 0x03, 0x57, 0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x73, 0xca, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0xe2, 0x02, 0x19, 0x57, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73, 0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a, 0x3a,
	0x43, 0x61, 0x73, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_case_comment_proto_rawDescOnce sync.Once
	file_case_comment_proto_rawDescData []byte
)

func file_case_comment_proto_rawDescGZIP() []byte {
	file_case_comment_proto_rawDescOnce.Do(func() {
		file_case_comment_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_case_comment_proto_rawDesc), len(file_case_comment_proto_rawDesc)))
	})
	return file_case_comment_proto_rawDescData
}

var file_case_comment_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_case_comment_proto_goTypes = []any{
	(*CaseComment)(nil),           // 0: webitel.cases.CaseComment
	(*CaseCommentList)(nil),       // 1: webitel.cases.CaseCommentList
	(*InputCaseComment)(nil),      // 2: webitel.cases.InputCaseComment
	(*LocateCommentRequest)(nil),  // 3: webitel.cases.LocateCommentRequest
	(*UpdateCommentRequest)(nil),  // 4: webitel.cases.UpdateCommentRequest
	(*DeleteCommentRequest)(nil),  // 5: webitel.cases.DeleteCommentRequest
	(*ListCommentsRequest)(nil),   // 6: webitel.cases.ListCommentsRequest
	(*PublishCommentRequest)(nil), // 7: webitel.cases.PublishCommentRequest
	(*Lookup)(nil),                // 8: general.Lookup
}
var file_case_comment_proto_depIdxs = []int32{
	8,  // 0: webitel.cases.CaseComment.created_by:type_name -> general.Lookup
	8,  // 1: webitel.cases.CaseComment.updated_by:type_name -> general.Lookup
	8,  // 2: webitel.cases.CaseComment.author:type_name -> general.Lookup
	0,  // 3: webitel.cases.CaseCommentList.items:type_name -> webitel.cases.CaseComment
	2,  // 4: webitel.cases.UpdateCommentRequest.input:type_name -> webitel.cases.InputCaseComment
	2,  // 5: webitel.cases.PublishCommentRequest.input:type_name -> webitel.cases.InputCaseComment
	3,  // 6: webitel.cases.CaseComments.LocateComment:input_type -> webitel.cases.LocateCommentRequest
	4,  // 7: webitel.cases.CaseComments.UpdateComment:input_type -> webitel.cases.UpdateCommentRequest
	5,  // 8: webitel.cases.CaseComments.DeleteComment:input_type -> webitel.cases.DeleteCommentRequest
	6,  // 9: webitel.cases.CaseComments.ListComments:input_type -> webitel.cases.ListCommentsRequest
	7,  // 10: webitel.cases.CaseComments.PublishComment:input_type -> webitel.cases.PublishCommentRequest
	0,  // 11: webitel.cases.CaseComments.LocateComment:output_type -> webitel.cases.CaseComment
	0,  // 12: webitel.cases.CaseComments.UpdateComment:output_type -> webitel.cases.CaseComment
	0,  // 13: webitel.cases.CaseComments.DeleteComment:output_type -> webitel.cases.CaseComment
	1,  // 14: webitel.cases.CaseComments.ListComments:output_type -> webitel.cases.CaseCommentList
	0,  // 15: webitel.cases.CaseComments.PublishComment:output_type -> webitel.cases.CaseComment
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_case_comment_proto_init() }
func file_case_comment_proto_init() {
	if File_case_comment_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_case_comment_proto_rawDesc), len(file_case_comment_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_case_comment_proto_goTypes,
		DependencyIndexes: file_case_comment_proto_depIdxs,
		MessageInfos:      file_case_comment_proto_msgTypes,
	}.Build()
	File_case_comment_proto = out.File
	file_case_comment_proto_goTypes = nil
	file_case_comment_proto_depIdxs = nil
}
