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

	cfg.ConfigurationIsExhaustive = true
	cfg.Impl = optional.Some(
		config.Implementation[dogma.Application]{
			TypeName: typename.Of(app),
			Source:   optional.Some(app),
		},
	)

	app.Configure(&applicationConfigurer{&cfg})

	return cfg
}

type applicationConfigurer struct {
	cfg *config.Application
}

func (c *applicationConfigurer) Identity(name, key string) {
	c.cfg.ConfiguredIdentities = append(
		c.cfg.ConfiguredIdentities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.cfg.ConfiguredHandlers = append(c.cfg.ConfiguredHandlers, FromAggregate(h))
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.cfg.ConfiguredHandlers = append(c.cfg.ConfiguredHandlers, FromProcess(h))
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.cfg.ConfiguredHandlers = append(c.cfg.ConfiguredHandlers, FromIntegration(h))
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.cfg.ConfiguredHandlers = append(c.cfg.ConfiguredHandlers, FromProjection(h))
}
