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

func TestProjection_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		h := &ProjectionMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A")
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
			Handler dogma.ProjectionMessageHandler
		}{
			{
				"no identity",
				`projection:ProjectionMessageHandlerStub is invalid: no identity is configured`,
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
				`projection:ProjectionMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
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
				`projection:ProjectionMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
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

func TestProjection_RouteSet(t *testing.T) {
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
			cfg.RouteSet().MessageTypes(),
			map[message.Type]RouteDirection{
				message.TypeFor[EventStub[TypeA]](): InboundDirection,
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
				`projection is invalid: route is invalid: could not evaluate entire configuration`,
				&Route{},
			},
			{
				"",
				`projection is invalid: unexpected handles-command route for pkg.SomeCommandType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(HandlesCommandRouteType),
						MessageTypeName: optional.Some("pkg.SomeCommandType"),
					},
				},
			},
			{
				"",
				`projection is invalid: unexpected executes-command route for pkg.SomeCommandType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(ExecutesCommandRouteType),
						MessageTypeName: optional.Some("pkg.SomeCommandType"),
					},
				},
			},
			{
				"",
				`projection is invalid: unexpected records-event route for pkg.SomeEventType`,
				&Route{
					AsConfigured: RouteAsConfigured{
						RouteType:       optional.Some(RecordsEventRouteType),
						MessageTypeName: optional.Some("pkg.SomeEventType"),
					},
				},
			},
			{
				"",
				`projection is invalid: unexpected schedules-timeout route for pkg.SomeTimeoutType`,
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
				cfg := &Projection{
					AsConfigured: ProjectionAsConfigured{
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

func TestProjection_IsDisabled(t *testing.T) {
	disable := false

	h := &ProjectionMessageHandlerStub{
		ConfigureFunc: func(c dogma.ProjectionConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesEvent[EventStub[TypeA]](),
			)
			if disable {
				c.Disable()
			}
		},
	}

	cfg := runtimeconfig.FromProjection(h)

	if cfg.IsDisabled() {
		t.Fatal("did not expect handler to be disabled")
	}

	disable = true
	cfg = runtimeconfig.FromProjection(h)

	if !cfg.IsDisabled() {
		t.Fatal("expected handler to be disabled")
	}
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

func TestProjection_render(t *testing.T) {
	cases := []renderTestCase{
		{
			Name:             "complete",
			ExpectDescriptor: `projection:ProjectionMessageHandlerStub`,
			ExpectDetails: multiline(
				`valid projection *github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - broadcast delivery policy`,
			),
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
					c.DeliveryPolicy(
						dogma.BroadcastProjectionDeliveryPolicy{
							PrimaryFirst: true,
						},
					)
				},
			}),
		},
		{
			Name:             "disabled",
			ExpectDescriptor: `projection:ProjectionMessageHandlerStub`,
			ExpectDetails: multiline(
				`disabled valid projection *github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - unicast delivery policy`,
			),
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
					c.Disable()
				},
			}),
		},
		{
			Name:             "no runtime type information",
			ExpectDescriptor: `projection:SomeProjection`,
			ExpectDetails: multiline(
				`valid projection pkg.SomeProjection (runtime type unavailable)`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid handles-event route for pkg.SomeEvent (runtime type unavailable)`,
				`  - broadcast delivery policy (runtime type unavailable)`,
			),
			Component: configbuilder.Projection(func(b *configbuilder.ProjectionBuilder) {
				b.SetSourceTypeName("pkg.SomeProjection")
				b.SetDisabled(false)
				b.SetDeliveryPolicyTypeName("github.com/dogmatiq/dogma.BroadcastProjectionDeliveryPolicy")

				b.Identity(func(b *configbuilder.IdentityBuilder) {
					b.SetName("name")
					b.SetKey("19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				})

				b.Route(func(b *configbuilder.RouteBuilder) {
					b.SetRouteType(HandlesEventRouteType)
					b.SetMessageTypeName("pkg.SomeEvent")
				})
			}),
		},
		{
			Name:             "empty",
			ExpectDescriptor: `projection`,
			ExpectDetails: multiline(
				`incomplete projection`,
				`  - no identity is configured`,
				`  - no "handles-event" routes are configured`,
			),
			Component: &Projection{},
		},
		{
			Name:             "invalid",
			ExpectDescriptor: `projection:ProjectionMessageHandlerStub`,
			ExpectDetails: multiline(
				`invalid projection *github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub`,
				`  - no "handles-event" routes are configured`,
				`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - unicast delivery policy`,
			),
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				},
			}),
		},
		{
			Name:             "invalid sub-component",
			ExpectDescriptor: `projection:ProjectionMessageHandlerStub`,
			ExpectDetails: multiline(
				`valid projection *github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub`,
				`  - invalid identity name/non-uuid`,
				`      - invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				`  - valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - unicast delivery policy`,
			),
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
	}

	runRenderTests(t, cases)
}

func TestProjection_validation(t *testing.T) {
	cases := []validationTestCase{
		{
			Name:   "valid",
			Expect: ``,
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "no runtime type information",
			Expect: ``,
			Component: &Projection{
				AsConfigured: ProjectionAsConfigured{
					Source: Value[dogma.ProjectionMessageHandler]{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub"),
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
					},
					DeliveryPolicy: optional.Some(Value[dogma.ProjectionDeliveryPolicy]{
						TypeName: optional.Some("github.com/dogmatiq/dogma.UnicastProjectionDeliveryPolicy"),
					}),
				},
			},
		},
		{
			Name: "no runtime type information using WithRuntimeTypes() option",
			Expect: `projection:ProjectionMessageHandlerStub is invalid:` +
				"\n" + `- dogma.ProjectionMessageHandler value is not available` +
				"\n" + `- dogma.ProjectionDeliveryPolicy value is not available` +
				"\n" + `- route:handles-event:EventStub[TypeA] is invalid: message.Type value is not available`,
			Options: []NormalizeOption{
				WithRuntimeTypes(),
			},
			Component: &Projection{
				AsConfigured: ProjectionAsConfigured{
					Source: Value[dogma.ProjectionMessageHandler]{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub"),
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
					},
					DeliveryPolicy: optional.Some(Value[dogma.ProjectionDeliveryPolicy]{
						TypeName: optional.Some("github.com/dogmatiq/dogma.UnicastProjectionDeliveryPolicy"),
					}),
				},
			},
		},
		{
			Name: "nil projection",
			Expect: `projection is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- no "handles-event" routes are configured` +
				"\n" + `- could not evaluate entire configuration`,
			Component: runtimeconfig.FromProjection(nil),
		},
		{
			Name: "unconfigured projection",
			Expect: `projection:ProjectionMessageHandlerStub is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- no "handles-event" routes are configured`,
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{}),
		},
		{
			Name:   "projection identity must be valid",
			Expect: `projection:ProjectionMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("name", "non-uuid")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "projection must not have multiple identities",
			Expect: `projection:ProjectionMessageHandlerStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
				},
			}),
		},
		{
			Name:   "projection must handle at least one event type",
			Expect: `projection:ProjectionMessageHandlerStub is invalid: no "handles-event" routes are configured`,
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					// <-- MISSING "handles-event" ROUTE
				},
			}),
		},
		{
			Name:   "projection must not have multiple routes for the same event type",
			Expect: `projection:ProjectionMessageHandlerStub is invalid: multiple "handles-event" routes are configured for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: runtimeconfig.FromProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
						dogma.HandlesEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
					)
				},
			}),
		},
	}

	runValidationTests(t, cases)
}
