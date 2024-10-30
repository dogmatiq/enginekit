package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/test"
)

type renderTestCases []struct {
	Name                string
	String, Description string
	Component           Component
}

func testDescribe(t *testing.T, cases renderTestCases) {
	t.Helper()

	t.Run("func Format()", func(t *testing.T) {
		t.Helper()

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				t.Helper()

				test.Expect(
					t,
					"unexpected short string representation",
					c.Component.String(),
					c.String,
				)

				test.Expect(
					t,
					"unexpected long string representation",
					Description(
						c.Component,
						WithValidationResult(Validate(c.Component)),
					),
					c.Description+"\n",
				)
			})
		}
	})
}
