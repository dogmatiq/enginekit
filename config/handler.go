package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/message"
)

// HandlerType is an enumeration of the types of handlers.
type HandlerType string

const (
	// AggregateHandlerType is the handler type for dogma.AggregateMessageHandler.
	AggregateHandlerType HandlerType = "aggregate"

	// ProcessHandlerType is the handler type for dogma.ProcessMessageHandler.
	ProcessHandlerType HandlerType = "process"

	// IntegrationHandlerType is the handler type for dogma.IntegrationMessageHandler.
	IntegrationHandlerType HandlerType = "integration"

	// ProjectionHandlerType is the handler type for dogma.ProjectionMessageHandler.
	ProjectionHandlerType HandlerType = "projection"
)

// HandlerTypes is a slice of the valid handler types.
var HandlerTypes = []HandlerType{
	AggregateHandlerType,
	ProcessHandlerType,
	IntegrationHandlerType,
	ProjectionHandlerType,
}

// Validate returns an error if r is not a valid message role.
func (t HandlerType) Validate() error {
	switch t {
	case AggregateHandlerType,
		ProcessHandlerType,
		IntegrationHandlerType,
		ProjectionHandlerType:
		return nil
	default:
		return fmt.Errorf("invalid handler type: %s", t)
	}
}

// MustValidate panics if r is not a valid message role.
func (t HandlerType) MustValidate() {
	if err := t.Validate(); err != nil {
		panic(err)
	}
}

// Is returns true if t is one of the given types.
func (t HandlerType) Is(types ...HandlerType) bool {
	t.MustValidate()

	for _, x := range types {
		x.MustValidate()

		if t == x {
			return true
		}
	}

	return false
}

// MustBe panics if t is not one of the given types.
func (t HandlerType) MustBe(types ...HandlerType) {
	if !t.Is(types...) {
		panic("unexpected type: " + t)
	}
}

// MustNotBe panics if t is one of the given types.
func (t HandlerType) MustNotBe(types ...HandlerType) {
	if t.Is(types...) {
		panic("unexpected type: " + t)
	}
}

// IsConsumerOf returns true if handlers of type t can consume messages with the
// given role.
func (t HandlerType) IsConsumerOf(r message.Role) bool {
	return r.Is(t.Consumes()...)
}

// IsProducerOf returns true if handlers of type t can produce messages with the
// given role.
func (t HandlerType) IsProducerOf(r message.Role) bool {
	return r.Is(t.Produces()...)
}

// Consumes returns the roles of messages that can be consumed by handlers of
// this type.
func (t HandlerType) Consumes() []message.Role {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []message.Role{message.CommandRole}
	case ProcessHandlerType:
		return []message.Role{message.EventRole, message.TimeoutRole}
	case IntegrationHandlerType:
		return []message.Role{message.CommandRole}
	default: // ProjectionHandlerType
		return []message.Role{message.EventRole}
	}
}

// Produces returns the roles of messages that can be produced by handlers of
// this type.
func (t HandlerType) Produces() []message.Role {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []message.Role{message.EventRole}
	case ProcessHandlerType:
		return []message.Role{message.CommandRole, message.TimeoutRole}
	case IntegrationHandlerType:
		return []message.Role{message.EventRole}
	default: // ProjectionHandlerType
		return nil
	}
}

// ShortString returns a short (3-character) representation of the handler type.
func (t HandlerType) ShortString() string {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return "agg"
	case ProcessHandlerType:
		return "prc"
	case IntegrationHandlerType:
		return "int"
	default: // ProjectionHandlerType
		return "prj"
	}
}

func (t HandlerType) String() string {
	return string(t)
}

// ConsumersOf returns the handler types that can consume messages with the
// given role.
func ConsumersOf(r message.Role) []HandlerType {
	r.MustValidate()

	switch r {
	case message.CommandRole:
		return []HandlerType{AggregateHandlerType, IntegrationHandlerType}
	case message.EventRole:
		return []HandlerType{ProcessHandlerType, ProjectionHandlerType}
	default: // message.TimeoutRole
		return []HandlerType{ProcessHandlerType}
	}
}

// ProducersOf returns the handler types that can produces messages with the
// given role.
func ProducersOf(r message.Role) []HandlerType {
	switch r {
	case message.CommandRole:
		return []HandlerType{ProcessHandlerType}
	case message.EventRole:
		return []HandlerType{AggregateHandlerType, IntegrationHandlerType}
	default: // message.TimeoutRole
		return []HandlerType{ProcessHandlerType}
	}
}
