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
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestProcess_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &ProcessMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A")
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
			&identitypb.Identity{
				Name: "name",
				Key:  uuidpb.MustParse("19cb98d5-dd17-4daf-ae00-1b413b7b899a"),
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
				`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: no identity is configured`,
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
				"multiple identities",
				`process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
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

func TestProcess_RouteSet(t *testing.T) {
	t.Run("it returns the normalized routes", func(t *testing.T) {
		h := &ProcessMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
					dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromProcess(h)

		Expect(
			t,
			"unexpected routes",
			cfg.RouteSet().MessageTypes(),
			map[message.Type]RouteDirection{
				message.TypeFor[EventStub[TypeA]]():   InboundDirection,
				message.TypeFor[CommandStub[TypeA]](): OutboundDirection,
				message.TypeFor[TimeoutStub[TypeA]](): InboundDirection | OutboundDirection,
			},
		)
	})

	t.Run("it panics if the routes are invalid", func(t *testing.T) {
		cases := []struct {
			Name  string
			Want  string
			Route *Route
		}{
			{
				"empty route",
				`process is invalid: route is invalid: could not evaluate entire configuration`,
				&Route{},
			},
			{
				"unexpected HandlesCommand route",
				`process is invalid: unexpected route: route:handles-command(pkg.SomeCommandType)`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(HandlesCommandRouteType),
						MessageTypeName: optional.Some("pkg.SomeCommandType"),
					},
				},
			},
			{
				"unexpected RecordsEvent route",
				`process is invalid: unexpected route: route:records-event(pkg.SomeEventType)`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(RecordsEventRouteType),
						MessageTypeName: optional.Some("pkg.SomeEventType"),
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := &Process{
					AsConfigured: ProcessAsConfigured{
						Routes: []*Route{c.Route},
					},
				}

				ExpectPanic(
					t,
					c.Want,
					func() {
						cfg.RouteSet()
					},
				)
			})
		}
	})
}

func TestProcess_IsDisabled(t *testing.T) {
	disable := false

	h := &ProcessMessageHandlerStub{
		ConfigureFunc: func(c dogma.ProcessConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesEvent[EventStub[TypeA]](),
				dogma.ExecutesCommand[CommandStub[TypeA]](),
			)
			if disable {
				c.Disable()
			}
		},
	}

	cfg := runtimeconfig.FromProcess(h)

	if cfg.IsDisabled() {
		t.Fatal("did not expect handler to be disabled")
	}

	disable = true
	cfg = runtimeconfig.FromProcess(h)

	if !cfg.IsDisabled() {
		t.Fatal("expected handler to be disabled")
	}
}

func TestProcess_Interface(t *testing.T) {
	h := &ProcessMessageHandlerStub{
		ConfigureFunc: func(c dogma.ProcessConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesEvent[EventStub[TypeA]](),
				dogma.ExecutesCommand[CommandStub[TypeA]](),
			)
		},
	}

	cfg := runtimeconfig.FromProcess(h)

	Expect(
		t,
		"unexpected result",
		cfg.Interface(),
		h,
	)
}

func TestProcess_validation(t *testing.T) {
	cases := []validationTestCase{
		{
			Name:   "valid",
			Expect: ``,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "missing implementations",
			Expect: ``,
			Component: &Process{
				AsConfigured: ProcessAsConfigured{
					Source: Value[dogma.ProcessMessageHandler]{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
					},
					Identities: []*Identity{
						{
							AsConfigured: IdentityAsConfigured{
								Name: optional.Some("name"),
								Key:  optional.Some("494157ef-6f91-45ec-ab19-df61bb96210a"),
							},
						},
					},
					Routes: []*Route{
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(HandlesEventRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(ExecutesCommandRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
					},
				},
			},
		},
		{
			Name: "missing implementations using WithRuntimeValues() option",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid:` +
				"\n" + `- missing implementation: dogma.ProcessMessageHandler value is not available` +
				"\n" + `- route:handles-event(github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]) is invalid: missing implementation: message.Type value is not available` +
				"\n" + `- route:executes-command(github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]) is invalid: missing implementation: message.Type value is not available`,
			Options: []NormalizeOption{
				WithRuntimeValues(),
			},
			Component: &Process{
				AsConfigured: ProcessAsConfigured{
					Source: Value[dogma.ProcessMessageHandler]{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
					},
					Identities: []*Identity{
						{
							AsConfigured: IdentityAsConfigured{
								Name: optional.Some("name"),
								Key:  optional.Some("494157ef-6f91-45ec-ab19-df61bb96210a"),
							},
						},
					},
					Routes: []*Route{
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(HandlesEventRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(ExecutesCommandRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
					},
				},
			},
		},
		{
			Name: "nil process",
			Expect: `process is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "handles-event" route` +
				"\n" + `- expected at least one "executes-command" route` +
				"\n" + `- could not evaluate entire configuration`,
			Component: runtimeconfig.FromProcess(nil),
		},
		{
			Name: "unconfigured process",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "handles-event" route` +
				"\n" + `- expected at least one "executes-command" route`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{}),
		},
		{
			Name:   "process identity must be valid",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "process must not have multiple identities",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "process must handle at least one event type",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: expected at least one "handles-event" route`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						// <-- MISSING "handles-event" ROUTE
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "process must execute at least one command type",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: expected at least one "executes-command" route`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						// <-- MISSING "executes-command" ROUTE
					)
				},
			}),
		},
		{
			Name:   "process must not have multiple routes for the same event type",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple "handles-event" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "process must not have multiple routes for the same command type",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple "executes-command" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.ExecutesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			}),
		},
		{
			Name:   "process must not have multiple routes for the same timeout type",
			Expect: `process:github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub is invalid: multiple "schedules-timeout" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			}),
		},
	}

	runValidationTests(t, cases)
}
