package config

import (
	"context"
	"reflect"
	"strings"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

// IntegrationConfig represents the configuration of an integration message handler.
type IntegrationConfig struct {
	// Handler is the handler that the configuration applies to.
	Handler dogma.IntegrationMessageHandler

	// HandlerName is the handler's name, as specified by its Configure() method.
	HandlerName string

	messageTypes MessageTypes
}

// NewIntegrationConfig returns an IntegrationConfig for the given handler.
func NewIntegrationConfig(h dogma.IntegrationMessageHandler) (*IntegrationConfig, error) {
	cfg := &IntegrationConfig{
		Handler: h,
		messageTypes: MessageTypes{
			AcceptedCommandTypes: message.TypeSet{},
			RecordedEventTypes:   message.TypeSet{},
		},
	}

	c := &integrationConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerName == "" {
		return nil, errorf(
			"%T.Configure() did not call IntegrationConfigurer.Name()",
			h,
		)
	}

	if len(c.cfg.messageTypes.AcceptedCommandTypes) == 0 {
		return nil, errorf(
			"%T.Configure() did not call IntegrationConfigurer.AcceptsCommandType()",
			h,
		)
	}

	// if len(c.cfg.messageTypes.RecordedEventTypes) == 0 {
	// 	return nil, errorf(
	// 		"%T.Configure() did not call IntegrationConfigurer.RecordsEventType()",
	// 		h,
	// 	)
	// }

	return cfg, nil
}

// Name returns the integration name.
func (c *IntegrationConfig) Name() string {
	return c.HandlerName
}

// HandlerType returns handler.IntegrationType.
func (c *IntegrationConfig) HandlerType() handler.Type {
	return handler.IntegrationType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *IntegrationConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// MessageTypes returns the message types used by the handler.
func (c *IntegrationConfig) MessageTypes() MessageTypes {
	return c.messageTypes
}

// Accept calls v.VisitIntegrationConfig(ctx, c).
func (c *IntegrationConfig) Accept(ctx context.Context, v Visitor) error {
	return v.VisitIntegrationConfig(ctx, c)
}

// integrationConfigurer is an implementation of dogma.integrationConfigurer
// that builds an IntegrationConfig value.
type integrationConfigurer struct {
	cfg *IntegrationConfig
}

func (c *integrationConfigurer) Name(n string) {
	if c.cfg.HandlerName != "" {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.Name(%#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
		)
	}

	if strings.TrimSpace(n) == "" {
		panicf(
			`%T.Configure() called IntegrationConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *integrationConfigurer) AcceptsCommandType(m dogma.Message) {
	if !c.cfg.messageTypes.AcceptedCommandTypes.AddM(m) {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.AcceptsCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *integrationConfigurer) RecordsEventType(m dogma.Message) {
	if !c.cfg.messageTypes.RecordedEventTypes.AddM(m) {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.RecordsEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
