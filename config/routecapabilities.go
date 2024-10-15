package config

import "github.com/dogmatiq/enginekit/message"

// RouteTypeCapability is describes how a [HandlerType] makes use of a
// particular [RouteType].
type RouteTypeCapability int

const (
	// RouteTypeDisallowed indicates that the [HandlerType] does not support the
	// [RouteType].
	RouteTypeDisallowed RouteTypeCapability = iota

	// RouteTypeAllowed indicates that the [HandlerType] supports the [RouteType],
	// but it is not required.
	RouteTypeAllowed

	// RouteTypeRequired indicates that the [HandlerType] requires at least one
	// route of the [RouteType].
	RouteTypeRequired
)

// RouteCapabilities is a map that describes how a [HandlerType] makes use of
// each [RouteType].
type RouteCapabilities struct {
	RouteTypes map[RouteType]RouteTypeCapability
}

// DirectionOf returns the direction in which messages of the given kind flow
// for the [HandlerType].
func (s RouteCapabilities) DirectionOf(k message.Kind) RouteDirection {
	var dir RouteDirection

	for rt, cap := range s.RouteTypes {
		if cap != RouteTypeDisallowed {
			if rt.MessageKind() == k {
				dir |= rt.Direction()
			}
		}
	}

	return dir
}
