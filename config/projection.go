package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/handler"
	"github.com/dogmatiq/enginekit/message"
)

// ProjectionConfig represents the configuration of an aggregate message handler.
type ProjectionConfig struct {
	// Handler is the handler that the configuration applies to.
	Handler dogma.ProjectionMessageHandler

	// HandlerName is the handler's name, as specified by its Configure() method.
	HandlerName string

	// HandlerKey is the handler's key, as specified by its Configure() method.
	HandlerKey string

	consumed message.RoleMap
}

// NewProjectionConfig returns an ProjectionConfig for the given handler.
func NewProjectionConfig(h dogma.ProjectionMessageHandler) (*ProjectionConfig, error) {
	cfg := &ProjectionConfig{
		Handler:  h,
		consumed: message.RoleMap{},
	}

	c := &projectionConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerName == "" {
		return nil, errorf(
			"%T.Configure() did not call ProjectionConfigurer.Identity()",
			h,
		)
	}

	if len(c.cfg.consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProjectionConfigurer.ConsumesEventType()",
			h,
		)
	}

	return cfg, nil
}

// Name returns the projection name.
func (c *ProjectionConfig) Name() string {
	return c.HandlerName
}

// Key returns the projection key.
func (c *ProjectionConfig) Key() string {
	return c.HandlerKey
}

// HandlerType returns handler.ProjectionType.
func (c *ProjectionConfig) HandlerType() handler.Type {
	return handler.ProjectionType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *ProjectionConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *ProjectionConfig) ConsumedMessageTypes() message.RoleMap {
	return c.consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *ProjectionConfig) ProducedMessageTypes() message.RoleMap {
	return nil
}

// Accept calls v.VisitProjectionConfig(ctx, c).
func (c *ProjectionConfig) Accept(ctx context.Context, v Visitor) error {
	return v.VisitProjectionConfig(ctx, c)
}

// projectionConfigurer is an implementation of dogma.ProjectionConfigurer
// that builds an ProjectionConfig value.
type projectionConfigurer struct {
	cfg *ProjectionConfig
}

func (c *projectionConfigurer) Identity(n, k string) {
	if c.cfg.HandlerName != "" {
		panicf(
			`%T.Configure() has already called ProjectionConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
			c.cfg.HandlerKey,
		)
	}

	if !IsValidName(n) {
		panicf(
			`%T.Configure() called ProjectionConfigurer.Identity() with an invalid name %#v`,
			c.cfg.Handler,
			n,
		)
	}

	if !IsValidKey(k) {
		panicf(
			`%T.Configure() called ProjectionConfigurer.Identity() with an invalid key %#v`,
			c.cfg.Handler,
			k,
		)
	}

	c.cfg.HandlerName = n
	c.cfg.HandlerKey = k
}

func (c *projectionConfigurer) ConsumesEventType(m dogma.Message) {
	if !c.cfg.consumed.AddM(m, message.EventRole) {
		panicf(
			`%T.Configure() has already called ProjectionConfigurer.ConsumesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
