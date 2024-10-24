package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromProcess returns a new [config.Process] that represents the configuration
// of the given [dogma.ProcessMessageHandler].
func FromProcess(h dogma.ProcessMessageHandler) *config.Process {
	return configbuilder.Process(func(b *configbuilder.ProcessBuilder) {
		if h == nil {
			b.UpdateFidelity(config.Incomplete)
		} else {
			buildProcess(b, h)
		}
	})
}

func buildProcess(b *configbuilder.ProcessBuilder, h dogma.ProcessMessageHandler) {
	b.SetSource(h)
	h.Configure(&handlerConfigurer[dogma.ProcessRoute]{b})
}
