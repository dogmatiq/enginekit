package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromApplication returns a new [config.Application] that represents the
// configuration of the given [dogma.Application].
func FromApplication(app dogma.Application) *config.Application {
	b := configbuilder.Application()

	if app == nil {
		return b.Done(config.Incomplete)
	}

	b.SetSource(app)
	app.Configure(&applicationConfigurer{b})

	return b.Done(config.Immaculate)
}

type applicationConfigurer struct {
	b *configbuilder.ApplicationBuilder
}

func (c *applicationConfigurer) Identity(name, key string) {
	c.b.
		AddIdentity().
		SetName(name).
		SetKey(key).
		Done()
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.b.BuildAggregate(
		func(b *configbuilder.AggregateBuilder) {
			b.SetSource(h)
			b.SetDisabled(false)
			h.Configure(&aggregateConfigurer{b})
		},
	)
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.b.BuildProcess(
		func(b *configbuilder.ProcessBuilder) {
			b.SetSource(h)
			b.SetDisabled(false)
			h.Configure(&processConfigurer{b})
		},
	)
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.b.BuildIntegration(
		func(b *configbuilder.IntegrationBuilder) {
			b.SetSource(h)
			b.SetDisabled(false)
			h.Configure(&integrationConfigurer{b})
		},
	)
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.b.BuildProjection(
		func(b *configbuilder.ProjectionBuilder) {
			b.SetSource(h)
			b.SetDisabled(false)
			b.SetDeliveryPolicy(dogma.UnicastProjectionDeliveryPolicy{})
			h.Configure(&projectionConfigurer{b})
		},
	)
}
