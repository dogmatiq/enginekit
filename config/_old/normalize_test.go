package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/test"
)

type validationTestCase struct {
	Name      string
	Expect    string
	Options   []NormalizeOption
	Component Component
}

func runValidationTests(
	t *testing.T,
	cases []validationTestCase,
) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Helper()

			var got string
			if _, err := Normalize(c.Component, c.Options...); err != nil {
				got = err.Error()
			}

			if c.Expect != got {
				test.Expect(
					t,
					"unexpected error",
					got,
					c.Expect,
				)
			}

			// Only test MustNormalize() if Normalize() passes, otherwise
			// we'll likely just get duplicate error messages.
			if c.Expect != "" {
				defer func() {
					r := recover()

					if r == nil {
						t.Fatal("expected MustNormalize() to panic")
					}

					err, ok := r.(error)
					if !ok {
						t.Fatalf("expected panic to be an error, got %T", r)
					}

					test.Expect(
						t,
						"unexpected error",
						err.Error(),
						c.Expect,
					)
				}()
			}

			MustNormalize(c.Component, c.Options...)
		})
	}
}
