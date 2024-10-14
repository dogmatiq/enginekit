package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromApplication returns a new [config.Application] that represents the
// configuration of the given [dogma.Application].
func FromApplication(app dogma.Application) *config.Application {
	cfg := &config.Application{}

	if app == nil {
		return cfg
	}

	cfg.AsConfigured.Source = optional.Some(
		config.Source[dogma.Application]{
			TypeName:  typename.Of(app),
			Interface: optional.Some(app),
		},
	)

	app.Configure(&applicationConfigurer{cfg})

	return cfg
}

type applicationConfigurer struct {
	cfg *config.Application
}

func (c *applicationConfigurer) Identity(name, key string) {
	c.cfg.AsConfigured.Identities = append(
		c.cfg.AsConfigured.Identities,
		config.Identity{
			AsConfigured: config.IdentityAsConfigured{
				Name: name,
				Key:  key,
			},
		},
	)
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.cfg.AsConfigured.Handlers = append(c.cfg.AsConfigured.Handlers, FromAggregate(h))
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.cfg.AsConfigured.Handlers = append(c.cfg.AsConfigured.Handlers, FromProcess(h))
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.cfg.AsConfigured.Handlers = append(c.cfg.AsConfigured.Handlers, FromIntegration(h))
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.cfg.AsConfigured.Handlers = append(c.cfg.AsConfigured.Handlers, FromProjection(h))
}
