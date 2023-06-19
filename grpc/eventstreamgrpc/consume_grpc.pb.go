// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
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
	ConsumeAPI_List_FullMethodName    = "/dogma.eventstream.consume.v1.ConsumeAPI/List"
	ConsumeAPI_Consume_FullMethodName = "/dogma.eventstream.consume.v1.ConsumeAPI/Consume"
)

// ConsumeAPIClient is the client API for ConsumeAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConsumeAPIClient interface {
	// List lists the streams that the server provides.
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
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
	// with the stream, as per the List operation. If any other event types are
	// requested the server MUST return an INVALID_ARGUMENT error with an attached
	// [UnrecognizedEventType] value for each unrecognized event type. See
	// [UnrecognizedEventTypeError].
	//
	// If no types are specified the server MUST return an INVALID_ARGUMENT error
	// with an attached [NoEventTypes] value. See [NoEventTypesError].
	//
	// If none of the requested media-types for a given event type are supported
	// the server MUST return an INVALID_ARGUMENT error with an attached
	// [NoRecognizedMediaTypes] value for each such event type. See
	// [NoRecognizedMediaTypesError].
	Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (ConsumeAPI_ConsumeClient, error)
}

type consumeAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewConsumeAPIClient(cc grpc.ClientConnInterface) ConsumeAPIClient {
	return &consumeAPIClient{cc}
}

func (c *consumeAPIClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, ConsumeAPI_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *consumeAPIClient) Consume(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (ConsumeAPI_ConsumeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConsumeAPI_ServiceDesc.Streams[0], ConsumeAPI_Consume_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &consumeAPIConsumeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConsumeAPI_ConsumeClient interface {
	Recv() (*ConsumeResponse, error)
	grpc.ClientStream
}

type consumeAPIConsumeClient struct {
	grpc.ClientStream
}

func (x *consumeAPIConsumeClient) Recv() (*ConsumeResponse, error) {
	m := new(ConsumeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConsumeAPIServer is the server API for ConsumeAPI service.
// All implementations should embed UnimplementedConsumeAPIServer
// for forward compatibility
type ConsumeAPIServer interface {
	// List lists the streams that the server provides.
	List(context.Context, *ListRequest) (*ListResponse, error)
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
	// with the stream, as per the List operation. If any other event types are
	// requested the server MUST return an INVALID_ARGUMENT error with an attached
	// [UnrecognizedEventType] value for each unrecognized event type. See
	// [UnrecognizedEventTypeError].
	//
	// If no types are specified the server MUST return an INVALID_ARGUMENT error
	// with an attached [NoEventTypes] value. See [NoEventTypesError].
	//
	// If none of the requested media-types for a given event type are supported
	// the server MUST return an INVALID_ARGUMENT error with an attached
	// [NoRecognizedMediaTypes] value for each such event type. See
	// [NoRecognizedMediaTypesError].
	Consume(*ConsumeRequest, ConsumeAPI_ConsumeServer) error
}

// UnimplementedConsumeAPIServer should be embedded to have forward compatible implementations.
type UnimplementedConsumeAPIServer struct {
}

func (UnimplementedConsumeAPIServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedConsumeAPIServer) Consume(*ConsumeRequest, ConsumeAPI_ConsumeServer) error {
	return status.Errorf(codes.Unimplemented, "method Consume not implemented")
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

func _ConsumeAPI_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConsumeAPIServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ConsumeAPI_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConsumeAPIServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConsumeAPI_Consume_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConsumeAPIServer).Consume(m, &consumeAPIConsumeServer{stream})
}

type ConsumeAPI_ConsumeServer interface {
	Send(*ConsumeResponse) error
	grpc.ServerStream
}

type consumeAPIConsumeServer struct {
	grpc.ServerStream
}

func (x *consumeAPIConsumeServer) Send(m *ConsumeResponse) error {
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
			MethodName: "List",
			Handler:    _ConsumeAPI_List_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Consume",
			Handler:       _ConsumeAPI_Consume_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "github.com/dogmatiq/enginekit/grpc/eventstreamgrpc/consume.proto",
}
