package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	// TypeName is the fully-qualified name of the Go type that implements
	// [dogma.ProcessMessageHandler], if available.
	TypeName optional.Optional[string]

	// Implementation is the value that produced the configuration, if
	// available.
	Implementation optional.Optional[dogma.ProcessMessageHandler]

	// Identity is the set of identities configured for the handler.
	Identities []Identity

	// Routes is the set of message routes to and from the handler.
	Routes []Route

	// IsDisabled is true if the handler was disabled via the configurer.
	IsDisabled bool
}

func (h Process) String() string {
	return stringify("process", h.TypeName, h.Identities)
}
