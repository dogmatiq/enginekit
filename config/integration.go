package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// IntegrationConfig represents the configuration of an integration message handler.
type IntegrationConfig struct {
	// Handler is the handler that the configuration applies to.
	// It is nil if the config was obtained via the config API.
	Handler dogma.IntegrationMessageHandler

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

// NewIntegrationConfig returns an IntegrationConfig for the given handler.
func NewIntegrationConfig(h dogma.IntegrationMessageHandler) (*IntegrationConfig, error) {
	cfg := &IntegrationConfig{
		Handler:  h,
		Consumed: MessageRoleMap{},
		Produced: MessageRoleMap{},
	}

	c := &integrationConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerIdentity.IsZero() {
		return nil, errorf(
			"%T.Configure() did not call IntegrationConfigurer.Identity()",
			h,
		)
	}

	if len(c.cfg.Consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call IntegrationConfigurer.ConsumesCommandType()",
			h,
		)
	}

	return cfg, nil
}

// Identity returns the integration identity.
func (c *IntegrationConfig) Identity() Identity {
	return c.HandlerIdentity
}

// HandlerType returns IntegrationHandlerType.
func (c *IntegrationConfig) HandlerType() HandlerType {
	return IntegrationHandlerType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *IntegrationConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *IntegrationConfig) ConsumedMessageTypes() MessageRoleMap {
	return c.Consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *IntegrationConfig) ProducedMessageTypes() MessageRoleMap {
	return c.Produced
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

func (c *integrationConfigurer) Identity(n, k string) {
	if !c.cfg.HandlerIdentity.IsZero() {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerIdentity.Name,
			c.cfg.HandlerIdentity.Key,
		)
	}

	i, err := NewIdentity(n, k)
	if err != nil {
		panicf(
			`%T.Configure() called IntegrationConfigurer.Identity() with an %s`,
			c.cfg.Handler,
			err,
		)
	}

	c.cfg.HandlerIdentity = i
}

func (c *integrationConfigurer) ConsumesCommandType(m dogma.Message) {
	if !c.cfg.Consumed.AddM(m, CommandMessageRole) {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.ConsumesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *integrationConfigurer) ProducesEventType(m dogma.Message) {
	if !c.cfg.Produced.AddM(m, EventMessageRole) {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.ProducesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
