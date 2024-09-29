package protobuf

import (
	"google.golang.org/protobuf/encoding/prototext"
)

// TextBasicMediaType is the type and subtype portion of the media-type used to
// identify data encoded in the protocol buffers text format.
const TextBasicMediaType = "text/vnd.google.protobuf"

// DefaultTextMarshaler is the marshaler used by DefaultTextCodec.
var DefaultTextMarshaler = prototext.MarshalOptions{
	Multiline: true,
}

// DefaultTextUnmarshaler is the unmarshaler used by DefaultTextCodec.
var DefaultTextUnmarshaler = prototext.UnmarshalOptions{
	DiscardUnknown: true,
}

// DefaultTextCodec is a [marshaler.Codec] that marshals protocol buffers
// messages in text format.
var DefaultTextCodec = Codec{
	MediaType:   TextBasicMediaType,
	Marshaler:   DefaultTextMarshaler,
	Unmarshaler: DefaultTextUnmarshaler,
}
