package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromProcess returns a new [config.Process] that represents the configuration
// of the given [dogma.ProcessMessageHandler].
func FromProcess(h dogma.ProcessMessageHandler) *config.Process {
	cfg := &config.Process{
		AsConfigured: config.ProcessAsConfigured{
			IsDisabled: optional.Some(false),
		},
	}

	if h == nil {
		return cfg
	}

	cfg.AsConfigured.Source = optional.Some(
		config.Value[dogma.ProcessMessageHandler]{
			TypeName: optional.Some(typename.Of(h)),
			Value:    optional.Some(h),
		},
	)

	h.Configure(&processConfigurer{cfg})

	return cfg
}

type processConfigurer struct {
	cfg *config.Process
}

func (c *processConfigurer) Identity(name, key string) {
	c.cfg.AsConfigured.Identities = append(
		c.cfg.AsConfigured.Identities,
		config.Identity{
			AsConfigured: config.IdentityAsConfigured{
				Name: optional.Some(name),
				Key:  optional.Some(key),
			},
		},
	)
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	for _, r := range routes {
		c.cfg.AsConfigured.Routes = append(c.cfg.AsConfigured.Routes, fromRoute(r))
	}
}

func (c *processConfigurer) Disable(...dogma.DisableOption) {
	c.cfg.AsConfigured.IsDisabled = optional.Some(true)
}
