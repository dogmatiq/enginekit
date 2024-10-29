package config

import (
	"github.com/dogmatiq/dogma"
)

// Aggregate is a [Handler] that represents the configuration of a
// [dogma.AggregateMessageHandler].
type Aggregate struct {
	HandlerCommon[dogma.AggregateMessageHandler]
}

// HandlerType returns the [HandlerType] of the handler.
func (h *Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the configuration does not specify a complete set of valid
// routes for the entity and its constituents.
func (h *Aggregate) RouteSet() RouteSet {
	h.validateRoutes(nil, h.HandlerType())
	return routeSetForHandler(h)
}

func (h *Aggregate) validate(ctx *validateContext) {
	h.HandlerCommon.validate(ctx, h.HandlerType())
}
