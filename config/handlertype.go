package config

import (
	"bytes"
	"fmt"
)

// HandlerType is an enumeration of the types of handlers.
type HandlerType byte

const (
	// AggregateHandlerType is the handler type for dogma.AggregateMessageHandler.
	AggregateHandlerType HandlerType = 'A'

	// ProcessHandlerType is the handler type for dogma.ProcessMessageHandler.
	ProcessHandlerType HandlerType = 'P'

	// IntegrationHandlerType is the handler type for dogma.IntegrationMessageHandler.
	IntegrationHandlerType HandlerType = 'I'

	// ProjectionHandlerType is the handler type for dogma.ProjectionMessageHandler.
	ProjectionHandlerType HandlerType = 'R'
)

// HandlerTypes is a slice of the valid handler types.
var HandlerTypes = []HandlerType{
	AggregateHandlerType,
	ProcessHandlerType,
	IntegrationHandlerType,
	ProjectionHandlerType,
}

const (
	aggregateHandlerTypeString   = "aggregate"
	processHandlerTypeString     = "process"
	integrationHandlerTypeString = "integration"
	projectionHandlerTypeString  = "projection"

	aggregateHandlerTypeShortString   = "AGG"
	processHandlerTypeShortString     = "PRC"
	integrationHandlerTypeShortString = "INT"
	projectionHandlerTypeShortString  = "PRJ"
)

var (
	aggregateHandlerTypeBytes   = []byte(aggregateHandlerTypeString)
	processHandlerTypeBytes     = []byte(processHandlerTypeString)
	integrationHandlerTypeBytes = []byte(integrationHandlerTypeString)
	projectionHandlerTypeBytes  = []byte(projectionHandlerTypeString)
)

// Validate returns an error if r is not a valid message role.
func (t HandlerType) Validate() error {
	switch t {
	case AggregateHandlerType,
		ProcessHandlerType,
		IntegrationHandlerType,
		ProjectionHandlerType:
		return nil
	default:
		return fmt.Errorf("invalid handler type: %#v", t)
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
		panic("unexpected handler type: " + t.String())
	}
}

// MustNotBe panics if t is one of the given types.
func (t HandlerType) MustNotBe(types ...HandlerType) {
	if t.Is(types...) {
		panic("unexpected handler type: " + t.String())
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
		return aggregateHandlerTypeShortString
	case ProcessHandlerType:
		return processHandlerTypeShortString
	case IntegrationHandlerType:
		return integrationHandlerTypeShortString
	default: // ProjectionHandlerType
		return projectionHandlerTypeShortString
	}
}

// String returns a string representation of the handler type.
func (t HandlerType) String() string {
	switch t {
	case AggregateHandlerType:
		return aggregateHandlerTypeString
	case ProcessHandlerType:
		return processHandlerTypeString
	case IntegrationHandlerType:
		return integrationHandlerTypeString
	case ProjectionHandlerType:
		return projectionHandlerTypeString
	default:
		return fmt.Sprintf("<invalid handler type %#v>", t)
	}
}

// MarshalText returns a UTF-8 representation of the handler type.
func (t HandlerType) MarshalText() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	switch t {
	case AggregateHandlerType:
		return aggregateHandlerTypeBytes, nil
	case ProcessHandlerType:
		return processHandlerTypeBytes, nil
	case IntegrationHandlerType:
		return integrationHandlerTypeBytes, nil
	default: // ProjectionHandlerType
		return projectionHandlerTypeBytes, nil
	}
}

// UnmarshalText unmarshals a type from its UTF-8 representation.
func (t *HandlerType) UnmarshalText(text []byte) error {
	if bytes.Equal(text, aggregateHandlerTypeBytes) {
		*t = AggregateHandlerType
	} else if bytes.Equal(text, processHandlerTypeBytes) {
		*t = ProcessHandlerType
	} else if bytes.Equal(text, integrationHandlerTypeBytes) {
		*t = IntegrationHandlerType
	} else if bytes.Equal(text, projectionHandlerTypeBytes) {
		*t = ProjectionHandlerType
	} else {
		return fmt.Errorf("invalid text representation of handler type: %s", text)
	}

	return nil
}

// MarshalBinary returns a binary representation of the handler type.
func (t HandlerType) MarshalBinary() ([]byte, error) {
	return []byte{byte(t)}, t.Validate()
}

// UnmarshalBinary unmarshals a type from its binary representation.
func (t *HandlerType) UnmarshalBinary(data []byte) error {
	if len(data) != 1 {
		return fmt.Errorf("invalid binary representation of handler type, expected exactly 1 byte")
	}

	*t = HandlerType(data[0])
	return t.Validate()
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
