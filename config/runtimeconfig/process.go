package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromProcess returns a new [config.Process] that represents the configuration
// of the given [dogma.ProcessMessageHandler].
func FromProcess(h dogma.ProcessMessageHandler) *config.Process {
	b := configbuilder.Process()

	if h == nil {
		return b.Done(config.Incomplete)
	}

	b.SetDisabled(false)
	b.SetSource(h)
	h.Configure(&processConfigurer{b})

	return b.Done(config.Immaculate)
}

type processConfigurer struct {
	b *configbuilder.ProcessBuilder
}

func (c *processConfigurer) Identity(name, key string) {
	c.b.
		AddIdentity().
		SetName(name).
		SetKey(key).
		Done(config.Immaculate)
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	c.b.Edit(
		func(cfg *config.ProcessAsConfigured) {
			for _, r := range routes {
				cfg.Routes = append(cfg.Routes, fromRoute(r))
			}
		},
	)
}

func (c *processConfigurer) Disable(...dogma.DisableOption) {
	c.b.SetDisabled(true)
}
