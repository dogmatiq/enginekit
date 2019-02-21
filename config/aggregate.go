package config

import (
	"context"
	"reflect"
	"strings"

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

	messageTypes MessageTypes
}

// NewAggregateConfig returns an AggregateConfig for the given handler.
func NewAggregateConfig(h dogma.AggregateMessageHandler) (*AggregateConfig, error) {
	cfg := &AggregateConfig{
		Handler: h,
		messageTypes: MessageTypes{
			AcceptedCommandTypes: message.TypeSet{},
			RecordedEventTypes:   message.TypeSet{},
		},
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

	if len(c.cfg.messageTypes.AcceptedCommandTypes) == 0 {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.AcceptsCommandType()",
			h,
		)
	}

	if len(c.cfg.messageTypes.RecordedEventTypes) == 0 {
		return nil, errorf(
			"%T.Configure() did not call AggregateConfigurer.RecordsEventType()",
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

// MessageTypes returns the message types used by the handler.
func (c *AggregateConfig) MessageTypes() MessageTypes {
	return c.messageTypes
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

	if strings.TrimSpace(n) == "" {
		panicf(
			`%T.Configure() called AggregateConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *aggregateConfigurer) AcceptsCommandType(m dogma.Message) {
	if !c.cfg.messageTypes.AcceptedCommandTypes.AddM(m) {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.AcceptsCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *aggregateConfigurer) RecordsEventType(m dogma.Message) {
	if !c.cfg.messageTypes.RecordedEventTypes.AddM(m) {
		panicf(
			`%T.Configure() has already called AggregateConfigurer.RecordsEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
