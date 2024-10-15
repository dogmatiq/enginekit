package runtimeconfig

import (
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// fromRoute returns a new [config.Route] that represents the configuration of
// the given [dogma.Route].
func fromRoute(r dogma.Route) *config.Route {
	cfg := &config.Route{}

	configure := func(
		rt config.RouteType,
		t reflect.Type,
	) {
		cfg.AsConfigured = config.RouteAsConfigured{
			RouteType:       optional.Some(rt),
			MessageTypeName: optional.Some(typename.Get(t)),
			MessageType:     optional.Some(message.TypeFromReflect(t)),
		}
	}

	if r != nil {
		switch r := r.(type) {
		case dogma.HandlesCommandRoute:
			configure(config.HandlesCommandRouteType, r.Type)
		case dogma.RecordsEventRoute:
			configure(config.RecordsEventRouteType, r.Type)
		case dogma.HandlesEventRoute:
			configure(config.HandlesEventRouteType, r.Type)
		case dogma.ExecutesCommandRoute:
			configure(config.ExecutesCommandRouteType, r.Type)
		case dogma.SchedulesTimeoutRoute:
			configure(config.SchedulesTimeoutRouteType, r.Type)
		}
	}

	return cfg
}
