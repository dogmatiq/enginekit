package runtimeconfig

import (
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

// FromRoute returns a new [config.Route] that represents the configuration of
// the given [dogma.Route].
func FromRoute(r dogma.Route) config.Route {
	cfg := config.Route{}

	if r != nil {
		switch r := r.(type) {
		case dogma.HandlesCommandRoute:
			cfg.Type = optional.Some(config.HandlesCommandRoute)
			cfg.MessageType = messageType(r.Type)
		case dogma.RecordsEventRoute:
			cfg.Type = optional.Some(config.RecordsEventRoute)
			cfg.MessageType = messageType(r.Type)
		case dogma.HandlesEventRoute:
			cfg.Type = optional.Some(config.HandlesEventRoute)
			cfg.MessageType = messageType(r.Type)
		case dogma.ExecutesCommandRoute:
			cfg.Type = optional.Some(config.ExecutesCommandRoute)
			cfg.MessageType = messageType(r.Type)
		case dogma.SchedulesTimeoutRoute:
			cfg.Type = optional.Some(config.SchedulesTimeoutRoute)
			cfg.MessageType = messageType(r.Type)
		}
	}

	return cfg
}

func messageType(r reflect.Type) optional.Optional[config.MessageType] {
	t := message.TypeFromReflect(r)

	return optional.Some(
		config.MessageType{
			TypeName: typename.Get(r),
			Kind:     t.Kind(),
			Type:     optional.Some(t),
		},
	)
}
