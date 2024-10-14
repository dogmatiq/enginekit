package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

type ApplicationSet struct {
	apps []*config.Application
}

func (b *ApplicationSet) Application() *Application {
	app := &Application{}
	b.apps = append(b.apps, &app.config)
	return app
}

func (b *ApplicationSet) Get() []*config.Application {
	return b.apps
}

// Application builds [config.Application] values.
type Application struct {
	config config.Application
}

func (b *Application) Configurer() dogma.ApplicationConfigurer {
	panic("not implemented")
}

func (b *Application) Aggregate() *Aggregate {
	h := &Aggregate{}
	b.config.ConfiguredHandlers = append(b.config.ConfiguredHandlers, &h.config)
	return h
}

func (b *Application) Integration() *Integration {
	h := &Integration{}
	b.config.ConfiguredHandlers = append(b.config.ConfiguredHandlers, &h.config)
	return h
}

func (b *Application) Process() *Process {
	h := &Process{}
	b.config.ConfiguredHandlers = append(b.config.ConfiguredHandlers, &h.config)
	return h
}

func (b *Application) Projection() *Projection {
	h := &Projection{}
	b.config.ConfiguredHandlers = append(b.config.ConfiguredHandlers, &h.config)
	return h
}

func (b *Application) Get() *config.Application {
	return &b.config
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
