package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Projection represents the (potentially invalid) configuration of a
// [dogma.ProjectionMessageHandler] implementation.
type Projection struct {
	// Implementation contains information about the type that produced the
	// configuration, if available.
	Implementation optional.Optional[Implementation[dogma.ProjectionMessageHandler]]

	// Identity is the set of identities configured for the handler.
	Identities []Identity

	// Routes is the set of message routes to and from the handler.
	Routes []Route

	// DeliveryPolicy is the delivery policy for the handler, if configured.
	DeliveryPolicy optional.Optional[ProjectionDeliveryPolicy]

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool

	// IsExhaustive indicates whether the complete configuration was loaded. It
	// is false when it cannot be guaranteed that the configuration is complete,
	// which is possible, for example, when attempting to load configuration by
	// static analysis.
	IsExhaustive bool
}

func (h Projection) String() string {
	return stringify("projection", h.Implementation, h.Identities)
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
