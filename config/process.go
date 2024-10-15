package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// ProcessAsConfigured contains the raw unvalidated properties of a
// [Process].
type ProcessAsConfigured struct {
	// Source describes the type and value that produced the configuration, if
	// available.
	Source optional.Optional[Value[dogma.ProcessMessageHandler]]

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

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	AsConfigured ProcessAsConfigured
}

func (h *Process) String() string {
	return renderEntity("process", h, h.AsConfigured.Source)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Process) Identity() *identitypb.Identity {
	return finalizeIdentity(newFinalizeContext(h), h)
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (h *Process) Fidelity() Fidelity {
	return h.AsConfigured.Fidelity
}

// HandlerType returns [HandlerType] of the handler.
func (h *Process) HandlerType() HandlerType {
	return ProcessHandlerType
}

// RouteSet returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h *Process) RouteSet() RouteSet {
	return finalizeRouteSet(newFinalizeContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Process) IsDisabled() bool {
	return h.AsConfigured.IsDisabled.Get()
}

// Interface returns the [dogma.ProcessMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Process) Interface() dogma.ProcessMessageHandler {
	return h.AsConfigured.Source.Get().Value.Get()
}

func (h *Process) normalize(ctx *normalizeContext) Component {
	h.AsConfigured.Fidelity, h.AsConfigured.Source = normalizeValue(ctx, h.AsConfigured.Fidelity, h.AsConfigured.Source)
	h.AsConfigured.Identities = normalizeIdentities(ctx, h)
	h.AsConfigured.Routes = normalizeRoutes(ctx, h)
	return h
}

func (h *Process) identitiesAsConfigured() []*Identity {
	return h.AsConfigured.Identities
}

func (h *Process) routesAsConfigured() []*Route {
	return h.AsConfigured.Routes
}
