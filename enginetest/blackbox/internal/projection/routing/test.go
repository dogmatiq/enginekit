package routing

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the projection event routing tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("it delivers the event to the handler", func(t *testing.T) {
		// Execute a command that causes an integration to record an event.
		// The projection handles that event and signals Received. If
		// Received is never closed, the engine did not route the event to
		// the projection handler.
		a := &app{
			Projection: projection{Received: make(chan struct{}, 1)},
		}
		x := setup(t, a)

		if err := x.ExecuteCommand(t.Context(), &triggerCommand{}); err != nil {
			t.Fatal(err)
		}

		select {
		case <-a.Projection.Received:
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for projection to receive the event")
		}
	})
}
