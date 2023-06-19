package eventstreamgrpc

import (
	uuidpb "github.com/dogmatiq/enginekit/protobuf/uuidpb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// UnrecognizedStreamError returns an error indicating that the given stream is
// not recognized by the server.
func UnrecognizedStreamError(id *uuidpb.UUID) error {
	s, err := status.
		New(codes.NotFound, "unrecognized stream").
		WithDetails(&UnrecognizedStream{
			StreamId: id,
		})

	if err != nil {
		panic(err)
	}

	return s.Err()
}

// UnrecognizedEventTypeError returns an error indicating that the given event
// type is not recognized by the server.
func UnrecognizedEventTypeError(portableName string) error {
	s, err := status.
		New(codes.InvalidArgument, "unrecognized event type").
		WithDetails(&UnrecognizedEventType{
			PortableName: portableName,
		})

	if err != nil {
		panic(err)
	}

	return s.Err()
}

// NoRecognizedMediaTypesError returns an error indicating that the server does
// not support any of the requested media-types for the given event type.
func NoRecognizedMediaTypesError(portableName string) error {
	s, err := status.
		New(codes.InvalidArgument, "no recognized media types").
		WithDetails(&NoRecognizedMediaTypes{
			PortableName: portableName,
		})

	if err != nil {
		panic(err)
	}

	return s.Err()
}
