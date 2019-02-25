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

	consumed map[message.Type]message.Role
	produced map[message.Type]message.Role
}

// NewProcessConfig returns an ProcessConfig for the given handler.
func NewProcessConfig(h dogma.ProcessMessageHandler) (*ProcessConfig, error) {
	cfg := &ProcessConfig{
		Handler:  h,
		consumed: map[message.Type]message.Role{},
		produced: map[message.Type]message.Role{},
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
func (c *ProcessConfig) ConsumedMessageTypes() map[message.Type]message.Role {
	return c.consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *ProcessConfig) ProducedMessageTypes() map[message.Type]message.Role {
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

	if strings.TrimSpace(n) == "" {
		panicf(
			`%T.Configure() called ProcessConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *processConfigurer) ConsumesEventType(m dogma.Message) {
	t := message.TypeOf(m)

	if _, ok := c.cfg.consumed[t]; ok {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.ConsumesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}

	c.cfg.consumed[t] = message.EventRole
}

func (c *processConfigurer) ProducesCommandType(m dogma.Message) {
	t := message.TypeOf(m)

	if _, ok := c.cfg.produced[t]; ok {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.ProducesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}

	c.cfg.produced[t] = message.CommandRole
}
