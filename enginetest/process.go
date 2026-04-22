package enginetest

import (
	"context"
	"testing"

	"github.com/dogmatiq/enginekit/enginetest/internal/action"
	"github.com/dogmatiq/enginekit/enginetest/internal/testapp"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func testProcess(_ context.Context, t *testing.T, e *engine) {
	t.Run("process message handlers", func(t *testing.T) {
		t.Parallel()

		t.Run("can execute commands", func(t *testing.T) {
			t.Parallel()

			expect := testapp.
				NewGenericEventBuilder().
				WithValue(uuidpb.Generate().AsString()).
				Build()

			e.RecordEvent(
				t,
				testapp.
					NewProcessEventABuilder().
					WithInstanceId(uuidpb.Generate().AsString()).
					WithActions(action.ExecuteCommand(
						testapp.
							NewDoActionsBuilder().
							WithActions(action.RecordEvent(expect)).
							Build(),
					)).
					Build(),
			)

			e.ExpectEvent(t, expect)
		})

		t.Run("can execute commands on existing un-ended instances", func(t *testing.T) {
			t.Parallel()

			processInstanceID := uuidpb.Generate().AsString()

			expect := testapp.
				NewGenericEventBuilder().
				WithValue(uuidpb.Generate().AsString()).
				Build()

			e.RecordEvent(
				t,
				testapp.
					NewProcessEventABuilder().
					WithInstanceId(processInstanceID).
					WithActions(action.ExecuteCommand(
						testapp.
							NewDoActionsBuilder().
							WithActions(action.RecordEvent(
								testapp.
									NewProcessEventABuilder().
									WithInstanceId(processInstanceID).
									WithActions(action.ExecuteCommand(
										testapp.
											NewDoActionsBuilder().
											WithActions(action.RecordEvent(expect)).
											Build(),
									)).
									Build(),
							)).
							Build(),
					)).
					Build(),
			)

			e.ExpectEvent(t, expect)
		})

		t.Run("do not handle events for ended process instances", func(t *testing.T) {
			t.Parallel()

			instanceID := uuidpb.Generate().AsString()

			e.RecordEvent(
				t,
				testapp.
					NewProcessEventABuilder().
					WithInstanceId(instanceID).
					WithActions(action.Sequence(
						action.ExecuteCommand(
							testapp.
								NewDoActionsBuilder().
								WithActions(action.RecordEvent(
									testapp.
										NewProcessEventABuilder().
										WithInstanceId(instanceID).
										WithActions(action.Fail("event handled for ended process instance")).
										Build(),
								)).
								Build(),
						),
						action.End(),
					)).
					Build(),
			)
		})
	})
}
