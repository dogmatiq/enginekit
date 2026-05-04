package timeout

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the process timeout tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("it delivers the timeout at the scheduled time", func(t *testing.T) {
		// startCommand → startIntegration records workflowStarted →
		// process handles workflowStarted, schedules workflowTimeout at time.Now() →
		// engine delivers workflowTimeout →
		// process handles workflowTimeout, signals Done.
		//
		// The engine must tick to deliver the timeout. The test setup is
		// expected to run the engine in the background so that timeouts are
		// delivered asynchronously.

		a := &app{Handler: handler{Done: make(chan struct{}, 1)}}
		x := setup(t, a)

		if err := x.ExecuteCommand(t.Context(), &startCommand{}); err != nil {
			t.Fatal(err)
		}

		select {
		case <-a.Handler.Done:
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for timeout to be delivered")
		}
	})
}
