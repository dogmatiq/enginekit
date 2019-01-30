package handler

import (
	"github.com/dogmatiq/enginekit/message"
)

// Type is an enumeration of the types of handlers.
type Type string

const (
	// AggregateType is the handler type for dogma.AggregateMessageHandler.
	AggregateType Type = "aggregate"

	// ProcessType is the handler type for dogma.ProcessMessageHandler.
	ProcessType Type = "process"

	// IntegrationType is the handler type for dogma.IntegrationMessageHandler.
	IntegrationType Type = "integration"

	// ProjectionType is the handler type for dogma.ProjectionMessageHandler.
	ProjectionType Type = "projection"
)

// Types is a slice of the valid handler types.
var Types = []Type{
	AggregateType,
	ProcessType,
	IntegrationType,
	ProjectionType,
}

// MustValidate panics if r is not a valid message role.
func (t Type) MustValidate() {
	switch t {
	case AggregateType:
	case ProcessType:
	case IntegrationType:
	case ProjectionType:
	default:
		panic("invalid type: " + t)
	}
}

// Is returns true if t is one of the given types.
func (t Type) Is(types ...Type) bool {
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
func (t Type) MustBe(types ...Type) {
	if !t.Is(types...) {
		panic("unexpected type: " + t)
	}
}

// MustNotBe panics if t is one of the given types.
func (t Type) MustNotBe(types ...Type) {
	if t.Is(types...) {
		panic("unexpected type: " + t)
	}
}

// IsConsumerOf returns true if handlers of type t can consume messages with the
// given role.
func (t Type) IsConsumerOf(r message.Role) bool {
	return r.Is(t.Consumes()...)
}

// IsProducerOf returns true if handlers of type t can produce messages with the
// given role.
func (t Type) IsProducerOf(r message.Role) bool {
	return r.Is(t.Produces()...)
}

// Consumes returns the roles of messages that can be consumed by handlers of
// this type.
func (t Type) Consumes() []message.Role {
	t.MustValidate()

	switch t {
	case AggregateType:
		return []message.Role{message.CommandRole}
	case ProcessType:
		return []message.Role{message.EventRole, message.TimeoutRole}
	case IntegrationType:
		return []message.Role{message.CommandRole}
	default: // ProjectionType
		return []message.Role{message.EventRole}
	}
}

// Produces returns the roles of messages that can be produced by handlers of
// this type.
func (t Type) Produces() []message.Role {
	t.MustValidate()

	switch t {
	case AggregateType:
		return []message.Role{message.EventRole}
	case ProcessType:
		return []message.Role{message.CommandRole, message.TimeoutRole}
	case IntegrationType:
		return []message.Role{message.EventRole}
	default: // ProjectionType
		return nil
	}
}

// ConsumersOf returns the handler types that can consume messages with the
// given role.
func ConsumersOf(r message.Role) []Type {
	r.MustValidate()

	switch r {
	case message.CommandRole:
		return []Type{AggregateType, IntegrationType}
	case message.EventRole:
		return []Type{ProcessType, ProjectionType}
	default: // message.TimeoutRole
		return []Type{ProcessType}
	}
}

// ProducersOf returns the handler types that can produces messages with the
// given role.
func ProducersOf(r message.Role) []Type {
	switch r {
	case message.CommandRole:
		return []Type{ProcessType}
	case message.EventRole:
		return []Type{AggregateType, IntegrationType}
	default: // message.TimeoutRole
		return []Type{ProcessType}
	}
}
