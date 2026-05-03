package idempotency

import (
	"context"
	"sync/atomic"

	"github.com/dogmatiq/dogma"
)

// app routes [command] to its [Handler].
type app struct {
	Handler handler
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "3bd30f3e-945f-4bf4-bdf2-c5a556cc0d96")
	c.Routes(dogma.ViaIntegration(&a.Handler))
}

// handler handles [command], records [handled], and counts invocations so the
// test can verify that deduplicated calls do not re-invoke the handler.
type handler struct {
	Calls atomic.Int32
}

func (*handler) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("handler", "a06b01b1-8c06-451a-8679-284e8bad111f")
	c.Routes(
		dogma.HandlesCommand[*command](),
		dogma.RecordsEvent[*handled](),
	)
}

func (h *handler) HandleCommand(
	_ context.Context,
	s dogma.IntegrationCommandScope,
	_ dogma.Command,
) error {
	h.Calls.Add(1)
	s.RecordEvent(&handled{})
	return nil
}
