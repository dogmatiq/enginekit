package config

import (
	"cmp"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// Route represents the configuration of a [dogma.Route].
type Route struct {
	ComponentCommon

	// RouteType is the type of route, if available.
	RouteType optional.Optional[RouteType]

	// MessageTypeName is the fully-qualified name of the Go type that
	// implements the [dogma.Message] interface, if available.
	MessageTypeName optional.Optional[string]

	// MessageType is the [message.Type], if available.
	MessageType optional.Optional[message.Type]
}

func (r *Route) String() string {
	var w strings.Builder

	w.WriteString("route")

	if rt, ok := r.RouteType.TryGet(); ok {
		w.WriteByte(':')
		w.WriteString(rt.String())
	}

	if mt, ok := r.MessageTypeName.TryGet(); ok {
		w.WriteByte(':')
		w.WriteString(typename.Unqualified(mt))
	}

	return w.String()
}

func (r *Route) key() (routeKey, bool) {
	if !r.RouteType.IsPresent() {
		return routeKey{}, false
	}

	if !r.MessageTypeName.IsPresent() {
		return routeKey{}, false
	}

	return routeKey{
		RouteType:       r.RouteType.Get(),
		MessageTypeName: r.MessageTypeName.Get(),
	}, true
}

func (r *Route) validate(ctx *validateContext) {
	validateComponent(
		ctx,
		func(ctx *validateContext) {
			routeType, hasRouteType := r.RouteType.TryGet()
			messageTypeName, hasMessageTypeName := r.MessageTypeName.TryGet()
			messageType, hasMessageType := r.MessageType.TryGet()

			if !hasRouteType {
				ctx.Invalid(UnknownRouteTypeError{})
			}

			if !hasMessageTypeName {
				ctx.Invalid(UnknownMessageTypeError{})
			}

			if hasMessageType {
				if hasRouteType {
					if routeType.MessageKind() != messageType.Kind() {
						ctx.Invalid(MessageKindMismatchError{routeType, messageType})
					}
				}

				if !hasMessageTypeName {
					ctx.Malformed("MessageType is present, but MessageTypeName is not")
				} else if messageTypeName != string(messageType.Name()) {
					ctx.Malformed(
						"MessageTypeName does not match MessageType: %q != %q",
						messageTypeName,
						messageType.Name(),
					)
				}
			}
		},
	)
}

func (r *Route) describe(ctx *describeContext) {
	ctx.DescribeFidelity()

	if rt, ok := r.RouteType.TryGet(); ok {
		ctx.Print(rt.String(), " ")
	}

	ctx.Print("route")

	if mt, ok := r.MessageTypeName.TryGet(); ok {
		ctx.Print(" for ", mt)

		if !r.MessageType.IsPresent() {
			ctx.Print(" (type unavailable)")
		}
	}

	ctx.Print("\n")
	ctx.DescribeErrors()
}

// routeKey is a [comparable] representation of a route's identity. No [Handler]
// may have two routes with the same key.
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

// UnknownRouteTypeError indicates that a [Route] does not specify a
// [RouteType].
type UnknownRouteTypeError struct{}

func (UnknownRouteTypeError) Error() string {
	return "unknown route type"
}

// UnknownMessageTypeError indicates that a [Route] does not specify a message
// type name.
type UnknownMessageTypeError struct{}

func (UnknownMessageTypeError) Error() string {
	return "unknown message type"
}

// MessageKindMismatchError indicates that a [Route] refers to a [message.Type]
// that has a different [message.Kind] than the route's [RouteType].
type MessageKindMismatchError struct {
	RouteType   RouteType
	MessageType message.Type
}

func (e MessageKindMismatchError) Error() string {
	return renderer.Inflect(
		"unexpected message kind: %s is a %s, expected a %s",
		typename.Get(e.MessageType.ReflectType()),
		e.MessageType.Kind(),
		e.RouteType.MessageKind(),
	)
}
