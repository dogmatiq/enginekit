package message

import (
	"fmt"

	"github.com/dogmatiq/dogma"
)

// SwitchByKind invokes one of the provided functions based on k.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [Kind] values are added in the future.
//
// It panics with a meaningful message if the function associated with k.
//
// It panics if the function associated with k is nil, or if k is not a valid
// [Kind].
func SwitchByKind(
	k Kind,
	command func(),
	event func(),
	timeout func(),
) {
	fn := MapByKind(k, command, event, timeout)

	if fn == nil {
		panic(fmt.Sprintf("no case function was provided for the %q kind", k))
	}

	fn()
}

// MapByKind maps k to a value of type T.
//
// It provides a compile-time guarantee that all possible values are handled,
// even if new [Kind] values are added in the future.
//
// It panics if k is not a valid [Kind].
func MapByKind[T any](k Kind, command, event, timeout T) (result T) {
	switch k {
	case CommandKind:
		return command
	case EventKind:
		return event
	case TimeoutKind:
		return timeout
	default:
		panic("invalid kind")
	}
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
			panic("no case function was provided for dogma.Command messages")
		}
		command(m)
	case dogma.Event:
		if event == nil {
			panic("no case function was provided for dogma.Event messages")
		}
		event(m)
	case dogma.Timeout:
		if timeout == nil {
			panic("no case function was provided for dogma.Timeout messages")
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
// It panics with a meaningful message if the function associated with m's kind
// is nil.
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
		mapFunc(command, &result),
		mapFunc(event, &result),
		mapFunc(timeout, &result),
	)

	return result
}

// mapFunc returns a function that invokes fn and assigns the result to *v.
func mapFunc[K dogma.Message, T any](
	fn func(K) T,
	v *T,
) func(K) {
	if fn == nil {
		return nil
	}
	return func(m K) {
		*v = fn(m)
	}
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
		mapWithErrFunc(command, &result, &err),
		mapWithErrFunc(event, &result, &err),
		mapWithErrFunc(timeout, &result, &err),
	)

	return result, err
}

// mapWithErrFunc returns a function that invokes fn and assigns the result and
// error value to *v and *err, respectively.
func mapWithErrFunc[K dogma.Message, T any](
	fn func(K) (T, error),
	v *T,
	err *error,
) func(K) {
	if fn == nil {
		return nil
	}
	return func(m K) {
		*v, *err = fn(m)
	}
}
