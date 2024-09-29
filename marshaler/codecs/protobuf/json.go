package protobuf

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// JSONBasicMediaType is the type and subtype portion of the media-type used to
// identify data encoded using JSON.
const JSONBasicMediaType = "application/vnd.google.protobuf+json"

// DefaultJSONMarshaler is the marshaler used by DefaultJSONCodec.
var DefaultJSONMarshaler = protojson.MarshalOptions{
	UseProtoNames: true,
}

// DefaultJSONUnmarshaler is the unmarshaler used by DefaultJSONCodec.
var DefaultJSONUnmarshaler = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

// DefaultJSONCodec is a [marshaler.Codec] that marshals protocol buffers
// messages in JSON format.
//
// See https://developers.google.com/protocol-buffers/docs/proto3#json.
var DefaultJSONCodec = Codec{
	MediaType:   JSONBasicMediaType,
	Marshaler:   DefaultJSONMarshaler,
	Unmarshaler: DefaultJSONUnmarshaler,
}
