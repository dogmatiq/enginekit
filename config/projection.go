package config

import (
	"context"
	"reflect"
	"strings"

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

	consumed map[message.Type]message.Role
}

// NewProjectionConfig returns an ProjectionConfig for the given handler.
func NewProjectionConfig(h dogma.ProjectionMessageHandler) (*ProjectionConfig, error) {
	cfg := &ProjectionConfig{
		Handler:  h,
		consumed: map[message.Type]message.Role{},
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
			"%T.Configure() did not call ProjectionConfigurer.Name()",
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

// HandlerType returns handler.ProjectionType.
func (c *ProjectionConfig) HandlerType() handler.Type {
	return handler.ProjectionType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *ProjectionConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *ProjectionConfig) ConsumedMessageTypes() map[message.Type]message.Role {
	return c.consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *ProjectionConfig) ProducedMessageTypes() map[message.Type]message.Role {
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

func (c *projectionConfigurer) Name(n string) {
	if c.cfg.HandlerName != "" {
		panicf(
			`%T.Configure() has already called ProjectionConfigurer.Name(%#v)`,
			c.cfg.Handler,
			c.cfg.HandlerName,
		)
	}

	if strings.TrimSpace(n) == "" {
		panicf(
			`%T.Configure() called ProjectionConfigurer.Name(%#v) with an invalid name`,
			c.cfg.Handler,
			n,
		)
	}

	c.cfg.HandlerName = n
}

func (c *projectionConfigurer) ConsumesEventType(m dogma.Message) {
	t := message.TypeOf(m)

	if _, ok := c.cfg.consumed[t]; ok {
		panicf(
			`%T.Configure() has already called ProjectionConfigurer.ConsumesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}

	c.cfg.consumed[t] = message.EventRole
}
