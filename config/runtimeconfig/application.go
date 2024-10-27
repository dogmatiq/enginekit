package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromApplication returns a new [config.Application] that represents the
// configuration of the given [dogma.Application].
func FromApplication(app dogma.Application) *config.Application {
	return configbuilder.Application(
		func(b *configbuilder.ApplicationBuilder) {
			if app != nil {
				b.Source(app)
				app.Configure(&applicationConfigurer{b})
			}
		},
	)
}

type applicationConfigurer struct {
	b *configbuilder.ApplicationBuilder
}

func (c *applicationConfigurer) Identity(name, key string) {
	c.b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.Name(name)
		b.Key(key)
	})
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.b.Aggregate(func(b *configbuilder.AggregateBuilder) {
		buildAggregate(b, h)
	})
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.b.Process(func(b *configbuilder.ProcessBuilder) {
		buildProcess(b, h)
	})
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.b.Integration(func(b *configbuilder.IntegrationBuilder) {
		buildIntegration(b, h)
	})
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.b.Projection(func(b *configbuilder.ProjectionBuilder) {
		buildProjection(b, h)
	})
}
