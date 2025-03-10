package message

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/internal/enum"
)

// SwitchByKind invokes one of the provided functions based on k.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [Kind] values are added in the future.
//
// It panics if the function associated with k is nil, or if k is not a valid
// [Kind].
func SwitchByKind(
	k Kind,
	command func(),
	event func(),
	timeout func(),
) {
	enum.Switch(k, command, event, timeout)
}

// MapByKind maps k to a value of type T.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [Kind] values are added in the future.
//
// It panics if k is not a valid [Kind].
func MapByKind[T any](k Kind, command, event, timeout T) T {
	return enum.Map(k, command, event, timeout)
}

// SwitchByKindOf invokes one of the provided functions based on the [Kind] of m.
//
// It provides a compile-time guarantee that all kinds are handled, even if new
// [Kind] values are added in the future.
//
// It panics if the function associated with m's kind is nil, or if m does not
// implement [dogma.Command], [dogma.Event] or [dogma.Timeout].
func SwitchByKindOf(
	m dogma.Message,
	command func(dogma.Command),
	event func(dogma.Event),
	timeout func(dogma.Timeout),
) {
	switch m := m.(type) {
	case dogma.Command:
		if command == nil {
			panic("no case function was provided for dogma.Command")
		}
		command(m)
	case dogma.Event:
		if event == nil {
			panic("no case function was provided for dogma.Event")
		}
		event(m)
	case dogma.Timeout:
		if timeout == nil {
			panic("no case function was provided for dogma.Timeout")
		}
		timeout(m)
	default:
		panic(fmt.Sprintf(
			"%T implements dogma.Message, but does not implement dogma.Command, dogma.Event or dogma.Timeout",
			m,
		))
	}
}

// MapByKindOf invokes one of the provided functions based on the [Kind] of m, and
// returns the result.
//
// It provides a compile-time guarantee that all kinds are handled, even if new
// [Kind] values are added in the future.
//
// It panics if the function associated with m's kind is nil, or if m does not
// implement [dogma.Command], [dogma.Event] or [dogma.Timeout].
func MapByKindOf[T any](
	m dogma.Message,
	command func(dogma.Command) T,
	event func(dogma.Event) T,
	timeout func(dogma.Timeout) T,
) (result T) {
	SwitchByKindOf(
		m,
		enum.AssignResult(command, &result),
		enum.AssignResult(event, &result),
		enum.AssignResult(timeout, &result),
	)

	return result
}

// MapByKindOfWithErr invokes one of the provided functions based on the [Kind] of m,
// and returns the result and error value.
//
// It provides a compile-time guarantee that all kinds are handled, even if new
// [Kind] values are added in the future.
//
// It panics if the function associated with m's kind is nil, or if m does not
// implement [dogma.Command], [dogma.Event] or [dogma.Timeout].
func MapByKindOfWithErr[T any](
	m dogma.Message,
	command func(dogma.Command) (T, error),
	event func(dogma.Event) (T, error),
	timeout func(dogma.Timeout) (T, error),
) (result T, err error) {
	SwitchByKindOf(
		m,
		enum.AssignResultErr(command, &result, &err),
		enum.AssignResultErr(event, &result, &err),
		enum.AssignResultErr(timeout, &result, &err),
	)

	return result, err
}
