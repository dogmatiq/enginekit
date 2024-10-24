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

// Disabled is the label for a [Flag] that indicates that a [Handler] has been
// disabled.
type Disabled struct{ Label }
