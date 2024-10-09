// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: github.com/dogmatiq/enginekit/protobuf/configpb/config.proto

package configpb

import (
	identitypb "github.com/dogmatiq/enginekit/protobuf/identitypb"
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

// MessageKind is an enumeration of the kinds of message, represented by the
// [dogma.Command], [dogma.Event], and [dogma.Timeout] interfaces.
type MessageKind int32

const (
	MessageKind_UNKNOWN_MESSAGE_KIND MessageKind = 0
	MessageKind_COMMAND              MessageKind = 1
	MessageKind_EVENT                MessageKind = 2
	MessageKind_TIMEOUT              MessageKind = 3
)

// Enum value maps for MessageKind.
var (
	MessageKind_name = map[int32]string{
		0: "UNKNOWN_MESSAGE_KIND",
		1: "COMMAND",
		2: "EVENT",
		3: "TIMEOUT",
	}
	MessageKind_value = map[string]int32{
		"UNKNOWN_MESSAGE_KIND": 0,
		"COMMAND":              1,
		"EVENT":                2,
		"TIMEOUT":              3,
	}
)

func (x MessageKind) Enum() *MessageKind {
	p := new(MessageKind)
	*p = x
	return p
}

func (x MessageKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageKind) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_enumTypes[0].Descriptor()
}

func (MessageKind) Type() protoreflect.EnumType {
	return &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_enumTypes[0]
}

func (x MessageKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageKind.Descriptor instead.
func (MessageKind) EnumDescriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescGZIP(), []int{0}
}

// HandlerType is an enumeration of the types of handlers that an application
// can contain.
type HandlerType int32

const (
	HandlerType_UNKNOWN_HANDLER_TYPE HandlerType = 0
	HandlerType_AGGREGATE            HandlerType = 1
	HandlerType_PROCESS              HandlerType = 2
	HandlerType_INTEGRATION          HandlerType = 3
	HandlerType_PROJECTION           HandlerType = 4
)

// Enum value maps for HandlerType.
var (
	HandlerType_name = map[int32]string{
		0: "UNKNOWN_HANDLER_TYPE",
		1: "AGGREGATE",
		2: "PROCESS",
		3: "INTEGRATION",
		4: "PROJECTION",
	}
	HandlerType_value = map[string]int32{
		"UNKNOWN_HANDLER_TYPE": 0,
		"AGGREGATE":            1,
		"PROCESS":              2,
		"INTEGRATION":          3,
		"PROJECTION":           4,
	}
)

func (x HandlerType) Enum() *HandlerType {
	p := new(HandlerType)
	*p = x
	return p
}

func (x HandlerType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HandlerType) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_enumTypes[1].Descriptor()
}

func (HandlerType) Type() protoreflect.EnumType {
	return &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_enumTypes[1]
}

func (x HandlerType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HandlerType.Descriptor instead.
func (HandlerType) EnumDescriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescGZIP(), []int{1}
}

// Application represents a Dogma application hosted by the engine on the
// server.
type Application struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Identity is the application's identity.
	Identity *identitypb.Identity `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	// GoType is the fully-qualified name of the Go type that provides as the
	// application's implementation.
	GoType string `protobuf:"bytes,2,opt,name=go_type,json=goType,proto3" json:"go_type,omitempty"`
	// Handlers is the set of handlers within the application.
	Handlers []*Handler `protobuf:"bytes,3,rep,name=handlers,proto3" json:"handlers,omitempty"`
	// MessageKinds is a map of each message type's fully-qualified Go type to its
	// the kind of message it implemented by that type.
	Messages map[string]MessageKind `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3,enum=dogma.protobuf.MessageKind"`
}

func (x *Application) Reset() {
	*x = Application{}
	mi := &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Application) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Application) ProtoMessage() {}

func (x *Application) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Application.ProtoReflect.Descriptor instead.
func (*Application) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescGZIP(), []int{0}
}

func (x *Application) GetIdentity() *identitypb.Identity {
	if x != nil {
		return x.Identity
	}
	return nil
}

func (x *Application) GetGoType() string {
	if x != nil {
		return x.GoType
	}
	return ""
}

func (x *Application) GetHandlers() []*Handler {
	if x != nil {
		return x.Handlers
	}
	return nil
}

func (x *Application) GetMessages() map[string]MessageKind {
	if x != nil {
		return x.Messages
	}
	return nil
}

// Handler is a message handler within an application.
type Handler struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Identity is the handler's identity.
	Identity *identitypb.Identity `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	// GoType is the fully-qualified name of the Go type that provides as the
	// handler's implementation.
	GoType string `protobuf:"bytes,2,opt,name=go_type,json=goType,proto3" json:"go_type,omitempty"`
	// Type is the handler's type.
	Type HandlerType `protobuf:"varint,3,opt,name=type,proto3,enum=dogma.protobuf.HandlerType" json:"type,omitempty"`
	// messages is the set of messages produced and consumed by this handler.
	//
	// The keys are the fully-qualified names of the message's Go type.
	Messages map[string]*MessageUsage `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// IsDisabled indicates whether the handler is disabled.
	IsDisabled bool `protobuf:"varint,5,opt,name=is_disabled,json=isDisabled,proto3" json:"is_disabled,omitempty"`
}

func (x *Handler) Reset() {
	*x = Handler{}
	mi := &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Handler) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Handler) ProtoMessage() {}

func (x *Handler) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Handler.ProtoReflect.Descriptor instead.
func (*Handler) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescGZIP(), []int{1}
}

func (x *Handler) GetIdentity() *identitypb.Identity {
	if x != nil {
		return x.Identity
	}
	return nil
}

func (x *Handler) GetGoType() string {
	if x != nil {
		return x.GoType
	}
	return ""
}

func (x *Handler) GetType() HandlerType {
	if x != nil {
		return x.Type
	}
	return HandlerType_UNKNOWN_HANDLER_TYPE
}

func (x *Handler) GetMessages() map[string]*MessageUsage {
	if x != nil {
		return x.Messages
	}
	return nil
}

func (x *Handler) GetIsDisabled() bool {
	if x != nil {
		return x.IsDisabled
	}
	return false
}

type MessageUsage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// IsProduced indicates whether the message is produced by the handler.
	IsProduced bool `protobuf:"varint,2,opt,name=is_produced,json=isProduced,proto3" json:"is_produced,omitempty"`
	// IsConsumed indicates whether the message is consumed by the handler.
	IsConsumed bool `protobuf:"varint,3,opt,name=is_consumed,json=isConsumed,proto3" json:"is_consumed,omitempty"`
}

func (x *MessageUsage) Reset() {
	*x = MessageUsage{}
	mi := &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageUsage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageUsage) ProtoMessage() {}

func (x *MessageUsage) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageUsage.ProtoReflect.Descriptor instead.
func (*MessageUsage) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescGZIP(), []int{2}
}

func (x *MessageUsage) GetIsProduced() bool {
	if x != nil {
		return x.IsProduced
	}
	return false
}

func (x *MessageUsage) GetIsConsumed() bool {
	if x != nil {
		return x.IsConsumed
	}
	return false
}

var File_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto protoreflect.FileDescriptor

var file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70,
	0x62, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e,
	0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x40,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61,
	0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x70,
	0x62, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xb2, 0x02, 0x0a, 0x0b, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x34, 0x0a, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x08, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x6f, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x33, 0x0a, 0x08, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x52, 0x08, 0x68, 0x61, 0x6e, 0x64,
	0x6c, 0x65, 0x72, 0x73, 0x12, 0x45, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x1a, 0x58, 0x0a, 0x0d, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x31,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e,
	0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xc8, 0x02, 0x0a, 0x07, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65,
	0x72, 0x12, 0x34, 0x0a, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x08, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x67, 0x6f, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67, 0x6f, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x2f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b,
	0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x41, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x64, 0x69, 0x73, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x44, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x64, 0x1a, 0x59, 0x0a, 0x0d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x50, 0x0a, 0x0c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d,
	0x65, 0x64, 0x2a, 0x4c, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4b, 0x69, 0x6e,
	0x64, 0x12, 0x18, 0x0a, 0x14, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x4d, 0x45, 0x53,
	0x53, 0x41, 0x47, 0x45, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43,
	0x4f, 0x4d, 0x4d, 0x41, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x56, 0x45, 0x4e,
	0x54, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x49, 0x4d, 0x45, 0x4f, 0x55, 0x54, 0x10, 0x03,
	0x2a, 0x64, 0x0a, 0x0b, 0x48, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x18, 0x0a, 0x14, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x48, 0x41, 0x4e, 0x44, 0x4c,
	0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x47, 0x47,
	0x52, 0x45, 0x47, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x52, 0x4f, 0x43,
	0x45, 0x53, 0x53, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x54, 0x45, 0x47, 0x52, 0x41,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x52, 0x4f, 0x4a, 0x45, 0x43,
	0x54, 0x49, 0x4f, 0x4e, 0x10, 0x04, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescData = file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDesc
)

func file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescData)
	})
	return file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_goTypes = []any{
	(MessageKind)(0),            // 0: dogma.protobuf.MessageKind
	(HandlerType)(0),            // 1: dogma.protobuf.HandlerType
	(*Application)(nil),         // 2: dogma.protobuf.Application
	(*Handler)(nil),             // 3: dogma.protobuf.Handler
	(*MessageUsage)(nil),        // 4: dogma.protobuf.MessageUsage
	nil,                         // 5: dogma.protobuf.Application.MessagesEntry
	nil,                         // 6: dogma.protobuf.Handler.MessagesEntry
	(*identitypb.Identity)(nil), // 7: dogma.protobuf.Identity
}
var file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_depIdxs = []int32{
	7, // 0: dogma.protobuf.Application.identity:type_name -> dogma.protobuf.Identity
	3, // 1: dogma.protobuf.Application.handlers:type_name -> dogma.protobuf.Handler
	5, // 2: dogma.protobuf.Application.messages:type_name -> dogma.protobuf.Application.MessagesEntry
	7, // 3: dogma.protobuf.Handler.identity:type_name -> dogma.protobuf.Identity
	1, // 4: dogma.protobuf.Handler.type:type_name -> dogma.protobuf.HandlerType
	6, // 5: dogma.protobuf.Handler.messages:type_name -> dogma.protobuf.Handler.MessagesEntry
	0, // 6: dogma.protobuf.Application.MessagesEntry.value:type_name -> dogma.protobuf.MessageKind
	4, // 7: dogma.protobuf.Handler.MessagesEntry.value:type_name -> dogma.protobuf.MessageUsage
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_init() }
func file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_init() {
	if File_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_depIdxs,
		EnumInfos:         file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_enumTypes,
		MessageInfos:      file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto = out.File
	file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_rawDesc = nil
	file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_protobuf_configpb_config_proto_depIdxs = nil
}
