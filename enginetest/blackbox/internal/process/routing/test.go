package routing

import (
	"context"
	"testing"

	"github.com/dogmatiq/dogma"
)

// Run runs the process event routing tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("it routes the event to the handler", func(t *testing.T) {
		// startCommand → startIntegration records workflowStarted →
		// process handles workflowStarted, executes actionCommand →
		// completionIntegration records workflowCompleted.
		x := setup(t, &app{})

		saw := false
		err := x.ExecuteCommand(
			t.Context(),
			&startCommand{},
			dogma.WithEventObserver(
				func(context.Context, *workflowCompleted) (bool, error) {
					saw = true
					return true, nil
				},
			),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !saw {
			t.Fatal("expected workflowCompleted event to be observed")
		}
	})
}
