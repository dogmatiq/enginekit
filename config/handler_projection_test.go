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

func TestProjection_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &ProjectionMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A") // note: non-canonical UUID
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromProjection(h)

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
			Handler dogma.ProjectionMessageHandler
		}{
			{
				"no identity",
				`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: no identity is configured`,
				&ProjectionMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
						)
					},
				},
			},
			{
				"invalid identity",
				`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				&ProjectionMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("name", "non-uuid")
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
						)
					},
				},
			},
			{
				"multiple identities",
				`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
				&ProjectionMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
						)
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := runtimeconfig.FromProjection(c.Handler)

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

func TestProjection_Routes(t *testing.T) {
	t.Run("it returns the normalized routes", func(t *testing.T) {
		h := &ProjectionMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromProjection(h)

		Expect(
			t,
			"unexpected routes",
			cfg.Routes(),
			RouteSet{
				{
					RouteType:       optional.Some(HandlesEventRouteType),
					MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
					MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
				},
			},
		)
	})

	t.Run("it panics if the routes are invalid", func(t *testing.T) {
		cfg := Projection{
			ConfiguredRoutes: []Route{
				{},
				{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommandType"),
				},
				{
					RouteType:       optional.Some(ExecutesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommandType"),
				},
				{
					RouteType:       optional.Some(RecordsEventRouteType),
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
			`partial projection is invalid:`+
				"\n"+`- route is invalid:`+
				"\n"+`  - missing route type`+
				"\n"+`  - missing message type`+
				"\n"+`- unexpected route: HandlesCommand:pkg.SomeCommandType`+
				"\n"+`- unexpected route: ExecutesCommand:pkg.SomeCommandType`+
				"\n"+`- unexpected route: RecordsEvent:pkg.SomeEventType`+
				"\n"+`- unexpected route: SchedulesTimeout:pkg.SomeTimeoutType`+
				"\n"+`- expected at least one "HandlesEvent" route`,
			func() {
				cfg.Routes()
			},
		)
	})
}

func TestProjection_Interface(t *testing.T) {
	h := &ProjectionMessageHandlerStub{
		ConfigureFunc: func(c dogma.ProjectionConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesEvent[EventStub[TypeA]](),
			)
		},
	}

	cfg := runtimeconfig.FromProjection(h)

	Expect(
		t,
		"unexpected result",
		cfg.Interface(),
		h,
	)
}

func TestProjection_validation(t *testing.T) {
	cases := []struct {
		Name    string
		Want    string
		Handler dogma.ProjectionMessageHandler
	}{
		{
			"valid",
			``, // no error
			&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"nil projection",
			`partial projection is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "HandlesEvent" route`,
			nil,
		},
		{
			"unconfigured projection",
			`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "HandlesEvent" route`,
			&ProjectionMessageHandlerStub{},
		},
		{
			"projection identity must be valid",
			`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"projection must not have multiple identities",
			`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			},
		},
		{
			"projection must handle at least one event type",
			`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: expected at least one "HandlesEvent" route`,
			&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					// <-- MISSING "HandlesEvent" ROUTE
				},
			},
		},
		{
			"projection must not have multiple routes for the same event type",
			`projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub is invalid: multiple "HandlesEvent" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cfg := runtimeconfig.FromProjection(c.Handler)

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
