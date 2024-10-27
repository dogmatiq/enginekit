package config

import (
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// Component is the "top-level" interface for the individual elements that form
// a complete configuration of a Dogma application or handler.
type Component interface {
	// Fidelity reports how faithfully the [Component] describes a complete
	// configuration that can be used to execute an application.
	Fidelity() Fidelity

	// Baggage is a collection of arbitrary data that is associated with the
	// [Component] by whatever system produced the configuration.
	Baggage() Baggage
}

// An Entity is a [Component] that that represents the configuration of a Dogma
// entity, which is a type that implements one of the following interfaces:
//
//   - [dogma.Application]
//   - [dogma.AggregateMessageHandler]
//   - [dogma.ProcessMessageHandler]
//   - [dogma.IntegrationMessageHandler]
//   - [dogma.ProjectionMessageHandler]
//
// See [Handler] for a more specific interface that represents the
// configuration of a Dogma handler.
type Entity interface {
	Component

	// Identity returns the entity's identity.
	//
	// It panics the configuration does not specify a singular valid identity.
	Identity() *identitypb.Identity

	// RouteSet returns the routes configured for the entity.
	//
	// It panics if the configuration does not specify a complete set of valid
	// routes for the entity and its constituents.
	RouteSet() RouteSet
}

// A Handler is an [Entity] that represents the configuration of a Dogma
// handler.
type Handler interface {
	Entity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	//
	// It panics if the configuration does not specify unambiguously whether the
	// handler is enabled or disabled.
	IsDisabled() bool
}

var (
	_ Component = (*Identity)(nil)

	_ Component = (*FlagModification)(nil)

	_ Entity = (*Application)(nil)

	_ Handler = (*Aggregate)(nil)
	_ Handler = (*Process)(nil)
	_ Handler = (*Integration)(nil)
	_ Handler = (*Projection)(nil)

// _ Component = (*Route)(nil)
)

// ComponentCommon is a partial implementation of [Component].
type ComponentCommon struct {
	ComponentFidelity Fidelity
	ComponentBaggage  Baggage
}

// Fidelity reports how faithfully the [Component] describes a complete
// configuration that can be used to execute an application.
func (c *ComponentCommon) Fidelity() Fidelity {
	return c.ComponentFidelity
}

// Baggage returns a collection of arbitrary data that is associated with the
// [Component] by whatever system produced the configuration.
func (c *ComponentCommon) Baggage() Baggage {
	return c.ComponentBaggage
}

// EntityCommon is a partial implementation of [Entity].
type EntityCommon[T any] struct {
	ComponentCommon

	SourceTypeName     string
	Source             optional.Optional[T]
	IdentityComponents []*Identity
}

// Identity returns the entity's identity.
//
// It panics the configuration does not specify a singular valid identity.
func (e *EntityCommon[T]) Identity() *identitypb.Identity {
	panic("not implemented")
}

// HandlerCommon is a partial implementation of [Handler].
type HandlerCommon[T any] struct {
	EntityCommon[T]

	RouteComponents []*Route
	DisabledFlag    Flag[Disabled]
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the configuration does not specify a complete set of valid
// routes for the entity and its constituents.
func (h *HandlerCommon[T]) RouteSet() RouteSet {
	panic("not implemented")
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the configuration does not specify unambiguously whether the
// handler is enabled or disabled.
func (h *HandlerCommon[T]) IsDisabled() bool {
	panic("not implemented")
}
