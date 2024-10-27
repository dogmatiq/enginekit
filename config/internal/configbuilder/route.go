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
		b.target.ComponentFidelity |= config.Incomplete
	}
}

// // SetRouteType sets the route type of the route.
// func (b *RouteBuilder) SetRouteType(t config.RouteType) {
// 	b.target.XRoute.RouteType = optional.Some(t)
// }

// // SetMessageTypeName sets the message type name of the route.
// func (b *RouteBuilder) SetMessageTypeName(name string) {
// 	b.target.XRoute.MessageTypeName = optional.Some(name)
// }

// // SetMessageType sets the message type of the route.
// func (b *RouteBuilder) SetMessageType(t message.Type) {
// 	b.target.XRoute.MessageTypeName = optional.Some(typename.Get(t.ReflectType()))
// 	b.target.XRoute.MessageType = optional.Some(t)
// }

// // Edit calls fn, which can apply arbitrary changes to the identity.
// func (b *RouteBuilder) Edit(fn func(*config.XRoute)) {
// 	fn(&b.target.XRoute)
// }

// // Fidelity returns the fidelity of the configuration.
// func (b *RouteBuilder) Fidelity() config.Fidelity {
// 	return b.target.XRoute.Fidelity
// }

// // UpdateFidelity merges f with the current fidelity of the configuration.
// func (b *RouteBuilder) UpdateFidelity(f config.Fidelity) {
// 	b.target.XRoute.Fidelity |= f
// }

// Done completes the configuration of the route.
func (b *RouteBuilder) Done() *config.Route {
	if !b.target.RouteType.IsPresent() || !b.target.MessageTypeName.IsPresent() {
		b.target.ComponentFidelity |= config.Incomplete
	}
	return &b.target
}
