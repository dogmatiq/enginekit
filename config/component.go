package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// A Component is some element of the configuration of a Dogma application.
type Component interface {
	fmt.Stringer

	// Fidelity returns information about how well the configuration represents
	// the actual configuration that would be used at runtime.
	Fidelity() Fidelity

	normalize(*normalizeContext) Component
}

// An Entity is a [Component] that represents the configuration of some
// configurable Dogma entity; that is, any type with a Configure() method that
// accepts one of the Dogma "configurer" interfaces.
type Entity interface {
	Component

	// Identity returns the entity's identity.
	//
	// It panics if no single valid identity is configured.
	Identity() *identitypb.Identity

	// RouteSet returns the routes configured for the entity.
	//
	// It panics if the route configuration is incomplete or invalid.
	RouteSet() RouteSet

	identitiesAsConfigured() []*Identity
}

// A Handler is a specialization of [Entity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	Entity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler was disabled via the configurer.
	IsDisabled() bool

	routesAsConfigured() []Route
}

// Fidelity describes how well a [Component] configuration represents the actual
// configuration that would be when running an application.
type Fidelity struct {
	// IsPartial is true if some configuration logic was not applied when
	// building the configuration.
	//
	// It is false if all of the _available_ configuration logic was applied.
	// This does not imply that all _mandatory_ configuration is present.
	IsPartial bool

	// IsSpeculative is true if the component is only included in the
	// configuration under certain conditions and those conditions could not be
	// evaluated at the time the configuration was built.
	IsSpeculative bool

	// IsUnresolved is true if any of the component's configuration values
	// could not be determined at the time the configuration was built.
	IsUnresolved bool
}

var (
	_ Component = (*Identity)(nil)
	_ Component = (*Route)(nil)

	_ Entity = (*Application)(nil)

	_ Handler = (*Aggregate)(nil)
	_ Handler = (*Process)(nil)
	_ Handler = (*Integration)(nil)
	_ Handler = (*Projection)(nil)
)
