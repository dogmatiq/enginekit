package enginetest

import (
	"context"
	"testing"

	"github.com/dogmatiq/enginekit/enginetest/internal/action"
	"github.com/dogmatiq/enginekit/enginetest/internal/testapp"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func testIntegration(_ context.Context, t *testing.T, e *engine) {
	t.Run("integration message handlers", func(t *testing.T) {
		t.Parallel()

		t.Run("can record events", func(t *testing.T) {
			t.Parallel()

			expect := testapp.
				NewIntegrationEventABuilder().
				WithValue(uuidpb.Generate().AsString()).
				Build()

			e.ExecuteCommand(
				t,
				testapp.
					NewIntegrationCommandABuilder().
					WithActions(action.RecordEvent(expect)).
					Build(),
			)

			e.ExpectEvent(t, expect)
		})
	})
}
