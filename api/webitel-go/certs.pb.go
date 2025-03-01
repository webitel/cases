// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: certs.proto

package api

import (
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

type Validity struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NotBefore     int64                  `protobuf:"varint,1,opt,name=not_before,json=notBefore,proto3" json:"not_before,omitempty"` // unix
	NotAfter      int64                  `protobuf:"varint,2,opt,name=not_after,json=notAfter,proto3" json:"not_after,omitempty"`    // unix
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Validity) Reset() {
	*x = Validity{}
	mi := &file_certs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Validity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Validity) ProtoMessage() {}

func (x *Validity) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Validity.ProtoReflect.Descriptor instead.
func (*Validity) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{0}
}

func (x *Validity) GetNotBefore() int64 {
	if x != nil {
		return x.NotBefore
	}
	return 0
}

func (x *Validity) GetNotAfter() int64 {
	if x != nil {
		return x.NotAfter
	}
	return 0
}

type License struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Serial        string                 `protobuf:"bytes,1,opt,name=serial,proto3" json:"serial,omitempty"`     // serial number assigned (global::License-ID)
	Scope         string                 `protobuf:"bytes,2,opt,name=scope,proto3" json:"scope,omitempty"`       // mandatory privilege codename, e.g.: DEVICE, MANAGER, OPERATOR
	Limit         uint32                 `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`      // required: limit count (maximum allowed usage quantity)
	Validity      *Validity              `protobuf:"bytes,4,opt,name=validity,proto3" json:"validity,omitempty"` // optional
	Competitive   bool                   `protobuf:"varint,5,opt,name=competitive,proto3" json:"competitive,omitempty"`
	Users         map[int64]string       `protobuf:"bytes,6,rep,name=users,proto3" json:"users,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"` // users.id (grantees)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *License) Reset() {
	*x = License{}
	mi := &file_certs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *License) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*License) ProtoMessage() {}

func (x *License) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use License.ProtoReflect.Descriptor instead.
func (*License) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{1}
}

func (x *License) GetSerial() string {
	if x != nil {
		return x.Serial
	}
	return ""
}

func (x *License) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

func (x *License) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *License) GetValidity() *Validity {
	if x != nil {
		return x.Validity
	}
	return nil
}

func (x *License) GetCompetitive() bool {
	if x != nil {
		return x.Competitive
	}
	return false
}

func (x *License) GetUsers() map[int64]string {
	if x != nil {
		return x.Users
	}
	return nil
}

type Certificate struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Serial        string                 `protobuf:"bytes,1,opt,name=serial,proto3" json:"serial,omitempty"`                            // serial number assigned (global::Customer-ID)
	Version       string                 `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`                          // version
	Validity      *Validity              `protobuf:"bytes,3,opt,name=validity,proto3" json:"validity,omitempty"`                        // int32 competitive = 4; // zero-based competitive sessions limit ?
	License       []*License             `protobuf:"bytes,5,rep,name=license,proto3" json:"license,omitempty"`                          // grants issued
	ThisUpdate    int64                  `protobuf:"varint,6,opt,name=this_update,json=thisUpdate,proto3" json:"this_update,omitempty"` // unix: issuer signed at
	NextUpdate    int64                  `protobuf:"varint,7,opt,name=next_update,json=nextUpdate,proto3" json:"next_update,omitempty"` // unix: signature expires; update required
	Valid         bool                   `protobuf:"varint,8,opt,name=valid,proto3" json:"valid,omitempty"`                             // validation status
	Domains       map[int64]string       `protobuf:"bytes,9,rep,name=domains,proto3" json:"domains,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Certificate) Reset() {
	*x = Certificate{}
	mi := &file_certs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Certificate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Certificate) ProtoMessage() {}

func (x *Certificate) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Certificate.ProtoReflect.Descriptor instead.
func (*Certificate) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{2}
}

func (x *Certificate) GetSerial() string {
	if x != nil {
		return x.Serial
	}
	return ""
}

func (x *Certificate) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Certificate) GetValidity() *Validity {
	if x != nil {
		return x.Validity
	}
	return nil
}

func (x *Certificate) GetLicense() []*License {
	if x != nil {
		return x.License
	}
	return nil
}

func (x *Certificate) GetThisUpdate() int64 {
	if x != nil {
		return x.ThisUpdate
	}
	return 0
}

func (x *Certificate) GetNextUpdate() int64 {
	if x != nil {
		return x.NextUpdate
	}
	return 0
}

func (x *Certificate) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

func (x *Certificate) GetDomains() map[int64]string {
	if x != nil {
		return x.Domains
	}
	return nil
}

// GET /certificate/{filter=**}
// GET /user/{userId}/certificate
// GET /domain/{domain}/certificate
type CertificateUsageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Domain        string                 `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`                // domain relation
	Serial        string                 `protobuf:"bytes,2,opt,name=serial,proto3" json:"serial,omitempty"`                // filter: serial
	UserId        int64                  `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // filter: grants.user.id grantee
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CertificateUsageRequest) Reset() {
	*x = CertificateUsageRequest{}
	mi := &file_certs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CertificateUsageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateUsageRequest) ProtoMessage() {}

func (x *CertificateUsageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateUsageRequest.ProtoReflect.Descriptor instead.
func (*CertificateUsageRequest) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{3}
}

func (x *CertificateUsageRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *CertificateUsageRequest) GetSerial() string {
	if x != nil {
		return x.Serial
	}
	return ""
}

func (x *CertificateUsageRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type CertificateUsageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Certificate   *Certificate           `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"` // detailed
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CertificateUsageResponse) Reset() {
	*x = CertificateUsageResponse{}
	mi := &file_certs_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CertificateUsageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CertificateUsageResponse) ProtoMessage() {}

func (x *CertificateUsageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CertificateUsageResponse.ProtoReflect.Descriptor instead.
func (*CertificateUsageResponse) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{4}
}

func (x *CertificateUsageResponse) GetCertificate() *Certificate {
	if x != nil {
		return x.Certificate
	}
	return nil
}

// PUT /certificate
type UpdateCertificateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Certificate   []byte                 `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"` // raw bytes
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCertificateRequest) Reset() {
	*x = UpdateCertificateRequest{}
	mi := &file_certs_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCertificateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCertificateRequest) ProtoMessage() {}

func (x *UpdateCertificateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCertificateRequest.ProtoReflect.Descriptor instead.
func (*UpdateCertificateRequest) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateCertificateRequest) GetCertificate() []byte {
	if x != nil {
		return x.Certificate
	}
	return nil
}

type UpdateCertificateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Certificate   *Certificate           `protobuf:"bytes,1,opt,name=certificate,proto3" json:"certificate,omitempty"` // detailed
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateCertificateResponse) Reset() {
	*x = UpdateCertificateResponse{}
	mi := &file_certs_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateCertificateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCertificateResponse) ProtoMessage() {}

func (x *UpdateCertificateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCertificateResponse.ProtoReflect.Descriptor instead.
func (*UpdateCertificateResponse) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateCertificateResponse) GetCertificate() *Certificate {
	if x != nil {
		return x.Certificate
	}
	return nil
}

// GET /certificates
type SearchCertificatesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Domain        string                 `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchCertificatesRequest) Reset() {
	*x = SearchCertificatesRequest{}
	mi := &file_certs_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchCertificatesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchCertificatesRequest) ProtoMessage() {}

func (x *SearchCertificatesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchCertificatesRequest.ProtoReflect.Descriptor instead.
func (*SearchCertificatesRequest) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{7}
}

func (x *SearchCertificatesRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

type SearchCertificatesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Certificates  []*Certificate         `protobuf:"bytes,1,rep,name=certificates,proto3" json:"certificates,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchCertificatesResponse) Reset() {
	*x = SearchCertificatesResponse{}
	mi := &file_certs_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchCertificatesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchCertificatesResponse) ProtoMessage() {}

func (x *SearchCertificatesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchCertificatesResponse.ProtoReflect.Descriptor instead.
func (*SearchCertificatesResponse) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{8}
}

func (x *SearchCertificatesResponse) GetCertificates() []*Certificate {
	if x != nil {
		return x.Certificates
	}
	return nil
}

// POST /user/{userId}/certificate/{serial}
type GrantCertificateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Domain        string                 `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`                // domain relation
	Serial        string                 `protobuf:"bytes,2,opt,name=serial,proto3" json:"serial,omitempty"`                // grants.serial
	UserId        int64                  `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // grantee
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GrantCertificateRequest) Reset() {
	*x = GrantCertificateRequest{}
	mi := &file_certs_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrantCertificateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantCertificateRequest) ProtoMessage() {}

func (x *GrantCertificateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantCertificateRequest.ProtoReflect.Descriptor instead.
func (*GrantCertificateRequest) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{9}
}

func (x *GrantCertificateRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *GrantCertificateRequest) GetSerial() string {
	if x != nil {
		return x.Serial
	}
	return ""
}

func (x *GrantCertificateRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GrantCertificateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GrantCertificateResponse) Reset() {
	*x = GrantCertificateResponse{}
	mi := &file_certs_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GrantCertificateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantCertificateResponse) ProtoMessage() {}

func (x *GrantCertificateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantCertificateResponse.ProtoReflect.Descriptor instead.
func (*GrantCertificateResponse) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{10}
}

// DELETE /user/{userId}/certificate/{serial}
type RevokeCertificateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Domain        string                 `protobuf:"bytes,1,opt,name=domain,proto3" json:"domain,omitempty"`                // domain relation
	Serial        string                 `protobuf:"bytes,2,opt,name=serial,proto3" json:"serial,omitempty"`                // grants.serial
	UserId        int64                  `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // grantee
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RevokeCertificateRequest) Reset() {
	*x = RevokeCertificateRequest{}
	mi := &file_certs_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RevokeCertificateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeCertificateRequest) ProtoMessage() {}

func (x *RevokeCertificateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeCertificateRequest.ProtoReflect.Descriptor instead.
func (*RevokeCertificateRequest) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{11}
}

func (x *RevokeCertificateRequest) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *RevokeCertificateRequest) GetSerial() string {
	if x != nil {
		return x.Serial
	}
	return ""
}

func (x *RevokeCertificateRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RevokeCertificateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RevokeCertificateResponse) Reset() {
	*x = RevokeCertificateResponse{}
	mi := &file_certs_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RevokeCertificateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeCertificateResponse) ProtoMessage() {}

func (x *RevokeCertificateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_certs_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeCertificateResponse.ProtoReflect.Descriptor instead.
func (*RevokeCertificateResponse) Descriptor() ([]byte, []int) {
	return file_certs_proto_rawDescGZIP(), []int{12}
}

var File_certs_proto protoreflect.FileDescriptor

var file_certs_proto_rawDesc = string([]byte{
	0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x69, 0x22, 0x46, 0x0a, 0x08, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x12, 0x1d,
	0x0a, 0x0a, 0x6e, 0x6f, 0x74, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x6e, 0x6f, 0x74, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x6e, 0x6f, 0x74, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x6e, 0x6f, 0x74, 0x41, 0x66, 0x74, 0x65, 0x72, 0x22, 0x83, 0x02, 0x0a, 0x07, 0x4c,
	0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x29, 0x0a, 0x08, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x52, 0x08, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x69, 0x74, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x65, 0x74, 0x69,
	0x74, 0x69, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70,
	0x65, 0x74, 0x69, 0x74, 0x69, 0x76, 0x65, 0x12, 0x2d, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x63,
	0x65, 0x6e, 0x73, 0x65, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0xdf, 0x02, 0x0a, 0x0b, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x29, 0x0a, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x69, 0x74, 0x79, 0x52, 0x08, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x12, 0x26, 0x0a,
	0x07, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x52, 0x07, 0x6c, 0x69,
	0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x68, 0x69, 0x73, 0x5f, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x68, 0x69, 0x73,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6e, 0x65, 0x78,
	0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x12, 0x37, 0x0a,
	0x07, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x2e, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x62, 0x0a, 0x17, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4e, 0x0a, 0x18, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x32, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x22, 0x3c, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x22, 0x4f, 0x0a, 0x19, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x32, 0x0a, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x22, 0x33, 0x0a, 0x19, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43,
	0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x52, 0x0a, 0x1a, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x0c, 0x63, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x52, 0x0c, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x22, 0x62,
	0x0a, 0x17, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x1a, 0x0a, 0x18, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x63,
	0x0a, 0x18, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x22, 0x1b, 0x0a, 0x19, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0xaf, 0x03, 0x0a, 0x02, 0x43, 0x41, 0x12, 0x51, 0x0a, 0x10, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x11, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12,
	0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x65, 0x72, 0x74,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x57, 0x0a, 0x12, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x12, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x10, 0x47, 0x72, 0x61,
	0x6e, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69,
	0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x54, 0x0a, 0x11,
	0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x12, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x43, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x55, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x42, 0x0a, 0x43,
	0x65, 0x72, 0x74, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x12, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2e, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69, 0xa2,
	0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x03, 0x41, 0x70, 0x69, 0xca, 0x02, 0x03, 0x41, 0x70,
	0x69, 0xe2, 0x02, 0x0f, 0x41, 0x70, 0x69, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x03, 0x41, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_certs_proto_rawDescOnce sync.Once
	file_certs_proto_rawDescData []byte
)

func file_certs_proto_rawDescGZIP() []byte {
	file_certs_proto_rawDescOnce.Do(func() {
		file_certs_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_certs_proto_rawDesc), len(file_certs_proto_rawDesc)))
	})
	return file_certs_proto_rawDescData
}

var file_certs_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_certs_proto_goTypes = []any{
	(*Validity)(nil),                   // 0: api.Validity
	(*License)(nil),                    // 1: api.License
	(*Certificate)(nil),                // 2: api.Certificate
	(*CertificateUsageRequest)(nil),    // 3: api.CertificateUsageRequest
	(*CertificateUsageResponse)(nil),   // 4: api.CertificateUsageResponse
	(*UpdateCertificateRequest)(nil),   // 5: api.UpdateCertificateRequest
	(*UpdateCertificateResponse)(nil),  // 6: api.UpdateCertificateResponse
	(*SearchCertificatesRequest)(nil),  // 7: api.SearchCertificatesRequest
	(*SearchCertificatesResponse)(nil), // 8: api.SearchCertificatesResponse
	(*GrantCertificateRequest)(nil),    // 9: api.GrantCertificateRequest
	(*GrantCertificateResponse)(nil),   // 10: api.GrantCertificateResponse
	(*RevokeCertificateRequest)(nil),   // 11: api.RevokeCertificateRequest
	(*RevokeCertificateResponse)(nil),  // 12: api.RevokeCertificateResponse
	nil,                                // 13: api.License.UsersEntry
	nil,                                // 14: api.Certificate.DomainsEntry
}
var file_certs_proto_depIdxs = []int32{
	0,  // 0: api.License.validity:type_name -> api.Validity
	13, // 1: api.License.users:type_name -> api.License.UsersEntry
	0,  // 2: api.Certificate.validity:type_name -> api.Validity
	1,  // 3: api.Certificate.license:type_name -> api.License
	14, // 4: api.Certificate.domains:type_name -> api.Certificate.DomainsEntry
	2,  // 5: api.CertificateUsageResponse.certificate:type_name -> api.Certificate
	2,  // 6: api.UpdateCertificateResponse.certificate:type_name -> api.Certificate
	2,  // 7: api.SearchCertificatesResponse.certificates:type_name -> api.Certificate
	3,  // 8: api.CA.CertificateUsage:input_type -> api.CertificateUsageRequest
	5,  // 9: api.CA.UpdateCertificate:input_type -> api.UpdateCertificateRequest
	7,  // 10: api.CA.SearchCertificates:input_type -> api.SearchCertificatesRequest
	9,  // 11: api.CA.GrantCertificate:input_type -> api.GrantCertificateRequest
	11, // 12: api.CA.RevokeCertificate:input_type -> api.RevokeCertificateRequest
	4,  // 13: api.CA.CertificateUsage:output_type -> api.CertificateUsageResponse
	6,  // 14: api.CA.UpdateCertificate:output_type -> api.UpdateCertificateResponse
	8,  // 15: api.CA.SearchCertificates:output_type -> api.SearchCertificatesResponse
	10, // 16: api.CA.GrantCertificate:output_type -> api.GrantCertificateResponse
	12, // 17: api.CA.RevokeCertificate:output_type -> api.RevokeCertificateResponse
	13, // [13:18] is the sub-list for method output_type
	8,  // [8:13] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_certs_proto_init() }
func file_certs_proto_init() {
	if File_certs_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_certs_proto_rawDesc), len(file_certs_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_certs_proto_goTypes,
		DependencyIndexes: file_certs_proto_depIdxs,
		MessageInfos:      file_certs_proto_msgTypes,
	}.Build()
	File_certs_proto = out.File
	file_certs_proto_goTypes = nil
	file_certs_proto_depIdxs = nil
}
