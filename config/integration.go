package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Integration represents the (potentially invalid) configuration of a
// [dogma.IntegrationMessageHandler] implementation.
type Integration struct {
	// Implementation contains information about the type that produced the
	// configuration, if available.
	Implementation optional.Optional[Implementation[dogma.IntegrationMessageHandler]]

	// Identity is the set of identities configured for the handler.
	Identities []Identity

	// Routes is the set of message routes to and from the handler.
	Routes []Route

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool

	// IsExhaustive indicates whether the complete configuration was loaded. It
	// is false when it cannot be guaranteed that the configuration is complete,
	// which is possible, for example, when attempting to load configuration by
	// static analysis.
	IsExhaustive bool
}

func (h Integration) String() string {
	return stringify("integration", h.Implementation, h.Identities)
}
