// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: cases/case_file.proto

package cases

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type CaseFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Author *Lookup `protobuf:"bytes,6,opt,name=author,proto3" json:"author,omitempty"` // contact of the creator
	File   *File   `protobuf:"bytes,7,opt,name=file,proto3" json:"file,omitempty"`
}

func (x *CaseFile) Reset() {
	*x = CaseFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_case_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CaseFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseFile) ProtoMessage() {}

func (x *CaseFile) ProtoReflect() protoreflect.Message {
	mi := &file_cases_case_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseFile.ProtoReflect.Descriptor instead.
func (*CaseFile) Descriptor() ([]byte, []int) {
	return file_cases_case_file_proto_rawDescGZIP(), []int{0}
}

func (x *CaseFile) GetAuthor() *Lookup {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *CaseFile) GetFile() *File {
	if x != nil {
		return x.File
	}
	return nil
}

type CaseFileList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int64       `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Next  bool        `protobuf:"varint,2,opt,name=next,proto3" json:"next,omitempty"`
	Items []*CaseFile `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *CaseFileList) Reset() {
	*x = CaseFileList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_case_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CaseFileList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseFileList) ProtoMessage() {}

func (x *CaseFileList) ProtoReflect() protoreflect.Message {
	mi := &file_cases_case_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_cases_case_file_proto_rawDescGZIP(), []int{1}
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

func (x *CaseFileList) GetItems() []*CaseFile {
	if x != nil {
		return x.Items
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` // storage.file id
	CreatedBy *Lookup `protobuf:"bytes,2,opt,name=created_by,json=createdBy,proto3" json:"created_by,omitempty"`
	CreatedAt int64   `protobuf:"varint,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"` // unixmilli
	Size      int64   `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Mime      string  `protobuf:"bytes,5,opt,name=mime,proto3" json:"mime,omitempty"` // MIME type
	Name      string  `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_case_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_cases_case_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_cases_case_file_proto_rawDescGZIP(), []int{2}
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

type ListFilesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page   int32    `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size   int32    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q      string   `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"` // covered by filters?
	Qin    []string `protobuf:"bytes,4,rep,name=qin,proto3" json:"qin,omitempty"`
	Sort   string   `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields []string `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	Filter string   `protobuf:"bytes,7,opt,name=filter,proto3" json:"filter,omitempty"`
	// simple filter query language [SFQL]
	// ( -> start of the filter or filter node
	// ) -> end of the filter or the filter node
	// last ( -> always start of the filter
	// [ -> start of the operator
	// ] -> end of the operator
	// operands - simple strings, first string should always be a valid field name, second should
	// operators:
	// [and] [or] -- only applied to the filter node, after them always should be a ( as filter or another filter node
	// [eq], [gte], [gt], [lte], [lt], [regex], [neq] -- only applied to the filter, value after them and to the ) symbol considered as string
	// ...?size=10&page=1&filter=((file[eq]fghj)[and](name[eq]yehor))[or])
	//
	//	oneof filters {
	//	  FilterNode node = 7;
	//	  Filter filter = 8;
	//	}
	CaseEtag string `protobuf:"bytes,9,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"`
}

func (x *ListFilesRequest) Reset() {
	*x = ListFilesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cases_case_file_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFilesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFilesRequest) ProtoMessage() {}

func (x *ListFilesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cases_case_file_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_cases_case_file_proto_rawDescGZIP(), []int{3}
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

func (x *ListFilesRequest) GetQin() []string {
	if x != nil {
		return x.Qin
	}
	return nil
}

func (x *ListFilesRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *ListFilesRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListFilesRequest) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *ListFilesRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

var File_cases_case_file_proto protoreflect.FileDescriptor

var file_cases_case_file_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x1a, 0x12, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x6c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x62, 0x0a,
	0x08, 0x43, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70,
	0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x27, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x04, 0x66, 0x69, 0x6c,
	0x65, 0x22, 0x65, 0x0a, 0x0c, 0x43, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x12, 0x2d, 0x0a, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xa7, 0x01, 0x0a, 0x04, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x34, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x69,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x69, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0xbb, 0x01, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73,
	0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12,
	0x0c, 0x0a, 0x01, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x10, 0x0a,
	0x03, 0x71, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x71, 0x69, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73,
	0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x66,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x61, 0x73, 0x65, 0x45, 0x74, 0x61, 0x67,
	0x32, 0x78, 0x0a, 0x09, 0x43, 0x61, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x6b, 0x0a,
	0x09, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1f, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65,
	0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x20, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1a,
	0x12, 0x18, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65,
	0x74, 0x61, 0x67, 0x7d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x3b, 0x63, 0x61, 0x73, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cases_case_file_proto_rawDescOnce sync.Once
	file_cases_case_file_proto_rawDescData = file_cases_case_file_proto_rawDesc
)

func file_cases_case_file_proto_rawDescGZIP() []byte {
	file_cases_case_file_proto_rawDescOnce.Do(func() {
		file_cases_case_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_cases_case_file_proto_rawDescData)
	})
	return file_cases_case_file_proto_rawDescData
}

var file_cases_case_file_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_cases_case_file_proto_goTypes = []any{
	(*CaseFile)(nil),         // 0: webitel.cases.CaseFile
	(*CaseFileList)(nil),     // 1: webitel.cases.CaseFileList
	(*File)(nil),             // 2: webitel.cases.File
	(*ListFilesRequest)(nil), // 3: webitel.cases.ListFilesRequest
	(*Lookup)(nil),           // 4: webitel.cases.Lookup
}
var file_cases_case_file_proto_depIdxs = []int32{
	4, // 0: webitel.cases.CaseFile.author:type_name -> webitel.cases.Lookup
	2, // 1: webitel.cases.CaseFile.file:type_name -> webitel.cases.File
	0, // 2: webitel.cases.CaseFileList.items:type_name -> webitel.cases.CaseFile
	4, // 3: webitel.cases.File.created_by:type_name -> webitel.cases.Lookup
	3, // 4: webitel.cases.CaseFiles.ListFiles:input_type -> webitel.cases.ListFilesRequest
	1, // 5: webitel.cases.CaseFiles.ListFiles:output_type -> webitel.cases.CaseFileList
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_cases_case_file_proto_init() }
func file_cases_case_file_proto_init() {
	if File_cases_case_file_proto != nil {
		return
	}
	file_cases_lookup_proto_init()
	file_cases_filters_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_cases_case_file_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CaseFile); i {
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
		file_cases_case_file_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CaseFileList); i {
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
		file_cases_case_file_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*File); i {
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
		file_cases_case_file_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ListFilesRequest); i {
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
			RawDescriptor: file_cases_case_file_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cases_case_file_proto_goTypes,
		DependencyIndexes: file_cases_case_file_proto_depIdxs,
		MessageInfos:      file_cases_case_file_proto_msgTypes,
	}.Build()
	File_cases_case_file_proto = out.File
	file_cases_case_file_proto_rawDesc = nil
	file_cases_case_file_proto_goTypes = nil
	file_cases_case_file_proto_depIdxs = nil
}
