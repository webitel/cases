// Specifies the syntax version of the protocol buffer.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
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
	unsafe "unsafe"
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
	state             protoimpl.MessageState `protogen:"open.v1"`
	Id                int64                  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`                                                       // Database ID of the communication.
	Ver               int32                  `protobuf:"varint,3,opt,name=ver,proto3" json:"ver,omitempty"`                                                     // Version of the communication record.
	Etag              string                 `protobuf:"bytes,4,opt,name=etag,proto3" json:"etag,omitempty"`                                                    // Version of the communication record.
	CommunicationType *Lookup                `protobuf:"bytes,5,opt,name=communication_type,json=communicationType,proto3" json:"communication_type,omitempty"` // Type of the communication (e.g., Chat, Call).
	CommunicationId   string                 `protobuf:"bytes,6,opt,name=communication_id,json=communicationId,proto3" json:"communication_id,omitempty"`       // External communication ID.
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

func (x *CaseCommunication) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CaseCommunication) GetVer() int32 {
	if x != nil {
		return x.Ver
	}
	return 0
}

func (x *CaseCommunication) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *CaseCommunication) GetCommunicationType() *Lookup {
	if x != nil {
		return x.CommunicationType
	}
	return nil
}

func (x *CaseCommunication) GetCommunicationId() string {
	if x != nil {
		return x.CommunicationId
	}
	return ""
}

// Represents input data for creating or linking a communication.
type InputCaseCommunication struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	CommunicationType *Lookup                `protobuf:"bytes,2,opt,name=communication_type,json=communicationType,proto3" json:"communication_type,omitempty"` // Type of the communication.
	CommunicationId   string                 `protobuf:"bytes,3,opt,name=communication_id,json=communicationId,proto3" json:"communication_id,omitempty"`       // External communication ID.
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

func (x *InputCaseCommunication) GetCommunicationType() *Lookup {
	if x != nil {
		return x.CommunicationType
	}
	return nil
}

func (x *InputCaseCommunication) GetCommunicationId() string {
	if x != nil {
		return x.CommunicationId
	}
	return ""
}

// Request message for linking communications to a case.
type LinkCommunicationRequest struct {
	state         protoimpl.MessageState  `protogen:"open.v1"`
	CaseEtag      string                  `protobuf:"bytes,1,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"` // Case identifier.
	Fields        []string                `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`                     // List of fields to include in the response.
	Input         *InputCaseCommunication `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`                       // Input data for the communications to link.
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

func (x *LinkCommunicationRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
	}
	return ""
}

func (x *LinkCommunicationRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *LinkCommunicationRequest) GetInput() *InputCaseCommunication {
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
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // Communication identifier.
	CaseEtag      string                 `protobuf:"bytes,2,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"`
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

func (x *UnlinkCommunicationRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
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
	CaseEtag      string                 `protobuf:"bytes,1,opt,name=case_etag,json=caseEtag,proto3" json:"case_etag,omitempty"` // Case identifier.
	Fields        []string               `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`                     // List of fields to include in the response.
	Q             string                 `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`                               // Query string for filtering results.
	Size          int32                  `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`                        // Number of records per page.
	Page          int32                  `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"`                        // Page number for pagination.
	Sort          string                 `protobuf:"bytes,6,opt,name=sort,proto3" json:"sort,omitempty"`                         // Sorting order.
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

func (x *ListCommunicationsRequest) GetCaseEtag() string {
	if x != nil {
		return x.CaseEtag
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

const file_case_communication_proto_rawDesc = "" +
	"\n" +
	"\x18case_communication.proto\x12\rwebitel.cases\x1a\rgeneral.proto\x1a\x1bgoogle/api/visibility.proto\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\x1a\x1aproto/webitel/option.proto\"\xb4\x01\n" +
	"\x11CaseCommunication\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\x03R\x02id\x12\x10\n" +
	"\x03ver\x18\x03 \x01(\x05R\x03ver\x12\x12\n" +
	"\x04etag\x18\x04 \x01(\tR\x04etag\x12>\n" +
	"\x12communication_type\x18\x05 \x01(\v2\x0f.general.LookupR\x11communicationType\x12)\n" +
	"\x10communication_id\x18\x06 \x01(\tR\x0fcommunicationId\"\x83\x01\n" +
	"\x16InputCaseCommunication\x12>\n" +
	"\x12communication_type\x18\x02 \x01(\v2\x0f.general.LookupR\x11communicationType\x12)\n" +
	"\x10communication_id\x18\x03 \x01(\tR\x0fcommunicationId\"\x8c\x01\n" +
	"\x18LinkCommunicationRequest\x12\x1b\n" +
	"\tcase_etag\x18\x01 \x01(\tR\bcaseEtag\x12\x16\n" +
	"\x06fields\x18\x02 \x03(\tR\x06fields\x12;\n" +
	"\x05input\x18\x03 \x01(\v2%.webitel.cases.InputCaseCommunicationR\x05input\"Q\n" +
	"\x19LinkCommunicationResponse\x124\n" +
	"\x04data\x18\x01 \x03(\v2 .webitel.cases.CaseCommunicationR\x04data\"a\n" +
	"\x1aUnlinkCommunicationRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1b\n" +
	"\tcase_etag\x18\x02 \x01(\tR\bcaseEtag\x12\x16\n" +
	"\x06fields\x18\x03 \x03(\tR\x06fields\"9\n" +
	"\x1bUnlinkCommunicationResponse\x12\x1a\n" +
	"\baffected\x18\x01 \x01(\x03R\baffected\"\x9a\x01\n" +
	"\x19ListCommunicationsRequest\x12\x1b\n" +
	"\tcase_etag\x18\x01 \x01(\tR\bcaseEtag\x12\x16\n" +
	"\x06fields\x18\x02 \x03(\tR\x06fields\x12\f\n" +
	"\x01q\x18\x03 \x01(\tR\x01q\x12\x12\n" +
	"\x04size\x18\x04 \x01(\x05R\x04size\x12\x12\n" +
	"\x04page\x18\x05 \x01(\x05R\x04page\x12\x12\n" +
	"\x04sort\x18\x06 \x01(\tR\x04sort\"z\n" +
	"\x1aListCommunicationsResponse\x124\n" +
	"\x04data\x18\x01 \x03(\v2 .webitel.cases.CaseCommunicationR\x04data\x12\x12\n" +
	"\x04page\x18\x05 \x01(\x05R\x04page\x12\x12\n" +
	"\x04next\x18\x06 \x01(\bR\x04next*o\n" +
	"\x17CaseCommunicationsTypes\x12\v\n" +
	"\aNO_TYPE\x10\x00\x12\x16\n" +
	"\x12COMMUNICATION_CHAT\x10\x01\x12\x16\n" +
	"\x12COMMUNICATION_CALL\x10\x02\x12\x17\n" +
	"\x13COMMUNICATION_EMAIL\x10\x032\xf9\x03\n" +
	"\x12CaseCommunications\x12\x9b\x01\n" +
	"\x11LinkCommunication\x12'.webitel.cases.LinkCommunicationRequest\x1a(.webitel.cases.LinkCommunicationResponse\"3\x90\xb5\x18\x01\x82\xd3\xe4\x93\x02):\x05input\" /cases/{case_etag}/communication\x12\x9f\x01\n" +
	"\x13UnlinkCommunication\x12).webitel.cases.UnlinkCommunicationRequest\x1a*.webitel.cases.UnlinkCommunicationResponse\"1\x90\xb5\x18\x02\x82\xd3\xe4\x93\x02'*%/cases/{case_etag}/communication/{id}\x12\x97\x01\n" +
	"\x12ListCommunications\x12(.webitel.cases.ListCommunicationsRequest\x1a).webitel.cases.ListCommunicationsResponse\",\x90\xb5\x18\x01\x82\xd3\xe4\x93\x02\"\x12 /cases/{case_etag}/communication\x1a\t\x8a\xb5\x18\x05casesB\xaa\x01\n" +
	"\x11com.webitel.casesB\x16CaseCommunicationProtoP\x01Z(github.com/webitel/cases/api/cases;cases\xa2\x02\x03WCX\xaa\x02\rWebitel.Cases\xca\x02\rWebitel\\Cases\xe2\x02\x19Webitel\\Cases\\GPBMetadata\xea\x02\x0eWebitel::Casesb\x06proto3"

var (
	file_case_communication_proto_rawDescOnce sync.Once
	file_case_communication_proto_rawDescData []byte
)

func file_case_communication_proto_rawDescGZIP() []byte {
	file_case_communication_proto_rawDescOnce.Do(func() {
		file_case_communication_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_case_communication_proto_rawDesc), len(file_case_communication_proto_rawDesc)))
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
	(*Lookup)(nil),                      // 9: general.Lookup
}
var file_case_communication_proto_depIdxs = []int32{
	9, // 0: webitel.cases.CaseCommunication.communication_type:type_name -> general.Lookup
	9, // 1: webitel.cases.InputCaseCommunication.communication_type:type_name -> general.Lookup
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
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_case_communication_proto_rawDesc), len(file_case_communication_proto_rawDesc)),
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
	file_case_communication_proto_goTypes = nil
	file_case_communication_proto_depIdxs = nil
}
