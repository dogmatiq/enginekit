package action

import (
	"errors"
	"time"

	"github.com/dogmatiq/dogma"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Sequence returns a new set of actions.
func Sequence(actions ...[]*Action) []*Action {
	var merged []*Action
	for _, a := range actions {
		merged = append(merged, a...)
	}
	return merged
}

// Fail returns an action that causes the handler to return an error.
func Fail(message string) []*Action {
	return []*Action{
		NewActionBuilder().
			WithFail(message).
			Build(),
	}
}

// Log returns an action that causes the handler to log a human-readable
// message.
func Log(message string) []*Action {
	return []*Action{
		NewActionBuilder().
			WithLog(message).
			Build(),
	}
}

// ExecuteCommand returns an action that causes the handler to execute a command
// message.
func ExecuteCommand(c dogma.Command) []*Action {
	return []*Action{
		NewActionBuilder().
			WithExecuteCommand(toAny(c)).
			Build(),
	}
}

// RecordEvent returns an action that causes the handler to record an event
// message.
func RecordEvent(e dogma.Message) []*Action {
	return []*Action{
		NewActionBuilder().
			WithRecordEvent(toAny(e)).
			Build(),
	}
}

// ScheduleTimeout returns an action that causes the handler to schedule a
// timeout message.
func ScheduleTimeout(t dogma.Timeout, at time.Time) []*Action {
	return []*Action{
		NewActionBuilder().
			WithScheduleTimeout(
				NewScheduleTimeoutDetailsBuilder().
					WithTimeout(toAny(t)).
					WithAt(timestamppb.New(at)).
					Build(),
			).
			Build(),
	}
}

// Destroy returns an action that causes the handler to destroy the aggregate
// instance.
func Destroy() []*Action {
	return []*Action{
		NewActionBuilder().
			WithDestroy(NewEmptyBuilder().Build()).
			Build(),
	}
}

// End returns an action that causes the handler to end the process instance.
func End() []*Action {
	return []*Action{
		NewActionBuilder().
			WithEnd(NewEmptyBuilder().Build()).
			Build(),
	}
}

func do(act *Action, s Scope) error {
	return MustMap_Action_Behavior(
		act,
		func(message string) error {
			return errors.New(message)
		},
		func(message string) error {
			s.Log("%s", message)
			return nil
		},
		func(c *anypb.Any) error {
			s.(executor).ExecuteCommand(fromAny[dogma.Command](c))
			return nil
		},
		func(e *anypb.Any) error {
			s.(recorder).RecordEvent(fromAny[dogma.Event](e))
			return nil
		},
		func(details *ScheduleTimeoutDetails) error {
			s.(scheduler).ScheduleTimeout(
				fromAny[dogma.Timeout](details.GetTimeout()),
				details.GetAt().AsTime(),
			)
			return nil
		},
		func(*Empty) error {
			s.(destroyer).Destroy()
			return nil
		},
		func(*Empty) error {
			s.(ender).End()
			return nil
		},
	)
}

func toAny(m dogma.Message) *anypb.Any {
	x, err := anypb.New(m.(proto.Message))
	if err != nil {
		panic(err)
	}
	return x
}

func fromAny[T dogma.Message](m *anypb.Any) T {
	x, err := m.UnmarshalNew()
	if err != nil {
		panic(err)
	}
	return x.(T)
}
