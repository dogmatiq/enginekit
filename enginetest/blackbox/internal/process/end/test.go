package end

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the process End tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("after End", func(t *testing.T) {
		t.Run("calling ExecuteCommand panics", func(t *testing.T) {
			// The process calls End() then ExecuteCommand(). The spec
			// requires ExecuteCommand to panic. The handler recovers the panic
			// and signals Panicked. The test asserts the panic occurred.

			a := &panicApp{Handler: panicHandler{Panicked: make(chan bool, 1)}}
			x := setup(t, a)

			if err := x.ExecuteCommand(t.Context(), &endCommand{}); err != nil {
				t.Fatal(err)
			}

			select {
			case panicked := <-a.Handler.Panicked:
				if !panicked {
					t.Fatal("expected ExecuteCommand to panic after End, but it did not")
				}
			case <-time.After(5 * time.Second):
				t.Fatal("timed out waiting for handler to report panic result")
			}
		})

		t.Run("future events targeting the same instance are ignored", func(t *testing.T) {
			// The process calls End() on first invocation. When the same
			// instance receives a second event, the engine must ignore it.
			// The handler's Calls counter must equal 1 after two dispatches.

			a := &replayApp{}
			x := setup(t, a)

			if err := x.ExecuteCommand(t.Context(), &replayCommand{}); err != nil {
				t.Fatal(err)
			}
			if err := x.ExecuteCommand(t.Context(), &replayCommand{}); err != nil {
				t.Fatal(err)
			}

			// Allow time for any asynchronous processing.
			time.Sleep(100 * time.Millisecond)

			if n := a.Handler.Calls.Load(); n != 1 {
				t.Fatalf("HandleEvent was called %d time(s), want 1", n)
			}
		})
	})
}
