package routing

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// app routes [command] to its [Handler].
type app struct {
	Handler handler
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "6c2f7cd6-9d1c-4ef0-a87c-210e21745d0b")
	c.Routes(dogma.ViaIntegration(&a.Handler))
}

// handler handles [command] and signals [Called] on each invocation.
type handler struct {
	Called chan struct{}
}

func (*handler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "4c6bddfb-c821-47c2-81d3-5c7611020786")
	c.Routes(dogma.HandlesCommand[*command]())
}

func (h *handler) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error {
	h.Called <- struct{}{}
	return nil
}
