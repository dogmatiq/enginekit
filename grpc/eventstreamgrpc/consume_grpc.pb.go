// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.0
// source: github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto

package eventstreamgrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ConsumeAPI_ListStreams_FullMethodName   = "/dogma.eventstream.consume.v1.ConsumeAPI/ListStreams"
	ConsumeAPI_ConsumeEvents_FullMethodName = "/dogma.eventstream.consume.v1.ConsumeAPI/ConsumeEvents"
)

// ConsumeAPIClient is the client API for ConsumeAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsumeAPIClient interface {
	// ListStreams lists the streams that the server provides.
	ListStreams(ctx context.Context, in *ListStreamsRequest, opts ...grpc.CallOption) (*ListStreamsResponse, error)
	// Consume starts consuming from a specific offset within an event stream.
	//
	// If the requested stream ID is unknown to the server it MUST return a
	// NOT_FOUND error with an attached [UnrecognizedStream] value. See
	// [UnrecognizedStreamError].
	//
	// If the requested offset is beyond the end of the stream, the server SHOULD
	// keep the stream open and send new events as they are written to the stream.
	//
	// The requested event types MUST be a subset of those event types associated
	// with the stream, as per the result of the ListStreams operation. If any
	// other event types are requested the server MUST return an INVALID_ARGUMENT
	// error with an attached [UnrecognizedEventType] value for each unrecognized
	// event type. See [UnrecognizedEventTypeError].
	//
	// If no types are specified the server MUST return an INVALID_ARGUMENT error
	// with an attached [NoEventTypes] value. See [NoEventTypesError].
	//
	// If none of the requested media-types for a given event type are supported
	// the server MUST return an INVALID_ARGUMENT error with an attached
	// [NoRecognizedMediaTypes] value for each such event type. See
	// [NoRecognizedMediaTypesError].
	ConsumeEvents(ctx context.Context, in *ConsumeEventsRequest, opts ...grpc.CallOption) (ConsumeAPI_ConsumeEventsClient, error)
}

type consumeAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewConsumeAPIClient(cc grpc.ClientConnInterface) ConsumeAPIClient {
	return &consumeAPIClient{cc}
}

func (c *consumeAPIClient) ListStreams(ctx context.Context, in *ListStreamsRequest, opts ...grpc.CallOption) (*ListStreamsResponse, error) {
	out := new(ListStreamsResponse)
	err := c.cc.Invoke(ctx, ConsumeAPI_ListStreams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consumeAPIClient) ConsumeEvents(ctx context.Context, in *ConsumeEventsRequest, opts ...grpc.CallOption) (ConsumeAPI_ConsumeEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConsumeAPI_ServiceDesc.Streams[0], ConsumeAPI_ConsumeEvents_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &consumeAPIConsumeEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConsumeAPI_ConsumeEventsClient interface {
	Recv() (*ConsumeEventsResponse, error)
	grpc.ClientStream
}

type consumeAPIConsumeEventsClient struct {
	grpc.ClientStream
}

func (x *consumeAPIConsumeEventsClient) Recv() (*ConsumeEventsResponse, error) {
	m := new(ConsumeEventsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConsumeAPIServer is the server API for ConsumeAPI service.
// All implementations should embed UnimplementedConsumeAPIServer
// for forward compatibility
type ConsumeAPIServer interface {
	// ListStreams lists the streams that the server provides.
	ListStreams(context.Context, *ListStreamsRequest) (*ListStreamsResponse, error)
	// Consume starts consuming from a specific offset within an event stream.
	//
	// If the requested stream ID is unknown to the server it MUST return a
	// NOT_FOUND error with an attached [UnrecognizedStream] value. See
	// [UnrecognizedStreamError].
	//
	// If the requested offset is beyond the end of the stream, the server SHOULD
	// keep the stream open and send new events as they are written to the stream.
	//
	// The requested event types MUST be a subset of those event types associated
	// with the stream, as per the result of the ListStreams operation. If any
	// other event types are requested the server MUST return an INVALID_ARGUMENT
	// error with an attached [UnrecognizedEventType] value for each unrecognized
	// event type. See [UnrecognizedEventTypeError].
	//
	// If no types are specified the server MUST return an INVALID_ARGUMENT error
	// with an attached [NoEventTypes] value. See [NoEventTypesError].
	//
	// If none of the requested media-types for a given event type are supported
	// the server MUST return an INVALID_ARGUMENT error with an attached
	// [NoRecognizedMediaTypes] value for each such event type. See
	// [NoRecognizedMediaTypesError].
	ConsumeEvents(*ConsumeEventsRequest, ConsumeAPI_ConsumeEventsServer) error
}

// UnimplementedConsumeAPIServer should be embedded to have forward compatible implementations.
type UnimplementedConsumeAPIServer struct {
}

func (UnimplementedConsumeAPIServer) ListStreams(context.Context, *ListStreamsRequest) (*ListStreamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStreams not implemented")
}
func (UnimplementedConsumeAPIServer) ConsumeEvents(*ConsumeEventsRequest, ConsumeAPI_ConsumeEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method ConsumeEvents not implemented")
}

// UnsafeConsumeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConsumeAPIServer will
// result in compilation errors.
type UnsafeConsumeAPIServer interface {
	mustEmbedUnimplementedConsumeAPIServer()
}

func RegisterConsumeAPIServer(s grpc.ServiceRegistrar, srv ConsumeAPIServer) {
	s.RegisterService(&ConsumeAPI_ServiceDesc, srv)
}

func _ConsumeAPI_ListStreams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStreamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsumeAPIServer).ListStreams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConsumeAPI_ListStreams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsumeAPIServer).ListStreams(ctx, req.(*ListStreamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsumeAPI_ConsumeEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeEventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConsumeAPIServer).ConsumeEvents(m, &consumeAPIConsumeEventsServer{stream})
}

type ConsumeAPI_ConsumeEventsServer interface {
	Send(*ConsumeEventsResponse) error
	grpc.ServerStream
}

type consumeAPIConsumeEventsServer struct {
	grpc.ServerStream
}

func (x *consumeAPIConsumeEventsServer) Send(m *ConsumeEventsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ConsumeAPI_ServiceDesc is the grpc.ServiceDesc for ConsumeAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConsumeAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dogma.eventstream.consume.v1.ConsumeAPI",
	HandlerType: (*ConsumeAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListStreams",
			Handler:    _ConsumeAPI_ListStreams_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConsumeEvents",
			Handler:       _ConsumeAPI_ConsumeEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto",
}
