package timeout

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
)

// app drives the timeout test:
// [startIntegration] records [workflowStarted] →
// [process] schedules [workflowTimeout] at time.Now() →
// engine tick delivers [workflowTimeout] →
// [process] handles [workflowTimeout] and signals [handler.Done].
type app struct {
	Handler handler
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "b3a00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&startIntegration{}),
		dogma.ViaProcess(&a.Handler),
	)
}

// startIntegration handles [startCommand] and records [workflowStarted].
type startIntegration struct{}

func (*startIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("start", "b3b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*startCommand](),
		dogma.RecordsEvent[*workflowStarted](),
	)
}

func (*startIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&workflowStarted{})
	return nil
}

// handler is the process under test. It handles [workflowStarted] by
// scheduling an immediate [workflowTimeout], then handles the timeout by
// sending to [Done].
type handler struct {
	Done chan struct{}
}

type root struct{}

func (*root) ProcessInstanceDescription(bool) string { return "" }
func (*root) MarshalBinary() ([]byte, error)         { return nil, nil }
func (*root) UnmarshalBinary([]byte) error           { return nil }

func (*handler) New() *root { return &root{} }

func (*handler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("process", "b3c00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesEvent[*workflowStarted](),
		dogma.SchedulesTimeout[*workflowTimeout](),
		dogma.ExecutesCommand[*startCommand](),
	)
}

func (*handler) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "instance", true, nil
}

func (*handler) HandleEvent(_ context.Context, _ *root, s dogma.ProcessEventScope[*root], _ dogma.Event) error {
	s.ScheduleTimeout(&workflowTimeout{}, time.Now())
	return nil
}

func (h *handler) HandleTimeout(_ context.Context, _ *root, _ dogma.ProcessTimeoutScope[*root], _ dogma.Timeout) error {
	h.Done <- struct{}{}
	return nil
}
