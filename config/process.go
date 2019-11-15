package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// ProcessConfig represents the configuration of an process message handler.
type ProcessConfig struct {
	// Handler is the handler that the configuration applies to.
	// It is nil if the config was obtained via the config API.
	Handler dogma.ProcessMessageHandler

	// HandlerIdentity is the handler's identity, as specified by its
	// Configure() method.
	HandlerIdentity Identity

	// Consumed is a map of message type to role for those message types
	// consumed by this handler.
	Consumed MessageRoleMap

	// Produced is a map of message type to role for those message types
	// produced by this handler.
	Produced MessageRoleMap
}

// NewProcessConfig returns an ProcessConfig for the given handler.
func NewProcessConfig(h dogma.ProcessMessageHandler) (*ProcessConfig, error) {
	cfg := &ProcessConfig{
		Handler:  h,
		Consumed: MessageRoleMap{},
		Produced: MessageRoleMap{},
	}

	c := &processConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerIdentity.IsZero() {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.Identity()",
			h,
		)
	}

	if len(c.cfg.Consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.ConsumesEventType()",
			h,
		)
	}

	if len(c.cfg.Produced) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProcessConfigurer.ProducesCommandType()",
			h,
		)
	}

	return cfg, nil
}

// Identity returns the process identity.
func (c *ProcessConfig) Identity() Identity {
	return c.HandlerIdentity
}

// HandlerType returns ProcessHandlerType.
func (c *ProcessConfig) HandlerType() HandlerType {
	return ProcessHandlerType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *ProcessConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *ProcessConfig) ConsumedMessageTypes() MessageRoleMap {
	return c.Consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *ProcessConfig) ProducedMessageTypes() MessageRoleMap {
	return c.Produced
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
	if !c.cfg.HandlerIdentity.IsZero() {
		panicf(
			`%T.Configure() has already called ProcessConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerIdentity.Name,
			c.cfg.HandlerIdentity.Key,
		)
	}

	i, err := NewIdentity(n, k)
	if err != nil {
		panicf(
			`%T.Configure() called ProcessConfigurer.Identity() with an %s`,
			c.cfg.Handler,
			err,
		)
	}

	c.cfg.HandlerIdentity = i
}

func (c *processConfigurer) ConsumesEventType(m dogma.Message) {
	c.add(c.cfg.Consumed, m, EventMessageRole)
}

func (c *processConfigurer) ProducesCommandType(m dogma.Message) {
	c.add(c.cfg.Produced, m, CommandMessageRole)
}

func (c *processConfigurer) SchedulesTimeoutType(m dogma.Message) {
	c.add(c.cfg.Consumed, m, TimeoutMessageRole)
	c.add(c.cfg.Produced, m, TimeoutMessageRole)
}

func (c *processConfigurer) add(rm MessageRoleMap, m dogma.Message, r MessageRole) {
	mt := MessageTypeOf(m)

	if rm.Add(mt, r) {
		return
	}

	var f string
	switch rm[mt] {
	case CommandMessageRole:
		f = "ProducesCommandType"
	case EventMessageRole:
		f = "ConsumesEventType"
	default: // TimeoutMessageRole
		f = "SchedulesTimeoutType"
	}

	panicf(
		`%T.Configure() has already called ProcessConfigurer.%s(%T)`,
		c.cfg.Handler,
		f,
		m,
	)
}
