// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: case_file.proto

package cases

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "github.com/webitel/webitel-go-kit/cmd/protoc-gen-go-webitel/gen/go/proto/webitel"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Contains a list of case files with pagination.
type CaseFileList struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Current page number.
	Page int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// Indicator if there is a next page.
	Next bool `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	// List of case files.
	Items         []*File `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CaseFileList) Reset() {
	*x = CaseFileList{}
	mi := &file_case_file_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CaseFileList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseFileList) ProtoMessage() {}

func (x *CaseFileList) ProtoReflect() protoreflect.Message {
	mi := &file_case_file_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseFileList.ProtoReflect.Descriptor instead.
func (*CaseFileList) Descriptor() ([]byte, []int) {
	return file_case_file_proto_rawDescGZIP(), []int{0}
}

func (x *CaseFileList) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *CaseFileList) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *CaseFileList) GetItems() []*File {
	if x != nil {
		return x.Items
	}
	return nil
}

// Metadata for a file associated with a case.
type File struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Storage file ID.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Creator of the file.
	CreatedBy *Lookup `protobuf:"bytes,2,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	// Creation timestamp in Unix milliseconds.
	CreatedAt int64 `protobuf:"varint,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// File size in bytes.
	Size int64 `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	// MIME type of the file.
	Mime string `protobuf:"bytes,5,opt,name=mime,proto3" json:"mime,omitempty"`
	// File name.
	Name          string  `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	Author        *Lookup `protobuf:"bytes,7,opt,name=author,proto3" json:"author,omitempty"`
	Url           string  `protobuf:"bytes,8,opt,name=url,proto3" json:"url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *File) Reset() {
	*x = File{}
	mi := &file_case_file_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_case_file_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_case_file_proto_rawDescGZIP(), []int{1}
}

func (x *File) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *File) GetCreatedBy() *Lookup {
	if x != nil {
		return x.CreatedBy
	}
	return nil
}

func (x *File) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *File) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *File) GetMime() string {
	if x != nil {
		return x.Mime
	}
	return ""
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetAuthor() *Lookup {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *File) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// Request message for listing files associated with a case.
type ListFilesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// ID of the case to fetch files for (required).
	CaseEtag string `protobuf:"bytes,1,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"`
	// The page number to retrieve.
	Page int32 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	// Number of items per page.
	Size int32 `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	// Search term.
	Q string `protobuf:"bytes,4,opt,name=q,proto3" json:"q,omitempty"`
	// Fields to include in the response.
	Fields []string `protobuf:"bytes,5,rep,name=fields,proto3" json:"fields,omitempty"`
	// Array of requested id.
	Ids           []string `protobuf:"bytes,6,rep,name=ids,proto3" json:"ids,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListFilesRequest) Reset() {
	*x = ListFilesRequest{}
	mi := &file_case_file_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListFilesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFilesRequest) ProtoMessage() {}

func (x *ListFilesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_file_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFilesRequest.ProtoReflect.Descriptor instead.
func (*ListFilesRequest) Descriptor() ([]byte, []int) {
	return file_case_file_proto_rawDescGZIP(), []int{2}
}

func (x *ListFilesRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

func (x *ListFilesRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListFilesRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListFilesRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListFilesRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListFilesRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

var File_case_file_proto protoreflect.FileDescriptor

var file_case_file_proto_rawDesc = string([]byte{
	0x0a, 0x0f, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x1a, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x61, 0x0a, 0x0c, 0x43, 0x61, 0x73,
	0x65, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78,
	0x74, 0x12, 0x29, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xdc, 0x01, 0x0a,
	0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x62, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x69, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x27, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75,
	0x70, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0xa2, 0x01, 0x0a, 0x10,
	0x4c, 0x69, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x73, 0x65, 0x45, 0x74, 0x61, 0x67, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x01, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69,
	0x64, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x3a, 0x11, 0x92,
	0x41, 0x0e, 0x0a, 0x0c, 0xd2, 0x01, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67,
	0x32, 0xbc, 0x01, 0x0a, 0x09, 0x43, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0xa3,
	0x01, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61,
	0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x58, 0x92, 0x41, 0x31, 0x12,
	0x2f, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x20, 0x61, 0x20, 0x6c, 0x69, 0x73, 0x74,
	0x20, 0x6f, 0x66, 0x20, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x20, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69,
	0x61, 0x74, 0x65, 0x64, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x61, 0x20, 0x63, 0x61, 0x73, 0x65,
	0x90, 0xb5, 0x18, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a, 0x12, 0x18, 0x2f, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x2f, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x1a, 0x09, 0x8a, 0xb5, 0x18, 0x05, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42,
	0xa1, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0x0d, 0x43, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x50,
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
	file_case_file_proto_rawDescOnce sync.Once
	file_case_file_proto_rawDescData []byte
)

func file_case_file_proto_rawDescGZIP() []byte {
	file_case_file_proto_rawDescOnce.Do(func() {
		file_case_file_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_case_file_proto_rawDesc), len(file_case_file_proto_rawDesc)))
	})
	return file_case_file_proto_rawDescData
}

var file_case_file_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_case_file_proto_goTypes = []any{
	(*CaseFileList)(nil),     // 0: webitel.cases.CaseFileList
	(*File)(nil),             // 1: webitel.cases.File
	(*ListFilesRequest)(nil), // 2: webitel.cases.ListFilesRequest
	(*Lookup)(nil),           // 3: general.Lookup
}
var file_case_file_proto_depIdxs = []int32{
	1, // 0: webitel.cases.CaseFileList.items:type_name -> webitel.cases.File
	3, // 1: webitel.cases.File.created_by:type_name -> general.Lookup
	3, // 2: webitel.cases.File.author:type_name -> general.Lookup
	2, // 3: webitel.cases.CaseFiles.ListFiles:input_type -> webitel.cases.ListFilesRequest
	0, // 4: webitel.cases.CaseFiles.ListFiles:output_type -> webitel.cases.CaseFileList
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_case_file_proto_init() }
func file_case_file_proto_init() {
	if File_case_file_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_case_file_proto_rawDesc), len(file_case_file_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_case_file_proto_goTypes,
		DependencyIndexes: file_case_file_proto_depIdxs,
		MessageInfos:      file_case_file_proto_msgTypes,
	}.Build()
	File_case_file_proto = out.File
	file_case_file_proto_goTypes = nil
	file_case_file_proto_depIdxs = nil
}
