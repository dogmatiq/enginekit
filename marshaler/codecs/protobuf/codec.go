package protobuf

import (
	"reflect"

	"google.golang.org/protobuf/proto"
)

// Codec is an implementation of marshalkit.Codec that encodes Protocol Buffers
// messages.
//
// It supports three common protocol buffers formats, that is, the native binary
// format, the JSON "mapping", and the text-based encoding scheme.
//
// See DefaultNativeCodec, DefaultJSONCodec and DefaultTextCodec, respectively.
type Codec struct {
	// MediaType is the type and subtype portion of the media-type used to
	// identify data encoded by this codec. If it is empty, NativeMediaType is
	// used.
	MediaType string

	// Marshaler is the Marshaler used to marshal messages.
	// If it is nil, DefaultNativeMarshaler is used.
	Marshaler Marshaler

	// Unmarshaler is the JSON Unmarshaler used to unmarshal messages.
	// If it is nil, DefaultNativeUnmarshaler is used.
	Unmarshaler Unmarshaler
}

// Marshaler is an interface for marshaling protocol buffers messages.
type Marshaler interface {
	Marshal(proto.Message) ([]byte, error)
}

// Unmarshaler is an interface for unmarshaling protocol buffers messages.
type Unmarshaler interface {
	Unmarshal([]byte, proto.Message) error
}

// PortableNames returns a map of type to its portable name for each of the
// given types that the codec supports.
func (c Codec) PortableNames(types []reflect.Type) map[reflect.Type]string {
	names := map[reflect.Type]string{}

	for _, rt := range types {
		m := reflect.Zero(rt).Interface()

		if m, ok := m.(proto.Message); ok {
			if n := proto.MessageName(m); n != "" {
				names[rt] = string(n)
			}
		}
	}

	return names
}

// BasicMediaType returns the type and subtype portion of the media-type used to
// identify data encoded by this codec.
func (c Codec) BasicMediaType() string {
	if c.MediaType == "" {
		return NativeBasicMediaType
	}

	return c.MediaType
}

// Marshal returns the binary representation of v.
func (c Codec) Marshal(v any) ([]byte, error) {
	m, err := cast(v)
	if err != nil {
		return nil, err
	}

	marshaler := c.Marshaler
	if marshaler == nil {
		marshaler = DefaultNativeMarshaler
	}

	return marshaler.Marshal(m)
}

// Unmarshal decodes a binary representation into v.
func (c Codec) Unmarshal(data []byte, v any) error {
	m, err := cast(v)
	if err != nil {
		return err
	}

	unmarshaler := c.Unmarshaler
	if unmarshaler == nil {
		unmarshaler = DefaultNativeUnmarshaler
	}

	return unmarshaler.Unmarshal(data, m)
}
