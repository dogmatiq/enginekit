package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromAggregate returns a new [config.Aggregate] that represents the
// configuration of the given [dogma.AggregateMessageHandler].
func FromAggregate(h dogma.AggregateMessageHandler) *config.Aggregate {
	cfg := &config.Aggregate{
		AsConfigured: config.AggregateAsConfigured{
			IsDisabled: optional.Some(false),
		},
	}

	if h == nil {
		return cfg
	}

	cfg.AsConfigured.Source.TypeName = optional.Some(typename.Of(h))
	cfg.AsConfigured.Source.Value = optional.Some(h)

	h.Configure(&aggregateConfigurer{cfg})

	return cfg
}

type aggregateConfigurer struct {
	cfg *config.Aggregate
}

func (c *aggregateConfigurer) Identity(name, key string) {
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

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	for _, r := range routes {
		c.cfg.AsConfigured.Routes = append(c.cfg.AsConfigured.Routes, fromRoute(r))
	}
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.AsConfigured.IsDisabled = optional.Some(true)
}
