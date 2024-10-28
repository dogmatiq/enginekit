package config

import (
	"cmp"

	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// Route represents the configuration of a [dogma.Route].
type Route struct {
	ComponentCommon

	// RouteType is the type of route, if available.
	RouteType optional.Optional[RouteType]

	// MessageTypeName is the fully-qualified name of the Go type that
	// implements the [dogma.Message] interface, if available.
	MessageTypeName optional.Optional[string]

	// MessageType is the [message.Type], if available.
	MessageType optional.Optional[message.Type]
}

func (r *Route) String() string {
	panic("not implemented")
}

func (r *Route) key() (routeKey, bool) {
	if !r.RouteType.IsPresent() {
		return routeKey{}, false
	}

	if !r.MessageTypeName.IsPresent() {
		return routeKey{}, false
	}

	return routeKey{
		RouteType:       r.RouteType.Get(),
		MessageTypeName: r.MessageTypeName.Get(),
	}, true
}

// routeKey is a [comparable] representation of a route's identity. No [Handler]
// may have two routes with the same key.
type routeKey struct {
	RouteType       RouteType
	MessageTypeName string
}

func (k routeKey) Compare(x routeKey) int {
	if c := cmp.Compare(k.RouteType, x.RouteType); c != 0 {
		return c
	}
	return cmp.Compare(k.MessageTypeName, x.MessageTypeName)
}
