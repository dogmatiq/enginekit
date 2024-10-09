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

	// DeliveryPolicy is the delivery policy for the handler, if configured.
	DeliveryPolicy optional.Optional[ProjectionDeliveryPolicy]

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool

	// IsExhaustive is true if the complete configuration was loaded. It may be
	// false, for example, when attempting to load configuration using static
	// analysis, but the code depends on runtime type information.
	IsExhaustive bool
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

// Routes returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h Projection) Routes() []Route {
	return normalizeRoutes(h)
}

func (h Projection) configuredIdentities() []Identity { return h.ConfiguredIdentities }
func (h Projection) configuredRoutes() []Route        { return h.ConfiguredRoutes }

func (h Projection) normalize(opts validationOptions) (_ Entity, errs error) {
	normalizeIdentitiesInPlace(opts, h, &errs, &h.ConfiguredIdentities)

	return h, errs
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
