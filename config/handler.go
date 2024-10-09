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

	routes() []Route
	routeSpec() routeSpec
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

// DuplicateRouteError indicates that a [Handler] is configured with multiple
// routes for the same [MessageType].
type DuplicateRouteError struct {
	MessageType     MessageType
	RouteType       RouteType
	DuplicateRoutes []Route
}

func (e DuplicateRouteError) Error() string {
	return fmt.Sprintf(
		"multiple %q routes are configured for %s",
		e.RouteType,
		e.MessageType,
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
		spec      = h.routeSpec()
		missing   maps.Ordered[RouteType, MissingRouteTypeError]
		duplicate maps.Ordered[string, DuplicateRouteError]
		filtered  []Route
	)

	for rt, req := range spec {
		if req == required {
			missing.Set(rt, MissingRouteTypeError{rt})
		}
	}

	for _, r := range slices.Clone(h.routes()) {
		rt, ok := r.RouteType.TryGet()
		if !ok {
			// TODO: invalid route
			continue
		}

		if len(filter) == 0 || slices.Contains(filter, rt) {
			filtered = append(filtered, r)
		}

		if spec[rt] == disallowed {
			ctx.Fail(UnexpectedRouteTypeError{r})
		} else {
			missing.Remove(rt)
		}

		mt, ok := r.MessageType.TryGet()
		if !ok {
			// TODO: invalid route
			continue
		}

		duplicate.Update(
			mt.TypeName,
			func(err *DuplicateRouteError) {
				err.MessageType = mt
				err.RouteType = rt
				err.DuplicateRoutes = append(err.DuplicateRoutes, r)
			},
		)
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
