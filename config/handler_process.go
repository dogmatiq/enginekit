package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	// ConfigurationSource contains information about the type and value that
	// produced the configuration, if available.
	ConfigurationSource optional.Optional[Source[dogma.ProcessMessageHandler]]

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

func (h Process) String() string {
	return stringify("process", h, h.ConfigurationSource)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h Process) Identity() Identity {
	return normalizedIdentity(h)
}

// IsExhaustive returns true if the entire configuration was loaded.
func (h Process) IsExhaustive() bool {
	return h.ConfigurationIsExhaustive
}

// HandlerType returns [HandlerType] of the handler.
func (h Process) HandlerType() HandlerType {
	return ProcessHandlerType
}

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Process) Routes() RouteSet {
	return normalizedRouteSet(h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h Process) IsDisabled() bool {
	return h.ConfiguredAsDisabled
}

// Interface returns the [dogma.ProcessMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h Process) Interface() dogma.ProcessMessageHandler {
	return h.ConfigurationSource.Get().Value.Get()
}

func (h Process) normalize(ctx *normalizationContext) Component {
	h.ConfiguredIdentities = normalizeIdentities(ctx, h)
	h.ConfiguredRoutes = normalizeRoutes(ctx, h)
	return h
}

func (h Process) identities() []Identity {
	return h.ConfiguredIdentities
}

func (h Process) routes() []Route {
	return h.ConfiguredRoutes
}
