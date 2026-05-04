package observer

import (
	"context"
	"errors"
	"testing"

	"github.com/dogmatiq/dogma"
)

// Run runs the WithEventObserver tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("WithEventObserver", func(t *testing.T) {
		t.Run("it blocks until the observer returns satisfied", func(t *testing.T) {
			// chainApp drives a three-handler causal chain:
			//   1. chainAggregate handles command, records triggered
			//   2. chainProcess handles triggered, executes relayed
			//   3. chainIntegration handles relayed, records observed
			//   4. Observer watches for observed (end of chain)
			x := setup(t, &chainApp{})

			saw := false
			err := x.ExecuteCommand(
				t.Context(),
				&command{},
				dogma.WithEventObserver(
					func(context.Context, *observed) (bool, error) {
						saw = true
						return true, nil
					},
				),
			)
			if err != nil {
				t.Fatal(err)
			}

			if !saw {
				t.Fatal("expected observer to be called")
			}
		})

		t.Run("when no further events can occur", func(t *testing.T) {
			t.Run("it returns ErrEventObserverNotSatisfied", func(t *testing.T) {
				// silentApp handles command but records no events, so the
				// observer waiting for triggered can never be satisfied.
				x := setup(t, &silentApp{})

				err := x.ExecuteCommand(
					t.Context(),
					&command{},
					dogma.WithEventObserver(
						func(context.Context, *triggered) (bool, error) {
							return false, nil
						},
					),
				)
				if !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
					t.Fatalf("got %v, want ErrEventObserverNotSatisfied", err)
				}
			})
		})
	})
}
