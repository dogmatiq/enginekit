package marshaler

import (
	"reflect"
)

// Codec is an interface for encoding and decoding values using a specific
// format.
type Codec interface {
	// PortableNames returns a map of type to its portable name for each of the
	// given types that the codec supports.
	PortableNames(types []reflect.Type) map[reflect.Type]string

	// BasicMediaType returns the type and subtype portion of the MIME
	// media-type used by this code. For example, "application/json".
	BasicMediaType() string

	// Marshal returns the portable representation of v.
	Marshal(v any) ([]byte, error)

	// Unmarshal decodes a portable representation into v.
	//
	// v must be a pointer to the type that the data represents.
	Unmarshal(data []byte, v any) error
}
