package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// ComponentBuilder an interface for builders that produce a [config.Component].
type ComponentBuilder[T config.Component] interface {
	Done() T
}

// EntityBuilder an interface for builders that produce a [config.Entity].
type EntityBuilder[T config.Entity, E any] interface {
	ComponentBuilder[T]

	// TypeName sets the name of the concrete type that implements the
	// entity.
	TypeName(n string)

	// Source sets the source value to h.
	Source(E)

	// Identity calls fn which configures a [config.Identity] that is added to
	// the handler.
	Identity(fn func(*IdentityBuilder))
}

// HandlerBuilder an interface for builders that produce a [config.Handler].
type HandlerBuilder[T config.Handler, H any] interface {
	EntityBuilder[T, H]
}

var (
	_ ComponentBuilder[*config.Identity]         = (*IdentityBuilder)(nil)
	_ ComponentBuilder[*config.FlagModification] = (*FlagBuilder)(nil)

	_ EntityBuilder[*config.Application, dogma.Application] = (*ApplicationBuilder)(nil)

	_ HandlerBuilder[*config.Aggregate, dogma.AggregateMessageHandler]     = (*AggregateBuilder)(nil)
	_ HandlerBuilder[*config.Process, dogma.ProcessMessageHandler]         = (*ProcessBuilder)(nil)
	_ HandlerBuilder[*config.Integration, dogma.IntegrationMessageHandler] = (*IntegrationBuilder)(nil)
	_ HandlerBuilder[*config.Projection, dogma.ProjectionMessageHandler]   = (*ProjectionBuilder)(nil)

	_ ComponentBuilder[*config.Route] = (*RouteBuilder)(nil)
)

func setTypeName[T any](
	typeName *optional.Optional[string],
	source *optional.Optional[T],
	n string,
) {
	if n == "" {
		// TODO: validate that this is actually a well-formed fully-qualified
		// type name, with optional asterisk prefix.
		panic("concrete type name must not be empty")
	}

	*typeName = optional.Some(n)
	*source = optional.None[T]()
}

func setSource[T any](
	typeName *optional.Optional[string],
	source *optional.Optional[T],
	v T,
) {
	if any(v) == nil {
		panic("source must not be nil")
	}

	*typeName = optional.Some(typename.Of(v))
	*source = optional.Some(v)
}
