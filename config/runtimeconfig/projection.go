package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromProjection returns a new [config.Projection] that represents the
// configuration of the given [dogma.ProjectionMessageHandler].
func FromProjection(h dogma.ProjectionMessageHandler) *config.Projection {
	cfg := &config.Projection{
		AsConfigured: config.ProjectionAsConfigured{
			IsDisabled: optional.Some(false),
		},
	}

	if h == nil {
		return cfg
	}

	cfg.AsConfigured.Source.TypeName = optional.Some(typename.Of(h))
	cfg.AsConfigured.Source.Value = optional.Some(h)

	h.Configure(&projectionConfigurer{cfg})

	return cfg
}

type projectionConfigurer struct {
	cfg *config.Projection
}

func (c *projectionConfigurer) Identity(name, key string) {
	c.cfg.AsConfigured.Identities = append(
		c.cfg.AsConfigured.Identities,
		&config.Identity{
			AsConfigured: config.IdentityAsConfigured{
				Name: optional.Some(name),
				Key:  optional.Some(key),
			},
		},
	)
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	for _, r := range routes {
		c.cfg.AsConfigured.Routes = append(c.cfg.AsConfigured.Routes, fromRoute(r))
	}
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	cfg := config.Value[dogma.ProjectionDeliveryPolicy]{}

	if p != nil {
		cfg.TypeName = optional.Some(typename.Of(p))
		cfg.Value = optional.Some(p)
	}

	c.cfg.AsConfigured.DeliveryPolicy = optional.Some(cfg)
}

func (c *projectionConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.AsConfigured.IsDisabled = optional.Some(true)
}
