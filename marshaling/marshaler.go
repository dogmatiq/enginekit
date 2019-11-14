package marshaling

import (
	"fmt"
	"mime"
	"reflect"
)

// Marshaler uses a set of codecs to marshal and unmarshal values.
type Marshaler struct {
	encoders map[reflect.Type]Codec
	decoders map[string]Codec
	names    map[reflect.Type]string
	types    map[string]reflect.Type
}

// NewMarshaler returns a new marshaler that uses the given set of codecs to
// marshal and unmarshal values.
//
// The codecs are given in order of preference.
func NewMarshaler(
	types []reflect.Type,
	codecs []Codec,
) (*Marshaler, error) {
	m := &Marshaler{
		encoders: map[reflect.Type]Codec{},
		decoders: map[string]Codec{},
		names:    map[reflect.Type]string{},
		types:    map[string]reflect.Type{},
	}

	// build a list of all of the "unsupported" types
	unsupported := map[reflect.Type]struct{}{}
	for _, rt := range types {
		unsupported[rt] = struct{}{}
	}

	// scan the codecs in order of preference
	for _, c := range codecs {
		caps := c.Query(types)

		if len(caps.Types) > 0 {
			for rt, n := range caps.Types {
				if _, ok := m.encoders[rt]; ok {
					// a higher priority codec has already "claimed" this type
					continue
				}

				if x, ok := m.types[n]; ok {
					return nil, fmt.Errorf(
						"the type name '%s' is used by both '%s' and '%s'",
						n,
						x,
						rt,
					)
				}

				m.encoders[rt] = c
				m.names[rt] = n
				m.types[n] = rt

				delete(unsupported, rt)
			}

			if _, ok := m.decoders[c.MediaType()]; ok {
				return nil, fmt.Errorf(
					"multiple codecs use the '%s' media-type",
					c.MediaType(),
				)
			}

			m.decoders[c.MediaType()] = c
		}
	}

	for rt := range unsupported {
		return nil, fmt.Errorf(
			"no codecs support the '%s' type",
			rt,
		)
	}

	return m, nil
}

// MarshalType marshals a type to its portable representation.
func (m *Marshaler) MarshalType(rt reflect.Type) (string, error) {
	if n, ok := m.names[rt]; ok {
		return n, nil
	}

	return "", fmt.Errorf(
		"no codecs support the '%s' type",
		rt,
	)
}

// UnmarshalType marshals a type from its portable representation.
func (m *Marshaler) UnmarshalType(n string) (reflect.Type, error) {
	if rt, ok := m.types[n]; ok {
		return rt, nil
	}

	return nil, fmt.Errorf(
		"the portable type name '%s' is not recognized",
		n,
	)
}

// UnmarshalTypeFromMediaType unmarshals a type from a MIME media-type.
func (m *Marshaler) UnmarshalTypeFromMediaType(mt string) (reflect.Type, error) {
	_, rt, err := m.unpackMediaType(mt)
	return rt, err
}

// Marshal returns a binary representation of v.
func (m *Marshaler) Marshal(v interface{}) (Packet, error) {
	rt := reflect.TypeOf(v)

	if c, ok := m.encoders[rt]; ok {
		data, err := c.Marshal(v)
		if err != nil {
			return Packet{}, err
		}

		return Packet{
			mime.FormatMediaType(
				c.MediaType(),
				map[string]string{"type": m.names[rt]},
			),
			data,
		}, nil
	}

	return Packet{}, fmt.Errorf(
		"no codecs support the '%s' type",
		rt,
	)
}

// Unmarshal produces a value from its binary representation.
func (m *Marshaler) Unmarshal(p Packet) (interface{}, error) {
	c, rt, err := m.unpackMediaType(p.MediaType)
	if err != nil {
		return nil, err
	}

	// If the type is already a pointer, we wan't to construct the element that
	// it points to, otherwise construct a new pointer to the actual type.
	var v reflect.Value
	if rt.Kind() == reflect.Ptr {
		v = reflect.New(rt.Elem())
	} else {
		v = reflect.New(rt)
	}

	if err := c.Unmarshal(p.Data, v.Interface()); err != nil {
		return nil, err
	}

	// Unwrap the pointer we created just to allow for unmarshalling.
	if rt.Kind() != reflect.Ptr {
		v = v.Elem()
	}

	return v.Interface(), nil
}

func (m *Marshaler) unpackMediaType(s string) (Codec, reflect.Type, error) {
	mt, params, err := mime.ParseMediaType(s)
	if err != nil {
		return nil, nil, err
	}

	c, ok := m.decoders[mt]
	if !ok {
		return nil, nil, fmt.Errorf(
			"no codecs support the '%s' media-type",
			mt,
		)
	}

	n, ok := params["type"]
	if !ok {
		return nil, nil, fmt.Errorf(
			"the media-type '%s' does not specify a 'type' parameter",
			mt,
		)
	}

	rt, ok := m.types[n]
	if !ok {
		return nil, nil, fmt.Errorf(
			"the portable type name '%s' is not recognized",
			n,
		)
	}

	return c, rt, nil
}
