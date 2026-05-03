package routing

import (
	"testing"
	"time"

	"github.com/dogmatiq/dogma"
)

// Run runs the integration routing tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("it routes the command to the handler", func(t *testing.T) {
		a := &app{Handler: handler{Called: make(chan struct{}, 1)}}
		x := setup(t, a)

		err := x.ExecuteCommand(t.Context(), &command{})
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-a.Handler.Called:
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for handler to be called")
		}
	})
}
