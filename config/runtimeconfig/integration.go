package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromIntegration returns a new [config.Integration] that represents the
// configuration of the given [dogma.IntegrationMessageHandler].
func FromIntegration(h dogma.IntegrationMessageHandler) *config.Integration {
	cfg := &config.Integration{
		AsConfigured: config.IntegrationAsConfigured{
			IsDisabled: optional.Some(false),
		},
	}

	if h == nil {
		return cfg
	}

	cfg.AsConfigured.Source.TypeName = optional.Some(typename.Of(h))
	cfg.AsConfigured.Source.Value = optional.Some(h)

	h.Configure(&integrationConfigurer{cfg})

	return cfg
}

type integrationConfigurer struct {
	cfg *config.Integration
}

func (c *integrationConfigurer) Identity(name, key string) {
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

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	for _, r := range routes {
		c.cfg.AsConfigured.Routes = append(c.cfg.AsConfigured.Routes, fromRoute(r))
	}
}

func (c *integrationConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.AsConfigured.IsDisabled = optional.Some(true)
}
