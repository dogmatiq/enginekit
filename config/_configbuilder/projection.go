package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// Projection builds [config.Projection] values.
type Projection struct {
	config config.Projection
}

func (b *Projection) Configurer() dogma.ProjectionConfigurer {
	panic("not implemented")
}

func (b *Projection) Get() *config.Projection {
	return &b.config
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

	c.cfg.ConfiguredDeliveryPolicy = optional.Some(cfg)
}

func (c *projectionConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.ConfiguredAsDisabled = true
}
