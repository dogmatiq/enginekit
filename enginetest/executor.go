package enginetest

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/enginetest/internal/testapp"
	"github.com/dogmatiq/enginekit/enginetest/stubs"
)

func testCommandExecutor(ctx context.Context, t *testing.T, e *engine) {
	t.Run("command executor", func(t *testing.T) {
		t.Parallel()

		type UnrecognizedCommand struct {
			stubs.CommandStub[string]
		}

		cases := []struct {
			Name    string
			Command dogma.Command
			Expect  string
		}{
			{
				Name:    "panics if passed an invalid command",
				Command: &testapp.IntegrationCommandA{IsInvalid: true},
				Expect:  testapp.ErrInvalidIntegrationMessage.Error(),
			},
			{
				Name:    "panics if passed an unrecognized command",
				Command: &UnrecognizedCommand{},
				Expect:  "UnrecognizedCommand",
			},
			{
				Name:    "panics if passed a nil command",
				Command: nil,
				Expect:  "nil",
			},
		}

		for _, c := range cases {
			c := c // capture loop variable

			t.Run(c.Name, func(t *testing.T) {
				t.Parallel()

				defer func() {
					actual := fmt.Sprint(recover())

					if !strings.Contains(actual, c.Expect) {
						t.Fatalf("got: %q, want: panic containing %q", actual, c.Expect)
					}
				}()

				e.Executor.ExecuteCommand(ctx, c.Command)
			})
		}
	})
}
