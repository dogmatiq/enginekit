package noroute

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// app routes [command] to [handler]. It has no route for [unroutedCommand].
type app struct{}

func (*app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "06dec7a5-a4df-44da-bda7-20afea15f1d3")
	c.Routes(dogma.ViaIntegration(&handler{}))
}

// handler handles [command] without recording any events.
type handler struct{}

func (*handler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "f2ffb4cb-53a2-4673-90ab-fa150901000c")
	c.Routes(dogma.HandlesCommand[*command]())
}

func (*handler) HandleCommand(context.Context, dogma.IntegrationCommandScope, dogma.Command) error {
	return nil
}
