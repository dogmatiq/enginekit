package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromIntegration returns a new [config.Integration] that represents the
// configuration of the given [dogma.IntegrationMessageHandler].
func FromIntegration(h dogma.IntegrationMessageHandler) config.Integration {
	var cfg config.Integration

	if h == nil {
		return cfg
	}

	cfg.TypeName = optional.Some(typename.Of(h))
	cfg.Implementation = optional.Some(h)

	h.Configure(&integrationConfigurer{&cfg})

	return cfg
}

type integrationConfigurer struct {
	cfg *config.Integration
}

func (c *integrationConfigurer) Identity(name, key string) {
	c.cfg.Identities = append(
		c.cfg.Identities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	for _, r := range routes {
		c.cfg.Routes = append(c.cfg.Routes, fromRoute(r))
	}
}

func (c *integrationConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.IsDisabled = true
}
