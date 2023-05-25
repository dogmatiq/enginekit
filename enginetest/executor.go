package enginetest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/enginekit/enginetest/internal/testapp"
)

func testCommandExecutor(ctx context.Context, t *testing.T, e *engine) {
	t.Run("command executor", func(t *testing.T) {
		t.Parallel()

		t.Run("panics if passed an invalid command", func(t *testing.T) {
			t.Parallel()

			defer func() {
				actual := fmt.Sprint(recover())
				expect := testapp.ErrInvalidIntegrationMessage.Error()

				if !strings.Contains(actual, expect) {
					t.Fatalf("got: %q, want: panic containing %q", actual, expect)
				}
			}()

			e.Executor.ExecuteCommand(
				ctx,
				&testapp.IntegrationCommandA{
					IsInvalid: true,
				},
			)
		})

		t.Run("panics if passed an unrecognized command", func(t *testing.T) {
			t.Parallel()

			defer func() {
				actual := fmt.Sprint(recover())
				expect := "MessageC" // the type name should appear in the message

				if !strings.Contains(actual, expect) {
					t.Fatalf("got: %q, want: panic containing %q", actual, expect)
				}
			}()

			e.Executor.ExecuteCommand(
				ctx,
				fixtures.MessageC1,
			)
		})
	})
}
