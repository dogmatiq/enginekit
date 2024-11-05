package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

func TestIntegration(t *testing.T) {
	testHandler(
		t,
		configbuilder.Integration,
		runtimeconfig.FromIntegration,
		func(fn func(dogma.IntegrationConfigurer)) dogma.IntegrationMessageHandler {
			return &IntegrationMessageHandlerStub{ConfigureFunc: fn}
		},
	)

	testValidate(
		t,
		validationTestCases{
			{
				Name:  "valid",
				Error: ``,
				Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeA]](),
						)
					},
				}),
			},
			{
				Name:  "no runtime values",
				Error: ``,
				Component: configbuilder.Integration(
					func(b *configbuilder.IntegrationBuilder) {
						b.TypeName("pkg.SomeIntegration")
						b.Identity(
							func(b *configbuilder.IdentityBuilder) {
								b.Name("name")
								b.Key("494157ef-6f91-45ec-ab19-df61bb96210a")
							},
						)
						b.Route(
							func(b *configbuilder.RouteBuilder) {
								b.RouteType(HandlesCommandRouteType)
								b.MessageTypeName("pkg.SomeCommand")
							},
						)
					},
				),
			},
			{
				Name: "no runtime values with ForExecution() option",
				Error: multiline(
					`integration:SomeIntegration is invalid:`,
					`  - dogma.IntegrationMessageHandler value is unavailable`,
					`  - route:handles-command:SomeCommand is invalid: message type is unavailable`,
				),
				Options: []ValidateOption{
					ForExecution(),
				},
				Component: configbuilder.Integration(
					func(b *configbuilder.IntegrationBuilder) {
						b.TypeName("pkg.SomeIntegration")
						b.Identity(
							func(b *configbuilder.IdentityBuilder) {
								b.Name("name")
								b.Key("494157ef-6f91-45ec-ab19-df61bb96210a")
							},
						)
						b.Route(
							func(b *configbuilder.RouteBuilder) {
								b.RouteType(HandlesCommandRouteType)
								b.MessageTypeName("pkg.SomeCommand")
							},
						)
					},
				),
			},
			{
				Name: "nil integration",
				Error: multiline(
					`integration is invalid:`,
					`  - could not evaluate entire configuration: handler is nil`,
					`  - no identity`,
					`  - no handles-command routes`,
				),
				Component: runtimeconfig.FromIntegration(nil),
			},
			{
				Name: "unconfigured integration",
				Error: multiline(
					`integration:IntegrationMessageHandlerStub is invalid:`,
					`  - no identity`,
					`  - no handles-command routes`,
				),
				Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{}),
			},
			{
				Name:  "integration identity must be valid",
				Error: `integration:IntegrationMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
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
				Name:  "integration must not have multiple identities",
				Error: `integration:IntegrationMessageHandlerStub is invalid: multiple identities`,
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
				Name:  "integration must handle at least one command type",
				Error: `integration:IntegrationMessageHandlerStub is invalid: no handles-command routes`,
				Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
						// <-- MISSING HandlesCommand ROUTE
					},
				}),
			},
			{
				Name:  "integration must not have multiple routes for the same command type",
				Error: `integration:IntegrationMessageHandlerStub is invalid: multiple handles-command routes for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
				Name:  "integration must not have multiple routes for the same event type",
				Error: `integration:IntegrationMessageHandlerStub is invalid: multiple records-event routes for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
			{
				Name:  "unsupported route type",
				Error: `integration:SomeIntegration is invalid: unsupported schedules-timeout route for pkg.SomeTimeout`,
				Component: configbuilder.Integration(
					func(b *configbuilder.IntegrationBuilder) {
						b.TypeName("pkg.SomeIntegration")
						b.Identity(
							func(b *configbuilder.IdentityBuilder) {
								b.Name("name")
								b.Key("494157ef-6f91-45ec-ab19-df61bb96210a")
							},
						)
						b.Route(
							func(b *configbuilder.RouteBuilder) {
								b.RouteType(HandlesCommandRouteType)
								b.MessageTypeName("pkg.SomeCommand")
							},
						)
						b.Route(
							func(b *configbuilder.RouteBuilder) {
								b.RouteType(SchedulesTimeoutRouteType) // <-- UNSUPPORTED ROUTE TYPE
								b.MessageTypeName("pkg.SomeTimeout")
							},
						)
					},
				),
			},
		},
	)

	testDescribe(
		t,
		describeTestCases{
			{
				Name:   "complete",
				String: `integration:IntegrationMessageHandlerStub`,
				Description: multiline(
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
				Name:   "disabled",
				String: `integration:IntegrationMessageHandlerStub`,
				Description: multiline(
					`valid integration *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
					`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					`  - valid disabled flag, set to true`,
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
				Name:   "no runtime values",
				String: `integration:SomeIntegration`,
				Description: multiline(
					`valid integration pkg.SomeIntegration (value unavailable)`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
					`  - valid handles-command route for pkg.SomeCommand (type unavailable)`,
				),
				Component: configbuilder.Integration(func(b *configbuilder.IntegrationBuilder) {
					b.TypeName("pkg.SomeIntegration")
					b.Identity(func(b *configbuilder.IdentityBuilder) {
						b.Name("name")
						b.Key("19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					})
					b.Route(func(b *configbuilder.RouteBuilder) {
						b.RouteType(HandlesCommandRouteType)
						b.MessageTypeName("pkg.SomeCommand")
					})
				}),
			},
			{
				Name:   "empty",
				String: `integration`,
				Description: multiline(
					`invalid integration`,
					`  - no identity`,
					`  - no handles-command routes`,
				),
				Component: &Integration{},
			},
			{
				Name:   "invalid",
				String: `integration:IntegrationMessageHandlerStub`,
				Description: multiline(
					`invalid integration *github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
					`  - no handles-command routes`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				),
				Component: runtimeconfig.FromIntegration(&IntegrationMessageHandlerStub{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					},
				}),
			},
			{
				Name:   "invalid sub-component",
				String: `integration:IntegrationMessageHandlerStub`,
				Description: multiline(
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
		},
	)

	t.Run("func RouteSet()", func(t *testing.T) {
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

			handler := runtimeconfig.FromIntegration(h)

			test.Expect(
				t,
				"unexpected routes",
				handler.RouteSet().MessageTypes(),
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
					`route is invalid: route type is unavailable`,
					&Route{},
				},
				{
					"unexpected ExecutesCommand route",
					`unsupported executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					&Route{
						RouteType:       optional.Some(ExecutesCommandRouteType),
						MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
						MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
					},
				},
				{
					"unexpected HandlesEvent route",
					`unsupported handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					&Route{
						RouteType:       optional.Some(HandlesEventRouteType),
						MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
						MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
					},
				},
				{
					"unexpected SchedulesTimeout route",
					`unsupported schedules-timeout route for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					&Route{
						RouteType:       optional.Some(SchedulesTimeoutRouteType),
						MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
						MessageType:     optional.Some(message.TypeFor[TimeoutStub[TypeA]]()),
					},
				},
			}

			for _, c := range cases {
				t.Run(c.Name, func(t *testing.T) {
					cfg := &Integration{
						HandlerCommon: HandlerCommon{
							RouteComponents: []*Route{c.Route},
						},
					}

					test.ExpectPanic(
						t,
						c.Want,
						func() {
							cfg.RouteSet()
						},
					)
				})
			}
		})
	})
}
