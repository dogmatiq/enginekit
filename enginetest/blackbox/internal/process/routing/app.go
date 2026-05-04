package routing

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// app drives a three-handler causal chain:
// [startIntegration] records [workflowStarted] →
// [process] executes [actionCommand] →
// [completionIntegration] records [workflowCompleted].
type app struct{}

func (*app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "b1a00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&startIntegration{}),
		dogma.ViaProcess(&process{}),
		dogma.ViaIntegration(&completionIntegration{}),
	)
}

// startIntegration handles [startCommand] and records [workflowStarted].
type startIntegration struct{}

func (*startIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("start", "b1b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*startCommand](),
		dogma.RecordsEvent[*workflowStarted](),
	)
}

func (*startIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&workflowStarted{})
	return nil
}

// process handles [workflowStarted] and executes [actionCommand].
type process struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutMessagesBehavior[dogma.StatelessProcessRoot]
}

func (*process) Configure(c dogma.ProcessConfigurer) {
	c.Identity("process", "b1c00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesEvent[*workflowStarted](),
		dogma.ExecutesCommand[*actionCommand](),
	)
}

func (*process) RouteEventToInstance(context.Context, dogma.Event) (string, bool, error) {
	return "instance", true, nil
}

func (*process) HandleEvent(_ context.Context, _ dogma.StatelessProcessRoot, s dogma.ProcessEventScope[dogma.StatelessProcessRoot], _ dogma.Event) error {
	s.ExecuteCommand(&actionCommand{})
	return nil
}

// completionIntegration handles [actionCommand] and records [workflowCompleted].
type completionIntegration struct{}

func (*completionIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("completion", "b1d00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*actionCommand](),
		dogma.RecordsEvent[*workflowCompleted](),
	)
}

func (*completionIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&workflowCompleted{})
	return nil
}
