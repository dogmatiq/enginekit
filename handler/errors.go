package handler

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/identity"
	"github.com/dogmatiq/enginekit/message"
)

// EmptyInstanceIDError indicates that an aggregate or process message handler has
// attempted to route a message to an instance with an empty ID.
type EmptyInstanceIDError struct {
	// Handler is the identity of the handler that caused the error.
	Handler identity.Identity

	// HandlerType is the type of handler that caused the error.
	HandlerType Type

	// Message is the message that was being handled when the error occurred.
	Message dogma.Message
}

func (e EmptyInstanceIDError) Error() string {
	return fmt.Sprintf(
		"the '%s' %s message handler attempted to route a %s message to an empty instance ID",
		e.Handler.Name,
		e.HandlerType,
		message.TypeOf(e.Message),
	)
}

// NilRootError indicates that an aggregate or process message handler has
// returned a nil "root" value from its New() method.
type NilRootError struct {
	// Handler is the identity of the handler that caused the error.
	Handler identity.Identity

	// HandlerType is the type of handler that caused the error.
	HandlerType Type
}

func (e NilRootError) Error() string {
	return fmt.Sprintf(
		"the '%s' %s message handler returned a nil root from New()",
		e.Handler.Name,
		e.HandlerType,
	)
}

// EventNotRecordedError indicates that an aggregate instance was created
// or destroyed without recording an event.
type EventNotRecordedError struct {
	// Handler is the identity of the handler that caused the error.
	Handler identity.Identity

	// WasDestroyed is true if the error occurred as a result of the description
	// of an aggregate, as opposed to creation.
	WasDestroyed bool

	// Message is the message that was being handled when the error occurred.
	Message dogma.Message

	// InstanceID is the aggregate instance ID that the message was routed to.
	InstanceID string
}

func (e EventNotRecordedError) Error() string {
	s := "created"

	if e.WasDestroyed {
		s = "destroyed"
	}

	return fmt.Sprintf(
		"the '%s' aggregate message handler %s the '%s' instance without recording an event while handling a %s command",
		e.Handler.Name,
		s,
		e.InstanceID,
		message.TypeOf(e.Message),
	)
}

// UnexpectedMessageError indicates that a message handler has panicked with a
// dogma.UnexpectedMessage value.
type UnexpectedMessageError struct {
	// Handler is the identity of the handler that caused the error.
	Handler identity.Identity

	// HandlerType is the type of handler that caused the error.
	HandlerType Type

	// Message is the message that was being handled when the error occurred.
	Message dogma.Message

	// InstanceID is the aggregate instance ID that the message was routed to.
	// It is empty if HandleType is not AggregateType or ProcessType, or
	// otherwise if the error occurred outside the context of any specific
	// instance.
	InstanceID string

	// StackTrace is the result of debug.Stack(), captured while recovering from
	// the panic.
	StackTrace []byte
}

func (e UnexpectedMessageError) Error() string {
	return fmt.Sprintf(
		"the '%s' %s message handler does not expect %s messages",
		e.Handler.Name,
		e.HandlerType,
		message.TypeOf(e.Message),
	)
}
