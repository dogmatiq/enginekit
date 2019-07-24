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

	// HandlerKey is the handler's key, as specified by its Configure() method.
	HandlerKey string

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
			"%T.Configure() did not call ProcessConfigurer.Identity()",
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

// Key returns the process key.
func (c *ProcessConfig) Key() string {
	return c.HandlerKey
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

func (c *processConfigurer) Identity(n, k string) {
	if c.cfg.HandlerName != "" {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
			c.cfg.HandlerKey,
		)
	}

	if !IsValidName(n) {
		panicf(
			`%T.Configure() called ProcessConfigurer.Identity() with an invalid name %#v`,
			c.cfg.Handler,
			n,
		)
	}

	if !IsValidKey(k) {
		panicf(
			`%T.Configure() called ProcessConfigurer.Identity() with an invalid key %#v`,
			c.cfg.Handler,
			k,
		)
	}

	c.cfg.HandlerName = n
	c.cfg.HandlerKey = k
}

func (c *processConfigurer) ConsumesEventType(m dogma.Message) {
	c.add(c.cfg.consumed, m, message.EventRole)
}

func (c *processConfigurer) ProducesCommandType(m dogma.Message) {
	c.add(c.cfg.produced, m, message.CommandRole)
}

func (c *processConfigurer) SchedulesTimeoutType(m dogma.Message) {
	c.add(c.cfg.consumed, m, message.TimeoutRole)
	c.add(c.cfg.produced, m, message.TimeoutRole)
}

func (c *processConfigurer) add(rm message.RoleMap, m dogma.Message, r message.Role) {
	mt := message.TypeOf(m)

	if rm.Add(mt, r) {
		return
	}

	var f string
	switch rm[mt] {
	case message.CommandRole:
		f = "ProducesCommandType"
	case message.EventRole:
		f = "ConsumesEventType"
	default: //case message.TimeoutRole
		f = "SchedulesTimeoutType"
	}

	panicf(
		`%T.Configure() has already called ProcessConfigurer.%s(%T)`,
		c.cfg.Handler,
		f,
		m,
	)
}
