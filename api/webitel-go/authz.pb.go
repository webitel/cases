// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: authz.proto

package api

import (
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

type UserinfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"` // string id_token = 2;
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserinfoRequest) Reset() {
	*x = UserinfoRequest{}
	mi := &file_authz_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserinfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserinfoRequest) ProtoMessage() {}

func (x *UserinfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authz_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserinfoRequest.ProtoReflect.Descriptor instead.
func (*UserinfoRequest) Descriptor() ([]byte, []int) {
	return file_authz_proto_rawDescGZIP(), []int{0}
}

func (x *UserinfoRequest) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

// AccessScope represents authorized access level
// to object class declaration within domain namespace
type Objclass struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Inherit::ObjectClass
	Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`      // class: object id
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`   // class: display common name
	Class string `protobuf:"bytes,3,opt,name=class,proto3" json:"class,omitempty"` // class: alphanumeric code name
	// Is [?]-[b]ased [A]ccess [C]ontrol *model enabled ?
	Abac bool `protobuf:"varint,4,opt,name=abac,proto3" json:"abac,omitempty"` // [A]ttribute-[b]ased;  not implemented; next releases ...
	Obac bool `protobuf:"varint,5,opt,name=obac,proto3" json:"obac,omitempty"` // [O]peration-[b]ased;  Mandatory: control access to object.class (alias: collection, section, etc.)
	Rbac bool `protobuf:"varint,6,opt,name=rbac,proto3" json:"rbac,omitempty"` // [R]ecord-[b]ased; Discretionary: control access to object.entry (alias: resource, record, etc.)
	// Extension: discretionary access control
	Access        string `protobuf:"bytes,7,opt,name=access,proto3" json:"access,omitempty"` // flags: [ CREATE | SELECT | UPDATE | DELETE ]
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Objclass) Reset() {
	*x = Objclass{}
	mi := &file_authz_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Objclass) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Objclass) ProtoMessage() {}

func (x *Objclass) ProtoReflect() protoreflect.Message {
	mi := &file_authz_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Objclass.ProtoReflect.Descriptor instead.
func (*Objclass) Descriptor() ([]byte, []int) {
	return file_authz_proto_rawDescGZIP(), []int{1}
}

func (x *Objclass) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Objclass) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Objclass) GetClass() string {
	if x != nil {
		return x.Class
	}
	return ""
}

func (x *Objclass) GetAbac() bool {
	if x != nil {
		return x.Abac
	}
	return false
}

func (x *Objclass) GetObac() bool {
	if x != nil {
		return x.Obac
	}
	return false
}

func (x *Objclass) GetRbac() bool {
	if x != nil {
		return x.Rbac
	}
	return false
}

func (x *Objclass) GetAccess() string {
	if x != nil {
		return x.Access
	}
	return ""
}

// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
type Userinfo struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	Dc                int64                  `protobuf:"varint,2,opt,name=dc,proto3" json:"dc,omitempty"`                          // current domain component id
	Domain            string                 `protobuf:"bytes,3,opt,name=domain,json=realm,proto3" json:"domain,omitempty"`        // session domain.name
	UserId            int64                  `protobuf:"varint,4,opt,name=user_id,json=sub,proto3" json:"user_id,omitempty"`       // current user.id
	Name              string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`                       // current user.id
	Username          string                 `protobuf:"bytes,6,opt,name=username,json=nickname,proto3" json:"username,omitempty"` // session end-user login
	PreferredUsername string                 `protobuf:"bytes,7,opt,name=preferred_username,proto3" json:"preferred_username,omitempty"`
	Extension         string                 `protobuf:"bytes,8,opt,name=extension,json=phone_number,proto3" json:"extension,omitempty"`
	Scope             []*Objclass            `protobuf:"bytes,10,rep,name=scope,proto3" json:"scope,omitempty"`                           // map[class]dac
	Roles             []*ObjectId            `protobuf:"bytes,11,rep,name=roles,proto3" json:"roles,omitempty"`                           // map[role]oid
	License           []*LicenseUser         `protobuf:"bytes,12,rep,name=license,proto3" json:"license,omitempty"`                       // map[key]details
	Permissions       []*Permission          `protobuf:"bytes,13,rep,name=permissions,proto3" json:"permissions,omitempty"`               //
	UpdatedAt         int64                  `protobuf:"varint,20,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"` // user last updated
	ExpiresAt         int64                  `protobuf:"varint,21,opt,name=expires_at,json=exp,proto3" json:"expires_at,omitempty"`       // unix
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *Userinfo) Reset() {
	*x = Userinfo{}
	mi := &file_authz_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Userinfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Userinfo) ProtoMessage() {}

func (x *Userinfo) ProtoReflect() protoreflect.Message {
	mi := &file_authz_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Userinfo.ProtoReflect.Descriptor instead.
func (*Userinfo) Descriptor() ([]byte, []int) {
	return file_authz_proto_rawDescGZIP(), []int{2}
}

func (x *Userinfo) GetDc() int64 {
	if x != nil {
		return x.Dc
	}
	return 0
}

func (x *Userinfo) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Userinfo) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Userinfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Userinfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Userinfo) GetPreferredUsername() string {
	if x != nil {
		return x.PreferredUsername
	}
	return ""
}

func (x *Userinfo) GetExtension() string {
	if x != nil {
		return x.Extension
	}
	return ""
}

func (x *Userinfo) GetScope() []*Objclass {
	if x != nil {
		return x.Scope
	}
	return nil
}

func (x *Userinfo) GetRoles() []*ObjectId {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *Userinfo) GetLicense() []*LicenseUser {
	if x != nil {
		return x.License
	}
	return nil
}

func (x *Userinfo) GetPermissions() []*Permission {
	if x != nil {
		return x.Permissions
	}
	return nil
}

func (x *Userinfo) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Userinfo) GetExpiresAt() int64 {
	if x != nil {
		return x.ExpiresAt
	}
	return 0
}

var File_authz_proto protoreflect.FileDescriptor

var file_authz_proto_rawDesc = string([]byte{
	0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x7a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x69, 0x1a, 0x09, 0x6f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0a, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x63, 0x75, 0x73, 0x74, 0x6f,
	0x6d, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x34, 0x0a, 0x0f, 0x55, 0x73,
	0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x98, 0x01, 0x0a, 0x08, 0x4f, 0x62, 0x6a, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x63, 0x6c, 0x61, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x62, 0x61, 0x63, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x61, 0x62, 0x61, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x6f,
	0x62, 0x61, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6f, 0x62, 0x61, 0x63, 0x12,
	0x12, 0x0a, 0x04, 0x72, 0x62, 0x61, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x72,
	0x62, 0x61, 0x63, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xa9, 0x03, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x64, 0x63, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61,
	0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x12,
	0x14, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x73, 0x75, 0x62, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63,
	0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x12, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72,
	0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x12, 0x70, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x09, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18,
	0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4f, 0x62, 0x6a, 0x63,
	0x6c, 0x61, 0x73, 0x73, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x05, 0x72,
	0x6f, 0x6c, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x52, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73,
	0x12, 0x2a, 0x0a, 0x07, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x18, 0x0c, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x07, 0x6c, 0x69, 0x63, 0x65, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x0b,
	0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x0b, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x17,
	0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61, 0x74, 0x18, 0x15, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x65, 0x78, 0x70, 0x32, 0x99, 0x01, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x4f, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66,
	0x6f, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x5a, 0x0b, 0x12, 0x09, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x09, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x69, 0x6e, 0x66,
	0x6f, 0x12, 0x40, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x11, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x0f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x09, 0x22, 0x07, 0x2f, 0x73, 0x69, 0x67,
	0x6e, 0x75, 0x70, 0x42, 0x55, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x42, 0x0a,
	0x41, 0x75, 0x74, 0x68, 0x7a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x12, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x67, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69,
	0xa2, 0x02, 0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x03, 0x41, 0x70, 0x69, 0xca, 0x02, 0x03, 0x41,
	0x70, 0x69, 0xe2, 0x02, 0x0f, 0x41, 0x70, 0x69, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x03, 0x41, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_authz_proto_rawDescOnce sync.Once
	file_authz_proto_rawDescData []byte
)

func file_authz_proto_rawDescGZIP() []byte {
	file_authz_proto_rawDescOnce.Do(func() {
		file_authz_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_authz_proto_rawDesc), len(file_authz_proto_rawDesc)))
	})
	return file_authz_proto_rawDescData
}

var file_authz_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_authz_proto_goTypes = []any{
	(*UserinfoRequest)(nil), // 0: api.UserinfoRequest
	(*Objclass)(nil),        // 1: api.Objclass
	(*Userinfo)(nil),        // 2: api.Userinfo
	(*ObjectId)(nil),        // 3: api.ObjectId
	(*LicenseUser)(nil),     // 4: api.LicenseUser
	(*Permission)(nil),      // 5: api.Permission
	(*LoginRequest)(nil),    // 6: api.LoginRequest
	(*LoginResponse)(nil),   // 7: api.LoginResponse
}
var file_authz_proto_depIdxs = []int32{
	1, // 0: api.Userinfo.scope:type_name -> api.Objclass
	3, // 1: api.Userinfo.roles:type_name -> api.ObjectId
	4, // 2: api.Userinfo.license:type_name -> api.LicenseUser
	5, // 3: api.Userinfo.permissions:type_name -> api.Permission
	0, // 4: api.Auth.UserInfo:input_type -> api.UserinfoRequest
	6, // 5: api.Auth.Signup:input_type -> api.LoginRequest
	2, // 6: api.Auth.UserInfo:output_type -> api.Userinfo
	7, // 7: api.Auth.Signup:output_type -> api.LoginResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_authz_proto_init() }
func file_authz_proto_init() {
	if File_authz_proto != nil {
		return
	}
	file_oid_proto_init()
	file_auth_proto_init()
	file_customers_proto_init()
	file_permission_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_authz_proto_rawDesc), len(file_authz_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_authz_proto_goTypes,
		DependencyIndexes: file_authz_proto_depIdxs,
		MessageInfos:      file_authz_proto_msgTypes,
	}.Build()
	File_authz_proto = out.File
	file_authz_proto_goTypes = nil
	file_authz_proto_depIdxs = nil
}
