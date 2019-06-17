package protobuf

import (
	"bytes"

	"github.com/golang/protobuf/proto"
)

// DefaultTextMarshaler is the text marshaler used by TextCodec if none is
// provided.
var DefaultTextMarshaler = &proto.TextMarshaler{
	Compact:   false,
	ExpandAny: true,
}

// TextCodec is an implementation of marshaling.Codec that marshals protocol
// buffers messages in text format.
type TextCodec struct {
	codec

	// Marshaler is the text marshaler used to marshal messages.
	// If it is nil, DefaultTextMarshaler is used.
	Marshaler *proto.TextMarshaler
}

// MediaType returns the media-type used to identify values encoded by this
// codec.
func (c *TextCodec) MediaType() string {
	return "text/vnd.google.protobuf"
}

// Marshal returns the binary representation of v.
func (c *TextCodec) Marshal(v interface{}) ([]byte, error) {
	m, err := cast(v)
	if err != nil {
		return nil, err
	}

	tm := c.Marshaler
	if tm == nil {
		tm = DefaultTextMarshaler
	}

	buf := &bytes.Buffer{}
	err = tm.Marshal(buf, m)
	return buf.Bytes(), err
}

// Unmarshal decodes a binary representation into v.
func (c *TextCodec) Unmarshal(data []byte, v interface{}) error {
	m, err := cast(v)
	if err != nil {
		return err
	}

	return proto.UnmarshalText(
		string(data),
		m,
	)
}
