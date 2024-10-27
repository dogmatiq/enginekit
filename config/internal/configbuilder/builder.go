package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// // ComponentBuilder is the interface shared by the builder types for all
// // [config.Component] types.
// type ComponentBuilder interface {
// 	// // Fidelity returns the fidelity of the configuration.
// 	// Fidelity() config.Fidelity

// 	// // UpdateFidelity merges f into the fidelity of the configuration.
// 	// UpdateFidelity(f config.Fidelity)
// }

// // EntityBuilder is a specialization of [ComponentBuilder] for building
// // [config.QEntity] configuration.
// type EntityBuilder interface {
// 	ComponentBuilder

// 	ConcreteTypeName(typeName string)

// 	// Identity calls fn which configures a [config.Identity] that is added to
// 	// the handler.
// 	Identity(fn func(*IdentityBuilder))
// }

// // HandlerBuilder is a specialization of [EntityBuilder] for building
// // [config.HandlerFor] configuration.
// type HandlerBuilder interface {
// 	EntityBuilder

// 	// Route calls fn which configures a [config.Route] that is added to the
// 	// handler.
// 	// Route(fn func(*RouteBuilder))

// 	// Disable calls fn which configures a [config.Flag] that indicates whether
// 	// the handler is disabled.
// 	// Disable(fn func(*FlagBuilder[config.IsDisabled]))
// }

// var (
// 	_ ComponentBuilder = (*IdentityBuilder)(nil)
// 	_ ComponentBuilder = (*FlagBuilder[config.IsDisabled])(nil)

// 	_ EntityBuilder = (*ApplicationBuilder)(nil)

// 	_ HandlerBuilder = (*AggregateBuilder)(nil)
// 	_ HandlerBuilder = (*ProcessBuilder)(nil)
// 	_ HandlerBuilder = (*IntegrationBuilder)(nil)
// 	_ HandlerBuilder = (*ProjectionBuilder)(nil)

// 	_ ComponentBuilder = (*RouteBuilder)(nil)
// )

func setSourceTypeName[T any](e *config.EntityCommon[T], n string) {
	if n == "" {
		// TODO: validate that this is actually a well-formed fully-qualified
		// type name, with optional asterisk prefix.
		panic("concrete type name must not be empty")
	}

	e.SourceTypeName = n
	e.Source = optional.None[T]()
}

func setSource[T any](e *config.EntityCommon[T], v T) {
	if any(v) == nil {
		panic("runtime value must not be nil")
	}

	e.SourceTypeName = typename.Of(v)
	e.Source = optional.Some(v)
}
