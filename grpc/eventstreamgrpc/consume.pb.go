// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto

package eventstreamgrpc

import (
	envelopepb "github.com/dogmatiq/enginekit/protobuf/envelopepb"
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
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

// ListStreamsRequest is the input to the ConsumeAPI.ListStreams method.
type ListStreamsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListStreamsRequest) Reset() {
	*x = ListStreamsRequest{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListStreamsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStreamsRequest) ProtoMessage() {}

func (x *ListStreamsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStreamsRequest.ProtoReflect.Descriptor instead.
func (*ListStreamsRequest) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{0}
}

// ListStreamsResponse is the output of the ConsumeAPI.ListStreams method.
type ListStreamsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Streams is a list of event streams that can be consumed from this server.
	Streams       []*Stream `protobuf:"bytes,1,rep,name=streams,proto3" json:"streams,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListStreamsResponse) Reset() {
	*x = ListStreamsResponse{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListStreamsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStreamsResponse) ProtoMessage() {}

func (x *ListStreamsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStreamsResponse.ProtoReflect.Descriptor instead.
func (*ListStreamsResponse) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{1}
}

func (x *ListStreamsResponse) GetStreams() []*Stream {
	if x != nil {
		return x.Streams
	}
	return nil
}

// Stream describes an offset-based ordered event stream.
type Stream struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// StreamId is a unique identifier for the stream.
	StreamId *uuidpb.UUID `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	// EventTypes is the set of event types that may appear on the stream.
	EventTypes    []*EventType `protobuf:"bytes,2,rep,name=event_types,json=eventTypes,proto3" json:"event_types,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Stream) Reset() {
	*x = Stream{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Stream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stream) ProtoMessage() {}

func (x *Stream) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[2]
	if x != nil {
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

// ConsumeEventsRequest is the input to the ConsumeAPI.ConsumeEvents method.
type ConsumeEventsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// StreamId is the ID from which events are consumed.
	StreamId *uuidpb.UUID `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	// Offset is the offset of the earliest event to be consumed.
	Offset uint64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	// EventTypes is a list of event types to be consumed. The consumer must be
	// explicit about the event types that it understands; there is no mechanism
	// to request all event types.
	EventTypes    []*EventType `protobuf:"bytes,3,rep,name=event_types,json=eventTypes,proto3" json:"event_types,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConsumeEventsRequest) Reset() {
	*x = ConsumeEventsRequest{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConsumeEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeEventsRequest) ProtoMessage() {}

func (x *ConsumeEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeEventsRequest.ProtoReflect.Descriptor instead.
func (*ConsumeEventsRequest) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{3}
}

func (x *ConsumeEventsRequest) GetStreamId() *uuidpb.UUID {
	if x != nil {
		return x.StreamId
	}
	return nil
}

func (x *ConsumeEventsRequest) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ConsumeEventsRequest) GetEventTypes() []*EventType {
	if x != nil {
		return x.EventTypes
	}
	return nil
}

// ConsumeResponse is the (streaming) output of the ConsumeAPI.ConsumeEvents
// method.
type ConsumeEventsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Operation:
	//
	//	*ConsumeEventsResponse_EventDelivery_
	Operation     isConsumeEventsResponse_Operation `protobuf_oneof:"operation"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConsumeEventsResponse) Reset() {
	*x = ConsumeEventsResponse{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConsumeEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeEventsResponse) ProtoMessage() {}

func (x *ConsumeEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeEventsResponse.ProtoReflect.Descriptor instead.
func (*ConsumeEventsResponse) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{4}
}

func (x *ConsumeEventsResponse) GetOperation() isConsumeEventsResponse_Operation {
	if x != nil {
		return x.Operation
	}
	return nil
}

func (x *ConsumeEventsResponse) GetEventDelivery() *ConsumeEventsResponse_EventDelivery {
	if x != nil {
		if x, ok := x.Operation.(*ConsumeEventsResponse_EventDelivery_); ok {
			return x.EventDelivery
		}
	}
	return nil
}

type isConsumeEventsResponse_Operation interface {
	isConsumeEventsResponse_Operation()
}

type ConsumeEventsResponse_EventDelivery_ struct {
	EventDelivery *ConsumeEventsResponse_EventDelivery `protobuf:"bytes,1,opt,name=event_delivery,json=eventDelivery,proto3,oneof"`
}

func (*ConsumeEventsResponse_EventDelivery_) isConsumeEventsResponse_Operation() {}

// EventType describes a type of event that may appear on a stream.
type EventType struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// PortableName is a name that uniquely identifies the event type across
	// process boundaries.
	PortableName string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
	// MediaTypes is the set of supported media-types that can be used to
	// represent events of this type, in order of preference.
	MediaTypes    []string `protobuf:"bytes,2,rep,name=media_types,json=mediaTypes,proto3" json:"media_types,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EventType) Reset() {
	*x = EventType{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EventType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventType) ProtoMessage() {}

func (x *EventType) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[5]
	if x != nil {
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// ApplicationKey is the ID of the unrecognized stream.
	StreamId      *uuidpb.UUID `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnrecognizedStream) Reset() {
	*x = UnrecognizedStream{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnrecognizedStream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnrecognizedStream) ProtoMessage() {}

func (x *UnrecognizedStream) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[6]
	if x != nil {
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

// NoEventTypes is an error-details value for INVALID_ARGUMENT errors that
// occurred because a client sent a consume request without specifying any event
// types.
type NoEventTypes struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoEventTypes) Reset() {
	*x = NoEventTypes{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoEventTypes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoEventTypes) ProtoMessage() {}

func (x *NoEventTypes) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoEventTypes.ProtoReflect.Descriptor instead.
func (*NoEventTypes) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{7}
}

// UnrecognizedEventType is an error-details value for INVALID_ARGUMENT errors
// that occurred because a specific event type was not recognized by the
// server.
type UnrecognizedEventType struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// PortableName is a name that uniquely identifies the event type across
	// process boundaries.
	PortableName  string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UnrecognizedEventType) Reset() {
	*x = UnrecognizedEventType{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UnrecognizedEventType) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnrecognizedEventType) ProtoMessage() {}

func (x *UnrecognizedEventType) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[8]
	if x != nil {
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
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{8}
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
	state protoimpl.MessageState `protogen:"open.v1"`
	// PortableName is a name that uniquely identifies the event type across
	// process boundaries.
	PortableName  string `protobuf:"bytes,1,opt,name=portable_name,json=portableName,proto3" json:"portable_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoRecognizedMediaTypes) Reset() {
	*x = NoRecognizedMediaTypes{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoRecognizedMediaTypes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoRecognizedMediaTypes) ProtoMessage() {}

func (x *NoRecognizedMediaTypes) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[9]
	if x != nil {
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
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{9}
}

func (x *NoRecognizedMediaTypes) GetPortableName() string {
	if x != nil {
		return x.PortableName
	}
	return ""
}

// EventDelivery represents the delivery of a single event to the consumer.
type ConsumeEventsResponse_EventDelivery struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Offset is the event's offset within the stream.
	Offset uint64 `protobuf:"varint,1,opt,name=offset,proto3" json:"offset,omitempty"`
	// Envelope is the envelope containing the event.
	Envelope      *envelopepb.Envelope `protobuf:"bytes,2,opt,name=envelope,proto3" json:"envelope,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConsumeEventsResponse_EventDelivery) Reset() {
	*x = ConsumeEventsResponse_EventDelivery{}
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConsumeEventsResponse_EventDelivery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsumeEventsResponse_EventDelivery) ProtoMessage() {}

func (x *ConsumeEventsResponse_EventDelivery) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsumeEventsResponse_EventDelivery.ProtoReflect.Descriptor instead.
func (*ConsumeEventsResponse_EventDelivery) Descriptor() ([]byte, []int) {
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP(), []int{4, 0}
}

func (x *ConsumeEventsResponse_EventDelivery) GetOffset() uint64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ConsumeEventsResponse_EventDelivery) GetEnvelope() *envelopepb.Envelope {
	if x != nil {
		return x.Envelope
	}
	return nil
}

var File_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto protoreflect.FileDescriptor

const file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc = "" +
	"\n" +
	"@github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto\x12\x1cdogma.eventstream.consume.v1\x1a@github.com/dogmatiq/enginekit/protobuf/envelopepb/envelope.proto\x1a8github.com/dogmatiq/enginekit/protobuf/uuidpb/uuid.proto\"\x14\n" +
	"\x12ListStreamsRequest\"U\n" +
	"\x13ListStreamsResponse\x12>\n" +
	"\astreams\x18\x01 \x03(\v2$.dogma.eventstream.consume.v1.StreamR\astreams\"\x85\x01\n" +
	"\x06Stream\x121\n" +
	"\tstream_id\x18\x01 \x01(\v2\x14.dogma.protobuf.UUIDR\bstreamId\x12H\n" +
	"\vevent_types\x18\x02 \x03(\v2'.dogma.eventstream.consume.v1.EventTypeR\n" +
	"eventTypes\"\xab\x01\n" +
	"\x14ConsumeEventsRequest\x121\n" +
	"\tstream_id\x18\x01 \x01(\v2\x14.dogma.protobuf.UUIDR\bstreamId\x12\x16\n" +
	"\x06offset\x18\x02 \x01(\x04R\x06offset\x12H\n" +
	"\vevent_types\x18\x03 \x03(\v2'.dogma.eventstream.consume.v1.EventTypeR\n" +
	"eventTypes\"\xef\x01\n" +
	"\x15ConsumeEventsResponse\x12j\n" +
	"\x0eevent_delivery\x18\x01 \x01(\v2A.dogma.eventstream.consume.v1.ConsumeEventsResponse.EventDeliveryH\x00R\reventDelivery\x1a]\n" +
	"\rEventDelivery\x12\x16\n" +
	"\x06offset\x18\x01 \x01(\x04R\x06offset\x124\n" +
	"\benvelope\x18\x02 \x01(\v2\x18.dogma.protobuf.EnvelopeR\benvelopeB\v\n" +
	"\toperation\"Q\n" +
	"\tEventType\x12#\n" +
	"\rportable_name\x18\x01 \x01(\tR\fportableName\x12\x1f\n" +
	"\vmedia_types\x18\x02 \x03(\tR\n" +
	"mediaTypes\"G\n" +
	"\x12UnrecognizedStream\x121\n" +
	"\tstream_id\x18\x01 \x01(\v2\x14.dogma.protobuf.UUIDR\bstreamId\"\x0e\n" +
	"\fNoEventTypes\"<\n" +
	"\x15UnrecognizedEventType\x12#\n" +
	"\rportable_name\x18\x01 \x01(\tR\fportableName\"=\n" +
	"\x16NoRecognizedMediaTypes\x12#\n" +
	"\rportable_name\x18\x01 \x01(\tR\fportableName2\xfc\x01\n" +
	"\n" +
	"ConsumeAPI\x12r\n" +
	"\vListStreams\x120.dogma.eventstream.consume.v1.ListStreamsRequest\x1a1.dogma.eventstream.consume.v1.ListStreamsResponse\x12z\n" +
	"\rConsumeEvents\x122.dogma.eventstream.consume.v1.ConsumeEventsRequest\x1a3.dogma.eventstream.consume.v1.ConsumeEventsResponse0\x01B4Z2github.com/dogmatiq/enginekit/grpc/eventstreamgrpcb\x06proto3"

var (
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescOnce sync.Once
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData []byte
)

func file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescGZIP() []byte {
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescOnce.Do(func() {
		file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc), len(file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc)))
	})
	return file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDescData
}

var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_goTypes = []any{
	(*ListStreamsRequest)(nil),                  // 0: dogma.eventstream.consume.v1.ListStreamsRequest
	(*ListStreamsResponse)(nil),                 // 1: dogma.eventstream.consume.v1.ListStreamsResponse
	(*Stream)(nil),                              // 2: dogma.eventstream.consume.v1.Stream
	(*ConsumeEventsRequest)(nil),                // 3: dogma.eventstream.consume.v1.ConsumeEventsRequest
	(*ConsumeEventsResponse)(nil),               // 4: dogma.eventstream.consume.v1.ConsumeEventsResponse
	(*EventType)(nil),                           // 5: dogma.eventstream.consume.v1.EventType
	(*UnrecognizedStream)(nil),                  // 6: dogma.eventstream.consume.v1.UnrecognizedStream
	(*NoEventTypes)(nil),                        // 7: dogma.eventstream.consume.v1.NoEventTypes
	(*UnrecognizedEventType)(nil),               // 8: dogma.eventstream.consume.v1.UnrecognizedEventType
	(*NoRecognizedMediaTypes)(nil),              // 9: dogma.eventstream.consume.v1.NoRecognizedMediaTypes
	(*ConsumeEventsResponse_EventDelivery)(nil), // 10: dogma.eventstream.consume.v1.ConsumeEventsResponse.EventDelivery
	(*uuidpb.UUID)(nil),                         // 11: dogma.protobuf.UUID
	(*envelopepb.Envelope)(nil),                 // 12: dogma.protobuf.Envelope
}
var file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_depIdxs = []int32{
	2,  // 0: dogma.eventstream.consume.v1.ListStreamsResponse.streams:type_name -> dogma.eventstream.consume.v1.Stream
	11, // 1: dogma.eventstream.consume.v1.Stream.stream_id:type_name -> dogma.protobuf.UUID
	5,  // 2: dogma.eventstream.consume.v1.Stream.event_types:type_name -> dogma.eventstream.consume.v1.EventType
	11, // 3: dogma.eventstream.consume.v1.ConsumeEventsRequest.stream_id:type_name -> dogma.protobuf.UUID
	5,  // 4: dogma.eventstream.consume.v1.ConsumeEventsRequest.event_types:type_name -> dogma.eventstream.consume.v1.EventType
	10, // 5: dogma.eventstream.consume.v1.ConsumeEventsResponse.event_delivery:type_name -> dogma.eventstream.consume.v1.ConsumeEventsResponse.EventDelivery
	11, // 6: dogma.eventstream.consume.v1.UnrecognizedStream.stream_id:type_name -> dogma.protobuf.UUID
	12, // 7: dogma.eventstream.consume.v1.ConsumeEventsResponse.EventDelivery.envelope:type_name -> dogma.protobuf.Envelope
	0,  // 8: dogma.eventstream.consume.v1.ConsumeAPI.ListStreams:input_type -> dogma.eventstream.consume.v1.ListStreamsRequest
	3,  // 9: dogma.eventstream.consume.v1.ConsumeAPI.ConsumeEvents:input_type -> dogma.eventstream.consume.v1.ConsumeEventsRequest
	1,  // 10: dogma.eventstream.consume.v1.ConsumeAPI.ListStreams:output_type -> dogma.eventstream.consume.v1.ListStreamsResponse
	4,  // 11: dogma.eventstream.consume.v1.ConsumeAPI.ConsumeEvents:output_type -> dogma.eventstream.consume.v1.ConsumeEventsResponse
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_init() }
func file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_init() {
	if File_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto != nil {
		return
	}
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes[4].OneofWrappers = []any{
		(*ConsumeEventsResponse_EventDelivery_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc), len(file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_goTypes,
		DependencyIndexes: file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_depIdxs,
		MessageInfos:      file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_msgTypes,
	}.Build()
	File_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto = out.File
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_goTypes = nil
	file_github_com_dogmatiq_enginekit_grpc_eventstreamgrpc_consume_proto_depIdxs = nil
}
