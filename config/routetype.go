package config

import (
	"github.com/dogmatiq/enginekit/message"
)

// RouteDirection is an enumeration of the "directions" in which a message flows
// for a specific [RouteType].
type RouteDirection int

const (
	// Inbound is a [RouteDirection] that indicates a message flowing into a
	// handler.
	Inbound RouteDirection = 1 << iota

	// Outbound is a [RouteDirection] that indicates a message flowing out of a
	// handler.
	Outbound
)

// Is returns true if d includes dir.
func (d RouteDirection) Is(dir RouteDirection) bool {
	return d&dir != 0
}

// RouteType is an enumeration of the types of message routes that can be
// configured on a [Handler].
type RouteType int

const (
	// HandlesCommandRouteType is a [RouteType] that indicates a [Handler]
	// handles a specific type of [dogma.Command] message.
	//
	// It is associated with routes configured by [dogma.HandlesCommand].
	HandlesCommandRouteType RouteType = iota

	// HandlesEventRouteType is a [RouteType] that indicates a [Handler] handles
	// a specific type of [dogma.Event] message.
	//
	// It is associated with routes configured by [dogma.HandlesEvent].
	HandlesEventRouteType

	// ExecutesCommandRouteType is a [RouteType] that indicates a [Handler]
	// executes a specific type of [dogma.Command] message.
	//
	// It is associated with routes configured by [dogma.ExecutesCommand].
	ExecutesCommandRouteType

	// RecordsEventRouteType is a [RouteType] that indicates a [Handler] records
	// a specific type of [dogma.Event] message.
	//
	// It is associated with routes configured by [dogma.RecordsEvent].
	RecordsEventRouteType

	// SchedulesTimeoutRouteType is a [RouteType] that indicates a [Handler]
	// schedules a specific type of [dogma.Timeout] message.
	//
	// It is associated with routes configured by [dogma.SchedulesTimeout].
	SchedulesTimeoutRouteType
)

// Direction returns the direction in which messages flow for the route type.
func (r RouteType) Direction() RouteDirection {
	switch r {
	case HandlesCommandRouteType, HandlesEventRouteType:
		return Inbound
	case ExecutesCommandRouteType, RecordsEventRouteType:
		return Outbound
	case SchedulesTimeoutRouteType:
		return Inbound | Outbound
	default:
		panic("unrecognized route type")
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
