package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Aggregate represents the (potentially invalid) configuration of a
// [dogma.AggregateMessageHandler] implementation.
type Aggregate struct {
	// TypeName is the fully-qualified name of the Go type that implements
	// [dogma.AggregateMessageHandler], if available.
	TypeName optional.Optional[string]

	// Implementation is the value that produced the configuration, if
	// available.
	Implementation optional.Optional[dogma.AggregateMessageHandler]

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

func (h Aggregate) String() string {
	return stringify("aggregate", h.TypeName, h.Identities)
}
