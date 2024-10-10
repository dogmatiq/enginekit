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
func fromRoute(r dogma.Route) config.Route {
	cfg := config.Route{}

	if r != nil {
		switch r := r.(type) {
		case dogma.HandlesCommandRoute:
			setupRoute(&cfg, config.HandlesCommandRoute, r.Type)
		case dogma.RecordsEventRoute:
			setupRoute(&cfg, config.RecordsEventRoute, r.Type)
		case dogma.HandlesEventRoute:
			setupRoute(&cfg, config.HandlesEventRoute, r.Type)
		case dogma.ExecutesCommandRoute:
			setupRoute(&cfg, config.ExecutesCommandRoute, r.Type)
		case dogma.SchedulesTimeoutRoute:
			setupRoute(&cfg, config.SchedulesTimeoutRoute, r.Type)
		}
	}

	return cfg
}

func setupRoute(
	cfg *config.Route,
	rt config.RouteType,
	t reflect.Type,
) {
	cfg.RouteType = optional.Some(rt)
	cfg.MessageTypeName = optional.Some(typename.Get(t))
	cfg.MessageType = optional.Some(message.TypeFromReflect(t))
}
