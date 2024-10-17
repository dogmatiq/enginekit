package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// ProjectionAsConfigured contains the raw unvalidated properties of a
// [Projection].
type ProjectionAsConfigured struct {
	// Source describes the type and value that produced the configuration.
	Source Value[dogma.ProjectionMessageHandler]

	// Identities is the list of identities configured for the handler.
	Identities []*Identity

	// Routes is the list of routes configured on the handler.
	Routes []*Route

	// IsDisabled is true if the handler was disabled via the configurer, if
	// known.
	IsDisabled optional.Optional[bool]

	// DeliveryPolicy is the delivery policy for the handler, if configured.
	DeliveryPolicy optional.Optional[Value[dogma.ProjectionDeliveryPolicy]]

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Projection represents the (potentially invalid) configuration of a
// [dogma.ProjectionMessageHandler] implementation.
type Projection struct {
	AsConfigured ProjectionAsConfigured
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Projection) Identity() *identitypb.Identity {
	return buildIdentity(strictContext(h), h.AsConfigured.Identities)
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (h *Projection) Fidelity() Fidelity {
	return h.AsConfigured.Fidelity
}

// HandlerType returns [HandlerType] of the handler.
func (h *Projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

// RouteSet returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h *Projection) RouteSet() RouteSet {
	return buildRouteSet(strictContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Projection) IsDisabled() bool {
	return h.AsConfigured.IsDisabled.Get()
}

// DeliveryPolicy returns the delivery policy for the handler.
func (h *Projection) DeliveryPolicy() dogma.ProjectionDeliveryPolicy {
	return h.AsConfigured.DeliveryPolicy.Get().Value.Get()
}

// Interface returns the [dogma.ProjectionMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Projection) Interface() dogma.ProjectionMessageHandler {
	return h.AsConfigured.Source.Value.Get()
}

func (h *Projection) String() string {
	return RenderDescriptor(h)
}

func (h *Projection) renderDescriptor(ren *renderer.Renderer) {
	renderEntityDescriptor(ren, "projection", h.AsConfigured.Source)
}

func (h *Projection) renderDetails(*renderer.Renderer) {
	panic("not implemented")
}

func (h *Projection) identities() []*Identity {
	return h.AsConfigured.Identities
}

func (h *Projection) routes() []*Route {
	return h.AsConfigured.Routes
}

func (h *Projection) clone() Component {
	clone := &Projection{h.AsConfigured}
	cloneInPlace(&clone.AsConfigured.Identities)
	cloneInPlace(&clone.AsConfigured.Routes)
	return clone
}

func (h *Projection) normalize(ctx *normalizationContext) {
	normalizeValue(ctx, &h.AsConfigured.Source, &h.AsConfigured.Fidelity)
	normalizeIdentities(ctx, h.AsConfigured.Identities)
	normalizeRoutes(ctx, h, h.AsConfigured.Routes)

	if p, ok := h.AsConfigured.DeliveryPolicy.TryGet(); ok {
		normalizeValue(ctx, &p, &h.AsConfigured.Fidelity)
		h.AsConfigured.DeliveryPolicy = optional.Some(p)
	}
}
