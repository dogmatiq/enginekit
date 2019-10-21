package message

import (
	"fmt"

	"github.com/dogmatiq/dogma"
)

// Describer is an interface implemented by messages that can describe
// themselves.
//
// It can be used to provide a more specific message description for message
// types that already implement fmt.Stringer, such as when the message
// implementations are generated Protocol Buffers structs.
type Describer interface {
	// MessageDescription returns a human-readable description of the message.
	//
	// Engines may include the message description in logs, error messages, etc.
	//
	// As a general rule, the description should be worded such that it
	// describes the behavior of the message to a non-technical person who is
	// familiar with the business domain.
	MessageDescription() string
}

// Description returns a string representation of m.
//
// If m implements Describer, it returns m.MessageDescription(). Otherwise, if m
// implements fmt.Stringer, it returns m.String().
//
// Finally, if m does not implement either of these interfaces, it returns
// the standard Go "%v" representation of the message.
func Description(m dogma.Message) string {
	if s, ok := m.(Describer); ok {
		return s.MessageDescription()
	}

	if s, ok := m.(fmt.Stringer); ok {
		return s.String()
	}

	return fmt.Sprintf("%v", m)
}
