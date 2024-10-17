package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// AggregateAsConfigured contains the raw unvalidated properties of an
// [Aggregate].
type AggregateAsConfigured struct {
	// Source describes the type and value that produced the configuration.
	Source Value[dogma.AggregateMessageHandler]

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

// Aggregate represents the (potentially invalid) configuration of a
// [dogma.AggregateMessageHandler] implementation.
type Aggregate struct {
	AsConfigured AggregateAsConfigured
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Aggregate) Identity() *identitypb.Identity {
	return buildIdentity(strictContext(h), h.AsConfigured.Identities)
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (h *Aggregate) Fidelity() Fidelity {
	return h.AsConfigured.Fidelity
}

// HandlerType returns [HandlerType] of the handler.
func (h *Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

// RouteSet returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h *Aggregate) RouteSet() RouteSet {
	return buildRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Aggregate) IsDisabled() bool {
	return h.AsConfigured.IsDisabled.Get()
}

// Interface returns the [dogma.AggregateMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Aggregate) Interface() dogma.AggregateMessageHandler {
	return h.AsConfigured.Source.Value.Get()
}

func (h *Aggregate) String() string {
	return RenderDescriptor(h)
}

func (h *Aggregate) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, "aggregate", h.AsConfigured.Source)
}

func (h *Aggregate) renderDetails(*renderer.Renderer) {
	panic("not implemented")
}

func (h *Aggregate) identities() []*Identity {
	return h.AsConfigured.Identities
}

func (h *Aggregate) routes() []*Route {
	return h.AsConfigured.Routes
}

func (h *Aggregate) clone() Component {
	clone := &Aggregate{h.AsConfigured}
	cloneInPlace(&clone.AsConfigured.Identities)
	cloneInPlace(&clone.AsConfigured.Routes)
	return clone
}

func (h *Aggregate) normalize(ctx *normalizationContext) {
	normalizeValue(ctx, &h.AsConfigured.Source, &h.AsConfigured.Fidelity)
	normalizeIdentities(ctx, h.AsConfigured.Identities)
	normalizeRoutes(ctx, h, h.AsConfigured.Routes)
}
