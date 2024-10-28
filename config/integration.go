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

func (h *Integration) validate(ctx *validationContext) {
	h.HandlerCommon.validate(ctx, h.HandlerType())
}
