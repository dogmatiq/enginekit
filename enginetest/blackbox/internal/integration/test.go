package integration

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/integration/atomicity"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/integration/routing"
)

// Run runs the Integration handler tests against the engine provided by
// setup.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("Integration", func(t *testing.T) {
		t.Run("func HandleCommand()", func(t *testing.T) {
			routing.Run(t, setup)
			atomicity.Run(t, setup)
		})
	})
}
