package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromProjection returns a new [config.Projection] that represents the
// configuration of the given [dogma.ProjectionMessageHandler].
func FromProjection(h dogma.ProjectionMessageHandler) *config.Projection {
	b := configbuilder.Projection()

	if h == nil {
		return b.Done(config.Incomplete)
	}

	b.SetDisabled(false)
	b.SetSource(h)
	b.SetDeliveryPolicy(dogma.UnicastProjectionDeliveryPolicy{})
	h.Configure(&projectionConfigurer{b})

	return b.Done(config.Immaculate)
}

type projectionConfigurer struct {
	b *configbuilder.ProjectionBuilder
}

func (c *projectionConfigurer) Identity(name, key string) {
	c.b.
		AddIdentity().
		SetName(name).
		SetKey(key).
		Done(config.Immaculate)
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	c.b.Edit(
		func(cfg *config.ProjectionAsConfigured) {
			for _, r := range routes {
				cfg.Routes = append(cfg.Routes, fromRoute(r))
			}
		},
	)
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	c.b.SetDeliveryPolicy(p)
}

func (c *projectionConfigurer) Disable(...dogma.DisableOption) {
	c.b.SetDisabled(true)
}
