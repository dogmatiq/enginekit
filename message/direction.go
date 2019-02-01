package message

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

// MustValidate panics if r is not a valid message direction.
func (d Direction) MustValidate() {
	switch d {
	case InboundDirection:
	case OutboundDirection:
	default:
		panic("invalid direction: " + d)
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
