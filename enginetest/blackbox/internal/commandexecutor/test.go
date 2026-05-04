package commandexecutor

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/commandexecutor/idempotency"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/commandexecutor/noroute"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/commandexecutor/observer"
)

// Run runs the CommandExecutor tests against the engine provided by
// setup.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("CommandExecutor", func(t *testing.T) {
		t.Run("func ExecuteCommand()", func(t *testing.T) {
			noroute.Run(t, setup)
			observer.Run(t, setup)
			idempotency.Run(t, setup)
		})
	})
}
