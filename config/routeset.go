package config

import (
	"iter"
	"slices"

	"github.com/dogmatiq/enginekit/message"
)

// RouteSet is the set of routes configured for a specific [Handler].
type RouteSet struct {
	byMessageType map[message.Type]map[RouteType]map[Handler]Route
}

// Routes returns a sequence that yields the routes in the set, and the handler
// it belongs to.
func (s RouteSet) Routes() iter.Seq2[Route, Handler] {
	return func(yield func(Route, Handler) bool) {
		for _, byRouteType := range s.byMessageType {
			for _, byHandler := range byRouteType {
				for h, r := range byHandler {
					if !yield(r, h) {
						return
					}
				}
			}
		}
	}
}

// MessageTypes returns a map all of the message types in the [RouteSet] and
// their respective [RouteDirection].
func (s RouteSet) MessageTypes() map[message.Type]RouteDirection {
	types := map[message.Type]RouteDirection{}

	for mt, byRouteType := range s.byMessageType {
		for rt := range byRouteType {
			types[mt] |= rt.Direction()
		}
	}

	return types
}

// DirectionOf returns the direction in which messages of the given type flow
// for the [Handler].
func (s RouteSet) DirectionOf(t message.Type) RouteDirection {
	var dir RouteDirection

	for rt := range s.byMessageType[t] {
		dir |= rt.Direction()
	}

	return dir
}

// Filter returns a new [RouteSet] that contains only the routes that match all
// of the given filters.
func (s RouteSet) Filter(filters ...RouteSetFilter) RouteSet {
	var filter routeSetFilters
	for _, f := range filters {
		f(&filter)
	}

	byMessageTypeFiltered := map[message.Type]map[RouteType]map[Handler]Route{}

	for mt, byRouteType := range s.byMessageType {
		byRouteTypeFiltered := map[RouteType]map[Handler]Route{}

		for rt, byHandler := range byRouteType {
			byHandlerFiltered := map[Handler]Route{}

			for h, r := range byHandler {
				if filter.TestRoute(r) {
					byHandlerFiltered[h] = r
				}
			}

			if len(byHandlerFiltered) != 0 {
				byRouteTypeFiltered[rt] = byHandlerFiltered
			}
		}

		if len(byRouteTypeFiltered) != 0 && filter.TestMessage(mt, byRouteTypeFiltered) {
			byMessageTypeFiltered[mt] = byRouteTypeFiltered
		}
	}

	return RouteSet{byMessageTypeFiltered}
}

func (s *RouteSet) merge(set RouteSet) {
	for mt, byRouteTypeSource := range set.byMessageType {
		if s.byMessageType == nil {
			s.byMessageType = map[message.Type]map[RouteType]map[Handler]Route{}
		}

		byRouteTypeTarget, ok := s.byMessageType[mt]
		if !ok {
			byRouteTypeTarget = map[RouteType]map[Handler]Route{}
			s.byMessageType[mt] = byRouteTypeTarget
		}

		for rt, byHandlerSource := range byRouteTypeSource {
			byHandlerTarget, ok := byRouteTypeTarget[rt]
			if !ok {
				byHandlerTarget = map[Handler]Route{}
				byRouteTypeTarget[rt] = byHandlerTarget
			}

			for h, r := range byHandlerSource {
				byHandlerTarget[h] = r
			}
		}
	}
}

// RouteSetFilter applies a filter to the routes within a [RouteSet].
type RouteSetFilter func(*routeSetFilters)

// WithRouteTypeFilter is a [RouteSetFilter] that limits results to routes with
// one of the given [RouteType] values.
func WithRouteTypeFilter(types ...RouteType) RouteSetFilter {
	return func(f *routeSetFilters) {
		f.routePredicates = append(
			f.routePredicates,
			func(r Route) bool {
				return slices.Contains(types, r.RouteType.Get())
			},
		)
	}
}

// WithRouteDirectionFilter is a [RouteSetFilter] that limits results to routes
// with a [RouteDirection] that matches one of the given directions bit-masks.
func WithRouteDirectionFilter(directions ...RouteDirection) RouteSetFilter {
	return func(f *routeSetFilters) {
		f.routePredicates = append(
			f.routePredicates,
			func(r Route) bool {
				for _, dir := range directions {
					if dir.Has(r.RouteType.Get().Direction()) {
						return true
					}
				}
				return false
			},
		)
	}
}

// WithMessageKindFilter is a [RouteSetFilter] that limits results to routes
// with a [message.Kind] that matches one of the given kinds.
func WithMessageKindFilter(kinds ...message.Kind) RouteSetFilter {
	return func(f *routeSetFilters) {
		f.routePredicates = append(
			f.routePredicates,
			func(r Route) bool {
				return slices.Contains(kinds, r.MessageType.Get().Kind())
			},
		)
	}
}

// WithMessageTypeFilter is a [RouteSetFilter] that limits results to routes
// with a [message.Type] that matches one of the given types.
func WithMessageTypeFilter(kinds ...message.Kind) RouteSetFilter {
	return func(f *routeSetFilters) {
		f.routePredicates = append(
			f.routePredicates,
			func(r Route) bool {
				return slices.Contains(kinds, r.RouteType.Get().MessageKind())
			},
		)
	}
}

// WithMessageDirectionFilter is a [RouteSetFilter] that limits results to
// routes for message types that have a [RouteDirection] that matches one of the
// given directions bit-masks, when considering all routes for that message
// type.
func WithMessageDirectionFilter(directions ...RouteDirection) RouteSetFilter {
	return func(f *routeSetFilters) {
		f.messagePredicates = append(
			f.messagePredicates,
			func(t message.Type, routes map[RouteType]map[Handler]Route) bool {
				var d RouteDirection

				for rt := range routes {
					d |= rt.Direction()
				}

				for _, dir := range directions {
					if dir.Has(d) {
						return true
					}
				}

				return false
			},
		)
	}
}

type routeSetFilters struct {
	routePredicates   []func(r Route) bool
	messagePredicates []func(t message.Type, routes map[RouteType]map[Handler]Route) bool
}

func (f routeSetFilters) TestRoute(r Route) bool {
	for _, p := range f.routePredicates {
		if !p(r) {
			return false
		}
	}
	return true
}

func (f routeSetFilters) TestMessage(t message.Type, routes map[RouteType]map[Handler]Route) bool {
	for _, p := range f.messagePredicates {
		if !p(t, routes) {
			return false
		}
	}
	return true
}

func finalizeRouteSet(ctx *normalizeContext, h Handler) RouteSet {
	routes := normalizeRoutes(ctx, h)

	set := RouteSet{}

	for _, r := range routes {
		rt := r.RouteType.Get()
		mt := r.MessageType.Get()

		if set.byMessageType == nil {
			set.byMessageType = map[message.Type]map[RouteType]map[Handler]Route{}
		}

		byRouteType, ok := set.byMessageType[mt]
		if !ok {
			byRouteType = map[RouteType]map[Handler]Route{}
			set.byMessageType[mt] = byRouteType
		}

		byHandler, ok := byRouteType[rt]
		if !ok {
			byHandler = map[Handler]Route{}
			byRouteType[rt] = byHandler
		}

		byHandler[h] = r
	}

	return set
}
