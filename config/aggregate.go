package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

// AggregateConfig represents the configuration of an aggregate message handler.
type AggregateConfig struct {
	// Handler is the handler that the configuration applies to.
	Handler dogma.AggregateMessageHandler

	// HandlerName is the handler's name, as specified by its Configure() method.
	HandlerName string

	consumed message.RoleMap
	produced message.RoleMap
}

// NewAggregateConfig returns an AggregateConfig for the given handler.
func NewAggregateConfig(h dogma.AggregateMessageHandler) (*AggregateConfig, error) {
	cfg := &AggregateConfig{
		Handler:  h,
		consumed: message.RoleMap{},
		produced: message.RoleMap{},
	}

	c := &aggregateConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerName == "" {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.Name()",
			h,
		)
	}

	if len(c.cfg.consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.ConsumesCommandType()",
			h,
		)
	}

	if len(c.cfg.produced) == 0 {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.ProducesEventType()",
			h,
		)
	}

	return cfg, nil
}

// Name returns the aggregate name.
func (c *AggregateConfig) Name() string {
	return c.HandlerName
}

// HandlerType returns handler.AggregateType.
func (c *AggregateConfig) HandlerType() handler.Type {
	return handler.AggregateType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *AggregateConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *AggregateConfig) ConsumedMessageTypes() message.RoleMap {
	return c.consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *AggregateConfig) ProducedMessageTypes() message.RoleMap {
	return c.produced
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

func (c *aggregateConfigurer) Name(n string) {
	if c.cfg.HandlerName != "" {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.Name(%#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
		)
	}

	if !IsValidName(n) {
		panicf(
			`%T.Configure() called AggregateConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *aggregateConfigurer) ConsumesCommandType(m dogma.Message) {
	if !c.cfg.consumed.AddM(m, message.CommandRole) {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.ConsumesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *aggregateConfigurer) ProducesEventType(m dogma.Message) {
	if !c.cfg.produced.AddM(m, message.EventRole) {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.ProducesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
