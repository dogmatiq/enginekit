package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// FromProjection returns a new [config.Projection] that represents the
// configuration of the given [dogma.ProjectionMessageHandler].
func FromProjection(h dogma.ProjectionMessageHandler) config.Projection {
	var cfg config.Projection

	if h == nil {
		return cfg
	}

	cfg.TypeName = optional.Some(typename.Of(h))
	cfg.Implementation = optional.Some(h)

	return cfg
}
