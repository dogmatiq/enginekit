package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

// Config is an interface for all configuration values.
type Config interface {
	// Name returns the name of the configured item.
	// For example, the application or handler name.
	Name() string

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

	// MessageTypes returns the message types used by the handler.
	MessageTypes() MessageTypes
}

// MessageTypes contains sets that indicate what types of messages are accepted
// and produced by a handler.
type MessageTypes struct {
	AcceptedCommandTypes message.TypeSet
	AcceptedEventTypes   message.TypeSet
	ExecutedCommandTypes message.TypeSet
	RecordedEventTypes   message.TypeSet
}
