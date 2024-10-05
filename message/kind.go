package message

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Kind is an enumeration of the different kinds of messages.
type Kind int

const (
	// CommandKind is the [Kind] for a "command" message, which is a request to make a
	// single atomic change to the application's state. Command messages
	// implement the [dogma.Command] interface.
	CommandKind Kind = iota

	// EventKind is the [Kind] for an "event" message, which indicates that the
	// application state has changed in some way. Event messages implement the
	// [dogma.Event] interface.
	EventKind

	// TimeoutKind is the [Kind] for a "timeout" message, which is a message
	// that model business logic that depends on the passage of time. Timeout
	// messages implement the [dogma.Timeout] interface.
	TimeoutKind
)

func (k Kind) String() string {
	return MapByKind(k, "command", "event", "timeout")
}

// Symbol returns a single-character symbol that represents the kind. It is
// often appended to a [Name] to make it easy to distinguish between different
// kinds of messages.
func (k Kind) Symbol() string {
	return MapByKind(k, "?", "!", "@")
}

// KindFor returns the [Kind] of the message with type T.
//
// It panics if T does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func KindFor[T dogma.Message]() Kind {
	return kindFromReflect(
		reflect.TypeFor[T](),
	)
}

// KindOf returns the [Kind] of m.
//
// It panics if m does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func KindOf(m dogma.Message) Kind {
	switch m.(type) {
	case dogma.Command:
		return CommandKind
	case dogma.Event:
		return EventKind
	case dogma.Timeout:
		return TimeoutKind
	default:
		panic(fmt.Sprintf("%T does not implement dogma.Command, dogma.Event or dogma.Timeout", m))
	}
}

// kindFromReflect returns the [Kind] of a message with type r.
//
// It panics if r does not implement [dogma.Command], [dogma.Event] or
// [dogma.Timeout].
func kindFromReflect(r reflect.Type) Kind {
	switch interfaceOf(r) {
	case commandInterface:
		return CommandKind
	case eventInterface:
		return EventKind
	case timeoutInterface:
		return TimeoutKind
	default:
		panic("unexpected message interface")
	}
}
