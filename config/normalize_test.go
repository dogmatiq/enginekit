package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
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
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var got string
			if _, err := Normalize(c.Component, c.Options...); err != nil {
				got = err.Error()
			}

			if c.Expect != got {
				t.Log("unexpected error:")
				t.Log("  got:  ", got)
				t.Log("  want: ", c.Expect)
				t.FailNow()
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

					if c.Expect != err.Error() {
						t.Log("unexpected error:")
						t.Log("  got:  ", err.Error())
						t.Log("  want: ", c.Expect)
						t.FailNow()
					}
				}()
			}

			MustNormalize(c.Component, c.Options...)
		})
	}
}
