package message

import "fmt"

// Direction is an enumeration of the "direction" a message is flowing, relative
// to a handler.
type Direction string

const (
	// InboundDirection means a message is being passed into a handler, that is,
	// being handled.
	InboundDirection Direction = "inbound"

	// OutboundDirection means a message was produced by a handler.
	OutboundDirection Direction = "outbound"
)

// Directions is a slice of the valid message direcations.
var Directions = []Direction{
	InboundDirection,
	OutboundDirection,
}

// Validate returns an error if r is not a valid message direction.
func (d Direction) Validate() error {
	switch d {
	case InboundDirection,
		OutboundDirection:
		return nil
	default:
		return fmt.Errorf("invalid direction: %s", d)
	}
}

// MustValidate panics if r is not a valid message direction.
func (d Direction) MustValidate() {
	if err := d.Validate(); err != nil {
		panic(err)
	}
}

// MustBe panics if r is not equal to x.
func (d Direction) MustBe(x Direction) {
	if d != x {
		panic("unexpected direction: " + d)
	}
}

// MustNotBe panics if r is equal to x.
func (d Direction) MustNotBe(x Direction) {
	if d == x {
		panic("unexpected role: " + d)
	}
}

func (d Direction) String() string {
	return string(d)
}
