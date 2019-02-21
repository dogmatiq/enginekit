package config

import (
	"context"
	"reflect"
	"strings"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

// ProcessConfig represents the configuration of an process message handler.
type ProcessConfig struct {
	// Handler is the handler that the configuration applies to.
	Handler dogma.ProcessMessageHandler

	// HandlerName is the handler's name, as specified by its Configure() method.
	HandlerName string

	messageTypes MessageTypes
}

// NewProcessConfig returns an ProcessConfig for the given handler.
func NewProcessConfig(h dogma.ProcessMessageHandler) (*ProcessConfig, error) {
	cfg := &ProcessConfig{
		Handler: h,
		messageTypes: MessageTypes{
			AcceptedEventTypes:   message.TypeSet{},
			ExecutedCommandTypes: message.TypeSet{},
		},
	}

	c := &processConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerName == "" {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.Name()",
			h,
		)
	}

	if len(c.cfg.messageTypes.AcceptedEventTypes) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.AcceptsEventType()",
			h,
		)
	}

	// if len(c.cfg.messageTypes.ExecutedCommandTypes) == 0 {
	// 	return nil, errorf(
	// 		"%T.Configure() did not call ProcessConfigurer.ExecutesCommandType()",
	// 		h,
	// 	)
	// }

	return cfg, nil
}

// Name returns the process name.
func (c *ProcessConfig) Name() string {
	return c.HandlerName
}

// HandlerType returns handler.ProcessType.
func (c *ProcessConfig) HandlerType() handler.Type {
	return handler.ProcessType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *ProcessConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// MessageTypes returns the message types used by the handler.
func (c *ProcessConfig) MessageTypes() MessageTypes {
	return c.messageTypes
}

// Accept calls v.VisitProcessConfig(ctx, c).
func (c *ProcessConfig) Accept(ctx context.Context, v Visitor) error {
	return v.VisitProcessConfig(ctx, c)
}

// processConfigurer is an implementation of dogma.ProcessConfigurer
// that builds an ProcessConfig value.
type processConfigurer struct {
	cfg *ProcessConfig
}

func (c *processConfigurer) Name(n string) {
	if c.cfg.HandlerName != "" {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.Name(%#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
		)
	}

	if strings.TrimSpace(n) == "" {
		panicf(
			`%T.Configure() called ProcessConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *processConfigurer) AcceptsEventType(m dogma.Message) {
	if !c.cfg.messageTypes.AcceptedEventTypes.AddM(m) {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.AcceptsEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *processConfigurer) ExecutesCommandType(m dogma.Message) {
	if !c.cfg.messageTypes.ExecutedCommandTypes.AddM(m) {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.ExecutesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
