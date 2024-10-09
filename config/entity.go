package config

import "github.com/dogmatiq/enginekit/optional"

// Entity represents the configuration of some [Configurable] Dogma entity.
type Entity interface {
	// Identity returns the entity's identity.
	//
	// It panics if no single valid identity is configured.
	Identity() Identity

	normalize(validationOptions) (Entity, error)
	configuredIdentities() []Identity
}

// Handler is an interface for configuration of a Dogma message handler.
type Handler interface {
	Entity

	// Routes returns the routes configured for the handler.
	//
	// It panics if the routes are incomplete or invalid.
	Routes() []Route

	configuredRoutes() []Route
}

// Implementation contains information about the implementation of the T
// interface.
type Implementation[T any] struct {
	// TypeName is the fully-qualified name of the Go type that implements T.
	TypeName string

	// Source is the value that produced the configuration, if available.
	Source optional.Optional[T]
}

var (
	_ Entity  = Application{}
	_ Handler = Aggregate{}
	_ Handler = Integration{}
	_ Handler = Process{}
	_ Handler = Projection{}
)
