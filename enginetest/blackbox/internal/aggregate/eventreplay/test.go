package eventreplay

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the aggregate event replay and state persistence tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("a second command observes state from the first command's events", func(t *testing.T) {
		// handler.Got receives the aggregate root's Value when checkCommand is
		// handled. The root's Value is set by ApplyEvent when the valueWritten
		// event is replayed. If the engine replays the event correctly, the
		// second command sees the value written by the first.
		//
		// The root type implements real MarshalBinary / UnmarshalBinary, so
		// the engine's snapshot path (if any) is exercised. The test does not
		// assert that snapshotting occurred — only that the final state is
		// correct.

		const want = "hello-from-first-command"

		a := &app{Handler: handler{Got: make(chan string, 1)}}
		x := setup(t, a)

		if err := x.ExecuteCommand(t.Context(), &writeCommand{Value: want}); err != nil {
			t.Fatal(err)
		}

		if err := x.ExecuteCommand(t.Context(), &checkCommand{}); err != nil {
			t.Fatal(err)
		}

		select {
		case got := <-a.Handler.Got:
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for handler to report the aggregate state")
		}
	})
}
