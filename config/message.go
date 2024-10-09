package config

import (
	"errors"
	"slices"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// MessageType represents a type of [dogma.Message] that is used by the handler.
type MessageType struct {
	// TypeName is the fully-qualified name of the Go type that implements the
	// [dogma.Message] interface.
	TypeName string

	// Kind is the kind of message represented by this type.
	Kind message.Kind

	// Type is the equivalent [message.Type], if available.
	Type optional.Optional[message.Type]
}

// Route represents a message route to or from a handler.
type Route struct {
	// RouteType is the type of route, if available.
	RouteType optional.Optional[RouteType]

	// MessageType is the type of message being routed, if available.
	MessageType optional.Optional[MessageType]
}

// RouteType is an enumeration of the types of message routes.
type RouteType int

const (
	// HandlesCommandRoute is the route type associated with
	// [dogma.HandlesCommandRoute].
	HandlesCommandRoute RouteType = iota

	// HandlesEventRoute is the route type associated with
	// [dogma.HandlesEventRoute].
	HandlesEventRoute

	// ExecutesCommandRoute is the route type associated with
	// [dogma.ExecutesCommandRoute].
	ExecutesCommandRoute

	// RecordsEventRoute is the route type associated with
	// [dogma.RecordsEventRoute].
	RecordsEventRoute

	// SchedulesTimeoutRoute is the route type associated with
	// [dogma.SchedulesTimeoutRoute].
	SchedulesTimeoutRoute
)

// IsConsume returns true if the route indicates that the handler consumes
// a message type.
func (r RouteType) IsConsume() bool {
	switch r {
	case HandlesCommandRoute, HandlesEventRoute, SchedulesTimeoutRoute:
		return true
	default:
		return false
	}
}

// IsProduce returns true if the route indicates that the handler produces
// a message type.
func (r RouteType) IsProduce() bool {
	switch r {
	case ExecutesCommandRoute, RecordsEventRoute, SchedulesTimeoutRoute:
		return true
	default:
		return false
	}
}

func (r RouteType) String() string {
	switch r {
	case HandlesCommandRoute:
		return "HandlesCommand"
	case HandlesEventRoute:
		return "HandlesEvent"
	case ExecutesCommandRoute:
		return "ExecutesCommand"
	case RecordsEventRoute:
		return "RecordsEvent"
	case SchedulesTimeoutRoute:
		return "SchedulesTimeout"
	default:
		panic("unrecognized route type")
	}
}

func normalizedRoutes(any) []Route {
	panic("not implemented")
}

func normalizeRoutesInPlace(
	h Handler,
	errs *error,
	routes *[]Route,
	types map[RouteType]bool,
) {
	*routes = slices.Clone(*routes)

	has := map[RouteType]struct{}{}

	for _, r := range *routes {
		t, ok := r.RouteType.TryGet()
		if !ok {
			// TODO: invalid route
			continue
		}

		has[t] = struct{}{}

		if _, ok := types[t]; !ok {
			*errs = errors.Join(*errs, UnexpectedRouteTypeError{h, r})
		}
	}

	for t, mandatory := range types {
		if _, ok := has[t]; ok {
			continue
		}

		if mandatory {
			*errs = errors.Join(*errs, MissingRouteError{h, t})
		}
	}
}
