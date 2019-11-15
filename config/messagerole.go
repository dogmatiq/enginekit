package config

import (
	"fmt"
)

// MessageRole is an enumeration of the roles a message can perform within an
// application.
type MessageRole string

const (
	// CommandMessageRole is the role for command messages.
	CommandMessageRole MessageRole = "command"

	// EventMessageRole is the role for event messages.
	EventMessageRole MessageRole = "event"

	// TimeoutMessageRole is the role for timeout messages.
	TimeoutMessageRole MessageRole = "timeout"
)

// Roles is a slice of the valid message roles.
var Roles = []MessageRole{
	CommandMessageRole,
	EventMessageRole,
	TimeoutMessageRole,
}

// Validate return an error if r is not a valid message role.
func (r MessageRole) Validate() error {
	switch r {
	case CommandMessageRole,
		EventMessageRole,
		TimeoutMessageRole:
		return nil
	default:
		return fmt.Errorf("invalid message role: %s", r)
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
		panic("unexpected role: " + r)
	}
}

// MustNotBe panics if r is one of the given roles.
func (r MessageRole) MustNotBe(roles ...MessageRole) {
	if r.Is(roles...) {
		panic("unexpected role: " + r)
	}
}

// Marker returns a character that identifies the message role when displaying
// message types.
func (r MessageRole) Marker() string {
	r.MustValidate()

	switch r {
	case CommandMessageRole:
		return "?"
	case EventMessageRole:
		return "!"
	default:
		return "@"
	}
}

func (r MessageRole) String() string {
	return string(r)
}
