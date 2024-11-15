// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        (unknown)
// source: google/api/visibility.proto

package visibility

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// `Visibility` restricts service consumer's access to service elements,
// such as whether an application can call a visibility-restricted method.
// The restriction is expressed by applying visibility labels on service
// elements. The visibility labels are elsewhere linked to service consumers.
//
// A service can define multiple visibility labels, but a service consumer
// should be granted at most one visibility label. Multiple visibility
// labels for a single service consumer are not supported.
//
// If an element and all its parents have no visibility label, its visibility
// is unconditionally granted.
//
// Example:
//
//	visibility:
//	  rules:
//	  - selector: google.calendar.Calendar.EnhancedSearch
//	    restriction: PREVIEW
//	  - selector: google.calendar.Calendar.Delegate
//	    restriction: INTERNAL
//
// Here, all methods are publicly visible except for the restricted methods
// EnhancedSearch and Delegate.
type Visibility struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of visibility rules that apply to individual API elements.
	//
	// **NOTE:** All service configuration rules follow "last one wins" order.
	Rules []*VisibilityRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *Visibility) Reset() {
	*x = Visibility{}
	mi := &file_google_api_visibility_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Visibility) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Visibility) ProtoMessage() {}

func (x *Visibility) ProtoReflect() protoreflect.Message {
	mi := &file_google_api_visibility_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Visibility.ProtoReflect.Descriptor instead.
func (*Visibility) Descriptor() ([]byte, []int) {
	return file_google_api_visibility_proto_rawDescGZIP(), []int{0}
}

func (x *Visibility) GetRules() []*VisibilityRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

// A visibility rule provides visibility configuration for an individual API
// element.
type VisibilityRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Selects methods, messages, fields, enums, etc. to which this rule applies.
	//
	// Refer to [selector][google.api.DocumentationRule.selector] for syntax
	// details.
	Selector string `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// A comma-separated list of visibility labels that apply to the `selector`.
	// Any of the listed labels can be used to grant the visibility.
	//
	// If a rule has multiple labels, removing one of the labels but not all of
	// them can break clients.
	//
	// Example:
	//
	//	visibility:
	//	  rules:
	//	  - selector: google.calendar.Calendar.EnhancedSearch
	//	    restriction: INTERNAL, PREVIEW
	//
	// Removing INTERNAL from this restriction will break clients that rely on
	// this method and only had access to it through INTERNAL.
	Restriction string `protobuf:"bytes,2,opt,name=restriction,proto3" json:"restriction,omitempty"`
}

func (x *VisibilityRule) Reset() {
	*x = VisibilityRule{}
	mi := &file_google_api_visibility_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VisibilityRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VisibilityRule) ProtoMessage() {}

func (x *VisibilityRule) ProtoReflect() protoreflect.Message {
	mi := &file_google_api_visibility_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VisibilityRule.ProtoReflect.Descriptor instead.
func (*VisibilityRule) Descriptor() ([]byte, []int) {
	return file_google_api_visibility_proto_rawDescGZIP(), []int{1}
}

func (x *VisibilityRule) GetSelector() string {
	if x != nil {
		return x.Selector
	}
	return ""
}

func (x *VisibilityRule) GetRestriction() string {
	if x != nil {
		return x.Restriction
	}
	return ""
}

var file_google_api_visibility_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.EnumOptions)(nil),
		ExtensionType: (*VisibilityRule)(nil),
		Field:         72295727,
		Name:          "google.api.enum_visibility",
		Tag:           "bytes,72295727,opt,name=enum_visibility",
		Filename:      "google/api/visibility.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*VisibilityRule)(nil),
		Field:         72295727,
		Name:          "google.api.value_visibility",
		Tag:           "bytes,72295727,opt,name=value_visibility",
		Filename:      "google/api/visibility.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*VisibilityRule)(nil),
		Field:         72295727,
		Name:          "google.api.field_visibility",
		Tag:           "bytes,72295727,opt,name=field_visibility",
		Filename:      "google/api/visibility.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*VisibilityRule)(nil),
		Field:         72295727,
		Name:          "google.api.message_visibility",
		Tag:           "bytes,72295727,opt,name=message_visibility",
		Filename:      "google/api/visibility.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*VisibilityRule)(nil),
		Field:         72295727,
		Name:          "google.api.method_visibility",
		Tag:           "bytes,72295727,opt,name=method_visibility",
		Filename:      "google/api/visibility.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*VisibilityRule)(nil),
		Field:         72295727,
		Name:          "google.api.api_visibility",
		Tag:           "bytes,72295727,opt,name=api_visibility",
		Filename:      "google/api/visibility.proto",
	},
}

// Extension fields to descriptorpb.EnumOptions.
var (
	// See `VisibilityRule`.
	//
	// optional google.api.VisibilityRule enum_visibility = 72295727;
	E_EnumVisibility = &file_google_api_visibility_proto_extTypes[0]
)

// Extension fields to descriptorpb.EnumValueOptions.
var (
	// See `VisibilityRule`.
	//
	// optional google.api.VisibilityRule value_visibility = 72295727;
	E_ValueVisibility = &file_google_api_visibility_proto_extTypes[1]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// See `VisibilityRule`.
	//
	// optional google.api.VisibilityRule field_visibility = 72295727;
	E_FieldVisibility = &file_google_api_visibility_proto_extTypes[2]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// See `VisibilityRule`.
	//
	// optional google.api.VisibilityRule message_visibility = 72295727;
	E_MessageVisibility = &file_google_api_visibility_proto_extTypes[3]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// See `VisibilityRule`.
	//
	// optional google.api.VisibilityRule method_visibility = 72295727;
	E_MethodVisibility = &file_google_api_visibility_proto_extTypes[4]
)

// Extension fields to descriptorpb.ServiceOptions.
var (
	// See `VisibilityRule`.
	//
	// optional google.api.VisibilityRule api_visibility = 72295727;
	E_ApiVisibility = &file_google_api_visibility_proto_extTypes[5]
)

var File_google_api_visibility_proto protoreflect.FileDescriptor

var file_google_api_visibility_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73,
	0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x0a, 0x56,
	0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x30, 0x0a, 0x05, 0x72, 0x75, 0x6c,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x22, 0x4e, 0x0a, 0x0e, 0x56,
	0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x75, 0x6c, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x73,
	0x74, 0x72, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x72, 0x65, 0x73, 0x74, 0x72, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x64, 0x0a, 0x0f, 0x65,
	0x6e, 0x75, 0x6d, 0x5f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xaf, 0xca, 0xbc,
	0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x75, 0x6c,
	0x65, 0x52, 0x0e, 0x65, 0x6e, 0x75, 0x6d, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x79, 0x3a, 0x6b, 0x0a, 0x10, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x76, 0x69, 0x73, 0x69, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xaf, 0xca, 0xbc, 0x22, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56,
	0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x3a, 0x67,
	0x0a, 0x10, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xaf, 0xca, 0xbc, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x69, 0x73,
	0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x3a, 0x6d, 0x0a, 0x12, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xaf,
	0xca, 0xbc, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52,
	0x75, 0x6c, 0x65, 0x52, 0x11, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x56, 0x69, 0x73, 0x69,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x3a, 0x6a, 0x0a, 0x11, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x5f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x1e, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xaf, 0xca, 0xbc, 0x22,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x75, 0x6c, 0x65,
	0x52, 0x10, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x56, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x3a, 0x65, 0x0a, 0x0e, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xaf, 0xca, 0xbc, 0x22, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x69, 0x73, 0x69,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x0d, 0x61, 0x70, 0x69, 0x56,
	0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x42, 0xae, 0x01, 0x0a, 0x0e, 0x63, 0x6f,
	0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x70, 0x69, 0x42, 0x0f, 0x56, 0x69,
	0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x3f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f,
	0x72, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x69, 0x73, 0x69, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x3b, 0x76, 0x69, 0x73, 0x69, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x47, 0x41, 0x58, 0xaa, 0x02, 0x0a, 0x47, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x41, 0x70, 0x69, 0xca, 0x02, 0x0a, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x5c,
	0x41, 0x70, 0x69, 0xe2, 0x02, 0x16, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x5c, 0x41, 0x70, 0x69,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x47,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x3a, 0x3a, 0x41, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_google_api_visibility_proto_rawDescOnce sync.Once
	file_google_api_visibility_proto_rawDescData = file_google_api_visibility_proto_rawDesc
)

func file_google_api_visibility_proto_rawDescGZIP() []byte {
	file_google_api_visibility_proto_rawDescOnce.Do(func() {
		file_google_api_visibility_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_api_visibility_proto_rawDescData)
	})
	return file_google_api_visibility_proto_rawDescData
}

var file_google_api_visibility_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_google_api_visibility_proto_goTypes = []any{
	(*Visibility)(nil),                    // 0: google.api.Visibility
	(*VisibilityRule)(nil),                // 1: google.api.VisibilityRule
	(*descriptorpb.EnumOptions)(nil),      // 2: google.protobuf.EnumOptions
	(*descriptorpb.EnumValueOptions)(nil), // 3: google.protobuf.EnumValueOptions
	(*descriptorpb.FieldOptions)(nil),     // 4: google.protobuf.FieldOptions
	(*descriptorpb.MessageOptions)(nil),   // 5: google.protobuf.MessageOptions
	(*descriptorpb.MethodOptions)(nil),    // 6: google.protobuf.MethodOptions
	(*descriptorpb.ServiceOptions)(nil),   // 7: google.protobuf.ServiceOptions
}
var file_google_api_visibility_proto_depIdxs = []int32{
	1,  // 0: google.api.Visibility.rules:type_name -> google.api.VisibilityRule
	2,  // 1: google.api.enum_visibility:extendee -> google.protobuf.EnumOptions
	3,  // 2: google.api.value_visibility:extendee -> google.protobuf.EnumValueOptions
	4,  // 3: google.api.field_visibility:extendee -> google.protobuf.FieldOptions
	5,  // 4: google.api.message_visibility:extendee -> google.protobuf.MessageOptions
	6,  // 5: google.api.method_visibility:extendee -> google.protobuf.MethodOptions
	7,  // 6: google.api.api_visibility:extendee -> google.protobuf.ServiceOptions
	1,  // 7: google.api.enum_visibility:type_name -> google.api.VisibilityRule
	1,  // 8: google.api.value_visibility:type_name -> google.api.VisibilityRule
	1,  // 9: google.api.field_visibility:type_name -> google.api.VisibilityRule
	1,  // 10: google.api.message_visibility:type_name -> google.api.VisibilityRule
	1,  // 11: google.api.method_visibility:type_name -> google.api.VisibilityRule
	1,  // 12: google.api.api_visibility:type_name -> google.api.VisibilityRule
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	7,  // [7:13] is the sub-list for extension type_name
	1,  // [1:7] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_google_api_visibility_proto_init() }
func file_google_api_visibility_proto_init() {
	if File_google_api_visibility_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_api_visibility_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 6,
			NumServices:   0,
		},
		GoTypes:           file_google_api_visibility_proto_goTypes,
		DependencyIndexes: file_google_api_visibility_proto_depIdxs,
		MessageInfos:      file_google_api_visibility_proto_msgTypes,
		ExtensionInfos:    file_google_api_visibility_proto_extTypes,
	}.Build()
	File_google_api_visibility_proto = out.File
	file_google_api_visibility_proto_rawDesc = nil
	file_google_api_visibility_proto_goTypes = nil
	file_google_api_visibility_proto_depIdxs = nil
}
