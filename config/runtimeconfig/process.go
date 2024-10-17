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
		b.UpdateFidelity(config.Incomplete)
	} else {
		b.SetDisabled(false)
		b.SetSource(h)
		h.Configure(&processConfigurer{b})
	}

	return b.Done()
}

type processConfigurer struct {
	b *configbuilder.ProcessBuilder
}

func (c *processConfigurer) Identity(name, key string) {
	c.b.
		AddIdentity().
		SetName(name).
		SetKey(key).
		Done()
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	for _, r := range routes {
		c.b.
			AddRoute().
			SetRoute(r).
			Done()
	}
}

func (c *processConfigurer) Disable(...dogma.DisableOption) {
	c.b.SetDisabled(true)
}
