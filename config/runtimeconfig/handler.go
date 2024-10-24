package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

type handlerConfigurer[R dogma.Route] struct {
	b configbuilder.HandlerBuilder
}

func (c *handlerConfigurer[R]) Identity(name, key string) {
	c.b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.SetName(name)
		b.SetKey(key)
	})
}

func (c *handlerConfigurer[R]) Routes(routes ...R) {
	for _, r := range routes {
		c.b.Route(func(b *configbuilder.RouteBuilder) {
			b.SetRoute(r)
		})
	}
}

func (c *handlerConfigurer[R]) Disable(...dogma.DisableOption) {
	c.b.Disable(func(*configbuilder.FlagBuilder[config.Disabled]) {})
}
