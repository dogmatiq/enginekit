package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

type handlerBuilder interface {
	Identity(func(*configbuilder.IdentityBuilder))
	Route(func(*configbuilder.RouteBuilder))
	Disabled(func(*configbuilder.FlagBuilder))
}

type handlerConfigurer[R dogma.Route] struct {
	b handlerBuilder
}

func (c *handlerConfigurer[R]) Identity(name, key string) {
	c.b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.Name(name)
		b.Key(key)
	})
}

func (c *handlerConfigurer[R]) Routes(routes ...R) {
	for _, r := range routes {
		c.b.Route(func(b *configbuilder.RouteBuilder) {
			b.AsPerRoute(r)
		})
	}
}

func (c *handlerConfigurer[R]) Disable(...dogma.DisableOption) {
	c.b.Disabled(func(b *configbuilder.FlagBuilder) {
		b.Value(true)
	})
}
