package statepersistence

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the process state persistence tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("the second invocation observes state from the first invocation", func(t *testing.T) {
		// firstCommand → firstIntegration records firstEvent →
		// process handles firstEvent, sets root.Value = "persisted-value".
		// secondCommand → secondIntegration records secondEvent →
		// process handles secondEvent, reads root.Value, sends to Got.
		//
		// The root implements real MarshalBinary / UnmarshalBinary. If the
		// engine persists state between invocations (via snapshots, event
		// replay, or any other mechanism), the second invocation must see the
		// value set by the first.

		const want = "persisted-value"

		a := &app{Handler: handler{Got: make(chan string, 1)}}
		x := setup(t, a)

		if err := x.ExecuteCommand(t.Context(), &firstCommand{}); err != nil {
			t.Fatal(err)
		}

		if err := x.ExecuteCommand(t.Context(), &secondCommand{}); err != nil {
			t.Fatal(err)
		}

		select {
		case got := <-a.Handler.Got:
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for process to report persisted state")
		}
	})
}
