package end

import (
	"context"
	"sync/atomic"

	"github.com/dogmatiq/dogma"
)

// --- "panics after End" app ---

// panicApp tests that calling ExecuteCommand after End panics.
type panicApp struct {
	Handler panicHandler
}

func (a *panicApp) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("panic-app", "b4a00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&endIntegration{}),
		dogma.ViaProcess(&a.Handler),
	)
}

// endIntegration handles [endCommand] and records [endTrigger].
type endIntegration struct{}

func (*endIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("end-integration", "b4b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*endCommand](),
		dogma.RecordsEvent[*endTrigger](),
	)
}

func (*endIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&endTrigger{})
	return nil
}

// panicHandler is a process that calls End() then ExecuteCommand(), expecting
// the latter to panic. The panic is recovered internally; Panicked is set to
// true if the expected panic occurred.
type panicHandler struct {
	Panicked chan bool
}

type panicRoot struct{}

func (*panicRoot) ProcessInstanceDescription(bool) string { return "" }
func (*panicRoot) MarshalBinary() ([]byte, error)         { return nil, nil }
func (*panicRoot) UnmarshalBinary([]byte) error           { return nil }

func (*panicHandler) New() *panicRoot { return &panicRoot{} }

func (*panicHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("panic-process", "b4c00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesEvent[*endTrigger](),
		dogma.ExecutesCommand[*continueCommand](),
	)
}

func (*panicHandler) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "instance", true, nil
}

func (h *panicHandler) HandleEvent(_ context.Context, _ *panicRoot, s dogma.ProcessEventScope[*panicRoot], _ dogma.Event) (err error) {
	panicked := false

	defer func() {
		h.Panicked <- panicked
	}()

	s.End()

	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()

	s.ExecuteCommand(&continueCommand{})

	return nil
}

func (*panicHandler) HandleTimeout(context.Context, *panicRoot, dogma.ProcessTimeoutScope[*panicRoot], dogma.Timeout) error {
	panic(dogma.UnexpectedMessage)
}

// --- "future events ignored after End" app ---

// replayApp tests that events targeting an ended process instance are ignored.
type replayApp struct {
	Handler replayHandler
}

func (a *replayApp) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("replay-app", "b4d00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&replayIntegration{}),
		dogma.ViaProcess(&a.Handler),
	)
}

// replayIntegration handles [replayCommand] and records [replayTrigger].
type replayIntegration struct{}

func (*replayIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("replay-integration", "b4e00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*replayCommand](),
		dogma.RecordsEvent[*replayTrigger](),
	)
}

func (*replayIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&replayTrigger{})
	return nil
}

// replayHandler is a process that calls End() on first invocation. Its Calls
// counter must equal 1 after two events targeting the same instance.
type replayHandler struct {
	dogma.NoTimeoutMessagesBehavior[*replayRoot]
	Calls atomic.Int32
}

type replayRoot struct{}

func (*replayRoot) ProcessInstanceDescription(bool) string { return "" }
func (*replayRoot) MarshalBinary() ([]byte, error)         { return nil, nil }
func (*replayRoot) UnmarshalBinary([]byte) error           { return nil }

func (*replayHandler) New() *replayRoot { return &replayRoot{} }

func (*replayHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("replay-process", "b4f00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesEvent[*replayTrigger](),
		dogma.ExecutesCommand[*replayCommand](),
	)
}

func (*replayHandler) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "instance", true, nil
}

func (h *replayHandler) HandleEvent(_ context.Context, _ *replayRoot, s dogma.ProcessEventScope[*replayRoot], _ dogma.Event) error {
	h.Calls.Add(1)
	s.End()
	return nil
}
