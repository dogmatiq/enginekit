package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// AggregateConfig represents the configuration of an aggregate message handler.
type AggregateConfig struct {
	// Handler is the handler that the configuration applies to.
	// It is nil if the config was obtained via the config API.
	Handler dogma.AggregateMessageHandler

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

// NewAggregateConfig returns an AggregateConfig for the given handler.
func NewAggregateConfig(h dogma.AggregateMessageHandler) (*AggregateConfig, error) {
	cfg := &AggregateConfig{
		Handler:  h,
		Consumed: MessageRoleMap{},
		Produced: MessageRoleMap{},
	}

	c := &aggregateConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerIdentity.IsZero() {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.Identity()",
			h,
		)
	}

	if len(c.cfg.Consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.ConsumesCommandType()",
			h,
		)
	}

	if len(c.cfg.Produced) == 0 {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.ProducesEventType()",
			h,
		)
	}

	return cfg, nil
}

// Identity returns the aggregate identity.
func (c *AggregateConfig) Identity() Identity {
	return c.HandlerIdentity
}

// HandlerType returns AggregateHandlerType.
func (c *AggregateConfig) HandlerType() HandlerType {
	return AggregateHandlerType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *AggregateConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *AggregateConfig) ConsumedMessageTypes() MessageRoleMap {
	return c.Consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *AggregateConfig) ProducedMessageTypes() MessageRoleMap {
	return c.Produced
}

// Accept calls v.VisitAggregateConfig(ctx, c).
func (c *AggregateConfig) Accept(ctx context.Context, v Visitor) error {
	return v.VisitAggregateConfig(ctx, c)
}

// aggregateConfigurer is an implementation of dogma.AggregateConfigurer
// that builds an AggregateConfig value.
type aggregateConfigurer struct {
	cfg *AggregateConfig
}

func (c *aggregateConfigurer) Identity(n, k string) {
	if !c.cfg.HandlerIdentity.IsZero() {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerIdentity.Name,
			c.cfg.HandlerIdentity.Key,
		)
	}

	i, err := NewIdentity(n, k)
	if err != nil {
		panicf(
			`%T.Configure() called AggregateConfigurer.Identity() with an %s`,
			c.cfg.Handler,
			err,
		)
	}

	c.cfg.HandlerIdentity = i
}

func (c *aggregateConfigurer) ConsumesCommandType(m dogma.Message) {
	if !c.cfg.Consumed.AddM(m, CommandMessageRole) {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.ConsumesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *aggregateConfigurer) ProducesEventType(m dogma.Message) {
	if !c.cfg.Produced.AddM(m, EventMessageRole) {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.ProducesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
