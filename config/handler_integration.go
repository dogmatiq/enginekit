package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
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

	// ConfigurationFidelity describes the configuration's accuracy in
	// comparison to the actual configuration that would be used at runtime.
	ConfigurationFidelity Fidelity
}

func (h *Integration) String() string {
	return renderEntity("integration", h, h.ConfigurationSource)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Integration) Identity() *identitypb.Identity {
	return finalizeIdentity(newFinalizeContext(h), h)
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (h *Integration) Fidelity() Fidelity {
	return h.ConfigurationFidelity
}

// HandlerType returns [HandlerType] of the handler.
func (h *Integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

// RouteSet returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h *Integration) RouteSet() RouteSet {
	return finalizeRouteSet(newFinalizeContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Integration) IsDisabled() bool {
	return h.ConfiguredAsDisabled
}

// Interface returns the [dogma.IntegrationMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Integration) Interface() dogma.IntegrationMessageHandler {
	return h.ConfigurationSource.Get().Interface.Get()
}

func (h *Integration) normalize(ctx *normalizeContext) Component {
	h.ConfiguredIdentities = normalizeIdentities(ctx, h)
	h.ConfiguredRoutes = normalizeRoutes(ctx, h)
	return h
}

func (h *Integration) identities() []Identity {
	return h.ConfiguredIdentities
}

func (h *Integration) routes() []Route {
	return h.ConfiguredRoutes
}
