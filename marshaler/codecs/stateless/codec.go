package stateless

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/dogmatiq/dogma"
)

// Codec is an implementation of [marshaler.Codec] that marshals
// [dogma.StatelessProcessRoot] values.
type Codec struct{}

// DefaultCodec is a [marshaler.Codec] that marshals
// [dogma.StatelessProcessRoot] values.
var DefaultCodec = Codec{}

var processRootType = reflect.TypeOf(dogma.StatelessProcessRoot)

// PortableNames returns a map of type to its portable name for each of the
// given types that the codec supports.
func (Codec) PortableNames(types []reflect.Type) map[reflect.Type]string {
	if slices.Contains(types, processRootType) {
		return map[reflect.Type]string{
			processRootType: "process",
		}
	}
	return nil
}

// BasicMediaType returns the type and subtype portion of the media-type used to
// identify data encoded by this codec.
func (Codec) BasicMediaType() string {
	return "application/x-empty"
}

// Marshal returns the binary representation of v.
func (Codec) Marshal(v any) ([]byte, error) {
	if v == dogma.StatelessProcessRoot {
		return nil, nil
	}
	return nil, fmt.Errorf("%T is not dogma.StatelessProcessRoot", v)
}

// Unmarshal decodes a binary representation into v.
func (Codec) Unmarshal(data []byte, v any) error {
	if v != &dogma.StatelessProcessRoot {
		return fmt.Errorf("%T is not a pointer to dogma.StatelessProcessRoot", v)
	}
	if len(data) != 0 {
		return fmt.Errorf("expected empty data, got %d byte(s)", len(data))
	}
	return nil
}
