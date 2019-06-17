package marshaling

import "reflect"

// Codec is an interface for encoding and decoding values.
type Codec interface {
	// Query returns the capabilities of the codec for the given types.
	Query(types []reflect.Type) CodecCapabilities

	// MediaType returns the media-type used to identify values encoded by this
	// codec.
	MediaType() string

	// Marshal returns the binary representation of v.
	Marshal(v interface{}) ([]byte, error)

	// Unmarshal decodes a binary representation into v.
	Unmarshal(data []byte, v interface{}) error
}

// CodecCapabilities describes the capabilities of a codec as it relates to
// specific Go types.
type CodecCapabilities struct {
	// Types is a map of the support types to their portable type name.
	Types map[reflect.Type]string
}
