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
	Routes(filter ...RouteType) []Route

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	routes() []Route
}

// HandlerType is an enumeration of the types of message handlers.
type HandlerType int

const (
	// AggregateHandlerType is the [HandlerType] for implementations of
	// [dogma.AggregateMessageHandler].
	AggregateHandlerType HandlerType = iota

	// IntegrationHandlerType is the [HandlerType] for implementations of
	// [dogma.IntegrationMessageHandler].
	IntegrationHandlerType

	// ProcessHandlerType is the [HandlerType] for implementations of
	// [dogma.ProcessMessageHandler].
	ProcessHandlerType

	// ProjectionHandlerType is the [HandlerType] for implementations of
	// [dogma.ProjectionMessageHandler].
	ProjectionHandlerType
)

// RouteSpec is describes how a [HandlerType] makes use of a particular
// [RouteType].
type RouteSpec int

const (
	// RouteTypeDisallowed indicates that the [HandlerType] does not support the
	// [RouteType].
	RouteTypeDisallowed RouteSpec = iota

	// RouteTypeAllowed indicates that the [HandlerType] supports the [RouteType],
	// but it is not required.
	RouteTypeAllowed

	// RouteTypeRequired indicates that the [HandlerType] requires at least one
	// route of the [RouteType].
	RouteTypeRequired
)

// RoutingSpec returns a map that describes how the [HandlerType] makes use of
// each [RouteType].
func (t HandlerType) RoutingSpec() map[RouteType]RouteSpec {
	switch t {
	case AggregateHandlerType:
		return map[RouteType]RouteSpec{
			HandlesCommandRouteType: RouteTypeRequired,
			RecordsEventRouteType:   RouteTypeRequired,
		}
	case IntegrationHandlerType:
		return map[RouteType]RouteSpec{
			HandlesCommandRouteType: RouteTypeRequired,
			RecordsEventRouteType:   RouteTypeAllowed,
		}
	case ProcessHandlerType:
		return map[RouteType]RouteSpec{
			HandlesEventRouteType:     RouteTypeRequired,
			ExecutesCommandRouteType:  RouteTypeRequired,
			SchedulesTimeoutRouteType: RouteTypeAllowed,
		}
	case ProjectionHandlerType:
		return map[RouteType]RouteSpec{
			HandlesEventRouteType: RouteTypeRequired,
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

func normalizedRoutes(h Handler, filter ...RouteType) []Route {
	ctx := &normalizationContext{
		Component: h,
	}

	routes := normalizeRoutes(ctx, h, filter...)

	if err := ctx.Err(); err != nil {
		panic(err)
	}

	return routes
}

func normalizeRoutes(ctx *normalizationContext, h Handler, filter ...RouteType) []Route {
	var (
		spec      = h.HandlerType().RoutingSpec()
		missing   maps.Ordered[RouteType, MissingRequiredRouteError]
		duplicate maps.OrderedByKey[routeKey, DuplicateRouteError]
		filtered  []Route
	)

	for rt, req := range spec {
		if req == RouteTypeRequired {
			missing.Set(rt, MissingRequiredRouteError{rt})
		}
	}

	for _, r := range slices.Clone(h.routes()) {
		r = normalize(ctx, r)

		rt, ok := r.RouteType.TryGet()
		if len(filter) == 0 || (ok && slices.Contains(filter, rt)) {
			filtered = append(filtered, r)
		}

		if !ok {
			continue
		}

		if spec[rt] == RouteTypeDisallowed {
			ctx.Fail(UnexpectedRouteError{r})
		} else {
			missing.Remove(rt)
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

	return filtered
}
