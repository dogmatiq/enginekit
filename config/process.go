package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Process represents the (potentially invalid) configuration of a
// [dogma.ProcessMessageHandler] implementation.
type Process struct {
	// Implementation contains information about the type that produced the
	// configuration, if available.
	Implementation optional.Optional[Implementation[dogma.ProcessMessageHandler]]

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

func (h Process) String() string {
	return stringify("process", h.Implementation, h.Identities)
}
