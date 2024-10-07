package config

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/optional"
)

// Application represents the (potentially invalid) configuration of a
// [dogma.Application] implementation.
type Application struct {
	// TypeName is the fully-qualified name of the Go type that implements
	// [dogma.Application], if available.
	TypeName optional.Optional[string]

	// Implementation is the value that produced the configuration, if
	// available.
	Implementation optional.Optional[dogma.Application]

	// Identity is the (potentially invalid) identity of the application, if
	// configured.
	Identity optional.Optional[Identity]

	// Aggregates is a list of [dogma.AggregateMessageHandler] implementations
	// that are registered with the application.
	Aggregates []Aggregate

	// Processes is a list of [dogma.ProcessMessageHandler] implementations that
	// are registered with the application.
	Processes []Process

	// Integrations is a list of [dogma.IntegrationMessageHandler]
	// implementations that are registered with the application.
	Integrations []Integration

	// Projections is a list of [dogma.ProjectionMessageHandler] implementations
	// that are registered with the application.
	Projections []Projection
}
