package config

import (
	"bytes"
	"fmt"
)

// MessageRole is an enumeration of the roles a message can perform within an
// application.
type MessageRole byte

const (
	// CommandMessageRole is the role for command messages.
	CommandMessageRole MessageRole = 'C'

	// EventMessageRole is the role for event messages.
	EventMessageRole MessageRole = 'E'

	// TimeoutMessageRole is the role for timeout messages.
	TimeoutMessageRole MessageRole = 'T'
)

// MessageRoles is a slice of the valid message roles.
var MessageRoles = []MessageRole{
	CommandMessageRole,
	EventMessageRole,
	TimeoutMessageRole,
}

const (
	commandMessageRoleString = "command"
	eventMessageRoleString   = "event"
	timeoutMessageRoleString = "timeout"

	commandMessageRoleShortString = "CMD"
	eventMessageRoleShortString   = "EVT"
	timeoutMessageRoleShortString = "TMO"

	commandMessageRoleMarker = "?"
	eventMessageRoleMarker   = "!"
	timeoutMessageRoleMarker = "@"
)

var (
	commandMessageRoleBytes = []byte(commandMessageRoleString)
	eventMessageRoleBytes   = []byte(eventMessageRoleString)
	timeoutMessageRoleBytes = []byte(timeoutMessageRoleString)
)

// Validate return an error if r is not a valid message role.
func (r MessageRole) Validate() error {
	switch r {
	case CommandMessageRole,
		EventMessageRole,
		TimeoutMessageRole:
		return nil
	default:
		return fmt.Errorf("invalid message role: %#v", r)
	}
}

// MustValidate panics if r is not a valid message role.
func (r MessageRole) MustValidate() {
	if err := r.Validate(); err != nil {
		panic(err)
	}
}

// Is returns true if r is one of the given roles.
func (r MessageRole) Is(roles ...MessageRole) bool {
	r.MustValidate()

	for _, x := range roles {
		x.MustValidate()

		if r == x {
			return true
		}
	}

	return false
}

// MustBe panics if r is not one of the given roles.
func (r MessageRole) MustBe(roles ...MessageRole) {
	if !r.Is(roles...) {
		panic("unexpected message role: " + r.String())
	}
}

// MustNotBe panics if r is one of the given roles.
func (r MessageRole) MustNotBe(roles ...MessageRole) {
	if r.Is(roles...) {
		panic("unexpected message role: " + r.String())
	}
}

// Marker returns a character that identifies the message role when displaying
// message types.
func (r MessageRole) Marker() string {
	r.MustValidate()

	switch r {
	case CommandMessageRole:
		return commandMessageRoleMarker
	case EventMessageRole:
		return eventMessageRoleMarker
	default: // TimeoutMessageRole
		return timeoutMessageRoleMarker
	}
}

// ShortString returns a short (3-character) representation of the handler type.
func (r MessageRole) ShortString() string {
	r.MustValidate()

	switch r {
	case CommandMessageRole:
		return commandMessageRoleShortString
	case EventMessageRole:
		return eventMessageRoleShortString
	default: // TimeoutMessageRole
		return timeoutMessageRoleShortString
	}
}

func (r MessageRole) String() string {
	switch r {
	case CommandMessageRole:
		return commandMessageRoleString
	case EventMessageRole:
		return eventMessageRoleString
	case TimeoutMessageRole:
		return timeoutMessageRoleString
	default:
		return fmt.Sprintf("<invalid message role %#v>", r)
	}
}

// MarshalText returns a UTF-8 representation of the message role.
func (r MessageRole) MarshalText() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}

	switch r {
	case CommandMessageRole:
		return commandMessageRoleBytes, nil
	case EventMessageRole:
		return eventMessageRoleBytes, nil
	default: // TimeoutMessageRole
		return timeoutMessageRoleBytes, nil
	}
}

// UnmarshalText unmarshals a role from its UTF-8 representation.
func (r *MessageRole) UnmarshalText(text []byte) error {
	if bytes.Equal(text, commandMessageRoleBytes) {
		*r = CommandMessageRole
	} else if bytes.Equal(text, eventMessageRoleBytes) {
		*r = EventMessageRole
	} else if bytes.Equal(text, timeoutMessageRoleBytes) {
		*r = TimeoutMessageRole
	} else {
		return fmt.Errorf("invalid text representation of message role: %s", text)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message role.
func (r MessageRole) MarshalBinary() ([]byte, error) {
	return []byte{byte(r)}, r.Validate()
}

// UnmarshalBinary unmarshals a role from its binary representation.
func (r *MessageRole) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("invalid binary representation of message role, expected exactly 1 byte")
	}

	*r = MessageRole(data[0])
	return r.Validate()
}
