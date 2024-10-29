package config

import "github.com/dogmatiq/dogma"

// Process is a [Handler] that represents the configuration of a
// [dogma.ProcessMessageHandler].
type Process struct {
	HandlerCommon[dogma.ProcessMessageHandler]
}

// HandlerType returns the [HandlerType] of the handler.
func (h *Process) HandlerType() HandlerType {
	return ProcessHandlerType
}

func (h *Process) validate(ctx *validateContext) {
	h.HandlerCommon.validate(ctx, h.HandlerType())
}
