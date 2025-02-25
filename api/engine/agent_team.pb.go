// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: agent_team.proto

package engine

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

type DeleteAgentTeamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId      int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteAgentTeamRequest) Reset() {
	*x = DeleteAgentTeamRequest{}
	mi := &file_agent_team_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteAgentTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAgentTeamRequest) ProtoMessage() {}

func (x *DeleteAgentTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAgentTeamRequest.ProtoReflect.Descriptor instead.
func (*DeleteAgentTeamRequest) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteAgentTeamRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DeleteAgentTeamRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

type UpdateAgentTeamRequest struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Id                  int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description         string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Strategy            string                 `protobuf:"bytes,4,opt,name=strategy,proto3" json:"strategy,omitempty"`
	MaxNoAnswer         int32                  `protobuf:"varint,5,opt,name=max_no_answer,json=maxNoAnswer,proto3" json:"max_no_answer,omitempty"`
	NoAnswerDelayTime   int32                  `protobuf:"varint,6,opt,name=no_answer_delay_time,json=noAnswerDelayTime,proto3" json:"no_answer_delay_time,omitempty"`
	WrapUpTime          int32                  `protobuf:"varint,7,opt,name=wrap_up_time,json=wrapUpTime,proto3" json:"wrap_up_time,omitempty"`
	CallTimeout         int32                  `protobuf:"varint,8,opt,name=call_timeout,json=callTimeout,proto3" json:"call_timeout,omitempty"`
	Admin               []*Lookup              `protobuf:"bytes,9,rep,name=admin,proto3" json:"admin,omitempty"`
	DomainId            int64                  `protobuf:"varint,10,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	InviteChatTimeout   int32                  `protobuf:"varint,11,opt,name=invite_chat_timeout,json=inviteChatTimeout,proto3" json:"invite_chat_timeout,omitempty"`
	TaskAcceptTimeout   int32                  `protobuf:"varint,12,opt,name=task_accept_timeout,json=taskAcceptTimeout,proto3" json:"task_accept_timeout,omitempty"`
	ForecastCalculation *Lookup                `protobuf:"bytes,13,opt,name=forecast_calculation,json=forecastCalculation,proto3" json:"forecast_calculation,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *UpdateAgentTeamRequest) Reset() {
	*x = UpdateAgentTeamRequest{}
	mi := &file_agent_team_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateAgentTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAgentTeamRequest) ProtoMessage() {}

func (x *UpdateAgentTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAgentTeamRequest.ProtoReflect.Descriptor instead.
func (*UpdateAgentTeamRequest) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateAgentTeamRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateAgentTeamRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateAgentTeamRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *UpdateAgentTeamRequest) GetMaxNoAnswer() int32 {
	if x != nil {
		return x.MaxNoAnswer
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetNoAnswerDelayTime() int32 {
	if x != nil {
		return x.NoAnswerDelayTime
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetWrapUpTime() int32 {
	if x != nil {
		return x.WrapUpTime
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetCallTimeout() int32 {
	if x != nil {
		return x.CallTimeout
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetAdmin() []*Lookup {
	if x != nil {
		return x.Admin
	}
	return nil
}

func (x *UpdateAgentTeamRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetInviteChatTimeout() int32 {
	if x != nil {
		return x.InviteChatTimeout
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetTaskAcceptTimeout() int32 {
	if x != nil {
		return x.TaskAcceptTimeout
	}
	return 0
}

func (x *UpdateAgentTeamRequest) GetForecastCalculation() *Lookup {
	if x != nil {
		return x.ForecastCalculation
	}
	return nil
}

type ReadAgentTeamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId      int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReadAgentTeamRequest) Reset() {
	*x = ReadAgentTeamRequest{}
	mi := &file_agent_team_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadAgentTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadAgentTeamRequest) ProtoMessage() {}

func (x *ReadAgentTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadAgentTeamRequest.ProtoReflect.Descriptor instead.
func (*ReadAgentTeamRequest) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{2}
}

func (x *ReadAgentTeamRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ReadAgentTeamRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

type SearchAgentTeamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Size          int32                  `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Q             string                 `protobuf:"bytes,3,opt,name=q,proto3" json:"q,omitempty"`
	Sort          string                 `protobuf:"bytes,5,opt,name=sort,proto3" json:"sort,omitempty"`
	Fields        []string               `protobuf:"bytes,6,rep,name=fields,proto3" json:"fields,omitempty"`
	Id            []uint32               `protobuf:"varint,7,rep,packed,name=id,proto3" json:"id,omitempty"`
	Strategy      []string               `protobuf:"bytes,8,rep,name=strategy,proto3" json:"strategy,omitempty"`
	AdminId       []uint32               `protobuf:"varint,9,rep,packed,name=admin_id,json=adminId,proto3" json:"admin_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchAgentTeamRequest) Reset() {
	*x = SearchAgentTeamRequest{}
	mi := &file_agent_team_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchAgentTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchAgentTeamRequest) ProtoMessage() {}

func (x *SearchAgentTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchAgentTeamRequest.ProtoReflect.Descriptor instead.
func (*SearchAgentTeamRequest) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{3}
}

func (x *SearchAgentTeamRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchAgentTeamRequest) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SearchAgentTeamRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *SearchAgentTeamRequest) GetSort() string {
	if x != nil {
		return x.Sort
	}
	return ""
}

func (x *SearchAgentTeamRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *SearchAgentTeamRequest) GetId() []uint32 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *SearchAgentTeamRequest) GetStrategy() []string {
	if x != nil {
		return x.Strategy
	}
	return nil
}

func (x *SearchAgentTeamRequest) GetAdminId() []uint32 {
	if x != nil {
		return x.AdminId
	}
	return nil
}

type CreateAgentTeamRequest struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Name                string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description         string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Strategy            string                 `protobuf:"bytes,3,opt,name=strategy,proto3" json:"strategy,omitempty"`
	MaxNoAnswer         int32                  `protobuf:"varint,4,opt,name=max_no_answer,json=maxNoAnswer,proto3" json:"max_no_answer,omitempty"`
	NoAnswerDelayTime   int32                  `protobuf:"varint,5,opt,name=no_answer_delay_time,json=noAnswerDelayTime,proto3" json:"no_answer_delay_time,omitempty"`
	WrapUpTime          int32                  `protobuf:"varint,6,opt,name=wrap_up_time,json=wrapUpTime,proto3" json:"wrap_up_time,omitempty"`
	CallTimeout         int32                  `protobuf:"varint,7,opt,name=call_timeout,json=callTimeout,proto3" json:"call_timeout,omitempty"`
	Admin               []*Lookup              `protobuf:"bytes,8,rep,name=admin,proto3" json:"admin,omitempty"`
	DomainId            int64                  `protobuf:"varint,9,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	InviteChatTimeout   int32                  `protobuf:"varint,10,opt,name=invite_chat_timeout,json=inviteChatTimeout,proto3" json:"invite_chat_timeout,omitempty"`
	TaskAcceptTimeout   int32                  `protobuf:"varint,11,opt,name=task_accept_timeout,json=taskAcceptTimeout,proto3" json:"task_accept_timeout,omitempty"`
	ForecastCalculation *Lookup                `protobuf:"bytes,12,opt,name=forecast_calculation,json=forecastCalculation,proto3" json:"forecast_calculation,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *CreateAgentTeamRequest) Reset() {
	*x = CreateAgentTeamRequest{}
	mi := &file_agent_team_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateAgentTeamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAgentTeamRequest) ProtoMessage() {}

func (x *CreateAgentTeamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAgentTeamRequest.ProtoReflect.Descriptor instead.
func (*CreateAgentTeamRequest) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{4}
}

func (x *CreateAgentTeamRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAgentTeamRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateAgentTeamRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *CreateAgentTeamRequest) GetMaxNoAnswer() int32 {
	if x != nil {
		return x.MaxNoAnswer
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetNoAnswerDelayTime() int32 {
	if x != nil {
		return x.NoAnswerDelayTime
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetWrapUpTime() int32 {
	if x != nil {
		return x.WrapUpTime
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetCallTimeout() int32 {
	if x != nil {
		return x.CallTimeout
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetAdmin() []*Lookup {
	if x != nil {
		return x.Admin
	}
	return nil
}

func (x *CreateAgentTeamRequest) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetInviteChatTimeout() int32 {
	if x != nil {
		return x.InviteChatTimeout
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetTaskAcceptTimeout() int32 {
	if x != nil {
		return x.TaskAcceptTimeout
	}
	return 0
}

func (x *CreateAgentTeamRequest) GetForecastCalculation() *Lookup {
	if x != nil {
		return x.ForecastCalculation
	}
	return nil
}

type AgentTeam struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	Id                  int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DomainId            int64                  `protobuf:"varint,2,opt,name=domain_id,json=domainId,proto3" json:"domain_id,omitempty"`
	Name                string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description         string                 `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Strategy            string                 `protobuf:"bytes,5,opt,name=strategy,proto3" json:"strategy,omitempty"`
	MaxNoAnswer         int32                  `protobuf:"varint,6,opt,name=max_no_answer,json=maxNoAnswer,proto3" json:"max_no_answer,omitempty"`
	NoAnswerDelayTime   int32                  `protobuf:"varint,7,opt,name=no_answer_delay_time,json=noAnswerDelayTime,proto3" json:"no_answer_delay_time,omitempty"`
	WrapUpTime          int32                  `protobuf:"varint,8,opt,name=wrap_up_time,json=wrapUpTime,proto3" json:"wrap_up_time,omitempty"`
	CallTimeout         int32                  `protobuf:"varint,9,opt,name=call_timeout,json=callTimeout,proto3" json:"call_timeout,omitempty"`
	UpdatedAt           int64                  `protobuf:"varint,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Admin               []*Lookup              `protobuf:"bytes,11,rep,name=admin,proto3" json:"admin,omitempty"`
	InviteChatTimeout   int32                  `protobuf:"varint,12,opt,name=invite_chat_timeout,json=inviteChatTimeout,proto3" json:"invite_chat_timeout,omitempty"`
	TaskAcceptTimeout   int32                  `protobuf:"varint,13,opt,name=task_accept_timeout,json=taskAcceptTimeout,proto3" json:"task_accept_timeout,omitempty"`
	ForecastCalculation *Lookup                `protobuf:"bytes,14,opt,name=forecast_calculation,json=forecastCalculation,proto3" json:"forecast_calculation,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *AgentTeam) Reset() {
	*x = AgentTeam{}
	mi := &file_agent_team_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AgentTeam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AgentTeam) ProtoMessage() {}

func (x *AgentTeam) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AgentTeam.ProtoReflect.Descriptor instead.
func (*AgentTeam) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{5}
}

func (x *AgentTeam) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AgentTeam) GetDomainId() int64 {
	if x != nil {
		return x.DomainId
	}
	return 0
}

func (x *AgentTeam) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AgentTeam) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *AgentTeam) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *AgentTeam) GetMaxNoAnswer() int32 {
	if x != nil {
		return x.MaxNoAnswer
	}
	return 0
}

func (x *AgentTeam) GetNoAnswerDelayTime() int32 {
	if x != nil {
		return x.NoAnswerDelayTime
	}
	return 0
}

func (x *AgentTeam) GetWrapUpTime() int32 {
	if x != nil {
		return x.WrapUpTime
	}
	return 0
}

func (x *AgentTeam) GetCallTimeout() int32 {
	if x != nil {
		return x.CallTimeout
	}
	return 0
}

func (x *AgentTeam) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *AgentTeam) GetAdmin() []*Lookup {
	if x != nil {
		return x.Admin
	}
	return nil
}

func (x *AgentTeam) GetInviteChatTimeout() int32 {
	if x != nil {
		return x.InviteChatTimeout
	}
	return 0
}

func (x *AgentTeam) GetTaskAcceptTimeout() int32 {
	if x != nil {
		return x.TaskAcceptTimeout
	}
	return 0
}

func (x *AgentTeam) GetForecastCalculation() *Lookup {
	if x != nil {
		return x.ForecastCalculation
	}
	return nil
}

type ListAgentTeam struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Next          bool                   `protobuf:"varint,1,opt,name=next,proto3" json:"next,omitempty"`
	Items         []*AgentTeam           `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListAgentTeam) Reset() {
	*x = ListAgentTeam{}
	mi := &file_agent_team_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListAgentTeam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAgentTeam) ProtoMessage() {}

func (x *ListAgentTeam) ProtoReflect() protoreflect.Message {
	mi := &file_agent_team_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAgentTeam.ProtoReflect.Descriptor instead.
func (*ListAgentTeam) Descriptor() ([]byte, []int) {
	return file_agent_team_proto_rawDescGZIP(), []int{6}
}

func (x *ListAgentTeam) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

func (x *ListAgentTeam) GetItems() []*AgentTeam {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_agent_team_proto protoreflect.FileDescriptor

var file_agent_team_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0b, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41,
	0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x22, 0xfa, 0x03, 0x0a,
	0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x22, 0x0a, 0x0d, 0x6d, 0x61, 0x78,
	0x5f, 0x6e, 0x6f, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0b, 0x6d, 0x61, 0x78, 0x4e, 0x6f, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x2f, 0x0a,
	0x14, 0x6e, 0x6f, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x6e, 0x6f, 0x41,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0c, 0x77, 0x72, 0x61, 0x70, 0x5f, 0x75, 0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x77, 0x72, 0x61, 0x70, 0x55, 0x70, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x12, 0x24, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x09, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b,
	0x75, 0x70, 0x52, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d,
	0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65,
	0x5f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x11, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x43, 0x68, 0x61, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x61,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x11, 0x74, 0x61, 0x73, 0x6b, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x41, 0x0a, 0x14, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61,
	0x73, 0x74, 0x5f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x13, 0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x43, 0x61,
	0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x43, 0x0a, 0x14, 0x52, 0x65, 0x61,
	0x64, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x22, 0xc1,
	0x01, 0x0a, 0x16, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12,
	0x12, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73,
	0x6f, 0x72, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x49, 0x64, 0x22, 0xea, 0x03, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65,
	0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12,
	0x22, 0x0a, 0x0d, 0x6d, 0x61, 0x78, 0x5f, 0x6e, 0x6f, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6d, 0x61, 0x78, 0x4e, 0x6f, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x14, 0x6e, 0x6f, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x5f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x11, 0x6e, 0x6f, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x61, 0x79,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x77, 0x72, 0x61, 0x70, 0x5f, 0x75, 0x70, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x77, 0x72, 0x61, 0x70,
	0x55, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x61,
	0x6c, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x24, 0x0a, 0x05, 0x61, 0x64, 0x6d,
	0x69, 0x6e, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12,
	0x1b, 0x0a, 0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13,
	0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x69, 0x6e, 0x76, 0x69, 0x74,
	0x65, 0x43, 0x68, 0x61, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x2e, 0x0a, 0x13,
	0x74, 0x61, 0x73, 0x6b, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x74, 0x61, 0x73, 0x6b, 0x41,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x41, 0x0a, 0x14,
	0x66, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x5f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x13, 0x66, 0x6f, 0x72, 0x65,
	0x63, 0x61, 0x73, 0x74, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x8c, 0x04, 0x0a, 0x09, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a,
	0x09, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x22, 0x0a, 0x0d,
	0x6d, 0x61, 0x78, 0x5f, 0x6e, 0x6f, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0b, 0x6d, 0x61, 0x78, 0x4e, 0x6f, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x12, 0x2f, 0x0a, 0x14, 0x6e, 0x6f, 0x5f, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x5f, 0x64, 0x65,
	0x6c, 0x61, 0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11,
	0x6e, 0x6f, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0c, 0x77, 0x72, 0x61, 0x70, 0x5f, 0x75, 0x70, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x77, 0x72, 0x61, 0x70, 0x55, 0x70, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x61, 0x6c, 0x6c, 0x54,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x24, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x0b,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x52, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x12, 0x2e, 0x0a, 0x13, 0x69,
	0x6e, 0x76, 0x69, 0x74, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65,
	0x43, 0x68, 0x61, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x74,
	0x61, 0x73, 0x6b, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x74, 0x61, 0x73, 0x6b, 0x41, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x41, 0x0a, 0x14, 0x66,
	0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x5f, 0x63, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x52, 0x13, 0x66, 0x6f, 0x72, 0x65, 0x63,
	0x61, 0x73, 0x74, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4c,
	0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e,
	0x65, 0x78, 0x74, 0x12, 0x27, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0x91, 0x04, 0x0a,
	0x10, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x63, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x54, 0x65, 0x61, 0x6d, 0x12, 0x1e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x3a,
	0x01, 0x2a, 0x22, 0x12, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72,
	0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x64, 0x0a, 0x0f, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x1e, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d,
	0x22, 0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x12, 0x12, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f,
	0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x12, 0x61, 0x0a, 0x0d,
	0x52, 0x65, 0x61, 0x64, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x1c, 0x2e,
	0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x22, 0x1f,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x12, 0x17, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12,
	0x68, 0x0a, 0x0f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65,
	0x61, 0x6d, 0x12, 0x1e, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x54, 0x65, 0x61, 0x6d, 0x22, 0x22, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1c, 0x3a, 0x01, 0x2a,
	0x1a, 0x17, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x74,
	0x65, 0x61, 0x6d, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x65, 0x0a, 0x0f, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x12, 0x1e, 0x2e, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x54, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2e, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x22,
	0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x2a, 0x17, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x5f, 0x63,
	0x65, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d,
	0x42, 0x76, 0x0a, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x42, 0x0e,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x54, 0x65, 0x61, 0x6d, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0xa2, 0x02, 0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0xca, 0x02, 0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0xe2, 0x02, 0x12, 0x45, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x06, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_agent_team_proto_rawDescOnce sync.Once
	file_agent_team_proto_rawDescData []byte
)

func file_agent_team_proto_rawDescGZIP() []byte {
	file_agent_team_proto_rawDescOnce.Do(func() {
		file_agent_team_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_agent_team_proto_rawDesc), len(file_agent_team_proto_rawDesc)))
	})
	return file_agent_team_proto_rawDescData
}

var file_agent_team_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_agent_team_proto_goTypes = []any{
	(*DeleteAgentTeamRequest)(nil), // 0: engine.DeleteAgentTeamRequest
	(*UpdateAgentTeamRequest)(nil), // 1: engine.UpdateAgentTeamRequest
	(*ReadAgentTeamRequest)(nil),   // 2: engine.ReadAgentTeamRequest
	(*SearchAgentTeamRequest)(nil), // 3: engine.SearchAgentTeamRequest
	(*CreateAgentTeamRequest)(nil), // 4: engine.CreateAgentTeamRequest
	(*AgentTeam)(nil),              // 5: engine.AgentTeam
	(*ListAgentTeam)(nil),          // 6: engine.ListAgentTeam
	(*Lookup)(nil),                 // 7: engine.Lookup
}
var file_agent_team_proto_depIdxs = []int32{
	7,  // 0: engine.UpdateAgentTeamRequest.admin:type_name -> engine.Lookup
	7,  // 1: engine.UpdateAgentTeamRequest.forecast_calculation:type_name -> engine.Lookup
	7,  // 2: engine.CreateAgentTeamRequest.admin:type_name -> engine.Lookup
	7,  // 3: engine.CreateAgentTeamRequest.forecast_calculation:type_name -> engine.Lookup
	7,  // 4: engine.AgentTeam.admin:type_name -> engine.Lookup
	7,  // 5: engine.AgentTeam.forecast_calculation:type_name -> engine.Lookup
	5,  // 6: engine.ListAgentTeam.items:type_name -> engine.AgentTeam
	4,  // 7: engine.AgentTeamService.CreateAgentTeam:input_type -> engine.CreateAgentTeamRequest
	3,  // 8: engine.AgentTeamService.SearchAgentTeam:input_type -> engine.SearchAgentTeamRequest
	2,  // 9: engine.AgentTeamService.ReadAgentTeam:input_type -> engine.ReadAgentTeamRequest
	1,  // 10: engine.AgentTeamService.UpdateAgentTeam:input_type -> engine.UpdateAgentTeamRequest
	0,  // 11: engine.AgentTeamService.DeleteAgentTeam:input_type -> engine.DeleteAgentTeamRequest
	5,  // 12: engine.AgentTeamService.CreateAgentTeam:output_type -> engine.AgentTeam
	6,  // 13: engine.AgentTeamService.SearchAgentTeam:output_type -> engine.ListAgentTeam
	5,  // 14: engine.AgentTeamService.ReadAgentTeam:output_type -> engine.AgentTeam
	5,  // 15: engine.AgentTeamService.UpdateAgentTeam:output_type -> engine.AgentTeam
	5,  // 16: engine.AgentTeamService.DeleteAgentTeam:output_type -> engine.AgentTeam
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_agent_team_proto_init() }
func file_agent_team_proto_init() {
	if File_agent_team_proto != nil {
		return
	}
	file_const_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_agent_team_proto_rawDesc), len(file_agent_team_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_agent_team_proto_goTypes,
		DependencyIndexes: file_agent_team_proto_depIdxs,
		MessageInfos:      file_agent_team_proto_msgTypes,
	}.Build()
	File_agent_team_proto = out.File
	file_agent_team_proto_goTypes = nil
	file_agent_team_proto_depIdxs = nil
}
