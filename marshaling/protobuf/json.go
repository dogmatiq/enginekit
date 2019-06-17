package protobuf

import (
	"bytes"

	"github.com/golang/protobuf/jsonpb"
)

// DefaultJSONMarshaler is the text marshaler used by JSONCodec if none is
// provided.
var DefaultJSONMarshaler = &jsonpb.Marshaler{
	OrigName: true,
}

// DefaultJSONUnmarshaler is the text marshaler used by JSONCodec if none is
// provided.
var DefaultJSONUnmarshaler = &jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

// JSONCodec is an implementation of marshaling.Codec that marshals protocol
// buffers messages in text format.
type JSONCodec struct {
	codec

	// Marshaler is the JSON marshaler used to marshal messages.
	// If it is nil, DefaultJSONMarshaler is used.
	Marshaler *jsonpb.Marshaler

	// Unmarshaler is the JSON unmarshaler used to unmarshal messages.
	// If it is nil, DefaultJSONUnmarshaler is used.
	Unmarshaler *jsonpb.Unmarshaler
}

// MediaType returns the media-type used to identify values encoded by this
// codec.
func (c *JSONCodec) MediaType() string {
	return "application/vnd.google.protobuf+json"
}

// Marshal returns the binary representation of v.
func (c *JSONCodec) Marshal(v interface{}) ([]byte, error) {
	m, err := cast(v)
	if err != nil {
		return nil, err
	}

	jm := c.Marshaler
	if jm == nil {
		jm = DefaultJSONMarshaler
	}

	buf := &bytes.Buffer{}
	err = jm.Marshal(buf, m)
	return buf.Bytes(), err
}

// Unmarshal decodes a binary representation into v.
func (c *JSONCodec) Unmarshal(data []byte, v interface{}) error {
	m, err := cast(v)
	if err != nil {
		return err
	}

	ju := c.Unmarshaler
	if ju == nil {
		ju = DefaultJSONUnmarshaler
	}

	return ju.Unmarshal(
		bytes.NewReader(data),
		m,
	)
}
