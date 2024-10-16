package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromIntegration returns a new [config.Integration] that represents the
// configuration of the given [dogma.IntegrationMessageHandler].
func FromIntegration(h dogma.IntegrationMessageHandler) *config.Integration {
	b := configbuilder.Integration()

	if h == nil {
		b.UpdateFidelity(config.Incomplete)
	} else {
		b.SetDisabled(false)
		b.SetSource(h)
		h.Configure(&integrationConfigurer{b})
	}

	return b.Done()
}

type integrationConfigurer struct {
	b *configbuilder.IntegrationBuilder
}

func (c *integrationConfigurer) Identity(name, key string) {
	c.b.
		AddIdentity().
		SetName(name).
		SetKey(key).
		Done()
}

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	c.b.Edit(
		func(cfg *config.IntegrationAsConfigured) {
			for _, r := range routes {
				cfg.Routes = append(cfg.Routes, fromRoute(r))
			}
		},
	)
}

func (c *integrationConfigurer) Disable(...dogma.DisableOption) {
	c.b.SetDisabled(true)
}
