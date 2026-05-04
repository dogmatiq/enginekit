package ignored

import (
	"context"
	"sync/atomic"

	"github.com/dogmatiq/dogma"
)

// app routes [triggerCommand] through [triggerIntegration] to produce
// [processTrigger], which [process] is registered to handle but ignores.
type app struct {
	Handler handler
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "b2a00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&triggerIntegration{}),
		dogma.ViaProcess(&a.Handler),
	)
}

// triggerIntegration handles [triggerCommand] and records [processTrigger].
type triggerIntegration struct{}

func (*triggerIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("trigger", "b2b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*triggerCommand](),
		dogma.RecordsEvent[*processTrigger](),
	)
}

func (*triggerIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&processTrigger{})
	return nil
}

// handler is a process that always ignores [processTrigger] events by
// returning ok == false from RouteEventToInstance. HandleEvent must never be
// called.
type handler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutMessagesBehavior[dogma.StatelessProcessRoot]
	Calls atomic.Int32
}

func (*handler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("process", "b2c00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesEvent[*processTrigger](),
		dogma.ExecutesCommand[*triggerCommand](),
	)
}

func (*handler) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "", false, nil
}

func (h *handler) HandleEvent(context.Context, dogma.StatelessProcessRoot, dogma.ProcessEventScope[dogma.StatelessProcessRoot], dogma.Event) error {
	h.Calls.Add(1)
	return nil
}
