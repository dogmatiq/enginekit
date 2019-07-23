package config

import (
	"context"
	"reflect"

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

	// HandlerKey is the handler's key, as specified by its Configure() method.
	HandlerKey string

	consumed message.RoleMap
	produced message.RoleMap
}

// NewIntegrationConfig returns an IntegrationConfig for the given handler.
func NewIntegrationConfig(h dogma.IntegrationMessageHandler) (*IntegrationConfig, error) {
	cfg := &IntegrationConfig{
		Handler:  h,
		consumed: message.RoleMap{},
		produced: message.RoleMap{},
	}

	c := &integrationConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerName == "" || c.cfg.HandlerKey == "" {
		return nil, errorf(
			"%T.Configure() did not call IntegrationConfigurer.Identity()",
			h,
		)
	}

	if len(c.cfg.consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call IntegrationConfigurer.ConsumesCommandType()",
			h,
		)
	}

	return cfg, nil
}

// Name returns the integration name.
func (c *IntegrationConfig) Name() string {
	return c.HandlerName
}

// Key returns the integration key.
func (c *IntegrationConfig) Key() string {
	return c.HandlerKey
}

// HandlerType returns handler.IntegrationType.
func (c *IntegrationConfig) HandlerType() handler.Type {
	return handler.IntegrationType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *IntegrationConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *IntegrationConfig) ConsumedMessageTypes() message.RoleMap {
	return c.consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *IntegrationConfig) ProducedMessageTypes() message.RoleMap {
	return c.produced
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

func (c *integrationConfigurer) Identity(name, key string) {
	if c.cfg.HandlerName != "" && c.cfg.HandlerKey != "" {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
			c.cfg.HandlerKey,
		)
	}

	if !IsValidName(name) {
		panicf(
			`%T.Configure() called IntegrationConfigurer.Identity() with an invalid name %#v`,
			c.cfg.Handler,
			name,
		)
	}

	if !IsValidKey(key) {
		panicf(
			`%T.Configure() called IntegrationConfigurer.Identity() with an invalid key %#v`,
			c.cfg.Handler,
			key,
		)
	}

	c.cfg.HandlerName = name
	c.cfg.HandlerKey = key
}

func (c *integrationConfigurer) ConsumesCommandType(m dogma.Message) {
	if !c.cfg.consumed.AddM(m, message.CommandRole) {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.ConsumesCommandType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}

func (c *integrationConfigurer) ProducesEventType(m dogma.Message) {
	if !c.cfg.produced.AddM(m, message.EventRole) {
		panicf(
			`%T.Configure() has already called IntegrationConfigurer.ProducesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
