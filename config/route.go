package config

import (
	"cmp"
	"fmt"

	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// Route represents a message route to or from a handler.
type Route struct {
	// RouteType is the type of route, if available.
	RouteType optional.Optional[RouteType]

	// MessageTypeName is the fully-qualified name of the Go type that
	// implements the [dogma.Message] interface, if available.
	MessageTypeName optional.Optional[string]

	// MessageType is the [message.Type], if available.
	MessageType optional.Optional[message.Type]
}

func (r Route) String() string {
	s := "route"

	if rt, ok := r.RouteType.TryGet(); ok {
		s = rt.String()
	}

	if mt, ok := r.MessageTypeName.TryGet(); ok {
		s += ":" + mt
	}

	return s
}

func (r Route) normalize(ctx *normalizationContext) Component {
	routeType, hasRouteType := r.RouteType.TryGet()
	typeName, hasTypeName := r.MessageTypeName.TryGet()
	messageType, hasMessageType := r.MessageType.TryGet()

	if !hasRouteType {
		ctx.Fail(MissingRouteTypeError{})
	}

	if !hasTypeName {
		ctx.Fail(MissingMessageTypeError{})
	}

	if hasMessageType {
		if hasRouteType && routeType.Kind() != messageType.Kind() {
			ctx.Fail(MessageKindMismatchError{routeType, messageType})
		}

		actualTypeName := typename.Get(messageType.ReflectType())
		if hasTypeName && typeName != actualTypeName {
			ctx.Fail(ImplementationTypeNameMismatchError{actualTypeName, typeName})
		}

		r.MessageTypeName = optional.Some(actualTypeName)
	} else if ctx.Options.RequireImplementations {
		ctx.Fail(MissingImplementationError{})
	}

	return r
}

// routeKey is the components of a [Route] that uniquely identify it.
type routeKey struct {
	RouteType       RouteType
	MessageTypeName string
}

func (k routeKey) Compare(x routeKey) int {
	if c := cmp.Compare(k.RouteType, x.RouteType); c != 0 {
		return c
	}
	return cmp.Compare(k.MessageTypeName, x.MessageTypeName)
}

func (r Route) key() (routeKey, bool) {
	rt, rtOK := r.RouteType.TryGet()
	mt, mtOK := r.MessageTypeName.TryGet()
	return routeKey{rt, mt}, rtOK && mtOK
}

// MissingRouteTypeError indicates that a [Route] is missing its [RouteType].
type MissingRouteTypeError struct{}

func (e MissingRouteTypeError) Error() string {
	return "missing route type"
}

// MissingMessageTypeError indicates that a [Route] is missing information about
// the message type.
type MissingMessageTypeError struct{}

func (e MissingMessageTypeError) Error() string {
	return "missing message type"
}

// MessageKindMismatchError indicates that a [Route] refers to a [message.Type]
// that has a different [message.Kind] than the route's [RouteType].
type MessageKindMismatchError struct {
	RouteType   RouteType
	MessageType message.Type
}

func (e MessageKindMismatchError) Error() string {
	return fmt.Sprintf(
		"message kind mismatch: %s expects %q, but %s is %q",
		e.RouteType,
		e.RouteType.Kind(),
		typename.Get(e.MessageType.ReflectType()),
		e.MessageType.Kind(),
	)
}
