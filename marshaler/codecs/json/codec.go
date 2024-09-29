package json

import (
	"encoding/json"
	"reflect"
)

// Codec is an implementation of [marshaler.Codec] that uses Go's standard JSON
// implementation.
type Codec struct{}

// DefaultCodec is a [marshaler.Codec] that marshals values using Go's standard
// JSON implementation.
var DefaultCodec = Codec{}

// PortableNames returns a map of type to its portable name for each of the
// given types that the codec supports.
func (Codec) PortableNames(types []reflect.Type) map[reflect.Type]string {
	names := map[reflect.Type]string{}

	for _, rt := range types {
		if n, ok := portableName(rt); ok {
			names[rt] = n
		}
	}

	return names
}

// BasicMediaType returns the type and subtype portion of the media-type used to
// identify data encoded by this codec.
func (Codec) BasicMediaType() string {
	return "application/json"
}

// Marshal returns the binary representation of v.
func (Codec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes a binary representation into v.
func (Codec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
