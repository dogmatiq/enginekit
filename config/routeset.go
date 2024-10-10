package config

import (
	"iter"
	"maps"

	"github.com/dogmatiq/enginekit/message"
)

// RouteSet is the set of routes configured for a specific [Handler].
type RouteSet []Route

// MessageTypes yields all of the message types in the [RouteSet] and their
// respective [RouteDirection].
func (s RouteSet) MessageTypes() iter.Seq2[message.Type, RouteDirection] {
	types := map[message.Type]RouteDirection{}

	for _, r := range s {
		types[r.MessageType.Get()] |= r.RouteType.Get().Direction()
	}

	return maps.All(types)
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
