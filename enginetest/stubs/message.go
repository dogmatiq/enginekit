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
