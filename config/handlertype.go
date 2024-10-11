package config

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

// HandlerTypes returns a list of all [HandlerType] values.
func HandlerTypes() []HandlerType {
	return []HandlerType{
		AggregateHandlerType,
		ProcessHandlerType,
		IntegrationHandlerType,
		ProjectionHandlerType,
	}
}

func (t HandlerType) String() string {
	switch t {
	case AggregateHandlerType:
		return "aggregate"
	case ProcessHandlerType:
		return "process"
	case IntegrationHandlerType:
		return "integration"
	case ProjectionHandlerType:
		return "projection"
	default:
		panic("invalid handler type")
	}
}

// SwitchByHandlerTypeOf invokes one of the provided functions based on the
// [HandlerType] of h.
func SwitchByHandlerTypeOf(
	h Handler,
	aggregate func(*Aggregate),
	process func(*Process),
	integration func(*Integration),
	projection func(*Projection),
) {
	switch h := h.(type) {
	case *Aggregate:
		if aggregate == nil {
			panic("no case function was provided for aggregate handlers")
		}
		aggregate(h)
	case *Process:
		if process == nil {
			panic("no case function was provided for process handlers")
		}
		process(h)
	case *Integration:
		if integration == nil {
			panic("no case function was provided for integration handlers")
		}
		integration(h)
	case *Projection:
		if projection == nil {
			panic("no case function was provided for projection handlers")
		}
		projection(h)
	default:
		panic("invalid handler type")
	}
}

// RouteCapabilities returns a value that describes the routing capabilities of
// the handler type.
func (t HandlerType) RouteCapabilities() RouteCapabilities {
	switch t {
	case AggregateHandlerType:
		return RouteCapabilities{
			map[RouteType]RouteTypeCapability{
				HandlesCommandRouteType: RouteTypeRequired,
				RecordsEventRouteType:   RouteTypeRequired,
			},
		}
	case IntegrationHandlerType:
		return RouteCapabilities{
			map[RouteType]RouteTypeCapability{
				HandlesCommandRouteType: RouteTypeRequired,
				RecordsEventRouteType:   RouteTypeAllowed,
			},
		}
	case ProcessHandlerType:
		return RouteCapabilities{
			map[RouteType]RouteTypeCapability{
				HandlesEventRouteType:     RouteTypeRequired,
				ExecutesCommandRouteType:  RouteTypeRequired,
				SchedulesTimeoutRouteType: RouteTypeAllowed,
			},
		}
	case ProjectionHandlerType:
		return RouteCapabilities{
			map[RouteType]RouteTypeCapability{
				HandlesEventRouteType: RouteTypeRequired,
			},
		}
	default:
		panic("invalid handler type")
	}
}
