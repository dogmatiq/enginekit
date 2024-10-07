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

	// Identity is the (potentially invalid) identity of the handler, if
	// configured.
	Identity optional.Optional[Identity]

	// Routes is the set of message routes to and from the handler.
	Routes []Route

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool
}
