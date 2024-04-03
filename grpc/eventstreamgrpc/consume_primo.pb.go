// Code generated by protoc-gen-go-primo. DO NOT EDIT.
// versions:
// 	protoc-gen-go-primo v
// 	protoc              v5.26.0
// source: github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto

package eventstreamgrpc

import uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"

type ListStreamsRequestBuilder struct {
	prototype ListStreamsRequest
}

// NewListStreamsRequestBuilder returns a builder that constructs [ListStreamsRequest] messages.
func NewListStreamsRequestBuilder() *ListStreamsRequestBuilder {
	return &ListStreamsRequestBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ListStreamsRequestBuilder) From(x *ListStreamsRequest) *ListStreamsRequestBuilder {
	return b
}

// Build returns a new [ListStreamsRequest] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ListStreamsRequestBuilder) Build() *ListStreamsRequest {
	return &ListStreamsRequest{}
}

type ListStreamsResponseBuilder struct {
	prototype ListStreamsResponse
}

// NewListStreamsResponseBuilder returns a builder that constructs [ListStreamsResponse] messages.
func NewListStreamsResponseBuilder() *ListStreamsResponseBuilder {
	return &ListStreamsResponseBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ListStreamsResponseBuilder) From(x *ListStreamsResponse) *ListStreamsResponseBuilder {
	b.prototype.Streams = x.Streams
	return b
}

// Build returns a new [ListStreamsResponse] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ListStreamsResponseBuilder) Build() *ListStreamsResponse {
	return &ListStreamsResponse{
		Streams: b.prototype.Streams,
	}
}

// WithStreams configures the builder to set the Streams field to v,
// then returns b.
func (b *ListStreamsResponseBuilder) WithStreams(v []*Stream) *ListStreamsResponseBuilder {
	b.prototype.Streams = v
	return b
}

type StreamBuilder struct {
	prototype Stream
}

// NewStreamBuilder returns a builder that constructs [Stream] messages.
func NewStreamBuilder() *StreamBuilder {
	return &StreamBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *StreamBuilder) From(x *Stream) *StreamBuilder {
	b.prototype.StreamId = x.StreamId
	b.prototype.EventTypes = x.EventTypes
	return b
}

// Build returns a new [Stream] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *StreamBuilder) Build() *Stream {
	return &Stream{
		StreamId:   b.prototype.StreamId,
		EventTypes: b.prototype.EventTypes,
	}
}

// WithStreamId configures the builder to set the StreamId field to v,
// then returns b.
func (b *StreamBuilder) WithStreamId(v *uuidpb.UUID) *StreamBuilder {
	b.prototype.StreamId = v
	return b
}

// WithEventTypes configures the builder to set the EventTypes field to v,
// then returns b.
func (b *StreamBuilder) WithEventTypes(v []*EventType) *StreamBuilder {
	b.prototype.EventTypes = v
	return b
}

type ConsumeEventsRequestBuilder struct {
	prototype ConsumeEventsRequest
}

// NewConsumeEventsRequestBuilder returns a builder that constructs [ConsumeEventsRequest] messages.
func NewConsumeEventsRequestBuilder() *ConsumeEventsRequestBuilder {
	return &ConsumeEventsRequestBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ConsumeEventsRequestBuilder) From(x *ConsumeEventsRequest) *ConsumeEventsRequestBuilder {
	b.prototype.StreamId = x.StreamId
	b.prototype.Offset = x.Offset
	b.prototype.EventTypes = x.EventTypes
	return b
}

// Build returns a new [ConsumeEventsRequest] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ConsumeEventsRequestBuilder) Build() *ConsumeEventsRequest {
	return &ConsumeEventsRequest{
		StreamId:   b.prototype.StreamId,
		Offset:     b.prototype.Offset,
		EventTypes: b.prototype.EventTypes,
	}
}

// WithStreamId configures the builder to set the StreamId field to v,
// then returns b.
func (b *ConsumeEventsRequestBuilder) WithStreamId(v *uuidpb.UUID) *ConsumeEventsRequestBuilder {
	b.prototype.StreamId = v
	return b
}

// WithOffset configures the builder to set the Offset field to v,
// then returns b.
func (b *ConsumeEventsRequestBuilder) WithOffset(v uint64) *ConsumeEventsRequestBuilder {
	b.prototype.Offset = v
	return b
}

// WithEventTypes configures the builder to set the EventTypes field to v,
// then returns b.
func (b *ConsumeEventsRequestBuilder) WithEventTypes(v []*EventType) *ConsumeEventsRequestBuilder {
	b.prototype.EventTypes = v
	return b
}

type ConsumeEventsResponseBuilder struct {
	prototype ConsumeEventsResponse
}

// NewConsumeEventsResponseBuilder returns a builder that constructs [ConsumeEventsResponse] messages.
func NewConsumeEventsResponseBuilder() *ConsumeEventsResponseBuilder {
	return &ConsumeEventsResponseBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *ConsumeEventsResponseBuilder) From(x *ConsumeEventsResponse) *ConsumeEventsResponseBuilder {
	b.prototype.Operation = x.Operation
	return b
}

// Build returns a new [ConsumeEventsResponse] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *ConsumeEventsResponseBuilder) Build() *ConsumeEventsResponse {
	return &ConsumeEventsResponse{
		Operation: b.prototype.Operation,
	}
}

// WithEventDelivery configures the builder to set the Operation field to a
// [ConsumeEventsResponse_EventDelivery_] value containing v, then returns b
func (b *ConsumeEventsResponseBuilder) WithEventDelivery(v *ConsumeEventsResponse_EventDelivery) *ConsumeEventsResponseBuilder {
	b.prototype.Operation = &ConsumeEventsResponse_EventDelivery_{EventDelivery: v}
	return b
}

type EventTypeBuilder struct {
	prototype EventType
}

// NewEventTypeBuilder returns a builder that constructs [EventType] messages.
func NewEventTypeBuilder() *EventTypeBuilder {
	return &EventTypeBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *EventTypeBuilder) From(x *EventType) *EventTypeBuilder {
	b.prototype.PortableName = x.PortableName
	b.prototype.MediaTypes = x.MediaTypes
	return b
}

// Build returns a new [EventType] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *EventTypeBuilder) Build() *EventType {
	return &EventType{
		PortableName: b.prototype.PortableName,
		MediaTypes:   b.prototype.MediaTypes,
	}
}

// WithPortableName configures the builder to set the PortableName field to v,
// then returns b.
func (b *EventTypeBuilder) WithPortableName(v string) *EventTypeBuilder {
	b.prototype.PortableName = v
	return b
}

// WithMediaTypes configures the builder to set the MediaTypes field to v,
// then returns b.
func (b *EventTypeBuilder) WithMediaTypes(v []string) *EventTypeBuilder {
	b.prototype.MediaTypes = v
	return b
}

type UnrecognizedStreamBuilder struct {
	prototype UnrecognizedStream
}

// NewUnrecognizedStreamBuilder returns a builder that constructs [UnrecognizedStream] messages.
func NewUnrecognizedStreamBuilder() *UnrecognizedStreamBuilder {
	return &UnrecognizedStreamBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *UnrecognizedStreamBuilder) From(x *UnrecognizedStream) *UnrecognizedStreamBuilder {
	b.prototype.StreamId = x.StreamId
	return b
}

// Build returns a new [UnrecognizedStream] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *UnrecognizedStreamBuilder) Build() *UnrecognizedStream {
	return &UnrecognizedStream{
		StreamId: b.prototype.StreamId,
	}
}

// WithStreamId configures the builder to set the StreamId field to v,
// then returns b.
func (b *UnrecognizedStreamBuilder) WithStreamId(v *uuidpb.UUID) *UnrecognizedStreamBuilder {
	b.prototype.StreamId = v
	return b
}

type NoEventTypesBuilder struct {
	prototype NoEventTypes
}

// NewNoEventTypesBuilder returns a builder that constructs [NoEventTypes] messages.
func NewNoEventTypesBuilder() *NoEventTypesBuilder {
	return &NoEventTypesBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *NoEventTypesBuilder) From(x *NoEventTypes) *NoEventTypesBuilder {
	return b
}

// Build returns a new [NoEventTypes] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *NoEventTypesBuilder) Build() *NoEventTypes {
	return &NoEventTypes{}
}

type UnrecognizedEventTypeBuilder struct {
	prototype UnrecognizedEventType
}

// NewUnrecognizedEventTypeBuilder returns a builder that constructs [UnrecognizedEventType] messages.
func NewUnrecognizedEventTypeBuilder() *UnrecognizedEventTypeBuilder {
	return &UnrecognizedEventTypeBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *UnrecognizedEventTypeBuilder) From(x *UnrecognizedEventType) *UnrecognizedEventTypeBuilder {
	b.prototype.PortableName = x.PortableName
	return b
}

// Build returns a new [UnrecognizedEventType] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *UnrecognizedEventTypeBuilder) Build() *UnrecognizedEventType {
	return &UnrecognizedEventType{
		PortableName: b.prototype.PortableName,
	}
}

// WithPortableName configures the builder to set the PortableName field to v,
// then returns b.
func (b *UnrecognizedEventTypeBuilder) WithPortableName(v string) *UnrecognizedEventTypeBuilder {
	b.prototype.PortableName = v
	return b
}

type NoRecognizedMediaTypesBuilder struct {
	prototype NoRecognizedMediaTypes
}

// NewNoRecognizedMediaTypesBuilder returns a builder that constructs [NoRecognizedMediaTypes] messages.
func NewNoRecognizedMediaTypesBuilder() *NoRecognizedMediaTypesBuilder {
	return &NoRecognizedMediaTypesBuilder{}
}

// From configures the builder to use x as the prototype for new messages,
// then returns b.
//
// It performs a shallow copy of x, such that any changes made via the builder
// do not modify x. It does not make a copy of the field values themselves.
func (b *NoRecognizedMediaTypesBuilder) From(x *NoRecognizedMediaTypes) *NoRecognizedMediaTypesBuilder {
	b.prototype.PortableName = x.PortableName
	return b
}

// Build returns a new [NoRecognizedMediaTypes] containing the values configured via the builder.
//
// Each call returns a new message, such that future changes to the builder do
// not modify previously constructed messages.
func (b *NoRecognizedMediaTypesBuilder) Build() *NoRecognizedMediaTypes {
	return &NoRecognizedMediaTypes{
		PortableName: b.prototype.PortableName,
	}
}

// WithPortableName configures the builder to set the PortableName field to v,
// then returns b.
func (b *NoRecognizedMediaTypesBuilder) WithPortableName(v string) *NoRecognizedMediaTypesBuilder {
	b.prototype.PortableName = v
	return b
}

// Switch_ConsumeEventsResponse_Operation invokes one of the given functions based on
// the value of x.Operation.
//
// It panics if x.Operation is nil.
func Switch_ConsumeEventsResponse_Operation(
	x *ConsumeEventsResponse,
	caseEventDelivery func(*ConsumeEventsResponse_EventDelivery),
) {
	switch v := x.Operation.(type) {
	case *ConsumeEventsResponse_EventDelivery_:
		caseEventDelivery(v.EventDelivery)
	default:
		panic("Switch_ConsumeEventsResponse_Operation: x.Operation is nil")
	}
}

// Map_ConsumeEventsResponse_Operation maps x.Operation to a value of type T by invoking
// one of the given functions.
//
// It invokes the function that corresponds to the value of x.Operation,
// and returns that function's result. It panics if x.Operation is nil.
func Map_ConsumeEventsResponse_Operation[T any](
	x *ConsumeEventsResponse,
	caseEventDelivery func(*ConsumeEventsResponse_EventDelivery) T,
) T {
	switch v := x.Operation.(type) {
	case *ConsumeEventsResponse_EventDelivery_:
		return caseEventDelivery(v.EventDelivery)
	default:
		panic("Map_ConsumeEventsResponse_Operation: x.Operation is nil")
	}
}

// SetStreams sets the x.Streams field to v, then returns x.
func (x *ListStreamsResponse) SetStreams(v []*Stream) {
	x.Streams = v
}

// SetStreamId sets the x.StreamId field to v, then returns x.
func (x *Stream) SetStreamId(v *uuidpb.UUID) {
	x.StreamId = v
}

// SetEventTypes sets the x.EventTypes field to v, then returns x.
func (x *Stream) SetEventTypes(v []*EventType) {
	x.EventTypes = v
}

// SetStreamId sets the x.StreamId field to v, then returns x.
func (x *ConsumeEventsRequest) SetStreamId(v *uuidpb.UUID) {
	x.StreamId = v
}

// SetOffset sets the x.Offset field to v, then returns x.
func (x *ConsumeEventsRequest) SetOffset(v uint64) {
	x.Offset = v
}

// SetEventTypes sets the x.EventTypes field to v, then returns x.
func (x *ConsumeEventsRequest) SetEventTypes(v []*EventType) {
	x.EventTypes = v
}

// SetEventDelivery sets the x.Operation field to a [ConsumeEventsResponse_EventDelivery_] value containing v,
// then returns x.
func (x *ConsumeEventsResponse) SetEventDelivery(v *ConsumeEventsResponse_EventDelivery) {
	x.Operation = &ConsumeEventsResponse_EventDelivery_{EventDelivery: v}
}

// SetPortableName sets the x.PortableName field to v, then returns x.
func (x *EventType) SetPortableName(v string) {
	x.PortableName = v
}

// SetMediaTypes sets the x.MediaTypes field to v, then returns x.
func (x *EventType) SetMediaTypes(v []string) {
	x.MediaTypes = v
}

// SetStreamId sets the x.StreamId field to v, then returns x.
func (x *UnrecognizedStream) SetStreamId(v *uuidpb.UUID) {
	x.StreamId = v
}

// SetPortableName sets the x.PortableName field to v, then returns x.
func (x *UnrecognizedEventType) SetPortableName(v string) {
	x.PortableName = v
}

// SetPortableName sets the x.PortableName field to v, then returns x.
func (x *NoRecognizedMediaTypes) SetPortableName(v string) {
	x.PortableName = v
}