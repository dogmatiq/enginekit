package enum

import "fmt"

// SwitchCase is a case for an enum value of type E.
type SwitchCase[E Enum, T any] struct {
	Case  E
	Value T
}

// Case creates a switch case.
func Case[E Enum, T any](in E, out T) SwitchCase[E, T] {
	return SwitchCase[E, T]{in, out}
}

// Switch calls a function associated with an enum value.
func Switch[E Enum](
	v E,
	cases ...SwitchCase[E, func()],
) {
	fn := Map(v, cases...)

	if fn == nil {
		panic(fmt.Sprintf("no case function was provided for %q", v))
	}

	fn()
}

// Map maps an enum value to a value of type T.
func Map[E Enum, T any](
	v E,
	cases ...SwitchCase[E, T],
) T {
	for _, c := range cases {
		if v == c.Case {
			return c.Value
		}
	}

	panic(fmt.Sprintf("invalid %T", v))
}
