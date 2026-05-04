package routing

import (
	"context"
	"testing"

	"github.com/dogmatiq/dogma"
)

// Run runs the aggregate routing tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("it routes the command to the handler", func(t *testing.T) {
		// Execute a command and use WithEventObserver to confirm the handler
		// was invoked. The handler records an event on every HandleCommand call;
		// if the observer fires, the command reached the right handler.
		x := setup(t, &app{})

		saw := false
		err := x.ExecuteCommand(
			t.Context(),
			&command{},
			dogma.WithEventObserver(
				func(context.Context, *recorded) (bool, error) {
					saw = true
					return true, nil
				},
			),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !saw {
			t.Fatal("expected recorded event to be observed")
		}
	})
}
