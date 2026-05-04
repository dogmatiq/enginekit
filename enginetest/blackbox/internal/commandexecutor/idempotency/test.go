package idempotency

import (
	"context"
	"errors"
	"testing"

	"github.com/dogmatiq/dogma"
)

// Run runs the WithIdempotencyKey tests.
func Run(
	t *testing.T,
	setup func(t *testing.T, app dogma.Application) dogma.CommandExecutor,
) {
	t.Helper()

	t.Run("WithIdempotencyKey", func(t *testing.T) {
		t.Run("it deduplicates commands", func(t *testing.T) {
			a := &app{}
			x := setup(t, a)

			// First call: the handler is invoked and the observer is satisfied.
			err := x.ExecuteCommand(
				t.Context(),
				&command{},
				dogma.WithIdempotencyKey("dedup-key"),
				dogma.WithEventObserver(
					func(context.Context, *handled) (bool, error) {
						return true, nil
					},
				),
			)
			if err != nil {
				t.Fatalf("first call: %v", err)
			}

			// Second call with the same key: the command is deduplicated, so no
			// events are produced and the observer cannot be satisfied.
			err = x.ExecuteCommand(
				t.Context(),
				&command{},
				dogma.WithIdempotencyKey("dedup-key"),
				dogma.WithEventObserver(
					func(context.Context, *handled) (bool, error) {
						return true, nil
					},
				),
			)
			if !errors.Is(err, dogma.ErrEventObserverNotSatisfied) {
				t.Fatalf("got %v, want ErrEventObserverNotSatisfied on second call", err)
			}

			if n := a.Handler.Calls.Load(); n != 1 {
				t.Fatalf("got %d handler invocations, want 1", n)
			}
		})
	})
}
