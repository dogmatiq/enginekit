package config

import "github.com/dogmatiq/dogma"

// Integration is a [Handler] that represents the configuration of a
// [dogma.IntegrationMessageHandler].
type Integration struct {
	HandlerCommon[dogma.IntegrationMessageHandler]
}

// HandlerType returns the [HandlerType] of the handler.
func (h *Integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

// RouteSet returns the routes configured for the entity.
//
// It panics if the configuration does not specify a complete set of valid
// routes for the entity and its constituents.
func (h *Integration) RouteSet() RouteSet {
	h.validateRoutes(nil, h.HandlerType())
	return routeSetForHandler(h)
}

func (h *Integration) validate(ctx *validateContext) {
	h.HandlerCommon.validate(ctx, h.HandlerType())
}
