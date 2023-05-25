package enginetest

import (
	"context"
	"testing"

	"github.com/dogmatiq/enginekit/enginetest/internal/action"
	"github.com/dogmatiq/enginekit/enginetest/internal/testapp"
	"github.com/google/uuid"
)

func testIntegration(ctx context.Context, t *testing.T, e *engine) {
	t.Run("integration message handler", func(t *testing.T) {
		t.Parallel()

		t.Run("it can record events", func(t *testing.T) {
			t.Parallel()

			id := uuid.NewString()
			expect := &testapp.IntegrationEventA{
				Value: id,
			}

			e.ExecuteCommand(
				t,
				&testapp.IntegrationCommandA{
					Actions: action.RecordEvent(expect),
				},
			)

			e.ExpectEvent(t, expect)
		})
	})
}
