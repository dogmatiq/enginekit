package config_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/test"
)

type validationTestCases []struct {
	Name      string
	Expect    string
	Options   []ValidateOption
	Component Component
}

func runValidationTests(t *testing.T, cases validationTestCases) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Helper()

			var got string
			if err := Validate(c.Component, c.Options...); err != nil {
				got = err.Error()
			}

			test.Expect(
				t,
				"unexpected error",
				got,
				c.Expect,
			)
		})
	}
}

func multiline(lines ...string) string {
	return strings.Join(lines, "\n")
}
