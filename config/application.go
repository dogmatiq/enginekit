package config

import "github.com/dogmatiq/dogma"

// Application is an [Entity] that represents the configuration of a
// [dogma.Application].
type Application struct {
	EntityCommon[dogma.Application]

	// HandlerComponents is the set of [Handler] components that have been
	// registered with the application.
	HandlerComponents []Handler
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the configuration does not specify a complete set of valid
// routes for the entity and its constituents.
func (a *Application) RouteSet() RouteSet {
	panic("not implemented")
}
