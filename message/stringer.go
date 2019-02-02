package message

import (
	"fmt"

	"github.com/dogmatiq/dogma"
)

// Stringer is an interface implemented by messages that can describe themselves.
//
// It can be used to provide a more specific message description for message
// types that already implement fmt.Stringer, such as when the message
// implementations are generated Protocol Buffers structs.
type Stringer interface {
	// MessageString returns a human-readable description of the message.
	MessageString() string
}

// ToString returns a string representation of m.
//
// If m implements Stringer, it returns m.MessageString(). Otherwise, if m
// implements fmt.Stringer, it returns m.String().
//
// Finally, if m does not implement either of these interfaces, it returns
// the standard Go "%v" representation of the message.
func ToString(m dogma.Message) string {
	if s, ok := m.(Stringer); ok {
		return s.MessageString()
	}

	if s, ok := m.(fmt.Stringer); ok {
		return s.String()
	}

	return fmt.Sprintf("%v", m)
}
