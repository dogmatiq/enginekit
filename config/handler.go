package config

import (
	"fmt"
	"slices"

	"github.com/dogmatiq/enginekit/collections/maps"
)

// A Handler is a specialization of [Entity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	Entity

	// Routes returns the routes configured for the handler.
	//
	// It panics if the routes are incomplete or invalid.
	Routes() RouteSet

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler was disabled via the configurer.
	IsDisabled() bool

	routes() []Route
}

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
	aggregate func(Aggregate),
	process func(Process),
	integration func(Integration),
	projection func(Projection),
) {
	switch h := h.(type) {
	case Aggregate:
		if aggregate == nil {
			panic("no case function was provided for aggregate handlers")
		}
		aggregate(h)
	case Process:
		if process == nil {
			panic("no case function was provided for process handlers")
		}
		process(h)
	case Integration:
		if integration == nil {
			panic("no case function was provided for integration handlers")
		}
		integration(h)
	case Projection:
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

// MissingRequiredRouteError indicates that a [Handler] is missing one of its
// required route types.
type MissingRequiredRouteError struct {
	RouteType RouteType
}

func (e MissingRequiredRouteError) Error() string {
	return fmt.Sprintf("expected at least one %q route", e.RouteType)
}

// UnexpectedRouteError indicates that a [Handler] is configured with a [Route]
// with a [RouteType] that is not allowed for that handler type.
type UnexpectedRouteError struct {
	UnexpectedRoute Route
}

func (e UnexpectedRouteError) Error() string {
	return fmt.Sprintf("unexpected route: %s", e.UnexpectedRoute)
}

// DuplicateRouteError indicates that a [Handler] is configured with multiple
// routes for the same [MessageType].
type DuplicateRouteError struct {
	RouteType       RouteType
	MessageTypeName string
	DuplicateRoutes []Route
}

func (e DuplicateRouteError) Error() string {
	return fmt.Sprintf(
		"multiple %q routes are configured for %s",
		e.RouteType,
		e.MessageTypeName,
	)
}

func normalizedRouteSet(h Handler) RouteSet {
	ctx := &normalizationContext{
		Component: h,
	}

	routes := normalizeRoutes(ctx, h)

	if err := ctx.Err(); err != nil {
		panic(err)
	}

	return routes
}

func normalizeRoutes(ctx *normalizationContext, h Handler) []Route {
	var (
		capabilities = h.HandlerType().RouteCapabilities()
		missing      maps.Ordered[RouteType, MissingRequiredRouteError]
		duplicate    maps.OrderedByKey[routeKey, DuplicateRouteError]
	)

	for rt, req := range capabilities.RouteTypes {
		if req == RouteTypeRequired {
			missing.Set(rt, MissingRequiredRouteError{rt})
		}
	}

	routes := slices.Clone(h.routes())

	for i, r := range routes {
		r = normalize(ctx, r)
		routes[i] = r

		if rt, ok := r.RouteType.TryGet(); ok {
			if capabilities.RouteTypes[rt] == RouteTypeDisallowed {
				ctx.Fail(UnexpectedRouteError{r})
			} else {
				missing.Remove(rt)
			}
		}

		if k, ok := r.key(); ok {
			duplicate.Update(
				k,
				func(err *DuplicateRouteError) {
					err.RouteType = k.RouteType
					err.MessageTypeName = k.MessageTypeName
					err.DuplicateRoutes = append(err.DuplicateRoutes, r)
				},
			)
		}
	}

	for err := range missing.Values() {
		ctx.Fail(err)
	}

	for err := range duplicate.Values() {
		if len(err.DuplicateRoutes) > 1 {
			ctx.Fail(err)
		}
	}

	return routes
}
