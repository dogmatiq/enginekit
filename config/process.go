package config

import (
	"context"
	"reflect"

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

	consumed message.RoleMap
	produced message.RoleMap
}

// NewProcessConfig returns an ProcessConfig for the given handler.
func NewProcessConfig(h dogma.ProcessMessageHandler) (*ProcessConfig, error) {
	cfg := &ProcessConfig{
		Handler:  h,
		consumed: message.RoleMap{},
		produced: message.RoleMap{},
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

	if len(c.cfg.consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.ConsumesEventType()",
			h,
		)
	}

	if len(c.cfg.produced) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.ProducesCommandType()",
			h,
		)
	}

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

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *ProcessConfig) ConsumedMessageTypes() message.RoleMap {
	return c.consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *ProcessConfig) ProducedMessageTypes() message.RoleMap {
	return c.produced
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

	if !IsValidName(n) {
		panicf(
			`%T.Configure() called ProcessConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *processConfigurer) ConsumesEventType(m dogma.Message) {
	if !c.cfg.consumed.AddM(m, message.EventRole) {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.ConsumesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *processConfigurer) ProducesCommandType(m dogma.Message) {
	if !c.cfg.produced.AddM(m, message.CommandRole) {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.ProducesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
