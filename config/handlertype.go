package config

import (
	"iter"

	"github.com/dogmatiq/enginekit/internal/enum"
)

// HandlerType is an enumeration of the types of message handlers.
type HandlerType int

const (
	// AggregateHandlerType is the [HandlerType] for implementations of
	// [dogma.AggregateMessageHandler].
	AggregateHandlerType HandlerType = iota

	// ProcessHandlerType is the [HandlerType] for implementations of
	// [dogma.ProcessMessageHandler].
	ProcessHandlerType

	// IntegrationHandlerType is the [HandlerType] for implementations of
	// [dogma.IntegrationMessageHandler].
	IntegrationHandlerType

	// ProjectionHandlerType is the [HandlerType] for implementations of
	// [dogma.ProjectionMessageHandler].
	ProjectionHandlerType
)

// HandlerTypes returns a sequence that yields all valid [HandlerType] values.
func HandlerTypes() iter.Seq[HandlerType] {
	return enum.Range(AggregateHandlerType, ProjectionHandlerType)
}

func (t HandlerType) String() string {
	return enum.String(t, "aggregate", "process", "integration", "projection")
}

// RouteCapabilities returns a value that describes the routing capabilities of
// the handler type.
func (t HandlerType) RouteCapabilities() RouteCapabilities {
	return RouteCapabilities{
		MapByHandlerType(
			t,
			map[RouteType]RouteTypeCapability{
				HandlesCommandRouteType: RouteTypeRequired,
				RecordsEventRouteType:   RouteTypeRequired,
			},
			map[RouteType]RouteTypeCapability{
				HandlesEventRouteType:     RouteTypeRequired,
				ExecutesCommandRouteType:  RouteTypeRequired,
				SchedulesTimeoutRouteType: RouteTypeAllowed,
			},
			map[RouteType]RouteTypeCapability{
				HandlesCommandRouteType: RouteTypeRequired,
				RecordsEventRouteType:   RouteTypeAllowed,
			},
			map[RouteType]RouteTypeCapability{
				HandlesEventRouteType: RouteTypeRequired,
			},
		),
	}
}
