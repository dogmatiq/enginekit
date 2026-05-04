package projection

import (
	"testing"

	"github.com/dogmatiq/dogma"
	projrouting "github.com/dogmatiq/enginekit/enginetest/blackbox/internal/projection/routing"
)

// Run runs the projection handler tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("Projection", func(t *testing.T) {
		t.Run("func HandleEvent()", func(t *testing.T) {
			projrouting.Run(t, setup)
		})
	})
}
