package config

import "github.com/dogmatiq/enginekit/message"

// RouteSet is the complete set of valid and complete [Route] components for a
// specific [Entity].
type RouteSet struct {
	byMessageType map[message.Type]map[RouteType]map[Handler]*Route
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

func (s *RouteSet) merge(x RouteSet) {
	for mt, byRouteTypeSource := range x.byMessageType {
		if s.byMessageType == nil {
			s.byMessageType = map[message.Type]map[RouteType]map[Handler]*Route{}
		}

		byRouteTypeTarget, ok := s.byMessageType[mt]
		if !ok {
			byRouteTypeTarget = map[RouteType]map[Handler]*Route{}
			s.byMessageType[mt] = byRouteTypeTarget
		}

		for rt, byHandlerSource := range byRouteTypeSource {
			byHandlerTarget, ok := byRouteTypeTarget[rt]
			if !ok {
				byHandlerTarget = map[Handler]*Route{}
				byRouteTypeTarget[rt] = byHandlerTarget
			}

			for h, r := range byHandlerSource {
				byHandlerTarget[h] = r
			}
		}
	}
}
