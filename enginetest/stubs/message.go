package stubs

import (
	"errors"
	"fmt"
)

// CommandStub is a test implementation of [dogma.Command].
type CommandStub[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (c CommandStub[T]) MessageDescription() string {
	validity := "valid"
	if c.Invalid != "" {
		validity = "invalid: " + c.Invalid
	}
	return fmt.Sprintf(
		"command(%T:%v, %s)",
		c.Content,
		c.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (c CommandStub[T]) Validate() error {
	if c.Invalid != "" {
		return errors.New(c.Invalid)
	}
	return nil
}

// EventStub is a test implementation of [dogma.Event].
type EventStub[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (e EventStub[T]) MessageDescription() string {
	validity := "valid"
	if e.Invalid != "" {
		validity = "invalid: " + e.Invalid
	}
	return fmt.Sprintf(
		"event(%T:%v, %s)",
		e.Content,
		e.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (e EventStub[T]) Validate() error {
	if e.Invalid != "" {
		return errors.New(e.Invalid)
	}
	return nil
}

// TimeoutStub is a test implementation of [dogma.Test].
type TimeoutStub[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (t TimeoutStub[T]) MessageDescription() string {
	validity := "valid"
	if t.Invalid != "" {
		validity = "invalid: " + t.Invalid
	}
	return fmt.Sprintf(
		"timeout(%T:%v, %s)",
		t.Content,
		t.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (t TimeoutStub[T]) Validate() error {
	if t.Invalid != "" {
		return errors.New(t.Invalid)
	}
	return nil
}

type (
	// TypeA is a named type used to differentiate test messages.
	TypeA string
	// TypeB is a named type used to differentiate test messages.
	TypeB string
	// TypeC is a named type used to differentiate test messages.
	TypeC string
	// TypeD is a named type used to differentiate test messages.
	TypeD string
	// TypeE is a named type used to differentiate test messages.
	TypeE string
	// TypeF is a named type used to differentiate test messages.
	TypeF string
	// TypeG is a named type used to differentiate test messages.
	TypeG string
	// TypeH is a named type used to differentiate test messages.
	TypeH string
	// TypeI is a named type used to differentiate test messages.
	TypeI string
	// TypeJ is a named type used to differentiate test messages.
	TypeJ string
	// TypeK is a named type used to differentiate test messages.
	TypeK string
	// TypeL is a named type used to differentiate test messages.
	TypeL string
	// TypeM is a named type used to differentiate test messages.
	TypeM string
	// TypeN is a named type used to differentiate test messages.
	TypeN string
	// TypeO is a named type used to differentiate test messages.
	TypeO string
	// TypeP is a named type used to differentiate test messages.
	TypeP string
	// TypeQ is a named type used to differentiate test messages.
	TypeQ string
	// TypeR is a named type used to differentiate test messages.
	TypeR string
	// TypeS is a named type used to differentiate test messages.
	TypeS string
	// TypeT is a named type used to differentiate test messages.
	TypeT string
	// TypeU is a named type used to differentiate test messages.
	TypeU string
	// TypeV is a named type used to differentiate test messages.
	TypeV string
	// TypeW is a named type used to differentiate test messages.
	TypeW string
	// TypeX is a named type used to differentiate test messages.
	TypeX string
	// TypeY is a named type used to differentiate test messages.
	TypeY string
	// TypeZ is a named type used to differentiate test messages.
	TypeZ string
)
