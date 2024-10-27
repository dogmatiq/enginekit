package config

import (
	"github.com/dogmatiq/dogma"
)

// Disabled is the [Symbol] for a [Flag] that indicates whether or not a
// [Handler] has been disabled.
type Disabled struct{ symbol }

// Aggregate is a [Handler] that represents the configuration of a
// [dogma.AggregateMessageHandler].
type Aggregate struct {
	HandlerCommon[dogma.AggregateMessageHandler]
}

// HandlerType returns the [HandlerType] of the handler.
func (a *Aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

// Process is a [Handler] that represents the configuration of a
// [dogma.ProcessMessageHandler].
type Process struct {
	HandlerCommon[dogma.ProcessMessageHandler]
}

// HandlerType returns the [HandlerType] of the handler.
func (p *Process) HandlerType() HandlerType {
	return ProcessHandlerType
}

// Integration is a [Handler] that represents the configuration of a
// [dogma.IntegrationMessageHandler].
type Integration struct {
	HandlerCommon[dogma.IntegrationMessageHandler]
}

// HandlerType returns the [HandlerType] of the handler.
func (i *Integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}
