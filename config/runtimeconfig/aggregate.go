package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromAggregate returns a new [config.Aggregate] that represents the
// configuration of the given [dogma.AggregateMessageHandler].
func FromAggregate(h dogma.AggregateMessageHandler) config.Aggregate {
	var cfg config.Aggregate

	if h == nil {
		return cfg
	}

	cfg.ConfigurationIsExhaustive = true
	cfg.ConfigurationSource = optional.Some(
		config.Source[dogma.AggregateMessageHandler]{
			TypeName: typename.Of(h),
			Value:    optional.Some(h),
		},
	)

	h.Configure(&aggregateConfigurer{&cfg})

	return cfg
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
		c.cfg.ConfiguredRoutes = append(c.cfg.ConfiguredRoutes, fromRoute(r))
	}
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.ConfiguredAsDisabled = true
}