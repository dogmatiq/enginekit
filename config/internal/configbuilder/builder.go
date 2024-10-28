package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// ComponentBuilder an interface for builders that produce a [config.Component].
type ComponentBuilder[T config.Component] interface {
	Done() T
}

// EntityBuilder an interface for builders that produce a [config.Entity].
type EntityBuilder[T config.Entity] interface {
	ComponentBuilder[T]

	// Identity calls fn which configures a [config.Identity] that is added to
	// the handler.
	Identity(fn func(*IdentityBuilder))
}

// HandlerBuilder an interface for builders that produce a [config.Handler].
type HandlerBuilder[T config.Handler] interface {
	EntityBuilder[T]
}

var (
	_ ComponentBuilder[*config.Identity]         = (*IdentityBuilder)(nil)
	_ ComponentBuilder[*config.FlagModification] = (*FlagBuilder)(nil)

	_ EntityBuilder[*config.Application] = (*ApplicationBuilder)(nil)

	_ HandlerBuilder[*config.Aggregate]   = (*AggregateBuilder)(nil)
	_ HandlerBuilder[*config.Process]     = (*ProcessBuilder)(nil)
	_ HandlerBuilder[*config.Integration] = (*IntegrationBuilder)(nil)
	_ HandlerBuilder[*config.Projection]  = (*ProjectionBuilder)(nil)

	_ ComponentBuilder[*config.Route] = (*RouteBuilder)(nil)
)

func setSourceTypeName[T any](e *config.EntityCommon[T], n string) {
	if n == "" {
		// TODO: validate that this is actually a well-formed fully-qualified
		// type name, with optional asterisk prefix.
		panic("concrete type name must not be empty")
	}

	e.SourceTypeName = optional.Some(n)
	e.Source = optional.None[T]()
}

func setSource[T any](e *config.EntityCommon[T], v T) {
	if any(v) == nil {
		panic("runtime value must not be nil")
	}

	e.SourceTypeName = optional.Some(typename.Of(v))
	e.Source = optional.Some(v)
}
