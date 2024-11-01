package configbuilder

import (
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// Route returns a new [config.Route] as configured by fn.
func Route(fn func(*RouteBuilder)) *config.Route {
	x := &RouteBuilder{}
	fn(x)
	return x.Done()
}

// RouteBuilder constructs a [config.Route].
type RouteBuilder struct {
	target config.Route
}

// AsPerRoute configures the builder to use the same properties as r.
func (b *RouteBuilder) AsPerRoute(r dogma.Route) {
	set := func(
		rt config.RouteType,
		t reflect.Type,
	) {
		b.target.RouteType = optional.Some(rt)
		b.target.MessageTypeName = optional.Some(typename.Get(t))
		b.target.MessageType = optional.Some(message.TypeFromReflect(t))
	}

	switch r := r.(type) {
	case dogma.HandlesCommandRoute:
		set(config.HandlesCommandRouteType, r.Type)
	case dogma.RecordsEventRoute:
		set(config.RecordsEventRouteType, r.Type)
	case dogma.HandlesEventRoute:
		set(config.HandlesEventRouteType, r.Type)
	case dogma.ExecutesCommandRoute:
		set(config.ExecutesCommandRouteType, r.Type)
	case dogma.SchedulesTimeoutRoute:
		set(config.SchedulesTimeoutRouteType, r.Type)
	default:
		panic("unsupported route type")
	}
}

// RouteType sets the route type of the route.
func (b *RouteBuilder) RouteType(t config.RouteType) {
	b.target.RouteType = optional.Some(t)
}

// MessageTypeName sets the message type name of the route.
func (b *RouteBuilder) MessageTypeName(name string) {
	b.target.MessageTypeName = optional.Some(name)
	b.target.MessageType = optional.None[message.Type]()
}

// MessageType sets the message type of the route.
func (b *RouteBuilder) MessageType(t message.Type) {
	b.target.MessageTypeName = optional.Some(typename.Get(t.ReflectType()))
	b.target.MessageType = optional.Some(t)
}

// Partial marks the compomnent as partially configured.
func (b *RouteBuilder) Partial() {
	b.target.IsPartial = true
}

// Speculative marks the component as speculative.
func (b *RouteBuilder) Speculative() {
	b.target.IsSpeculative = true
}

// Done completes the configuration of the route.
func (b *RouteBuilder) Done() *config.Route {
	return &b.target
}
