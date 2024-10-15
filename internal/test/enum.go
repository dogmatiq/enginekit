package test

import (
	"fmt"
	"iter"
	"reflect"
	"testing"

	"github.com/dogmatiq/enginekit/internal/enum"
)

// EnumParameters is a set of parameters for testing enumeration types.
type EnumParameters[E enum.Enum] struct {
	InclusiveBegin, InclusiveEnd E
	Switch                       any // func(E, ...func())
	MapToString                  any // func(E, ...func() string)
}

func (p EnumParameters[E]) cardinality() int {
	return 1 + int(p.InclusiveEnd-p.InclusiveBegin)
}

func (p EnumParameters[E]) values() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := p.InclusiveBegin; v <= p.InclusiveEnd; v++ {
			if !yield(v) {
				return
			}
		}
	}
}

// Enum is a helper for testing enumeration types.
func Enum[E enum.Enum](t *testing.T, params EnumParameters[E]) {
	if params.InclusiveBegin != 0 {
		t.Fatalf("unexpected begin value: got %d, want 0", params.InclusiveBegin)
	}

	if params.InclusiveBegin >= params.InclusiveEnd {
		t.Fatalf("unexpected end value: got %d, want > %d", params.InclusiveEnd, params.InclusiveBegin)
	}

	t.Run("switch", func(t *testing.T) { enumSwitch(t, params) })
	t.Run("map", func(t *testing.T) { enumMap(t, params) })
}

func enumSwitch[E enum.Enum](t *testing.T, params EnumParameters[E]) {
	fn := reflect.ValueOf(params.Switch)

	if fn.Kind() != reflect.Func {
		t.Fatalf("unexpected type for switch function: %T", params.Switch)
	}

	arity := 1 + params.cardinality()
	if fn.Type().NumIn() != arity {
		t.Fatalf("unexpected number of arguments for switch function: got %d, want %d", fn.Type().NumIn(), arity)
	}

	if fn.Type().NumOut() != 0 {
		t.Fatalf("switch function should not return a value")
	}

	t.Run("it calls the function associated with each case", func(t *testing.T) {
		for v := range params.values() {
			t.Run(fmt.Sprintf("case #%d", v), func(t *testing.T) {
				args := []reflect.Value{reflect.ValueOf(v)}
				called := false

				for x := range params.values() {
					arg := func() { t.Errorf("invoked unexpected case function: got %d, want %d", x, v) }
					if v == x {
						arg = func() { called = true }
					}
					args = append(args, reflect.ValueOf(arg))
				}

				fn.Call(args)

				if !called {
					t.Fatalf("case function for %d was not invoked", v)
				}
			})
		}
	})

	t.Run("it panics when the associated function is nil", func(t *testing.T) {
		for v := range params.values() {
			t.Run(fmt.Sprintf("case #%d", v), func(t *testing.T) {
				args := []reflect.Value{reflect.ValueOf(v)}

				for x := range params.values() {
					arg := func() { t.Errorf("invoked unexpected case function: got %d, want %d", x, v) }
					if v == x {
						arg = nil
					}
					args = append(args, reflect.ValueOf(arg))
				}

				defer func() {
					want := fmt.Sprintf("no case function was provided for %q", v.String())
					if got := recover(); got != want {
						t.Fatalf("unexpected panic: got %q, want %q", got, want)
					}
				}()

				fn.Call(args)
			})
		}
	})

	t.Run("it panics when the kind is invalid", func(t *testing.T) {
		cases := []E{
			params.InclusiveBegin - 1,
			params.InclusiveEnd + 1,
		}

		for _, c := range cases {
			t.Run(fmt.Sprintf("case #%d", c), func(t *testing.T) {
				args := []reflect.Value{reflect.ValueOf(c)}
				for range params.values() {
					var arg func()
					args = append(args, reflect.ValueOf(arg))
				}

				defer func() {
					want := fmt.Sprintf("invalid %T (%d)", c, c)
					if got := recover(); got != want {
						t.Fatalf("unexpected panic: got %q, want %q", got, want)
					}
				}()

				fn.Call(args)
			})
		}
	})
}

func enumMap[E enum.Enum](t *testing.T, params EnumParameters[E]) {
	fn := reflect.ValueOf(params.MapToString)

	if fn.Kind() != reflect.Func {
		t.Fatalf("unexpected type for map function: %T", params.MapToString)
	}

	arity := 1 + params.cardinality()
	if fn.Type().NumIn() != arity {
		t.Fatalf("unexpected number of arguments for map function: got %d, want %d", fn.Type().NumIn(), arity)
	}

	if fn.Type().NumOut() != 1 {
		t.Fatalf("unexpected number of return values for map function: got %d, want 1", fn.Type().NumOut())
	}

	t.Run("it returns the value associated with each case", func(t *testing.T) {
		for v := range params.values() {
			t.Run(fmt.Sprintf("case #%d", v), func(t *testing.T) {
				args := []reflect.Value{reflect.ValueOf(v)}

				for x := range params.values() {
					args = append(
						args,
						reflect.ValueOf(
							fmt.Sprintf("case #%d", x),
						),
					)
				}

				values := fn.Call(args)

				want := fmt.Sprintf("case #%d", v)
				got := values[0].Interface()

				if got != want {
					t.Fatalf("unexpected result: got %q, want %q", got, want)
				}
			})
		}
	})

	t.Run("it panics when the kind is invalid", func(t *testing.T) {
		cases := []E{
			params.InclusiveBegin - 1,
			params.InclusiveEnd + 1,
		}

		for _, c := range cases {
			t.Run(fmt.Sprintf("case #%d", c), func(t *testing.T) {
				args := []reflect.Value{reflect.ValueOf(c)}

				for x := range params.values() {
					args = append(
						args,
						reflect.ValueOf(
							fmt.Sprintf("case #%d", x),
						),
					)
				}

				defer func() {
					want := fmt.Sprintf("invalid %T (%d)", c, c)
					if got := recover(); got != want {
						t.Fatalf("unexpected panic: got %q, want %q", got, want)
					}
				}()

				fn.Call(args)
			})
		}
	})
}
