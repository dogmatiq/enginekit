package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Integration builds [config.Integration] values.
type Integration struct {
	config config.Integration
}

func (b *Integration) Configurer() dogma.IntegrationConfigurer {
	panic("not implemented")
}

func (b *Integration) Get() *config.Integration {
	return &b.config
}

type integrationConfigurer struct {
	cfg *config.Integration
}

func (c *integrationConfigurer) Identity(name, key string) {
	c.cfg.ConfiguredIdentities = append(
		c.cfg.ConfiguredIdentities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	for _, r := range routes {
		c.cfg.ConfiguredRoutes = append(c.cfg.ConfiguredRoutes, fromRoute(r))
	}
}

func (c *integrationConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.ConfiguredAsDisabled = true
}
