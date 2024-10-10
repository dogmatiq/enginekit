package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Projection represents the (potentially invalid) configuration of a
// [dogma.ProjectionMessageHandler] implementation.
type Projection struct {
	// Impl contains information about the type that produced the configuration,
	// if available.
	Impl optional.Optional[Implementation[dogma.ProjectionMessageHandler]]

	// ConfiguredIdentities is the list of (potentially invalid or duplicated)
	// identities configured for the handler.
	ConfiguredIdentities []Identity

	// ConfiguredRoutes is the list of (potentially invalid, incomplete or
	// duplicated) message routes configured for the handler.
	ConfiguredRoutes []Route

	// ConfiguredDeliveryPolicy is the delivery policy for the handler, if
	// configured.
	ConfiguredDeliveryPolicy optional.Optional[ProjectionDeliveryPolicy]

	// ConfiguredAsDisabled is true if the handler was disabled via the
	// configurer.
	ConfiguredAsDisabled bool

	// ConfigurationIsExhaustive is true if the entire configuration was loaded.
	ConfigurationIsExhaustive bool
}

func (h Projection) String() string {
	return stringify("projection", h.Impl, h.ConfiguredIdentities)
}

// Identity returns the entity's identity.
//
// It panics if no single valid identity is configured.
func (h Projection) Identity() Identity {
	return normalizedIdentity(h)
}

// IsExhaustive returns true if the entire configuration was loaded.
func (h Projection) IsExhaustive() bool {
	return h.ConfigurationIsExhaustive
}

// HandlerType returns [HandlerType] of the handler.
func (h Projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Projection) Routes(filter ...RouteType) []Route {
	return normalizedRoutes(h, filter...)
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h Projection) IsDisabled() bool {
	return h.ConfiguredAsDisabled
}

// DeliveryPolicy returns the delivery policy for the handler.
func (h Projection) DeliveryPolicy() dogma.ProjectionDeliveryPolicy {
	if p, ok := h.ConfiguredDeliveryPolicy.TryGet(); ok {
		return p.Implementation.Get()
	}
	return dogma.UnicastProjectionDeliveryPolicy{}
}

func (h Projection) normalize(ctx *normalizationContext) Component {
	h.ConfiguredIdentities = normalizeIdentities(ctx, h)
	h.ConfiguredRoutes = normalizeRoutes(ctx, h)
	return h
}

func (h Projection) identities() []Identity {
	return h.ConfiguredIdentities
}

func (h Projection) routes() []Route {
	return h.ConfiguredRoutes
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
