package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromProjection returns a new [config.Projection] that represents the
// configuration of the given [dogma.ProjectionMessageHandler].
func FromProjection(h dogma.ProjectionMessageHandler) config.Projection {
	var cfg config.Projection

	if h == nil {
		return cfg
	}

	cfg.IsExhaustive = true
	cfg.Impl = optional.Some(
		config.Implementation[dogma.ProjectionMessageHandler]{
			TypeName: typename.Of(h),
			Source:   optional.Some(h),
		},
	)

	h.Configure(&projectionConfigurer{&cfg})

	return cfg
}

type projectionConfigurer struct {
	cfg *config.Projection
}

func (c *projectionConfigurer) Identity(name, key string) {
	c.cfg.ConfiguredIdentities = append(
		c.cfg.ConfiguredIdentities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	for _, r := range routes {
		c.cfg.ConfiguredRoutes = append(c.cfg.ConfiguredRoutes, fromRoute(r))
	}
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	cfg := config.ProjectionDeliveryPolicy{}

	if p != nil {
		cfg.TypeName = optional.Some(typename.Of(p))
		cfg.Implementation = optional.Some(p)
	}

	c.cfg.DeliveryPolicy = optional.Some(cfg)
}

func (c *projectionConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.IsDisabled = true
}
