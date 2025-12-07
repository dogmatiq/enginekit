package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromAggregate returns a new [config.Aggregate] that represents the
// configuration of the given [dogma.AggregateMessageHandler].
func FromAggregate(h dogma.AggregateMessageHandler) *config.Aggregate {
	return configbuilder.Aggregate(
		func(b *configbuilder.AggregateBuilder) {
			buildAggregate(b, h)
		},
	)
}

func buildAggregate(b *configbuilder.AggregateBuilder, h dogma.AggregateMessageHandler) {
	if h == nil {
		b.Partial()
	} else {
		c := newHandlerConfigurer[dogma.AggregateRoute](b)
		b.Source(h)
		h.Configure(c)
	}
}
