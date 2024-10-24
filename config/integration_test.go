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

func TestIntegration_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &IntegrationMessageHandlerStub{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A")
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
			Handler dogma.IntegrationMessageHandler
		}{
			{
				"no identity",
				`integration:IntegrationMessageHandlerStub is invalid: no identity is configured`,
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
				`integration:IntegrationMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
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
				`integration:IntegrationMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
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

func TestIntegration_RouteSet(t *testing.T) {
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
			Route *Route
		}{
			{
				"empty route",
				`integration is invalid: route is invalid: could not evaluate entire configuration`,
				&Route{},
			},
			{
				"unexpected ExecutesCommand route",
				`integration is invalid: unexpected executes-command route for pkg.SomeCommandType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(ExecutesCommandRouteType),
						MessageTypeName: optional.Some("pkg.SomeCommandType"),
					},
				},
			},
			{
				"unexpected HandlesEvent route",
				`integration is invalid: unexpected handles-event route for pkg.SomeEventType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(HandlesEventRouteType),
						MessageTypeName: optional.Some("pkg.SomeEventType"),
					},
				},
			},
			{
				"unexpected SchedulesTimeout route",
				`integration is invalid: unexpected schedules-timeout route for pkg.SomeTimeoutType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(SchedulesTimeoutRouteType),
						MessageTypeName: optional.Some("pkg.SomeTimeoutType"),
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := &Integration{
					AsConfigured: IntegrationAsConfigured{
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

func TestIntegration_IsDisabled(t *testing.T) {
	disable := false

	h := &IntegrationMessageHandlerStub{
		ConfigureFunc: func(c dogma.IntegrationConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesCommand[CommandStub[TypeA]](),
			)
			if disable {
				c.Disable()
			}
		},
	}

	cfg := runtimeconfig.FromIntegration(h)

	if cfg.IsDisabled() {
		t.Fatal("did not expect handler to be disabled")
	}

	disable = true
	cfg = runtimeconfig.FromIntegration(h)

	if !cfg.IsDisabled() {
		t.Fatal("expected handler to be disabled")
	}
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

func TestIntegration_render(t *testing.T) {
	cases := []renderTestCase{
		{
			Name:             "complete",
			ExpectDescriptor: `integration:IntegrationMessageHandlerStub`,
			ExpectDetails: multiline(
				`valid integration *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:             "disabled",
			ExpectDescriptor: `integration:IntegrationMessageHandlerStub`,
			ExpectDetails: multiline(
				`disabled valid integration *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
					c.Disable()
				},
			}),
		},
		{
			Name:             "no runtime type information",
			ExpectDescriptor: `integration:SomeIntegration`,
			ExpectDetails: multiline(
				`valid integration pkg.SomeIntegration (runtime type unavailable)`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-command route for pkg.SomeCommand (runtime type unavailable)`,
			),
			Component: configbuilder.Integration(func(b *configbuilder.IntegrationBuilder) {
				b.SetSourceTypeName("pkg.SomeIntegration")
				b.SetDisabled(false)

				b.Identity(func(b *configbuilder.IdentityBuilder) {
					b.SetName("name")
					b.SetKey("19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				})

				b.Route(func(b *configbuilder.RouteBuilder) {
					b.SetRouteType(HandlesCommandRouteType)
					b.SetMessageTypeName("pkg.SomeCommand")
				})
			}),
		},
		{
			Name:             "empty",
			ExpectDescriptor: `integration`,
			ExpectDetails: multiline(
				`incomplete integration`,
				`  - no identity is configured`,
				`  - no "handles-command" routes are configured`,
			),
			Component: &Integration{},
		},
		{
			Name:             "invalid",
			ExpectDescriptor: `integration:IntegrationMessageHandlerStub`,
			ExpectDetails: multiline(
				`invalid integration *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
				`  - no "handles-command" routes are configured`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
			),
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				},
			}),
		},
		{
			Name:             "invalid sub-component",
			ExpectDescriptor: `integration:IntegrationMessageHandlerStub`,
			ExpectDetails: multiline(
				`valid integration *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
				`  - invalid identity name/non-uuid`,
				`      - invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
	}

	runRenderTests(t, cases)
}

func TestIntegration_validation(t *testing.T) {
	cases := []validationTestCase{
		{
			Name:   "valid",
			Expect: ``,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "no runtime type information",
			Expect: ``,
			Component: &Integration{
				AsConfigured: IntegrationAsConfigured{
					Source: Value[dogma.IntegrationMessageHandler]{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
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
								RouteType:       optional.Some(HandlesCommandRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
					},
				},
			},
		},
		{
			Name: "no runtime type information using WithRuntimeTypes() option",
			Expect: `integration:IntegrationMessageHandlerStub is invalid:` +
				"\n" + `- dogma.IntegrationMessageHandler value is not available` +
				"\n" + `- route:handles-command:CommandStub[TypeA] is invalid: message.Type value is not available`,
			Options: []NormalizeOption{
				WithRuntimeTypes(),
			},
			Component: &Integration{
				AsConfigured: IntegrationAsConfigured{
					Source: Value[dogma.IntegrationMessageHandler]{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
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
								RouteType:       optional.Some(HandlesCommandRouteType),
								MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							},
						},
					},
				},
			},
		},
		{
			Name: "nil integration",
			Expect: `integration is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- no "handles-command" routes are configured` +
				"\n" + `- could not evaluate entire configuration`,
			Component: runtimeconfig.FromIntegration(nil),
		},
		{
			Name: "unconfigured integration",
			Expect: `integration:IntegrationMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- no "handles-command" routes are configured`,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{}),
		},
		{
			Name:   "integration identity must be valid",
			Expect: `integration:IntegrationMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "integration must not have multiple identities",
			Expect: `integration:IntegrationMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "integration must handle at least one command type",
			Expect: `integration:IntegrationMessageHandlerStub is invalid: no "handles-command" routes are configured`,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					// <-- MISSING "handles-command" ROUTE
				},
			}),
		},
		{
			Name:   "integration must not have multiple routes for the same command type",
			Expect: `integration:IntegrationMessageHandlerStub is invalid: multiple "handles-command" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			}),
		},
		{
			Name:   "integration must not have multiple routes for the same event type",
			Expect: `integration:IntegrationMessageHandlerStub is invalid: multiple "records-event" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
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
