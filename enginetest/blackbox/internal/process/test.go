package process

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/process/end"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/process/ignored"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/process/routing"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/process/statepersistence"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/process/timeout"
)

// Run runs the process handler tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("Process", func(t *testing.T) {
		t.Run("func HandleEvent()", func(t *testing.T) {
			routing.Run(t, setup)
			ignored.Run(t, setup)
			statepersistence.Run(t, setup)
			end.Run(t, setup)
		})

		t.Run("func HandleTimeout()", func(t *testing.T) {
			timeout.Run(t, setup)
		})
	})
}
