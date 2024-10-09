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
			"nil aggregate",
			`aggregate is invalid:` +
				"\n" + `  - no identity is configured` +
				"\n" + `  - expected at least one "HandlesCommand" route` +
				"\n" + `  - expected at least one "RecordsEvent" route`,
			nil,
		},
		{
			"unconfigured aggregate",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid:` +
				"\n" + `  - no identity is configured` +
				"\n" + `  - expected at least one "HandlesCommand" route` +
				"\n" + `  - expected at least one "RecordsEvent" route`,
			&AggregateMessageHandlerStub{},
		},
		{
			"aggregate identity must be valid",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"aggregate must not have multiple identities",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"aggregate must handle at least one command",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: expected at least one "HandlesCommand" route`,
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						// <-- MISSING "HandlesCommand" ROUTE
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"aggregate must record at least one event type",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: expected at least one "RecordsEvent" route`,
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						// <-- MISSING "RecordEvent" ROUTE
					)
				},
			},
		},
		{
			"aggregate must not have multiple routes for the same command type",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple "HandlesCommand" routes are configured for command:github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"aggregate must not have multiple routes for the same event type",
			`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple "RecordsEvent" routes are configured for event:github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.RecordsEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			},
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
