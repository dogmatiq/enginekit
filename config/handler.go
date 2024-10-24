package config

// A Handler is a specialization of [Entity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	Entity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler was disabled via the configurer.
	IsDisabled() bool

	disabledFlags() FlagSet[Disabled]
	routes() []*Route
}

// HandlerTrait is a partial implementation of [Handler].
type HandlerTrait[T any] struct {
	EntityTrait[T]

	// Routes is the list of routes configured on the handler.
	Routes []*Route

	// DisabledFlags represents calls to [dogma.AggregateConfigurer.Disable].
	DisabledFlags FlagSet[Disabled]
}

// HandlerType returns [HandlerType] of the handler.
func (h HandlerTrait[T]) HandlerType() HandlerType {
	return HandlerTypeFor[T]()
}

// IsDisabled returns true if the handler was disabled via the configurer.
func (h HandlerTrait[T]) IsDisabled() bool {
	return h.DisabledFlags.resolve(h.F).Get()
}

// RouteSet returns the routes configured for the handler.
//
// It panics if the routes are incomplete or invalid.
func (h HandlerTrait[T]) RouteSet() RouteSet {
	return buildRouteSet(strictContext(), h)
}

// Disabled is the label for a [Flag] that indicates that a [Handler] has been
// disabled.
type Disabled struct{ Label }
