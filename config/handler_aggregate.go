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

	// ConfiguredAsDisabled is true if the handler was disabled via the
	// configurer.
	ConfiguredAsDisabled bool

	// ConfigurationIsExhaustive is true if the entire configuration was loaded.
	ConfigurationIsExhaustive bool
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

// IsExhaustive returns true if the entire configuration was loaded.
func (h Aggregate) IsExhaustive() bool {
	return h.ConfigurationIsExhaustive
}

// HandlerType returns [HandlerType] of the handler.
func (h Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Aggregate) Routes(filter ...RouteType) []Route {
	return normalizedRoutes(h, filter...)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h Aggregate) IsDisabled() bool {
	return h.ConfiguredAsDisabled
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
