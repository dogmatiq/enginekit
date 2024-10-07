package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromApplication returns a new [config.Application] that represents the
// configuration of the given [dogma.Application].
func FromApplication(app dogma.Application) config.Application {
	var cfg config.Application

	if app == nil {
		return cfg
	}

	cfg.TypeName = optional.Some(typename.Of(app))
	cfg.Implementation = optional.Some(app)
	cfg.IsExhaustive = true

	app.Configure(&applicationConfigurer{&cfg})

	return cfg
}

type applicationConfigurer struct {
	cfg *config.Application
}

func (c *applicationConfigurer) Identity(name, key string) {
	c.cfg.Identities = append(
		c.cfg.Identities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.cfg.Aggregates = append(c.cfg.Aggregates, FromAggregate(h))
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.cfg.Processes = append(c.cfg.Processes, FromProcess(h))
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.cfg.Integrations = append(c.cfg.Integrations, FromIntegration(h))
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.cfg.Projections = append(c.cfg.Projections, FromProjection(h))
}
