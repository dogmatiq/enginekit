package blackbox

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/aggregate"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/commandexecutor"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/integration"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/process"
	"github.com/dogmatiq/enginekit/enginetest/blackbox/internal/projection"
)

// Run runs the engine acceptance test suite.
//
// setup is called once per scenario. It starts the engine with the provided
// application, registers t.Cleanup to stop it, and returns a
// [dogma.CommandExecutor] ready for use.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()
	commandexecutor.Run(t, setup)
	integration.Run(t, setup)
	aggregate.Run(t, setup)
	process.Run(t, setup)
	projection.Run(t, setup)
}
