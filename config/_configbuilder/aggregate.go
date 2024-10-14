package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Aggregate builds [config.Aggregate] values.
type Aggregate struct {
	config config.Aggregate
}

// Configurer returns a [dogma.AggregateConfigurer] that
func (b *Aggregate) Configurer() dogma.AggregateConfigurer {
	return &aggregateConfigurer{&b.config}
}

// Get returns the configuration.
func (b *Aggregate) Get() *config.Aggregate {
	return &b.config
}

type aggregateConfigurer struct {
	cfg *config.Aggregate
}

func (c *aggregateConfigurer) Identity(name, key string) {
	c.cfg.ConfiguredIdentities = append(
		c.cfg.ConfiguredIdentities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	for _, r := range routes {
		b := &Route{}
		b.Set(r)

		c.cfg.ConfiguredRoutes = append(
			c.cfg.ConfiguredRoutes,
			b.Get(),
		)
	}
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.ConfiguredAsDisabled = true
}
