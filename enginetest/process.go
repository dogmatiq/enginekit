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

			expect := &testapp.GenericEvent{
				Value: uuidpb.Generate().AsString(),
			}

			e.RecordEvent(
				t,
				&testapp.ProcessEventA{
					InstanceId: uuidpb.Generate().AsString(),
					Actions: action.ExecuteCommand(
						&testapp.DoActions{
							Actions: action.RecordEvent(expect),
						},
					),
				},
			)

			e.ExpectEvent(t, expect)
		})

		t.Run("can execute commands on existing un-ended instances", func(t *testing.T) {
			t.Parallel()

			processInstanceID := uuidpb.Generate().AsString()

			expect := &testapp.GenericEvent{
				Value: uuidpb.Generate().AsString(),
			}

			e.RecordEvent(
				t,
				&testapp.ProcessEventA{
					InstanceId: processInstanceID,
					Actions: action.ExecuteCommand(
						&testapp.DoActions{
							Actions: action.RecordEvent(
								&testapp.ProcessEventA{
									InstanceId: processInstanceID,
									Actions: action.ExecuteCommand(
										&testapp.DoActions{
											Actions: action.RecordEvent(expect),
										},
									),
								},
							),
						},
					),
				},
			)

			e.ExpectEvent(t, expect)
		})

		t.Run("do not handle events for ended process instances", func(t *testing.T) {
			t.Parallel()

			instanceID := uuidpb.Generate().AsString()

			e.RecordEvent(
				t,
				&testapp.ProcessEventA{
					InstanceId: instanceID,
					Actions: action.Sequence(
						action.ExecuteCommand(
							&testapp.DoActions{
								Actions: action.RecordEvent(
									&testapp.ProcessEventA{
										InstanceId: instanceID,
										Actions:    action.Fail("event handled for ended process instance"),
									},
								),
							},
						),
						action.End(),
					),
				},
			)
		})
	})
}
