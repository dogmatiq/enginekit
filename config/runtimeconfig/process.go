package runtimeconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// FromProcess returns a new [config.Process] that represents the configuration
// of the given [dogma.ProcessMessageHandler].
func FromProcess(h dogma.ProcessMessageHandler) *config.Process {
	return configbuilder.Process(
		func(b *configbuilder.ProcessBuilder) {
			buildProcess(b, h)
		},
	)
}

func buildProcess(b *configbuilder.ProcessBuilder, h dogma.ProcessMessageHandler) {
	if h == nil {
		b.Partial()
	} else {
		c := newHandlerConfigurer[dogma.ProcessRoute](b)
		b.Source(h)
		h.Configure(c)
	}
}
