package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
)

type validationTestCase[T any] struct {
	Name  string
	Error string
	Input T
}

func runValidationTests[T any, C Component](
	t *testing.T,
	cases []validationTestCase[T],
	buildConfig func(T) C,
) {
	t.Run("Normalize()", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := buildConfig(c.Input)

				var got string
				if _, err := Normalize(cfg); err != nil {
					got = err.Error()
				}

				if c.Error != got {
					t.Log("unexpected error:")
					t.Log("  got:  ", got)
					t.Log("  want: ", c.Error)
					t.FailNow()
				}
			})
		}
	})

	t.Run("MustNormalize()", func(t *testing.T) {
		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := buildConfig(c.Input)

				if c.Error != "" {
					defer func() {
						r := recover()

						if r == nil {
							t.Fatal("expected MustNormalize() to panic")
						}

						err, ok := r.(error)
						if !ok {
							t.Fatalf("expected panic to be an error, got %T", r)
						}

						if c.Error != err.Error() {
							t.Log("unexpected error:")
							t.Log("  got:  ", err.Error())
							t.Log("  want: ", c.Error)
							t.FailNow()
						}
					}()
				}

				MustNormalize(cfg)
			})
		}
	})
}
