// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.1
// source: github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto

package eventstreamgrpc

import (
	envelopepb "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
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

type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{0}
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Streams []*Stream `protobuf:"bytes,1,rep,name=streams,proto3" json:"streams,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{1}
}

func (x *ListResponse) GetStreams() []*Stream {
	if x != nil {
		return x.Streams
	}
	return nil
}

type Stream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// StreamId is a unique identifier for the stream.
	StreamId *uuidpb.UUID `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	// EventTypes is the set of event types that may appear on the stream.
	EventTypes []*EventType `protobuf:"bytes,2,rep,name=event_types,json=eventTypes,proto3" json:"event_types,omitempty"`
}

func (x *Stream) Reset() {
	*x = Stream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stream) ProtoMessage() {}

func (x *Stream) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stream.ProtoReflect.Descriptor instead.
func (*Stream) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{2}
}

func (x *Stream) GetStreamId() *uuidpb.UUID {
	if x != nil {
		return x.StreamId
	}
	return nil
}

func (x *Stream) GetEventTypes() []*EventType {
	if x != nil {
		return x.EventTypes
	}
	return nil
}

type ConsumeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// StreamId is the ID from which events are consumed.
	StreamId *uuidpb.UUID `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	// Offset is the offset of the earliest event to be consumed.
	Offset uint64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	// EventTypes is a list of event types to be consumed. Consumers must be
	// explicit about the event types that it understands; there is no mechanism
	// to request all event types.
	EventTypes []*EventType `protobuf:"bytes,3,rep,name=event_types,json=eventTypes,proto3" json:"event_types,omitempty"`
}

func (x *ConsumeRequest) Reset() {
	*x = ConsumeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeRequest) ProtoMessage() {}

func (x *ConsumeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeRequest.ProtoReflect.Descriptor instead.
func (*ConsumeRequest) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{3}
}

func (x *ConsumeRequest) GetStreamId() *uuidpb.UUID {
	if x != nil {
		return x.StreamId
	}
	return nil
}

func (x *ConsumeRequest) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ConsumeRequest) GetEventTypes() []*EventType {
	if x != nil {
		return x.EventTypes
	}
	return nil
}

type ConsumeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Offset is the offset of the event within the stream.
	Offset uint64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	// Event is the envelope containing the event.
	Event *envelopepb.Envelope `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *ConsumeResponse) Reset() {
	*x = ConsumeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsumeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeResponse) ProtoMessage() {}

func (x *ConsumeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeResponse.ProtoReflect.Descriptor instead.
func (*ConsumeResponse) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{4}
}

func (x *ConsumeResponse) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ConsumeResponse) GetEvent() *envelopepb.Envelope {
	if x != nil {
		return x.Event
	}
	return nil
}

type EventType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PortableName is a name that uniquely identifies the event type across
	// process boundaries.
	PortableName string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
	// MediaTypes is the set of supported media-types that can be used to
	// represent events of this type, in order of preference.
	MediaTypes []string `protobuf:"bytes,2,rep,name=media_types,json=mediaTypes,proto3" json:"media_types,omitempty"`
}

func (x *EventType) Reset() {
	*x = EventType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventType) ProtoMessage() {}

func (x *EventType) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventType.ProtoReflect.Descriptor instead.
func (*EventType) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{5}
}

func (x *EventType) GetPortableName() string {
	if x != nil {
		return x.PortableName
	}
	return ""
}

func (x *EventType) GetMediaTypes() []string {
	if x != nil {
		return x.MediaTypes
	}
	return nil
}

// UnrecognizedStream is an error-details value for INVALID_ARGUMENT errors that
// occurred because a consumer requested an unrecognized stream ID.
type UnrecognizedStream struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ApplicationKey is the ID of the unrecognized stream.
	StreamId *uuidpb.UUID `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
}

func (x *UnrecognizedStream) Reset() {
	*x = UnrecognizedStream{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnrecognizedStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnrecognizedStream) ProtoMessage() {}

func (x *UnrecognizedStream) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnrecognizedStream.ProtoReflect.Descriptor instead.
func (*UnrecognizedStream) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{6}
}

func (x *UnrecognizedStream) GetStreamId() *uuidpb.UUID {
	if x != nil {
		return x.StreamId
	}
	return nil
}

// UnrecognizedEventType is an error-details value for INVALID_ARGUMENT errors
// that occurred because a specific event type was not recognized by the
// server.
type UnrecognizedEventType struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PortableName is a name that uniquely identifies the event type across
	// process boundaries.
	PortableName string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
}

func (x *UnrecognizedEventType) Reset() {
	*x = UnrecognizedEventType{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnrecognizedEventType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnrecognizedEventType) ProtoMessage() {}

func (x *UnrecognizedEventType) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnrecognizedEventType.ProtoReflect.Descriptor instead.
func (*UnrecognizedEventType) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{7}
}

func (x *UnrecognizedEventType) GetPortableName() string {
	if x != nil {
		return x.PortableName
	}
	return ""
}

// NoRecognizedMediaTypes is an error-details value for INVALID_ARGUMENT errors
// that occurred because a the server does not support any of the media-types
// requested by the client for a specific event type.
type NoRecognizedMediaTypes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// PortableName is a name that uniquely identifies the event type across
	// process boundaries.
	PortableName string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
}

func (x *NoRecognizedMediaTypes) Reset() {
	*x = NoRecognizedMediaTypes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoRecognizedMediaTypes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoRecognizedMediaTypes) ProtoMessage() {}

func (x *NoRecognizedMediaTypes) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoRecognizedMediaTypes.ProtoReflect.Descriptor instead.
func (*NoRecognizedMediaTypes) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{8}
}

func (x *NoRecognizedMediaTypes) GetPortableName() string {
	if x != nil {
		return x.PortableName
	}
	return ""
}

var File_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto protoreflect.FileDescriptor

var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc = []byte{
	0x0a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1c, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x2e, 0x76, 0x31,
	0x1a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x70, 0x62, 0x2f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64,
	0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69,
	0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x70,
	0x62, 0x2f, 0x75, 0x75, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4e, 0x0a, 0x0c, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64,
	0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x07, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x22, 0x85, 0x01, 0x0a, 0x06,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x31, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d,
	0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52,
	0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x48, 0x0a, 0x0b, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27,
	0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x73, 0x22, 0xa5, 0x01, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d,
	0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52,
	0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66,
	0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65,
	0x74, 0x12, 0x48, 0x0a, 0x0b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75,
	0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73, 0x22, 0x59, 0x0a, 0x0f, 0x43,
	0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06,
	0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x2e, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x52,
	0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x51, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x6f, 0x72, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x6d,
	0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x73, 0x22, 0x47, 0x0a, 0x12, 0x55, 0x6e, 0x72,
	0x65, 0x63, 0x6f, 0x67, 0x6e, 0x69, 0x7a, 0x65, 0x64, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12,
	0x31, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x49, 0x64, 0x22, 0x3c, 0x0a, 0x15, 0x55, 0x6e, 0x72, 0x65, 0x63, 0x6f, 0x67, 0x6e, 0x69, 0x7a,
	0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70,
	0x6f, 0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x3d, 0x0a, 0x16, 0x4e, 0x6f, 0x52, 0x65, 0x63, 0x6f, 0x67, 0x6e, 0x69, 0x7a, 0x65, 0x64,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x6f,
	0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x70, 0x6f, 0x72, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x32,
	0xd5, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x41, 0x50, 0x49, 0x12, 0x5d,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x29, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75,
	0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2a, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x68, 0x0a,
	0x07, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x12, 0x2c, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e,
	0x73, 0x75, 0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x63, 0x6f, 0x6e, 0x73, 0x75,
	0x6d, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x6b, 0x69, 0x74, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData = file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc
)

func file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData)
	})
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_goTypes = []interface{}{
	(*ListRequest)(nil),            // 0: dogma.eventstream.consume.v1.ListRequest
	(*ListResponse)(nil),           // 1: dogma.eventstream.consume.v1.ListResponse
	(*Stream)(nil),                 // 2: dogma.eventstream.consume.v1.Stream
	(*ConsumeRequest)(nil),         // 3: dogma.eventstream.consume.v1.ConsumeRequest
	(*ConsumeResponse)(nil),        // 4: dogma.eventstream.consume.v1.ConsumeResponse
	(*EventType)(nil),              // 5: dogma.eventstream.consume.v1.EventType
	(*UnrecognizedStream)(nil),     // 6: dogma.eventstream.consume.v1.UnrecognizedStream
	(*UnrecognizedEventType)(nil),  // 7: dogma.eventstream.consume.v1.UnrecognizedEventType
	(*NoRecognizedMediaTypes)(nil), // 8: dogma.eventstream.consume.v1.NoRecognizedMediaTypes
	(*uuidpb.UUID)(nil),            // 9: dogma.protobuf.UUID
	(*envelopepb.Envelope)(nil),    // 10: dogma.protobuf.Envelope
}
var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_depIdxs = []int32{
	2,  // 0: dogma.eventstream.consume.v1.ListResponse.streams:type_name -> dogma.eventstream.consume.v1.Stream
	9,  // 1: dogma.eventstream.consume.v1.Stream.stream_id:type_name -> dogma.protobuf.UUID
	5,  // 2: dogma.eventstream.consume.v1.Stream.event_types:type_name -> dogma.eventstream.consume.v1.EventType
	9,  // 3: dogma.eventstream.consume.v1.ConsumeRequest.stream_id:type_name -> dogma.protobuf.UUID
	5,  // 4: dogma.eventstream.consume.v1.ConsumeRequest.event_types:type_name -> dogma.eventstream.consume.v1.EventType
	10, // 5: dogma.eventstream.consume.v1.ConsumeResponse.event:type_name -> dogma.protobuf.Envelope
	9,  // 6: dogma.eventstream.consume.v1.UnrecognizedStream.stream_id:type_name -> dogma.protobuf.UUID
	0,  // 7: dogma.eventstream.consume.v1.ConsumeAPI.List:input_type -> dogma.eventstream.consume.v1.ListRequest
	3,  // 8: dogma.eventstream.consume.v1.ConsumeAPI.Consume:input_type -> dogma.eventstream.consume.v1.ConsumeRequest
	1,  // 9: dogma.eventstream.consume.v1.ConsumeAPI.List:output_type -> dogma.eventstream.consume.v1.ListResponse
	4,  // 10: dogma.eventstream.consume.v1.ConsumeAPI.Consume:output_type -> dogma.eventstream.consume.v1.ConsumeResponse
	9,  // [9:11] is the sub-list for method output_type
	7,  // [7:9] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_init() }
func file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_init() {
	if File_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRequest); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListResponse); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stream); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeRequest); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsumeResponse); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventType); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnrecognizedStream); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnrecognizedEventType); i {
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
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoRecognizedMediaTypes); i {
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
			RawDescriptor: file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto = out.File
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc = nil
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_depIdxs = nil
}
