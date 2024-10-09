package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
)

func TestProcess_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &ProcessMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A") // note: non-canonical UUID
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromProcess(h)

		Expect(
			t,
			"unexpected identity",
			cfg.Identity(),
			Identity{
				Name: "name",
				Key:  "19cb98d5-dd17-4daf-ae00-1b413b7b899a", // note: canonicalized
			},
		)
	})

	t.Run("it panics if there is no singular valid identity", func(t *testing.T) {
		cases := []struct {
			Name    string
			Want    string
			Handler dogma.ProcessMessageHandler
		}{
			{
				"no identity",
				`no identity is configured`,
				&ProcessMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
							dogma.ExecutesCommand[CommandStub[TypeA]](),
						)
					},
				},
			},
			{
				"invalid identity",
				`identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				&ProcessMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("name", "non-uuid")
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
							dogma.ExecutesCommand[CommandStub[TypeA]](),
						)
					},
				},
			},
			{
				"multiple identities",
				`multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
				&ProcessMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
							dogma.ExecutesCommand[CommandStub[TypeA]](),
						)
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := runtimeconfig.FromProcess(c.Handler)

				ExpectPanic(
					t,
					c.Want,
					func() {
						cfg.Identity()
					},
				)
			})
		}
	})
}

func TestProcess_validation(t *testing.T) {
	cases := []struct {
		Name    string
		Want    string
		Handler dogma.ProcessMessageHandler
	}{
		{
			"valid",
			``, // no error
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					)
				},
			},
		},
		{
			"nil process",
			`process is invalid:` +
				"\n" + `  - no identity is configured` +
				"\n" + `  - expected at least one "HandlesEvent" route` +
				"\n" + `  - expected at least one "ExecutesCommand" route`,
			nil,
		},
		{
			"unconfigured process",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid:` +
				"\n" + `  - no identity is configured` +
				"\n" + `  - expected at least one "HandlesEvent" route` +
				"\n" + `  - expected at least one "ExecutesCommand" route`,
			&ProcessMessageHandlerStub{},
		},
		{
			"process identity must be valid",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"process must not have multiple identities",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"process must handle at least one event type",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: expected at least one "HandlesEvent" route`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						// <-- MISSING "HandlesEvent" ROUTE
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"process must execute at least one command type",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: expected at least one "ExecutesCommand" route`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						// <-- MISSING "ExecutesCommand" ROUTE
					)
				},
			},
		},
		{
			"process must not have multiple routes for the same event type",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple "HandlesEvent" routes are configured for event:github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"process must not have multiple routes for the same command type",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple "ExecutesCommand" routes are configured for command:github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.ExecutesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			},
		},
		{
			"process must not have multiple routes for the same timeout type",
			`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple "SchedulesTimeout" routes are configured for timeout:github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cfg := runtimeconfig.FromProcess(c.Handler)

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
