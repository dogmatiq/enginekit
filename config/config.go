package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

// Config is an interface for all configuration values.
type Config interface {
	// Identity returns the identity of the configured item.
	// For example, the application or handler identity.
	Identity() Identity

	// Accept calls the appropriate method on v for this configuration type.
	Accept(ctx context.Context, v Visitor) error
}

// HandlerConfig is an interface for configuration values that refer to a
// specific message handler.
type HandlerConfig interface {
	Config

	// HandleType returns the type of handler.
	HandlerType() handler.Type

	// HandlerReflectType returns the reflect.Type of the handler.
	HandlerReflectType() reflect.Type

	// ConsumedMessageTypes returns the message types consumed by the handler.
	ConsumedMessageTypes() message.RoleMap

	// ProducedMessageTypes returns the message types produced by the handler.
	ProducedMessageTypes() message.RoleMap
}
