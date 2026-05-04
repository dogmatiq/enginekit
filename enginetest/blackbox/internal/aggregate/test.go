package aggregate

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/aggregate/atomicity"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/aggregate/eventreplay"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/aggregate/routing"
)

// Run runs the aggregate handler tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("Aggregate", func(t *testing.T) {
		t.Run("func HandleCommand()", func(t *testing.T) {
			routing.Run(t, setup)
			eventreplay.Run(t, setup)
			atomicity.Run(t, setup)
		})
	})
}
