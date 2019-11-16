package config

import (
	"fmt"
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
func (t HandlerType) IsConsumerOf(r MessageRole) bool {
	return r.Is(t.Consumes()...)
}

// IsProducerOf returns true if handlers of type t can produce messages with the
// given role.
func (t HandlerType) IsProducerOf(r MessageRole) bool {
	return r.Is(t.Produces()...)
}

// Consumes returns the roles of messages that can be consumed by handlers of
// this type.
func (t HandlerType) Consumes() []MessageRole {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []MessageRole{CommandMessageRole}
	case ProcessHandlerType:
		return []MessageRole{EventMessageRole, TimeoutMessageRole}
	case IntegrationHandlerType:
		return []MessageRole{CommandMessageRole}
	default: // ProjectionHandlerType
		return []MessageRole{EventMessageRole}
	}
}

// Produces returns the roles of messages that can be produced by handlers of
// this type.
func (t HandlerType) Produces() []MessageRole {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []MessageRole{EventMessageRole}
	case ProcessHandlerType:
		return []MessageRole{CommandMessageRole, TimeoutMessageRole}
	case IntegrationHandlerType:
		return []MessageRole{EventMessageRole}
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
func ConsumersOf(r MessageRole) []HandlerType {
	r.MustValidate()

	switch r {
	case CommandMessageRole:
		return []HandlerType{AggregateHandlerType, IntegrationHandlerType}
	case EventMessageRole:
		return []HandlerType{ProcessHandlerType, ProjectionHandlerType}
	default: // TimeoutMessageRole
		return []HandlerType{ProcessHandlerType}
	}
}

// ProducersOf returns the handler types that can produces messages with the
// given role.
func ProducersOf(r MessageRole) []HandlerType {
	switch r {
	case CommandMessageRole:
		return []HandlerType{ProcessHandlerType}
	case EventMessageRole:
		return []HandlerType{AggregateHandlerType, IntegrationHandlerType}
	default: // TimeoutMessageRole
		return []HandlerType{ProcessHandlerType}
	}
}
