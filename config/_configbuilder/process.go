package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Process builds [config.Process] values.
type Process struct {
	config config.Process
}

func (b *Process) Configurer() dogma.ProcessConfigurer {
	panic("not implemented")
}

func (b *Process) Get() *config.Process {
	return &b.config
}

type processConfigurer struct {
	cfg *config.Process
}

func (c *processConfigurer) Identity(name, key string) {
	c.cfg.ConfiguredIdentities = append(
		c.cfg.ConfiguredIdentities,
		config.Identity{
			Name: name,
			Key:  key,
		},
	)
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	for _, r := range routes {
		c.cfg.ConfiguredRoutes = append(c.cfg.ConfiguredRoutes, fromRoute(r))
	}
}

func (c *processConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.ConfiguredAsDisabled = true
}
