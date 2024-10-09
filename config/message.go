package config

import (
	"fmt"

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

func (m MessageType) String() string {
	return fmt.Sprintf("%s:%s", m.Kind, m.TypeName)
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

// IsInbound returns true if the route indicates that the handler consumes
// a message type.
func (r RouteType) IsInbound() bool {
	switch r {
	case HandlesCommandRoute, HandlesEventRoute, SchedulesTimeoutRoute:
		return true
	default:
		return false
	}
}

// IsOutbound returns true if the route indicates that the handler produces a
// message type.
func (r RouteType) IsOutbound() bool {
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

// routeSpec is a specification of the types of routes that can (and must) be
// configured for a specific handler type.
type routeSpec map[RouteType]requirement

// requirement is an enumeration of the "requirement level" of some value.
type requirement int

const (
	disallowed requirement = iota
	allowed
	required
)
