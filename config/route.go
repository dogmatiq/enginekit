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

func (r Route) key() (routeKey, bool) {
	rt, rtOK := r.RouteType.TryGet()
	mt, mtOK := r.MessageTypeName.TryGet()
	return routeKey{rt, mt}, rtOK && mtOK
}

// RouteType is an enumeration of the types of message routes.
type RouteType int

const (
	// HandlesCommandRouteType is the [RouteType] associated with
	// [dogma.HandlesCommandRouteType].
	HandlesCommandRouteType RouteType = iota

	// HandlesEventRouteType is the [RouteType] associated with
	// [dogma.HandlesEventRouteType].
	HandlesEventRouteType

	// ExecutesCommandRouteType is the [RouteType] associated with
	// [dogma.ExecutesCommandRouteType].
	ExecutesCommandRouteType

	// RecordsEventRouteType is the [RouteType] associated with
	// [dogma.RecordsEventRouteType].
	RecordsEventRouteType

	// SchedulesTimeoutRouteType is the [RouteType] associated with
	// [dogma.SchedulesTimeoutRouteType].
	SchedulesTimeoutRouteType
)

// IsInbound returns true if the route indicates that the handler consumes
// a message type.
func (r RouteType) IsInbound() bool {
	switch r {
	case HandlesCommandRouteType,
		HandlesEventRouteType,
		SchedulesTimeoutRouteType:
		return true
	default:
		return false
	}
}

// IsOutbound returns true if the route indicates that the handler produces a
// message type.
func (r RouteType) IsOutbound() bool {
	switch r {
	case ExecutesCommandRouteType,
		RecordsEventRouteType,
		SchedulesTimeoutRouteType:
		return true
	default:
		return false
	}
}

// Kind returns the kind of message that the route type is associated with.
func (r RouteType) Kind() message.Kind {
	switch r {
	case HandlesCommandRouteType, ExecutesCommandRouteType:
		return message.CommandKind
	case HandlesEventRouteType, RecordsEventRouteType:
		return message.EventKind
	case SchedulesTimeoutRouteType:
		return message.TimeoutKind
	default:
		panic("unrecognized route type")
	}
}

func (r RouteType) String() string {
	switch r {
	case HandlesCommandRouteType:
		return "HandlesCommand"
	case HandlesEventRouteType:
		return "HandlesEvent"
	case ExecutesCommandRouteType:
		return "ExecutesCommand"
	case RecordsEventRouteType:
		return "RecordsEvent"
	case SchedulesTimeoutRouteType:
		return "SchedulesTimeout"
	default:
		panic("unrecognized route type")
	}
}

// MessageKindMismatchError indicates that two components that refer to the same
// message type disagree on the kind of message they are referring to.
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

// MissingRouteTypeError indicates that a [Route] is missing the route type.
type MissingRouteTypeError struct{}

func (e MissingRouteTypeError) Error() string {
	return "missing route type"
}

// MissingMessageTypeError indicates that a [Route] is missing the message type.
type MissingMessageTypeError struct{}

func (e MissingMessageTypeError) Error() string {
	return "missing message type"
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
