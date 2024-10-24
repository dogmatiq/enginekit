package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
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

	// DisabledFlags represents calls to [dogma.AggregateConfigurer.Disable].
	DisabledFlags FlagSet[Disabled]

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	AsConfigured ProcessAsConfigured
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
	return h.AsConfigured.DisabledFlags.resolve(h).Get()
}

// Interface returns the [dogma.ProcessMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Process) Interface() dogma.ProcessMessageHandler {
	return h.AsConfigured.Source.Value.Get()
}

func (h *Process) String() string {
	return RenderDescriptor(h)
}

func (h *Process) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, "process", h.AsConfigured.Source)
}

func (h *Process) renderDetails(ren *renderer.Renderer) {
	renderHandlerDetails(ren, h, h.AsConfigured.Source)
}

func (h *Process) disabledFlags() FlagSet[Disabled] {
	return h.AsConfigured.DisabledFlags
}

func (h *Process) identities() []*Identity {
	return h.AsConfigured.Identities
}

func (h *Process) routes() []*Route {
	return h.AsConfigured.Routes
}

func (h *Process) clone() Component {
	clone := &Process{h.AsConfigured}
	cloneInPlace(&clone.AsConfigured.Identities)
	cloneInPlace(&clone.AsConfigured.Routes)
	return clone
}

func (h *Process) normalize(ctx *normalizationContext) {
	normalizeValue(ctx, &h.AsConfigured.Source, &h.AsConfigured.Fidelity)

	normalizeChildren(ctx, h.AsConfigured.Identities)
	normalizeChildren(ctx, h.AsConfigured.Routes)

	reportIdentityErrors(ctx, h.AsConfigured.Identities)
	reportRouteErrors(ctx, h, h.AsConfigured.Routes)
}
