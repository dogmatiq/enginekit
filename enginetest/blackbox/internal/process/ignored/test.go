package ignored

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the process event-ignored tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("when RouteEventToInstance returns ok == false", func(t *testing.T) {
		t.Run("it ignores the event", func(t *testing.T) {
			// The process's RouteEventToInstance always returns ok=false.
			// HandleEvent must never be called.

			a := &app{}
			x := setup(t, a)

			if err := x.ExecuteCommand(t.Context(), &triggerCommand{}); err != nil {
				t.Fatal(err)
			}

			// Allow time for any asynchronous processing.
			time.Sleep(100 * time.Millisecond)

			if n := a.Handler.Calls.Load(); n != 0 {
				t.Fatalf("HandleEvent was called %d time(s), want 0", n)
			}
		})
	})
}
