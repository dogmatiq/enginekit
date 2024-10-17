package configbuilder

import (
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// Route returns an [RouteBuilder] that builds a new [config.Route].
func Route() *RouteBuilder {
	return &RouteBuilder{}
}

// RouteBuilder constructs a [config.Route].
type RouteBuilder struct {
	target   config.Route
	appendTo *[]*config.Route
}

// SetRoute populates the route type and message type of the route from r.
func (b *RouteBuilder) SetRoute(r dogma.Route) *RouteBuilder {
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

	return b
}

// SetRouteType sets the route type of the route.
func (b *RouteBuilder) SetRouteType(t config.RouteType) *RouteBuilder {
	b.target.AsConfigured.RouteType = optional.Some(t)
	return b
}

// SetMessageTypeName sets the message type name of the route.
func (b *RouteBuilder) SetMessageTypeName(name string) *RouteBuilder {
	b.target.AsConfigured.MessageTypeName = optional.Some(name)
	return b
}

// SetMessageType sets the message type of the route.
func (b *RouteBuilder) SetMessageType(t message.Type) *RouteBuilder {
	b.target.AsConfigured.MessageTypeName = optional.Some(typename.Get(t.ReflectType()))
	b.target.AsConfigured.MessageType = optional.Some(t)
	return b
}

// Edit calls fn, which can apply arbitrary changes to the identity.
func (b *RouteBuilder) Edit(fn func(*config.RouteAsConfigured)) *RouteBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// UpdateFidelity merges f with the current fidelity of the identity.
func (b *RouteBuilder) UpdateFidelity(f config.Fidelity) *RouteBuilder {
	b.target.AsConfigured.Fidelity |= f
	return b
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

	if b.appendTo != nil {
		*b.appendTo = append(*b.appendTo, &b.target)
		b.appendTo = nil
	}

	return &b.target
}
