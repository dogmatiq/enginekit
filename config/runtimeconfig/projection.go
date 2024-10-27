package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromProjection returns a new [config.Projection] that represents the
// configuration of the given [dogma.ProjectionMessageHandler].
func FromProjection(h dogma.ProjectionMessageHandler) *config.Projection {
	return configbuilder.Projection(func(b *configbuilder.ProjectionBuilder) {
		buildProjection(b, h)
	})
}

func buildProjection(b *configbuilder.ProjectionBuilder, h dogma.ProjectionMessageHandler) {
	if h != nil {
		b.Source(h)
		h.Configure(&projectionConfigurer{
			handlerConfigurer[dogma.ProjectionRoute]{b},
			b,
		})
	}
}

type projectionConfigurer struct {
	handlerConfigurer[dogma.ProjectionRoute]
	b *configbuilder.ProjectionBuilder
}

func (c *projectionConfigurer) DeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	c.b.DeliveryPolicy(
		func(b *configbuilder.ProjectionDeliveryPolicyBuilder) {
			b.AsPerDeliveryPolicy(p)
		},
	)
}
