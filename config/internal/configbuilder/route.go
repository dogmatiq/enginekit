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

// SetRoute populates the route type and message type of the route from r.
func (b *RouteBuilder) SetRoute(r dogma.Route) {
	configure := func(
		rt config.RouteType,
		t reflect.Type,
	) {
		b.target.AsConfigured = config.RouteAsConfigured{
			RouteType:       optional.Some(rt),
			MessageTypeName: optional.Some(typename.Get(t)),
			MessageType:     optional.Some(message.TypeFromReflect(t)),
		}
	}

	switch r := r.(type) {
	case dogma.HandlesCommandRoute:
		configure(config.HandlesCommandRouteType, r.Type)
	case dogma.RecordsEventRoute:
		configure(config.RecordsEventRouteType, r.Type)
	case dogma.HandlesEventRoute:
		configure(config.HandlesEventRouteType, r.Type)
	case dogma.ExecutesCommandRoute:
		configure(config.ExecutesCommandRouteType, r.Type)
	case dogma.SchedulesTimeoutRoute:
		configure(config.SchedulesTimeoutRouteType, r.Type)
	default:
		b.target.AsConfigured.Fidelity |= config.Incomplete
	}
}

// SetRouteType sets the route type of the route.
func (b *RouteBuilder) SetRouteType(t config.RouteType) {
	b.target.AsConfigured.RouteType = optional.Some(t)
}

// SetMessageTypeName sets the message type name of the route.
func (b *RouteBuilder) SetMessageTypeName(name string) {
	b.target.AsConfigured.MessageTypeName = optional.Some(name)
}

// SetMessageType sets the message type of the route.
func (b *RouteBuilder) SetMessageType(t message.Type) {
	b.target.AsConfigured.MessageTypeName = optional.Some(typename.Get(t.ReflectType()))
	b.target.AsConfigured.MessageType = optional.Some(t)
}

// Edit calls fn, which can apply arbitrary changes to the identity.
func (b *RouteBuilder) Edit(fn func(*config.RouteAsConfigured)) {
	fn(&b.target.AsConfigured)
}

// UpdateFidelity merges f with the current fidelity of the identity.
func (b *RouteBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.AsConfigured.Fidelity |= f
}

// Done completes the configuration of the identity.
func (b *RouteBuilder) Done() *config.Route {
	if b.target.AsConfigured.Fidelity&config.Incomplete == 0 {
		if !b.target.AsConfigured.RouteType.IsPresent() {
			panic("route must have a route type or be marked as incomplete")
		}
		if !b.target.AsConfigured.MessageTypeName.IsPresent() {
			panic("route must have a message type name or be marked as incomplete")
		}
	}

	return &b.target
}
