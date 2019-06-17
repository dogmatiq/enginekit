package protobuf

import (
	"github.com/golang/protobuf/proto"
)

// NativeCodec is an implementation of marshaling.Codec that marshals protocol
// buffers messages in the native binary format.
type NativeCodec struct {
	codec
}

// MediaType returns the media-type used to identify values encoded by this
// codec.
func (c *NativeCodec) MediaType() string {
	return "application/vnd.google.protobuf"
}

// Marshal returns the binary representation of v.
func (c *NativeCodec) Marshal(v interface{}) ([]byte, error) {
	m, err := cast(v)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(m)
}

// Unmarshal decodes a binary representation into v.
func (c *NativeCodec) Unmarshal(data []byte, v interface{}) error {
	m, err := cast(v)
	if err != nil {
		return err
	}

	return proto.Unmarshal(data, m)
}
