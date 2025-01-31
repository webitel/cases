// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: messages/history.proto

package messages

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

// ChatMessages dataset
type ChatMessages struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Dataset page of messages.
	Messages []*Message `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	// List of chats mentioned in messages. [FROM]
	Chats []*Chat `protobuf:"bytes,2,rep,name=chats,proto3" json:"chats,omitempty"`
	// List of peers mentioned in messages. [FROM]
	Peers []*Peer `protobuf:"bytes,3,rep,name=peers,proto3" json:"peers,omitempty"`
	// Dataset page number.
	Page int32 `protobuf:"varint,5,opt,name=page,proto3" json:"page,omitempty"`
	// Next page is available ?
	Next          bool `protobuf:"varint,6,opt,name=next,proto3" json:"next,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatMessages) Reset() {
	*x = ChatMessages{}
	mi := &file_messages_history_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessages) ProtoMessage() {}

func (x *ChatMessages) ProtoReflect() protoreflect.Message {
	mi := &file_messages_history_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessages.ProtoReflect.Descriptor instead.
func (*ChatMessages) Descriptor() ([]byte, []int) {
	return file_messages_history_proto_rawDescGZIP(), []int{0}
}

func (x *ChatMessages) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *ChatMessages) GetChats() []*Chat {
	if x != nil {
		return x.Chats
	}
	return nil
}

func (x *ChatMessages) GetPeers() []*Peer {
	if x != nil {
		return x.Peers
	}
	return nil
}

func (x *ChatMessages) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ChatMessages) GetNext() bool {
	if x != nil {
		return x.Next
	}
	return false
}

type ChatMessagesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Offset messages.
	Offset *ChatMessagesRequest_Offset `protobuf:"bytes,1,opt,name=offset,proto3" json:"offset,omitempty"`
	// Number of messages to return.
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// Search term: message.text
	Q string `protobuf:"bytes,5,opt,name=q,proto3" json:"q,omitempty"`
	// Fields to return into result.
	Fields []string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty"`
	// Input peer identity
	//
	// Types that are valid to be assigned to Chat:
	//
	//	*ChatMessagesRequest_ChatId
	//	*ChatMessagesRequest_Peer
	Chat isChatMessagesRequest_Chat `protobuf_oneof:"chat"`
	// Includes the history of ONLY those dialogs
	// whose member channel(s) contain
	// a specified set of variables.
	Group         map[string]string `protobuf:"bytes,8,rep,name=group,proto3" json:"group,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatMessagesRequest) Reset() {
	*x = ChatMessagesRequest{}
	mi := &file_messages_history_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessagesRequest) ProtoMessage() {}

func (x *ChatMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_history_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessagesRequest.ProtoReflect.Descriptor instead.
func (*ChatMessagesRequest) Descriptor() ([]byte, []int) {
	return file_messages_history_proto_rawDescGZIP(), []int{1}
}

func (x *ChatMessagesRequest) GetOffset() *ChatMessagesRequest_Offset {
	if x != nil {
		return x.Offset
	}
	return nil
}

func (x *ChatMessagesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ChatMessagesRequest) GetQ() string {
	if x != nil {
		return x.Q
	}
	return ""
}

func (x *ChatMessagesRequest) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *ChatMessagesRequest) GetChat() isChatMessagesRequest_Chat {
	if x != nil {
		return x.Chat
	}
	return nil
}

func (x *ChatMessagesRequest) GetChatId() string {
	if x != nil {
		if x, ok := x.Chat.(*ChatMessagesRequest_ChatId); ok {
			return x.ChatId
		}
	}
	return ""
}

func (x *ChatMessagesRequest) GetPeer() *Peer {
	if x != nil {
		if x, ok := x.Chat.(*ChatMessagesRequest_Peer); ok {
			return x.Peer
		}
	}
	return nil
}

func (x *ChatMessagesRequest) GetGroup() map[string]string {
	if x != nil {
		return x.Group
	}
	return nil
}

type isChatMessagesRequest_Chat interface {
	isChatMessagesRequest_Chat()
}

type ChatMessagesRequest_ChatId struct {
	// Unique chat dialog
	ChatId string `protobuf:"bytes,6,opt,name=chat_id,json=chatId,proto3,oneof"`
}

type ChatMessagesRequest_Peer struct {
	// Unique peer contact
	Peer *Peer `protobuf:"bytes,7,opt,name=peer,proto3,oneof"`
}

func (*ChatMessagesRequest_ChatId) isChatMessagesRequest_Chat() {}

func (*ChatMessagesRequest_Peer) isChatMessagesRequest_Chat() {}

// Offset options
type ChatMessagesRequest_Offset struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Messages ONLY starting from the specified message ID
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Messages ONLY been sent before the specified epochtime(milli).
	Date          int64 `protobuf:"varint,2,opt,name=date,proto3" json:"date,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatMessagesRequest_Offset) Reset() {
	*x = ChatMessagesRequest_Offset{}
	mi := &file_messages_history_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessagesRequest_Offset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessagesRequest_Offset) ProtoMessage() {}

func (x *ChatMessagesRequest_Offset) ProtoReflect() protoreflect.Message {
	mi := &file_messages_history_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessagesRequest_Offset.ProtoReflect.Descriptor instead.
func (*ChatMessagesRequest_Offset) Descriptor() ([]byte, []int) {
	return file_messages_history_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ChatMessagesRequest_Offset) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ChatMessagesRequest_Offset) GetDate() int64 {
	if x != nil {
		return x.Date
	}
	return 0
}

var File_messages_history_proto protoreflect.FileDescriptor

var file_messages_history_proto_rawDesc = string([]byte{
	0x0a, 0x16, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x68, 0x69, 0x73, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x1a, 0x13, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2f, 0x70, 0x65, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x16, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbd, 0x01, 0x0a, 0x0c, 0x43, 0x68, 0x61,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x08, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x05,
	0x63, 0x68, 0x61, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x52,
	0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x04, 0x6e, 0x65, 0x78, 0x74, 0x22, 0x8c, 0x03, 0x0a, 0x13, 0x43, 0x68, 0x61,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x40, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x28, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x71, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x01, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x19,
	0x0a, 0x07, 0x63, 0x68, 0x61, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x06, 0x63, 0x68, 0x61, 0x74, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x04, 0x70, 0x65, 0x65,
	0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65,
	0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x48, 0x00, 0x52, 0x04, 0x70,
	0x65, 0x65, 0x72, 0x12, 0x42, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x08, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x68, 0x61,
	0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x1a, 0x2c, 0x0a, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x38, 0x0a, 0x0a, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42,
	0x06, 0x0a, 0x04, 0x63, 0x68, 0x61, 0x74, 0x42, 0x9a, 0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e,
	0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x42, 0x0c, 0x48, 0x69,
	0x73, 0x74, 0x6f, 0x72, 0x79, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x27, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x77, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0xa2, 0x02, 0x03, 0x57, 0x43, 0x58, 0xaa, 0x02, 0x0c, 0x57, 0x65,
	0x62, 0x69, 0x74, 0x65, 0x6c, 0x2e, 0x43, 0x68, 0x61, 0x74, 0xca, 0x02, 0x0c, 0x57, 0x65, 0x62,
	0x69, 0x74, 0x65, 0x6c, 0x5c, 0x43, 0x68, 0x61, 0x74, 0xe2, 0x02, 0x18, 0x57, 0x65, 0x62, 0x69,
	0x74, 0x65, 0x6c, 0x5c, 0x43, 0x68, 0x61, 0x74, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x57, 0x65, 0x62, 0x69, 0x74, 0x65, 0x6c, 0x3a, 0x3a,
	0x43, 0x68, 0x61, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_messages_history_proto_rawDescOnce sync.Once
	file_messages_history_proto_rawDescData []byte
)

func file_messages_history_proto_rawDescGZIP() []byte {
	file_messages_history_proto_rawDescOnce.Do(func() {
		file_messages_history_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_messages_history_proto_rawDesc), len(file_messages_history_proto_rawDesc)))
	})
	return file_messages_history_proto_rawDescData
}

var file_messages_history_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_messages_history_proto_goTypes = []any{
	(*ChatMessages)(nil),               // 0: webitel.chat.ChatMessages
	(*ChatMessagesRequest)(nil),        // 1: webitel.chat.ChatMessagesRequest
	(*ChatMessagesRequest_Offset)(nil), // 2: webitel.chat.ChatMessagesRequest.Offset
	nil,                                // 3: webitel.chat.ChatMessagesRequest.GroupEntry
	(*Message)(nil),                    // 4: webitel.chat.Message
	(*Chat)(nil),                       // 5: webitel.chat.Chat
	(*Peer)(nil),                       // 6: webitel.chat.Peer
}
var file_messages_history_proto_depIdxs = []int32{
	4, // 0: webitel.chat.ChatMessages.messages:type_name -> webitel.chat.Message
	5, // 1: webitel.chat.ChatMessages.chats:type_name -> webitel.chat.Chat
	6, // 2: webitel.chat.ChatMessages.peers:type_name -> webitel.chat.Peer
	2, // 3: webitel.chat.ChatMessagesRequest.offset:type_name -> webitel.chat.ChatMessagesRequest.Offset
	6, // 4: webitel.chat.ChatMessagesRequest.peer:type_name -> webitel.chat.Peer
	3, // 5: webitel.chat.ChatMessagesRequest.group:type_name -> webitel.chat.ChatMessagesRequest.GroupEntry
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_messages_history_proto_init() }
func file_messages_history_proto_init() {
	if File_messages_history_proto != nil {
		return
	}
	file_messages_peer_proto_init()
	file_messages_chat_proto_init()
	file_messages_message_proto_init()
	file_messages_history_proto_msgTypes[1].OneofWrappers = []any{
		(*ChatMessagesRequest_ChatId)(nil),
		(*ChatMessagesRequest_Peer)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_messages_history_proto_rawDesc), len(file_messages_history_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_history_proto_goTypes,
		DependencyIndexes: file_messages_history_proto_depIdxs,
		MessageInfos:      file_messages_history_proto_msgTypes,
	}.Build()
	File_messages_history_proto = out.File
	file_messages_history_proto_goTypes = nil
	file_messages_history_proto_depIdxs = nil
}
