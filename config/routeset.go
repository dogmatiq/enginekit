package config

import (
	"github.com/dogmatiq/enginekit/message"
)

// RouteSet is the set of routes configured for a specific [Handler].
type RouteSet []Route

// MessageTypes returns a map all of the message types in the [RouteSet] and
// their respective [RouteDirection].
func (s RouteSet) MessageTypes(filter ...RouteDirection) map[message.Type]RouteDirection {
	types := map[message.Type]RouteDirection{}

	for _, r := range s {
		types[r.MessageType.Get()] |= r.RouteType.Get().Direction()
	}

	filterByDirection(types, filter)

	return types
}

// DirectionOf returns the direction in which messages of the given type flow
// for the [Handler].
func (s RouteSet) DirectionOf(t message.Type) RouteDirection {
	var dir RouteDirection

	for _, r := range s {
		if r.MessageType.Get() == t {
			dir |= r.RouteType.Get().Direction()
		}
	}

	return dir
}

func filterByDirection(types map[message.Type]RouteDirection, filter []RouteDirection) {
	if len(filter) == 0 {
		return
	}

	match := func(d RouteDirection) bool {
		for _, f := range filter {
			if d.Has(f) {
				return true
			}
		}
		return false
	}

	for t, d := range types {
		if !match(d) {
			delete(types, t)
		}
	}
}
