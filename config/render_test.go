package config_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/test"
)

type renderTestCase struct {
	Name             string
	ExpectDescriptor string
	ExpectDetails    string
	Component        Component
}

func runRenderTests(t *testing.T, cases []renderTestCase) {
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			test.Expect(
				t,
				"unexpected descriptor",
				RenderDescriptor(c.Component),
				c.ExpectDescriptor,
			)

			test.Expect(
				t,
				"unexpected descriptor from String() method",
				c.Component.String(),
				c.ExpectDescriptor,
			)

			test.Expect(
				t,
				"unexpected details",
				RenderDetails(c.Component),
				c.ExpectDetails+"\n",
			)
		})
	}
}

func multiline(lines ...string) string {
	return strings.Join(lines, "\n")
}
