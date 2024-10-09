package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestAggregate_validation(t *testing.T) {
	cases := []struct {
		Name    string
		Want    string
		Handler dogma.AggregateMessageHandler
	}{
		{
			"valid",
			``, // no error
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"nil handler",
			`aggregate(?) is invalid:` +
				"\n" + `  - no identity is configured` +
				"\n" + `  - expected at least one "HandlesCommand" route` +
				"\n" + `  - expected at least one "RecordsEvent" route`,
			nil,
		},
		{
			"unconfigured handler",
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) is invalid:` +
				"\n" + `  - no identity is configured` +
				"\n" + `  - expected at least one "HandlesCommand" route` +
				"\n" + `  - expected at least one "RecordsEvent" route`,
			&AggregateMessageHandlerStub{},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cfg := runtimeconfig.FromAggregate(c.Handler)

			var got string
			if _, err := Normalize(cfg); err != nil {
				got = err.Error()
			}

			if c.Want != got {
				t.Log("unexpected error:")
				t.Log("  got:  ", got)
				t.Log("  want: ", c.Want)
				t.FailNow()
			}
		})
	}
}
