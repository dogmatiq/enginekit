package config

import (
	"fmt"
	"slices"
)

// A Handler is a specialization of [Entity] that represents configuration of a
// Dogma message handler.
type Handler interface {
	Entity

	// Routes returns the routes configured for the handler.
	//
	// It panics if the routes are incomplete or invalid.
	Routes(filter ...RouteType) []Route

	routes() []Route
	routeTypes() map[RouteType]bool
}

// MissingRouteTypeError indicates that a [Handler] is missing one of its mandatory
// route types.
type MissingRouteTypeError struct {
	RouteType RouteType
}

func (e MissingRouteTypeError) Error() string {
	return fmt.Sprintf("expected at least one %q route", e.RouteType)
}

// UnexpectedRouteTypeError indicates that a [Handler] is configured with a
// [Route] with a [RouteType] that is not allowed for that handler type.
type UnexpectedRouteTypeError struct {
	UnexpectedRoute Route
}

func (e UnexpectedRouteTypeError) Error() string {
	return fmt.Sprintf("%q routes are not allowed for this handler type", e.UnexpectedRoute.RouteType.Get())
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
	allowed := h.routeTypes()
	present := map[RouteType]struct{}{}
	var routes []Route

	for _, r := range slices.Clone(h.routes()) {
		t, ok := r.RouteType.TryGet()
		if !ok {
			// TODO: invalid route
			continue
		}

		present[t] = struct{}{}

		if _, ok := allowed[t]; !ok {
			ctx.Fail(UnexpectedRouteTypeError{r})
		}

		if len(filter) == 0 || slices.Contains(filter, t) {
			routes = append(routes, r)
		}
	}

	for _, t := range []RouteType{
		HandlesCommandRoute,
		HandlesEventRoute,
		ExecutesCommandRoute,
		RecordsEventRoute,
		SchedulesTimeoutRoute,
	} {
		if _, ok := present[t]; ok {
			continue
		}

		if allowed[t] {
			ctx.Fail(MissingRouteTypeError{t})
		}
	}

	return routes
}
