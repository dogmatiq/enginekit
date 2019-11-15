package config

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// ProjectionConfig represents the configuration of an aggregate message handler.
type ProjectionConfig struct {
	// Handler is the handler that the configuration applies to.
	// It is nil if the config was obtained via the config API.
	Handler dogma.ProjectionMessageHandler

	// HandlerIdentity is the handler's identity, as specified by its
	// Configure() method.
	HandlerIdentity Identity

	// Consumed is a map of message type to role for those message types
	// consumed by this handler.
	Consumed MessageRoleMap
}

// NewProjectionConfig returns an ProjectionConfig for the given handler.
func NewProjectionConfig(h dogma.ProjectionMessageHandler) (*ProjectionConfig, error) {
	cfg := &ProjectionConfig{
		Handler:  h,
		Consumed: MessageRoleMap{},
	}

	c := &projectionConfigurer{
		cfg: cfg,
	}

	if err := catch(func() {
		h.Configure(c)
	}); err != nil {
		return nil, err
	}

	if c.cfg.HandlerIdentity.IsZero() {
		return nil, errorf(
			"%T.Configure() did not call ProjectionConfigurer.Identity()",
			h,
		)
	}

	if len(c.cfg.Consumed) == 0 {
		return nil, errorf(
			"%T.Configure() did not call ProjectionConfigurer.ConsumesEventType()",
			h,
		)
	}

	return cfg, nil
}

// Identity returns the projection identity.
func (c *ProjectionConfig) Identity() Identity {
	return c.HandlerIdentity
}

// HandlerType returns ProjectionHandlerType.
func (c *ProjectionConfig) HandlerType() HandlerType {
	return ProjectionHandlerType
}

// HandlerReflectType returns the reflect.Type of the handler.
func (c *ProjectionConfig) HandlerReflectType() reflect.Type {
	return reflect.TypeOf(c.Handler)
}

// ConsumedMessageTypes returns the message types consumed by the handler.
func (c *ProjectionConfig) ConsumedMessageTypes() MessageRoleMap {
	return c.Consumed
}

// ProducedMessageTypes returns the message types produced by the handler.
func (c *ProjectionConfig) ProducedMessageTypes() MessageRoleMap {
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
	if !c.cfg.HandlerIdentity.IsZero() {
		panicf(
			`%T.Configure() has already called ProjectionConfigurer.Identity(%#v, %#v)`,
			c.cfg.Handler,
			c.cfg.HandlerIdentity.Name,
			c.cfg.HandlerIdentity.Key,
		)
	}

	i, err := NewIdentity(n, k)
	if err != nil {
		panicf(
			`%T.Configure() called ProjectionConfigurer.Identity() with an %s`,
			c.cfg.Handler,
			err,
		)
	}

	c.cfg.HandlerIdentity = i
}

func (c *projectionConfigurer) ConsumesEventType(m dogma.Message) {
	if !c.cfg.Consumed.AddM(m, EventMessageRole) {
		panicf(
			`%T.Configure() has already called ProjectionConfigurer.ConsumesEventType(%T)`,
			c.cfg.Handler,
			m,
		)
	}
}
