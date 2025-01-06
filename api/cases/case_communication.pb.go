// Specifies the syntax version of the protocol buffer.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: case_communication.proto

// Defines the package for the generated Go files.

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

// Enum to define types of case communications.
type CaseCommunicationsTypes int32

const (
	CaseCommunicationsTypes_NO_TYPE             CaseCommunicationsTypes = 0 // Default value, no type specified.
	CaseCommunicationsTypes_COMMUNICATION_CHAT  CaseCommunicationsTypes = 1 // Communication type: Chat.
	CaseCommunicationsTypes_COMMUNICATION_CALL  CaseCommunicationsTypes = 2 // Communication type: Call.
	CaseCommunicationsTypes_COMMUNICATION_EMAIL CaseCommunicationsTypes = 3 // Communication type: Email.
)

// Enum value maps for CaseCommunicationsTypes.
var (
	CaseCommunicationsTypes_name = map[int32]string{
		0: "NO_TYPE",
		1: "COMMUNICATION_CHAT",
		2: "COMMUNICATION_CALL",
		3: "COMMUNICATION_EMAIL",
	}
	CaseCommunicationsTypes_value = map[string]int32{
		"NO_TYPE":             0,
		"COMMUNICATION_CHAT":  1,
		"COMMUNICATION_CALL":  2,
		"COMMUNICATION_EMAIL": 3,
	}
)

func (x CaseCommunicationsTypes) Enum() *CaseCommunicationsTypes {
	p := new(CaseCommunicationsTypes)
	*p = x
	return p
}

func (x CaseCommunicationsTypes) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CaseCommunicationsTypes) Descriptor() protoreflect.EnumDescriptor {
	return file_case_communication_proto_enumTypes[0].Descriptor()
}

func (CaseCommunicationsTypes) Type() protoreflect.EnumType {
	return &file_case_communication_proto_enumTypes[0]
}

func (x CaseCommunicationsTypes) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CaseCommunicationsTypes.Descriptor instead.
func (CaseCommunicationsTypes) EnumDescriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{0}
}

// Represents a single case communication.
type CaseCommunication struct {
	state             protoimpl.MessageState  `protogen:"open.v1"`
	Id                string                  `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`                                                                                                    // Database ID of the communication.
	Ver               int32                   `protobuf:"varint,3,opt,name=ver,proto3" json:"ver,omitempty"`                                                                                                 // Version of the communication record.
	CommunicationType CaseCommunicationsTypes `protobuf:"varint,4,opt,name=communication_type,json=communicationType,proto3,enum=webitel.cases.CaseCommunicationsTypes" json:"communication_type,omitempty"` // Type of the communication (e.g., Chat, Call).
	CommunicationId   string                  `protobuf:"bytes,5,opt,name=communication_id,json=communicationId,proto3" json:"communication_id,omitempty"`                                                   // External communication ID.
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *CaseCommunication) Reset() {
	*x = CaseCommunication{}
	mi := &file_case_communication_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CaseCommunication) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaseCommunication) ProtoMessage() {}

func (x *CaseCommunication) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaseCommunication.ProtoReflect.Descriptor instead.
func (*CaseCommunication) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{0}
}

func (x *CaseCommunication) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CaseCommunication) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *CaseCommunication) GetCommunicationType() CaseCommunicationsTypes {
	if x != nil {
		return x.CommunicationType
	}
	return CaseCommunicationsTypes_NO_TYPE
}

func (x *CaseCommunication) GetCommunicationId() string {
	if x != nil {
		return x.CommunicationId
	}
	return ""
}

// Represents input data for creating or linking a communication.
type InputCaseCommunication struct {
	state             protoimpl.MessageState  `protogen:"open.v1"`
	CommunicationType CaseCommunicationsTypes `protobuf:"varint,2,opt,name=communication_type,json=communicationType,proto3,enum=webitel.cases.CaseCommunicationsTypes" json:"communication_type,omitempty"` // Type of the communication.
	CommunicationId   string                  `protobuf:"bytes,3,opt,name=communication_id,json=communicationId,proto3" json:"communication_id,omitempty"`                                                   // External communication ID.
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *InputCaseCommunication) Reset() {
	*x = InputCaseCommunication{}
	mi := &file_case_communication_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InputCaseCommunication) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputCaseCommunication) ProtoMessage() {}

func (x *InputCaseCommunication) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputCaseCommunication.ProtoReflect.Descriptor instead.
func (*InputCaseCommunication) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{1}
}

func (x *InputCaseCommunication) GetCommunicationType() CaseCommunicationsTypes {
	if x != nil {
		return x.CommunicationType
	}
	return CaseCommunicationsTypes_NO_TYPE
}

func (x *InputCaseCommunication) GetCommunicationId() string {
	if x != nil {
		return x.CommunicationId
	}
	return ""
}

// Request message for linking communications to a case.
type LinkCommunicationRequest struct {
	state         protoimpl.MessageState    `protogen:"open.v1"`
	CaseId        string                    `protobuf:"bytes,1,opt,name=case_id,json=caseId,proto3" json:"case_id,omitempty"` // Case identifier.
	Fields        []string                  `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`               // List of fields to include in the response.
	Input         []*InputCaseCommunication `protobuf:"bytes,3,rep,name=input,proto3" json:"input,omitempty"`                 // Input data for the communications to link.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LinkCommunicationRequest) Reset() {
	*x = LinkCommunicationRequest{}
	mi := &file_case_communication_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LinkCommunicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LinkCommunicationRequest) ProtoMessage() {}

func (x *LinkCommunicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LinkCommunicationRequest.ProtoReflect.Descriptor instead.
func (*LinkCommunicationRequest) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{2}
}

func (x *LinkCommunicationRequest) GetCaseId() string {
	if x != nil {
		return x.CaseId
	}
	return ""
}

func (x *LinkCommunicationRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *LinkCommunicationRequest) GetInput() []*InputCaseCommunication {
	if x != nil {
		return x.Input
	}
	return nil
}

// Response message after linking communications to a case.
type LinkCommunicationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*CaseCommunication   `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"` // List of linked communications.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LinkCommunicationResponse) Reset() {
	*x = LinkCommunicationResponse{}
	mi := &file_case_communication_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LinkCommunicationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LinkCommunicationResponse) ProtoMessage() {}

func (x *LinkCommunicationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LinkCommunicationResponse.ProtoReflect.Descriptor instead.
func (*LinkCommunicationResponse) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{3}
}

func (x *LinkCommunicationResponse) GetData() []*CaseCommunication {
	if x != nil {
		return x.Data
	}
	return nil
}

// Request message for unlinking communications from a case.
type UnlinkCommunicationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`         // Communication identifier.
	Fields        []string               `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"` // List of fields to include in the response.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnlinkCommunicationRequest) Reset() {
	*x = UnlinkCommunicationRequest{}
	mi := &file_case_communication_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnlinkCommunicationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlinkCommunicationRequest) ProtoMessage() {}

func (x *UnlinkCommunicationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlinkCommunicationRequest.ProtoReflect.Descriptor instead.
func (*UnlinkCommunicationRequest) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{4}
}

func (x *UnlinkCommunicationRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UnlinkCommunicationRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

// Response message after unlinking a communication from a case.
type UnlinkCommunicationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Affected      int64                  `protobuf:"varint,1,opt,name=affected,proto3" json:"affected,omitempty"` // Affected rows.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnlinkCommunicationResponse) Reset() {
	*x = UnlinkCommunicationResponse{}
	mi := &file_case_communication_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnlinkCommunicationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlinkCommunicationResponse) ProtoMessage() {}

func (x *UnlinkCommunicationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlinkCommunicationResponse.ProtoReflect.Descriptor instead.
func (*UnlinkCommunicationResponse) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{5}
}

func (x *UnlinkCommunicationResponse) GetAffected() int64 {
	if x != nil {
		return x.Affected
	}
	return 0
}

// Request message for listing communications linked to a case.
type ListCommunicationsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CaseId        string                 `protobuf:"bytes,1,opt,name=case_id,json=caseId,proto3" json:"case_id,omitempty"` // Case identifier.
	Fields        []string               `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`               // List of fields to include in the response.
	Q             string                 `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`                         // Query string for filtering results.
	Size          int32                  `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`                  // Number of records per page.
	Page          int32                  `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"`                  // Page number for pagination.
	Sort          string                 `protobuf:"bytes,6,opt,name=sort,proto3" json:"sort,omitempty"`                   // Sorting order.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCommunicationsRequest) Reset() {
	*x = ListCommunicationsRequest{}
	mi := &file_case_communication_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCommunicationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommunicationsRequest) ProtoMessage() {}

func (x *ListCommunicationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommunicationsRequest.ProtoReflect.Descriptor instead.
func (*ListCommunicationsRequest) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{6}
}

func (x *ListCommunicationsRequest) GetCaseId() string {
	if x != nil {
		return x.CaseId
	}
	return ""
}

func (x *ListCommunicationsRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ListCommunicationsRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ListCommunicationsRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ListCommunicationsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCommunicationsRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

// Response message for listing communications linked to a case.
type ListCommunicationsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*CaseCommunication   `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`  // List of communications.
	Page          int32                  `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"` // Current page number.
	Next          bool                   `protobuf:"varint,6,opt,name=next,proto3" json:"next,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListCommunicationsResponse) Reset() {
	*x = ListCommunicationsResponse{}
	mi := &file_case_communication_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListCommunicationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListCommunicationsResponse) ProtoMessage() {}

func (x *ListCommunicationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_case_communication_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListCommunicationsResponse.ProtoReflect.Descriptor instead.
func (*ListCommunicationsResponse) Descriptor() ([]byte, []int) {
	return file_case_communication_proto_rawDescGZIP(), []int{7}
}

func (x *ListCommunicationsResponse) GetData() []*CaseCommunication {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *ListCommunicationsResponse) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListCommunicationsResponse) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

var File_case_communication_proto protoreflect.FileDescriptor

var file_case_communication_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x77, 0x65, 0x62, 0x69,
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
	0xb7, 0x01, 0x0a, 0x11, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x76, 0x65, 0x72, 0x12, 0x55, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x75,
	0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x26, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x11, 0x63, 0x6f, 0x6d,
	0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x29,
	0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x9a, 0x01, 0x0a, 0x16, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x55, 0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x26, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73,
	0x2e, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52, 0x11, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x63,
	0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x88, 0x01, 0x0a, 0x18, 0x4c, 0x69, 0x6e, 0x6b, 0x43,
	0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x12, 0x3b, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61,
	0x73, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x22, 0x51, 0x0a, 0x19, 0x4c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73,
	0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x44, 0x0a, 0x1a, 0x55, 0x6e, 0x6c, 0x69, 0x6e, 0x6b, 0x43, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x39, 0x0a, 0x1b, 0x55, 0x6e,
	0x6c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x66, 0x66,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x61, 0x66, 0x66,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x22, 0x96, 0x01, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x01, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6f,
	0x72, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x22, 0x7a,
	0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x43, 0x61, 0x73, 0x65, 0x43,
	0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x2a, 0x6f, 0x0a, 0x17, 0x43, 0x61,
	0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x4e, 0x4f, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f, 0x4d, 0x4d, 0x55, 0x4e, 0x49, 0x43, 0x41, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x48, 0x41, 0x54, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x4f,
	0x4d, 0x4d, 0x55, 0x4e, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x41, 0x4c, 0x4c,
	0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x4f, 0x4d, 0x4d, 0x55, 0x4e, 0x49, 0x43, 0x41, 0x54,
	0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0x10, 0x03, 0x32, 0xf4, 0x03, 0x0a, 0x12,
	0x43, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x94, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75,
	0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d,
	0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x28, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2c, 0x90, 0xb5, 0x18,
	0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x22, 0x20, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f,
	0x7b, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0xa1, 0x01, 0x0a, 0x13, 0x55, 0x6e,
	0x6c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x29, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65,
	0x73, 0x2e, 0x55, 0x6e, 0x6c, 0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x55, 0x6e, 0x6c,
	0x69, 0x6e, 0x6b, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x33, 0x90, 0xb5, 0x18, 0x02, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x29, 0x2a, 0x27, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x63, 0x61,
	0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x7b, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x12, 0x97, 0x01,
	0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x28, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63,
	0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29,
	0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2c, 0x90, 0xb5, 0x18, 0x01, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x22, 0x12, 0x20, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x7b, 0x63,
	0x61, 0x73, 0x65, 0x5f, 0x65, 0x74, 0x61, 0x67, 0x7d, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x09, 0x8a, 0xb5, 0x18, 0x05, 0x63, 0x61, 0x73,
	0x65, 0x73, 0x42, 0xaa, 0x01, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74,
	0x65, 0x6c, 0x2e, 0x63, 0x61, 0x73, 0x65, 0x73, 0x42, 0x16, 0x43, 0x61, 0x73, 0x65, 0x43, 0x6f,
	0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77,
	0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x3b, 0x63, 0x61, 0x73, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x57,
	0x43, 0x58, 0xaa, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x43, 0x61, 0x73,
	0x65, 0x73, 0xca, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73,
	0x65, 0x73, 0xe2, 0x02, 0x19, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x61, 0x73,
	0x65, 0x73, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x0e, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a, 0x3a, 0x43, 0x61, 0x73, 0x65, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_case_communication_proto_rawDescOnce sync.Once
	file_case_communication_proto_rawDescData = file_case_communication_proto_rawDesc
)

func file_case_communication_proto_rawDescGZIP() []byte {
	file_case_communication_proto_rawDescOnce.Do(func() {
		file_case_communication_proto_rawDescData = protoimpl.X.CompressGZIP(file_case_communication_proto_rawDescData)
	})
	return file_case_communication_proto_rawDescData
}

var file_case_communication_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_case_communication_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_case_communication_proto_goTypes = []any{
	(CaseCommunicationsTypes)(0),        // 0: webitel.cases.CaseCommunicationsTypes
	(*CaseCommunication)(nil),           // 1: webitel.cases.CaseCommunication
	(*InputCaseCommunication)(nil),      // 2: webitel.cases.InputCaseCommunication
	(*LinkCommunicationRequest)(nil),    // 3: webitel.cases.LinkCommunicationRequest
	(*LinkCommunicationResponse)(nil),   // 4: webitel.cases.LinkCommunicationResponse
	(*UnlinkCommunicationRequest)(nil),  // 5: webitel.cases.UnlinkCommunicationRequest
	(*UnlinkCommunicationResponse)(nil), // 6: webitel.cases.UnlinkCommunicationResponse
	(*ListCommunicationsRequest)(nil),   // 7: webitel.cases.ListCommunicationsRequest
	(*ListCommunicationsResponse)(nil),  // 8: webitel.cases.ListCommunicationsResponse
}
var file_case_communication_proto_depIdxs = []int32{
	0, // 0: webitel.cases.CaseCommunication.communication_type:type_name -> webitel.cases.CaseCommunicationsTypes
	0, // 1: webitel.cases.InputCaseCommunication.communication_type:type_name -> webitel.cases.CaseCommunicationsTypes
	2, // 2: webitel.cases.LinkCommunicationRequest.input:type_name -> webitel.cases.InputCaseCommunication
	1, // 3: webitel.cases.LinkCommunicationResponse.data:type_name -> webitel.cases.CaseCommunication
	1, // 4: webitel.cases.ListCommunicationsResponse.data:type_name -> webitel.cases.CaseCommunication
	3, // 5: webitel.cases.CaseCommunications.LinkCommunication:input_type -> webitel.cases.LinkCommunicationRequest
	5, // 6: webitel.cases.CaseCommunications.UnlinkCommunication:input_type -> webitel.cases.UnlinkCommunicationRequest
	7, // 7: webitel.cases.CaseCommunications.ListCommunications:input_type -> webitel.cases.ListCommunicationsRequest
	4, // 8: webitel.cases.CaseCommunications.LinkCommunication:output_type -> webitel.cases.LinkCommunicationResponse
	6, // 9: webitel.cases.CaseCommunications.UnlinkCommunication:output_type -> webitel.cases.UnlinkCommunicationResponse
	8, // 10: webitel.cases.CaseCommunications.ListCommunications:output_type -> webitel.cases.ListCommunicationsResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_case_communication_proto_init() }
func file_case_communication_proto_init() {
	if File_case_communication_proto != nil {
		return
	}
	file_general_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_case_communication_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_case_communication_proto_goTypes,
		DependencyIndexes: file_case_communication_proto_depIdxs,
		EnumInfos:         file_case_communication_proto_enumTypes,
		MessageInfos:      file_case_communication_proto_msgTypes,
	}.Build()
	File_case_communication_proto = out.File
	file_case_communication_proto_rawDesc = nil
	file_case_communication_proto_goTypes = nil
	file_case_communication_proto_depIdxs = nil
}
