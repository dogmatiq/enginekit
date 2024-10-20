package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
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
				`process:ProcessMessageHandlerStub is invalid: no identity is configured`,
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
				`process:ProcessMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
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
				`process:ProcessMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
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
				`process is invalid: unexpected handles-command route for pkg.SomeCommandType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(HandlesCommandRouteType),
						MessageTypeName: optional.Some("pkg.SomeCommandType"),
					},
				},
			},
			{
				"unexpected RecordsEvent route",
				`process is invalid: unexpected records-event route for pkg.SomeEventType`,
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

func TestProcess_render(t *testing.T) {
	cases := []renderTestCase{
		{
			Name:             "complete",
			ExpectDescriptor: `process:ProcessMessageHandlerStub`,
			ExpectDetails: multiline(
				`valid process *github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - valid executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - valid schedules-timeout route for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
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
			Name:             "disabled",
			ExpectDescriptor: `process:ProcessMessageHandlerStub`,
			ExpectDetails: multiline(
				`disabled valid process *github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - valid executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
					)
					c.Disable()
				},
			}),
		},
		{
			Name:             "no runtime type information",
			ExpectDescriptor: `process:SomeProcess`,
			ExpectDetails: multiline(
				`valid process pkg.SomeProcess (runtime type unavailable)`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-event route for pkg.SomeEvent (runtime type unavailable)`,
				`  - valid executes-command route for pkg.SomeCommand (runtime type unavailable)`,
			),
			Component: configbuilder.Process(func(b *configbuilder.ProcessBuilder) {
				b.SetSourceTypeName("pkg.SomeProcess")
				b.SetDisabled(false)

				b.Identity(func(b *configbuilder.IdentityBuilder) {
					b.SetName("name")
					b.SetKey("19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				})

				b.Route(func(b *configbuilder.RouteBuilder) {
					b.SetRouteType(HandlesEventRouteType)
					b.SetMessageTypeName("pkg.SomeEvent")
				})

				b.Route(func(b *configbuilder.RouteBuilder) {
					b.SetRouteType(ExecutesCommandRouteType)
					b.SetMessageTypeName("pkg.SomeCommand")
				})
			}),
		},
		{
			Name:             "empty",
			ExpectDescriptor: `process`,
			ExpectDetails: multiline(
				`incomplete process`,
				`  - no identity is configured`,
				`  - no "handles-event" routes are configured`,
				`  - no "executes-command" routes are configured`,
			),
			Component: &Process{},
		},
		{
			Name:             "invalid",
			ExpectDescriptor: `process:ProcessMessageHandlerStub`,
			ExpectDetails: multiline(
				`invalid process *github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub`,
				`  - no "handles-event" routes are configured`,
				`  - no "executes-command" routes are configured`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
			),
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				},
			}),
		},
		{
			Name:             "invalid sub-component",
			ExpectDescriptor: `process:ProcessMessageHandlerStub`,
			ExpectDetails: multiline(
				`valid process *github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub`,
				`  - invalid identity name/non-uuid`,
				`      - invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				`  - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - valid executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
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
	}

	runRenderTests(t, cases)
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
			Name:   "no runtime type information",
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
			Name: "no runtime type information using WithRuntimeTypes() option",
			Expect: `process:ProcessMessageHandlerStub is invalid:` +
				"\n" + `- dogma.ProcessMessageHandler value is not available` +
				"\n" + `- route:handles-event:EventStub[TypeA] is invalid: message.Type value is not available` +
				"\n" + `- route:executes-command:CommandStub[TypeA] is invalid: message.Type value is not available`,
			Options: []NormalizeOption{
				WithRuntimeTypes(),
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
				"\n" + `- no "handles-event" routes are configured` +
				"\n" + `- no "executes-command" routes are configured` +
				"\n" + `- could not evaluate entire configuration`,
			Component: runtimeconfig.FromProcess(nil),
		},
		{
			Name: "unconfigured process",
			Expect: `process:ProcessMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- no "handles-event" routes are configured` +
				"\n" + `- no "executes-command" routes are configured`,
			Component: runtimeconfig.FromProcess(&ProcessMessageHandlerStub{}),
		},
		{
			Name:   "process identity must be valid",
			Expect: `process:ProcessMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
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
			Expect: `process:ProcessMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
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
			Expect: `process:ProcessMessageHandlerStub is invalid: no "handles-event" routes are configured`,
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
			Expect: `process:ProcessMessageHandlerStub is invalid: no "executes-command" routes are configured`,
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
			Expect: `process:ProcessMessageHandlerStub is invalid: multiple "handles-event" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
			Expect: `process:ProcessMessageHandlerStub is invalid: multiple "executes-command" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
			Expect: `process:ProcessMessageHandlerStub is invalid: multiple "schedules-timeout" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
