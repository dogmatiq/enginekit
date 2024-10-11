package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// A Component is some element of the configuration of a Dogma application.
type Component interface {
	fmt.Stringer

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

	// Routes returns the routes configured for the entity.
	//
	// It panics if the route configuration is incomplete or invalid.
	Routes() RouteSet

	// IsExhaustive returns true if the complete configuration was loaded. It
	// may be false, for example, when attempting to load configuration using
	// static analysis, but the code depends on runtime type information.
	IsExhaustive() bool

	identities() []Identity
}

// A Handler is a specialization of [Entity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	Entity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler was disabled via the configurer.
	IsDisabled() bool

	routes() []Route
}
