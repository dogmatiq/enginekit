package message

import (
	"github.com/dogmatiq/dogma"
)

// TypeContainer is an interface for contains of message types.
type TypeContainer interface {
	// Has returns true if the set contains t.
	Has(t Type) bool

	// HasM returns true if the set contains TypeOf(m).
	HasM(m dogma.Message) bool
}
