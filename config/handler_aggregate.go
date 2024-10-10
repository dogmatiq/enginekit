package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Aggregate represents the (potentially invalid) configuration of a
// [dogma.AggregateMessageHandler] implementation.
type Aggregate struct {
	// Impl contains information about the type that produced the configuration,
	// if available.
	Impl optional.Optional[Implementation[dogma.AggregateMessageHandler]]

	// ConfiguredIdentities is the list of (potentially invalid or duplicated)
	// identities configured for the handler.
	ConfiguredIdentities []Identity

	// ConfiguredRoutes is the list of (potentially invalid, incomplete or
	// duplicated) message routes configured for the handler.
	ConfiguredRoutes []Route

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool

	// IsExhaustive is true if the complete configuration was loaded. It may be
	// false, for example, when attempting to load configuration using static
	// analysis, but the code depends on runtime type information.
	IsExhaustive bool
}

func (h Aggregate) String() string {
	return stringify("aggregate", h.Impl, h.ConfiguredIdentities)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h Aggregate) Identity() Identity {
	return normalizedIdentity(h)
}

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Aggregate) Routes(filter ...RouteType) []Route {
	return normalizedRoutes(h, filter...)
}

// HandlerType returns [HandlerType] of the handler.
func (h Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

func (h Aggregate) normalize(ctx *normalizationContext) Component {
	h.ConfiguredIdentities = normalizeIdentities(ctx, h)
	h.ConfiguredRoutes = normalizeRoutes(ctx, h)
	return h
}

func (h Aggregate) identities() []Identity {
	return h.ConfiguredIdentities
}

func (h Aggregate) routes() []Route {
	return h.ConfiguredRoutes
}
