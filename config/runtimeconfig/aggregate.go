package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromAggregate returns a new [config.Aggregate] that represents the
// configuration of the given [dogma.AggregateMessageHandler].
func FromAggregate(h dogma.AggregateMessageHandler) *config.Aggregate {
	return configbuilder.Aggregate(func(b *configbuilder.AggregateBuilder) {
		if h == nil {
			b.UpdateFidelity(config.Incomplete)
		} else {
			buildAggregate(b, h)
		}
	})
}

func buildAggregate(b *configbuilder.AggregateBuilder, h dogma.AggregateMessageHandler) {
	b.SetSource(h)
	b.SetDisabled(false)
	h.Configure(&handlerConfigurer[dogma.AggregateRoute]{b})
}
