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

func TestAggregate_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &AggregateMessageHandlerStub{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A")
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromAggregate(h)

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
			Handler dogma.AggregateMessageHandler
		}{
			{
				"no identity",
				`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: no identity is configured`,
				&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeA]](),
						)
					},
				},
			},
			{
				"invalid identity",
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
				"multiple identities",
				`aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
				&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeA]](),
						)
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := runtimeconfig.FromAggregate(c.Handler)

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

func TestAggregate_RouteSet(t *testing.T) {
	t.Run("it returns the normalized routes", func(t *testing.T) {
		h := &AggregateMessageHandlerStub{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromAggregate(h)

		Expect(
			t,
			"unexpected routes",
			cfg.RouteSet().MessageTypes(),
			map[message.Type]RouteDirection{
				message.TypeFor[CommandStub[TypeA]](): InboundDirection,
				message.TypeFor[EventStub[TypeA]]():   OutboundDirection,
			},
		)
	})

	t.Run("it panics if the routes are invalid", func(t *testing.T) {
		cases := []struct {
			Name  string
			Want  string
			Route Route
		}{
			{
				"empty route",
				`aggregate is invalid: route is invalid: could not evaluate entire configuration`,
				Route{},
			},
			{
				"unexpected ExecutesCommand route",
				`aggregate is invalid: unexpected route: executes-command(pkg.SomeCommandType)`,
				Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(ExecutesCommandRouteType),
						MessageTypeName: optional.Some("pkg.SomeCommandType"),
					},
				},
			},
			{
				"unexpected HandlesEvent route",
				`aggregate is invalid: unexpected route: handles-event(pkg.SomeEventType)`,
				Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(HandlesEventRouteType),
						MessageTypeName: optional.Some("pkg.SomeEventType"),
					},
				},
			},
			{
				"unexpected SchedulesTimeout route",
				`aggregate is invalid: unexpected route: schedules-timeout(pkg.SomeTimeoutType)`,
				Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(SchedulesTimeoutRouteType),
						MessageTypeName: optional.Some("pkg.SomeTimeoutType"),
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := &Aggregate{
					AsConfigured: AggregateAsConfigured{
						Routes: []Route{c.Route},
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

func TestAggregate_IsDisabled(t *testing.T) {
	disable := false

	h := &AggregateMessageHandlerStub{
		ConfigureFunc: func(c dogma.AggregateConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesCommand[CommandStub[TypeA]](),
				dogma.RecordsEvent[EventStub[TypeA]](),
			)
			if disable {
				c.Disable()
			}
		},
	}

	cfg := runtimeconfig.FromAggregate(h)

	if cfg.IsDisabled() {
		t.Fatal("did not expect handler to be disabled")
	}

	disable = true
	cfg = runtimeconfig.FromAggregate(h)

	if !cfg.IsDisabled() {
		t.Fatal("expected handler to be disabled")
	}
}

func TestAggregate_Interface(t *testing.T) {
	h := &AggregateMessageHandlerStub{
		ConfigureFunc: func(c dogma.AggregateConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesCommand[CommandStub[TypeA]](),
				dogma.RecordsEvent[EventStub[TypeA]](),
			)
		},
	}

	cfg := runtimeconfig.FromAggregate(h)

	Expect(
		t,
		"unexpected result",
		cfg.Interface(),
		h,
	)
}

func TestAggregate_validation(t *testing.T) {
	cases := []validationTestCase{
		{
			Name:   "valid",
			Expect: ``,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "missing implementations",
			Expect: ``,
			Component: &Aggregate{
				AsConfigured: AggregateAsConfigured{
					Source: optional.Some(
						Value[dogma.AggregateMessageHandler]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"),
						},
					),
					Identities: []Identity{
						{
							AsConfigured: IdentityAsConfigured{
								Name: optional.Some("name"),
								Key:  optional.Some("494157ef-6f91-45ec-ab19-df61bb96210a"),
							},
						},
					},
					Routes: []Route{
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(HandlesCommandRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(RecordsEventRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
					},
				},
			},
		},
		{
			Name: "missing implementations using WithRuntimeValues() option",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid:` +
				"\n" + `- missing implementation: dogma.AggregateMessageHandler value is not available` +
				"\n" + `- handles-command(github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]) is invalid: missing implementation: message.Type value is not available` +
				"\n" + `- records-event(github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]) is invalid: missing implementation: message.Type value is not available`,
			Options: []NormalizeOption{
				WithRuntimeValues(),
			},
			Component: &Aggregate{
				AsConfigured: AggregateAsConfigured{
					Source: optional.Some(
						Value[dogma.AggregateMessageHandler]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"),
						},
					),
					Identities: []Identity{
						{
							AsConfigured: IdentityAsConfigured{
								Name: optional.Some("name"),
								Key:  optional.Some("494157ef-6f91-45ec-ab19-df61bb96210a"),
							},
						},
					},
					Routes: []Route{
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(HandlesCommandRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
						{
							AsConfigured: RouteAsConfigured{
								RouteType:       optional.Some(RecordsEventRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
					},
				},
			},
		},
		{
			Name: "nil aggregate",
			Expect: `aggregate is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "handles-command" route` +
				"\n" + `- expected at least one "records-event" route`,
			Component: runtimeconfig.FromAggregate(nil),
		},
		{
			Name: "unconfigured aggregate",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- expected at least one "handles-command" route` +
				"\n" + `- expected at least one "records-event" route`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{}),
		},
		{
			Name:   "aggregate identity must be valid",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "aggregate must not have multiple identities",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "aggregate must handle at least one command type",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: expected at least one "handles-command" route`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						// <-- MISSING "handles-command" ROUTE
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "aggregate must record at least one event type",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: expected at least one "records-event" route`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						// <-- MISSING "RecordEvent" ROUTE
					)
				},
			}),
		},
		{
			Name:   "aggregate must not have multiple routes for the same command type",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple "handles-command" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "aggregate must not have multiple routes for the same event type",
			Expect: `aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: multiple "records-event" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.RecordsEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			}),
		},
	}

	runValidationTests(t, cases)
}
