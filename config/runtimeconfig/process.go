package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromProcess returns a new [config.Process] that represents the configuration
// of the given [dogma.ProcessMessageHandler].
func FromProcess(h dogma.ProcessMessageHandler) config.Process {
	var cfg config.Process

	if h == nil {
		return cfg
	}

	cfg.ConfigurationIsExhaustive = true
	cfg.Impl = optional.Some(
		config.Implementation[dogma.ProcessMessageHandler]{
			TypeName: typename.Of(h),
			Source:   optional.Some(h),
		},
	)

	h.Configure(&processConfigurer{&cfg})

	return cfg
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
