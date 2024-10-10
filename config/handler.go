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
		spec      = h.routeSpec()
		missing   maps.Ordered[RouteType, MissingRequiredRouteError]
		duplicate maps.OrderedByKey[routeKey, DuplicateRouteError]
		filtered  []Route
	)

	for rt, req := range spec {
		if req == required {
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

		if spec[rt] == disallowed {
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
