package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// ComponentBuilder is the interface shared by the builder types for all
// [config.Component] types.
type ComponentBuilder interface {
	// Fidelity returns the fidelity of the configuration.
	Fidelity() config.Fidelity

	// UpdateFidelity merges f into the fidelity of the configuration.
	UpdateFidelity(f config.Fidelity)
}

// EntityBuilder is a specialization of [ComponentBuilder] for building
// [config.Entity] configuration.
type EntityBuilder interface {
	ComponentBuilder

	// SetSourceTypeName sets the source of the configuration.
	SetSourceTypeName(typeName string)

	// Identity calls fn which configures a [config.Identity] that is added to
	// the handler.
	Identity(fn func(*IdentityBuilder))
}

// HandlerBuilder is a specialization of [EntityBuilder] for building
// [config.Handler] configuration.
type HandlerBuilder interface {
	EntityBuilder

	// Route calls fn which configures a [config.Route] that is added to the
	// handler.
	Route(fn func(*RouteBuilder))

	// IsDisabled returns whether the handler is disabled or not.
	IsDisabled() optional.Optional[bool]

	// SetDisabled sets whether the handler is disabled or not.
	SetDisabled(disabled bool)
}

var (
	_ ComponentBuilder = (*IdentityBuilder)(nil)
	_ ComponentBuilder = (*RouteBuilder)(nil)

	_ EntityBuilder = (*ApplicationBuilder)(nil)

	_ HandlerBuilder = (*AggregateBuilder)(nil)
	_ HandlerBuilder = (*ProcessBuilder)(nil)
	_ HandlerBuilder = (*IntegrationBuilder)(nil)
	_ HandlerBuilder = (*ProjectionBuilder)(nil)
)
