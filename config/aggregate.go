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

func (h *Aggregate) validate(ctx *validationContext) {
	h.HandlerCommon.validate(ctx, h.HandlerType())
}
