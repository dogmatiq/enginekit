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
		{
			Behavior: &Action_Fail{message},
		},
	}
}

func (x *Action_Fail) do(Scope) error {
	return errors.New(x.Fail)
}

// Log returns an action that causes the handler to log a human-readable
// message.
func Log(message string) []*Action {
	return []*Action{
		{
			Behavior: &Action_Log{message},
		},
	}
}

func (x *Action_Log) do(s Scope) error {
	s.Log("%s", x.Log)
	return nil
}

// ExecuteCommand returns an action that causes the handler to execute a command
// message.
func ExecuteCommand(c dogma.Command) []*Action {
	return []*Action{
		{
			Behavior: &Action_ExecuteCommand{toAny(c)},
		},
	}
}

func (x *Action_ExecuteCommand) do(s Scope) error {
	s.(executor).ExecuteCommand(fromAny[dogma.Command](x.ExecuteCommand))
	return nil
}

// RecordEvent returns an action that causes the handler to record an event
// message.
func RecordEvent(e dogma.Message) []*Action {
	return []*Action{
		{
			Behavior: &Action_RecordEvent{toAny(e)},
		},
	}
}

func (x *Action_RecordEvent) do(s Scope) error {
	s.(recorder).RecordEvent(fromAny[dogma.Event](x.RecordEvent))
	return nil
}

// ScheduleTimeout returns an action that causes the handler to schedule a
// timeout message.
func ScheduleTimeout(t dogma.Timeout, at time.Time) []*Action {
	return []*Action{
		{
			Behavior: &Action_ScheduleTimeout{
				&ScheduleTimeoutDetails{
					Timeout: toAny(t),
					At:      timestamppb.New(at),
				},
			},
		},
	}
}

func (x *Action_ScheduleTimeout) do(s Scope) error {
	s.(scheduler).ScheduleTimeout(
		fromAny[dogma.Timeout](x.ScheduleTimeout.Timeout),
		x.ScheduleTimeout.At.AsTime(),
	)
	return nil
}

// Destroy returns an action that causes the handler to destroy the aggregate
// instance.
func Destroy() []*Action {
	return []*Action{
		{
			Behavior: &Action_Destroy{},
		},
	}
}

func (x *Action_Destroy) do(s Scope) error {
	s.(destroyer).Destroy()
	return nil
}

// End returns an action that causes the handler to end the process instance.
func End() []*Action {
	return []*Action{
		{
			Behavior: &Action_End{},
		},
	}
}

func (x *Action_End) do(s Scope) error {
	s.(ender).End()
	return nil
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
