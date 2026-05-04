package noroute

import (
	"testing"

	"github.com/dogmatiq/dogma"
)

// Run runs the no-route tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("when the command has no route", func(t *testing.T) {
		t.Run("it panics", func(t *testing.T) {
			x := setup(t, &app{})

			// unroutedCommand has no route in app. The engine must panic —
			// this is a programming error on the caller's part.
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic when executing a command with no route")
				}
			}()

			err := x.ExecuteCommand(t.Context(), &unroutedCommand{})
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
