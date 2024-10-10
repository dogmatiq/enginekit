package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

func TestIntegration_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &IntegrationMessageHandlerStub{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A") // note: non-canonical UUID
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromIntegration(h)

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
			Handler dogma.IntegrationMessageHandler
		}{
			{
				"no identity",
				`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: no identity is configured`,
				&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
						)
					},
				},
			},
			{
				"invalid identity",
				`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("name", "non-uuid")
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
						)
					},
				},
			},
			{
				"multiple identities",
				`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
				&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
						)
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := runtimeconfig.FromIntegration(c.Handler)

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

func TestIntegration_Routes(t *testing.T) {
	t.Run("it returns the normalized routes", func(t *testing.T) {
		h := &IntegrationMessageHandlerStub{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromIntegration(h)

		Expect(
			t,
			"unexpected routes",
			cfg.Routes(),
			RouteSet{
				{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
					MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
				},
				{
					RouteType:       optional.Some(RecordsEventRouteType),
					MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
					MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
				},
			},
		)
	})

	t.Run("it panics if the routes are invalid", func(t *testing.T) {
		cfg := Integration{
			ConfiguredRoutes: []Route{
				{},
				{
					RouteType:       optional.Some(ExecutesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommandType"),
				},
				{
					RouteType:       optional.Some(HandlesEventRouteType),
					MessageTypeName: optional.Some("pkg.SomeEventType"),
				},
				{
					RouteType:       optional.Some(SchedulesTimeoutRouteType),
					MessageTypeName: optional.Some("pkg.SomeTimeoutType"),
				},
			},
		}

		ExpectPanic(
			t,
			`partial integration is invalid:`+
				"\n"+`- route is invalid:`+
				"\n"+`  - missing route type`+
				"\n"+`  - missing message type`+
				"\n"+`- unexpected route: ExecutesCommand:pkg.SomeCommandType`+
				"\n"+`- unexpected route: HandlesEvent:pkg.SomeEventType`+
				"\n"+`- unexpected route: SchedulesTimeout:pkg.SomeTimeoutType`+
				"\n"+`- expected at least one "HandlesCommand" route`,
			func() {
				cfg.Routes()
			},
		)
	})
}

func TestIntegration_Interface(t *testing.T) {
	h := &IntegrationMessageHandlerStub{
		ConfigureFunc: func(c dogma.IntegrationConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesCommand[CommandStub[TypeA]](),
			)
		},
	}

	cfg := runtimeconfig.FromIntegration(h)

	Expect(
		t,
		"unexpected result",
		cfg.Interface(),
		h,
	)
}

func TestIntegration_validation(t *testing.T) {
	cases := []struct {
		Name    string
		Want    string
		Handler dogma.IntegrationMessageHandler
	}{
		{
			"valid",
			``, // no error
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"nil integration",
			`partial integration is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "HandlesCommand" route`,
			nil,
		},
		{
			"unconfigured integration",
			`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "HandlesCommand" route`,
			&IntegrationMessageHandlerStub{},
		},
		{
			"integration identity must be valid",
			`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"integration must not have multiple identities",
			`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			},
		},
		{
			"integration must handle at least one command type",
			`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: expected at least one "HandlesCommand" route`,
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					// <-- MISSING "HandlesCommand" ROUTE
				},
			},
		},
		{
			"integration must not have multiple routes for the same command type",
			`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: multiple "HandlesCommand" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			},
		},
		{
			"integration must not have multiple routes for the same event type",
			`integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub is invalid: multiple "RecordsEvent" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
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
			cfg := runtimeconfig.FromIntegration(c.Handler)

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