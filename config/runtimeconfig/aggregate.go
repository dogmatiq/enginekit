package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromAggregate returns a new [config.Aggregate] that represents the
// configuration of the given [dogma.AggregateMessageHandler].
func FromAggregate(h dogma.AggregateMessageHandler) *config.Aggregate {
	b := configbuilder.Aggregate()

	if h == nil {
		b.UpdateFidelity(config.Incomplete)
	} else {
		b.SetDisabled(false)
		b.SetSource(h)
		h.Configure(&aggregateConfigurer{b})
	}

	return b.Done()
}

type aggregateConfigurer struct {
	b *configbuilder.AggregateBuilder
}

func (c *aggregateConfigurer) Identity(name, key string) {
	c.b.
		AddIdentity().
		SetName(name).
		SetKey(key).
		Done()
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	c.b.Edit(
		func(cfg *config.AggregateAsConfigured) {
			for _, r := range routes {
				cfg.Routes = append(cfg.Routes, fromRoute(r))
			}
		},
	)
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.b.SetDisabled(true)
}
