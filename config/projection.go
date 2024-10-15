package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
)

// ProjectionAsConfigured contains the raw unvalidated properties of a
// [Projection].
type ProjectionAsConfigured struct {
	// Source describes the type and value that produced the configuration, if
	// available.
	Source optional.Optional[Source[dogma.ProjectionMessageHandler]]

	// Identities is the list of identities configured for the handler.
	Identities []Identity

	// Routes is the list of routes configured on the handler.
	Routes []Route

	// IsDisabled is true if the handler was disabled via the configurer, if
	// known.
	IsDisabled optional.Optional[bool]

	// DeliveryPolicy is the delivery policy for the handler, if configured.
	DeliveryPolicy optional.Optional[ProjectionDeliveryPolicy]

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Projection represents the (potentially invalid) configuration of a
// [dogma.ProjectionMessageHandler] implementation.
type Projection struct {
	AsConfigured ProjectionAsConfigured
}

func (h *Projection) String() string {
	return renderEntity("projection", h, h.AsConfigured.Source)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h *Projection) Identity() *identitypb.Identity {
	return finalizeIdentity(newFinalizeContext(h), h)
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
	return finalizeRouteSet(newFinalizeContext(h), h)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h *Projection) IsDisabled() bool {
	return h.AsConfigured.IsDisabled.Get()
}

// DeliveryPolicy returns the delivery policy for the handler.
func (h *Projection) DeliveryPolicy() dogma.ProjectionDeliveryPolicy {
	if p, ok := h.AsConfigured.DeliveryPolicy.TryGet(); ok {
		return p.Implementation.Get()
	}
	return dogma.UnicastProjectionDeliveryPolicy{}
}

// Interface returns the [dogma.ProjectionMessageHandler] instance that the
// configuration represents, or panics if it is not available.
func (h *Projection) Interface() dogma.ProjectionMessageHandler {
	return h.AsConfigured.Source.Get().Interface.Get()
}

func (h *Projection) normalize(ctx *normalizeContext) Component {
	h.AsConfigured.Identities = normalizeIdentities(ctx, h)
	h.AsConfigured.Routes = normalizeRoutes(ctx, h)
	return h
}

func (h *Projection) identitiesAsConfigured() []Identity {
	return h.AsConfigured.Identities
}

func (h *Projection) routesAsConfigured() []Route {
	return h.AsConfigured.Routes
}

// ProjectionDeliveryPolicy represents the (potentially invalid) configuration
// of a [dogma.ProjectionDeliveryPolicy].
type ProjectionDeliveryPolicy struct {
	// TypeName is the fully-qualified name of the Go type that implements
	// [dogma.DeliveryPolicy], if available.
	TypeName optional.Optional[string]

	// Implementation is the value that produced the configuration, if
	// available.
	Implementation optional.Optional[dogma.ProjectionDeliveryPolicy]
}