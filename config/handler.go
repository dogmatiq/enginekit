package config

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/enginekit/collections/maps"
	"github.com/dogmatiq/enginekit/message"
)

// A Handler is an [Entity] that represents the configuration of a Dogma
// handler.
type Handler interface {
	Entity

	// HandlerType returns [HandlerType] of the handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	//
	// It panics if the configuration does not specify unambiguously whether the
	// handler is enabled or disabled.
	IsDisabled() bool

	routes() []*Route
}

// HandlerCommon is a partial implementation of [Handler].
type HandlerCommon[T any] struct {
	EntityCommon[T]

	RouteComponents []*Route
	DisabledFlag    Flag[Disabled]
}

// IsDisabled returns true if the handler is disabled.
//
// It panics if the configuration does not specify unambiguously whether the
// handler is enabled or disabled.
func (h *HandlerCommon[T]) IsDisabled() bool {
	panic("not implemented")
}

func (h *HandlerCommon[T]) routes() []*Route {
	return h.RouteComponents
}

func (h *HandlerCommon[T]) validate(ctx *validateContext, t HandlerType) {
	h.EntityCommon.validate(ctx)
	h.validateRoutes(ctx, t)
}

func (h *HandlerCommon[T]) validateRoutes(ctx *validateContext, t HandlerType) {
	var (
		capabilities = t.RouteCapabilities()
		missing      maps.Ordered[RouteType, MissingRouteTypeError]
		duplicate    maps.OrderedByKey[routeKey, DuplicateRouteError]
	)

	for rt, req := range capabilities.RouteTypes {
		if req == RouteTypeRequired {
			missing.Set(rt, MissingRouteTypeError{rt})
		}
	}

	for _, r := range h.RouteComponents {
		if rt, ok := r.RouteType.TryGet(); ok {
			if capabilities.RouteTypes[rt] == RouteTypeDisallowed {
				ctx.Fail(UnexpectedRouteTypeError{r})
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
}

func (h *HandlerCommon[T]) describe(ctx *describeContext) {
	h.EntityCommon.describe(ctx)

	for _, r := range h.RouteComponents {
		ctx.DescribeChild(r)
	}
}

// Disabled is the [Symbol] for a [Flag] that indicates whether or not a
// [Handler] has been disabled.
type Disabled struct{ symbol }

// MissingRouteTypeError indicates that a [Handler] is missing one of its
// required route types.
type MissingRouteTypeError struct {
	RouteType RouteType
}

func (e MissingRouteTypeError) Error() string {
	return fmt.Sprintf("no %q routes configured", e.RouteType)
}

// UnexpectedRouteTypeError indicates that a [Handler] is configured with a
// [Route] of a [RouteType] that is not allowed for that handler type.
type UnexpectedRouteTypeError struct {
	UnexpectedRoute *Route
}

func (e UnexpectedRouteTypeError) Error() string {
	w := &strings.Builder{}

	fmt.Fprintf(w, "unexpected %s route", e.UnexpectedRoute.RouteType.Get())

	if name, ok := e.UnexpectedRoute.MessageTypeName.TryGet(); ok {
		fmt.Fprintf(w, " for %s", name)
	}

	return w.String()
}

// DuplicateRouteError indicates that a [Handler] is configured with multiple
// routes for the same [message.Type].
type DuplicateRouteError struct {
	RouteType       RouteType
	MessageTypeName string
	DuplicateRoutes []*Route
}

func (e DuplicateRouteError) Error() string {
	return fmt.Sprintf(
		"multiple %q routes configured for %s",
		e.RouteType,
		e.MessageTypeName,
	)
}

func routeSetForHandler(h Handler) RouteSet {
	set := RouteSet{}

	for _, r := range h.routes() {
		rt, ok := r.RouteType.TryGet()
		if !ok {
			continue
		}

		mt, ok := r.MessageType.TryGet()
		if !ok {
			continue
		}

		if set.byMessageType == nil {
			set.byMessageType = map[message.Type]map[RouteType]map[Handler]*Route{}
		}

		byRouteType, ok := set.byMessageType[mt]
		if !ok {
			byRouteType = map[RouteType]map[Handler]*Route{}
			set.byMessageType[mt] = byRouteType
		}

		byHandler, ok := byRouteType[rt]
		if !ok {
			byHandler = map[Handler]*Route{}
			byRouteType[rt] = byHandler
		}

		byHandler[h] = r
	}

	return set
}
