package json

import (
	"encoding/json"
	"reflect"

	"github.com/dogmatiq/enginekit/marshaling"
)

// Codec is an implementation of marshaling.Codec that uses Go's standard JSON
// implementation.
type Codec struct{}

// Query returns the capabilities of the codec for the given types.
func (*Codec) Query(types []reflect.Type) marshaling.CodecCapabilities {
	caps := marshaling.CodecCapabilities{
		Types: map[reflect.Type]string{},
	}

	for _, rt := range types {
		if n, ok := portableName(rt); ok {
			caps.Types[rt] = n
		}
	}

	return caps
}

// MediaType returns the media-type used to identify data encoded by this
// codec.
func (*Codec) MediaType() string {
	return "application/json"
}

// Marshal returns the binary representation of v.
func (*Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes a binary representation into v.
func (*Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// portableName returns the portable name to use for the given type.
func portableName(rt reflect.Type) (string, bool) {
	n := rt.Name()
	if n != "" {
		return n, true
	}

	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	n = rt.Name()
	return n, n != ""
}
