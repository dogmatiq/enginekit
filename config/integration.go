package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// IntegrationAsConfigured contains the raw unvalidated properties of an
// [Integration].
type IntegrationAsConfigured struct {
	// Source describes the type and value that produced the configuration, if
	// available.
	Source optional.Optional[Value[dogma.IntegrationMessageHandler]]

	// Identities is the list of identities configured for the handler.
	Identities []*Identity

	// Routes is the list of routes configured on the handler.
	Routes []*Route

	// IsDisabled is true if the handler was disabled via the configurer, if
	// known.
	IsDisabled optional.Optional[bool]

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Integration represents the (potentially invalid) configuration of a
// [dogma.IntegrationMessageHandler] implementation.
type Integration struct {
	AsConfigured IntegrationAsConfigured
}

func (h *Integration) String() string {
	return renderEntity("integration", h, h.AsConfigured.Source)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Integration) Identity() *identitypb.Identity {
	return buildIdentity(strictContext(h), h.AsConfigured.Identities)
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (h *Integration) Fidelity() Fidelity {
	return h.AsConfigured.Fidelity
}

// HandlerType returns [HandlerType] of the handler.
func (h *Integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

// RouteSet returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h *Integration) RouteSet() RouteSet {
	return buildRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Integration) IsDisabled() bool {
	return h.AsConfigured.IsDisabled.Get()
}

// Interface returns the [dogma.IntegrationMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Integration) Interface() dogma.IntegrationMessageHandler {
	return h.AsConfigured.Source.Get().Value.Get()
}

func (h *Integration) clone() Component {
	clone := &Integration{h.AsConfigured}
	cloneInPlace(&clone.AsConfigured.Identities)
	cloneInPlace(&clone.AsConfigured.Routes)
	return clone
}

func (h *Integration) normalize(ctx *normalizationContext) {
	normalizeValue(ctx, &h.AsConfigured.Source, &h.AsConfigured.Fidelity)
	normalizeIdentities(ctx, h.AsConfigured.Identities)
	normalizeRoutes(ctx, h, h.AsConfigured.Routes)
}

func (h *Integration) identities() []*Identity {
	return h.AsConfigured.Identities
}

func (h *Integration) routes() []*Route {
	return h.AsConfigured.Routes
}
