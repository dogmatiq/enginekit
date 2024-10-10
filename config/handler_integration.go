package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Integration represents the (potentially invalid) configuration of a
// [dogma.IntegrationMessageHandler] implementation.
type Integration struct {
	// ConfigurationSource contains information about the type and value that
	// produced the configuration, if available.
	ConfigurationSource optional.Optional[Source[dogma.IntegrationMessageHandler]]

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

func (h Integration) String() string {
	return stringify("integration", h, h.ConfigurationSource)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h Integration) Identity() Identity {
	return normalizedIdentity(h)
}

// IsExhaustive returns true if the entire configuration was loaded.
func (h Integration) IsExhaustive() bool {
	return h.ConfigurationIsExhaustive
}

// HandlerType returns [HandlerType] of the handler.
func (h Integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Integration) Routes() RouteSet {
	return normalizedRouteSet(h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h Integration) IsDisabled() bool {
	return h.ConfiguredAsDisabled
}

// Interface returns the [dogma.IntegrationMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h Integration) Interface() dogma.IntegrationMessageHandler {
	return h.ConfigurationSource.Get().Value.Get()
}

func (h Integration) normalize(ctx *normalizationContext) Component {
	h.ConfiguredIdentities = normalizeIdentities(ctx, h)
	h.ConfiguredRoutes = normalizeRoutes(ctx, h)
	return h
}

func (h Integration) identities() []Identity {
	return h.ConfiguredIdentities
}

func (h Integration) routes() []Route {
	return h.ConfiguredRoutes
}