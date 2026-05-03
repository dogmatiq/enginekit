package atomicity

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/dogmatiq/dogma"
)

// Run runs the integration scope atomicity tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("events recorded within a single scope are persisted atomically", func(t *testing.T) {
		// handler records two events in a single HandleCommand call. We observe
		// both via WithEventObserver to verify that all events from the scope
		// are persisted.
		//
		// NOTE: This tests the "all" direction of atomicity. The "none"
		// direction (rollback on handler error) is not testable at this
		// abstraction level — the engine may or may not expose failed scopes
		// externally.

		var sawStarted, sawCompleted atomic.Bool

		x := setup(t, &app{})

		err := x.ExecuteCommand(
			t.Context(),
			&command{},
			dogma.WithEventObserver(
				func(context.Context, *started) (bool, error) {
					sawStarted.Store(true)
					return false, nil
				},
			),
			dogma.WithEventObserver(
				func(context.Context, *completed) (bool, error) {
					sawCompleted.Store(true)
					return true, nil
				},
			),
		)
		if err != nil {
			t.Fatal(err)
		}

		if !sawStarted.Load() {
			t.Fatal("started event was not observed")
		}
		if !sawCompleted.Load() {
			t.Fatal("completed event was not observed")
		}
	})
}
