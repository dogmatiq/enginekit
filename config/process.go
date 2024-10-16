package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/internal/ioutil"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// ProcessAsConfigured contains the raw unvalidated properties of a
// [Process].
type ProcessAsConfigured struct {
	// Source describes the type and value that produced the configuration.
	Source Value[dogma.ProcessMessageHandler]

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
	return RenderDescriptor(h)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Process) Identity() *identitypb.Identity {
	return buildIdentity(strictContext(h), h.AsConfigured.Identities)
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
	return buildRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Process) IsDisabled() bool {
	return h.AsConfigured.IsDisabled.Get()
}

// Interface returns the [dogma.ProcessMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Process) Interface() dogma.ProcessMessageHandler {
	return h.AsConfigured.Source.Value.Get()
}

func (h *Process) clone() Component {
	clone := &Process{h.AsConfigured}
	cloneInPlace(&clone.AsConfigured.Identities)
	cloneInPlace(&clone.AsConfigured.Routes)
	return clone
}

func (h *Process) normalize(ctx *normalizationContext) {
	normalizeValue(ctx, &h.AsConfigured.Source, &h.AsConfigured.Fidelity)
	normalizeIdentities(ctx, h.AsConfigured.Identities)
	normalizeRoutes(ctx, h, h.AsConfigured.Routes)
}

func (h *Process) identities() []*Identity {
	return h.AsConfigured.Identities
}

func (h *Process) routes() []*Route {
	return h.AsConfigured.Routes
}

func (h *Process) renderDescriptor(ren *ioutil.Renderer) {
	renderEntityDescriptor(ren, "process", h, h.AsConfigured.Source)
}

func (h *Process) renderDetails(*ioutil.Renderer) {
	panic("not implemented")
}
