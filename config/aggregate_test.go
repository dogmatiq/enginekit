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

func TestAggregate(t *testing.T) {
	testHandler(
		t,
		configbuilder.Aggregate,
		runtimeconfig.FromAggregate,
		func(fn func(dogma.AggregateConfigurer)) dogma.AggregateMessageHandler {
			return &AggregateMessageHandlerStub{ConfigureFunc: fn}
		},
	)

	testValidate(
		t,
		validationTestCases{
			{
				Name:  "valid",
				Error: ``,
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
				Name:  "no runtime values",
				Error: ``,
				Component: configbuilder.Aggregate(
					func(b *configbuilder.AggregateBuilder) {
						b.TypeName("pkg.SomeAggregate")
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
								b.RouteType(RecordsEventRouteType)
								b.MessageTypeName("pkg.SomeEvent")
							},
						)
					},
				),
			},
			{
				Name: "no runtime values with ForExecution() option",
				Error: multiline(
					`aggregate:SomeAggregate is invalid:`,
					`  - dogma.AggregateMessageHandler value is unavailable`,
					`  - route:handles-command:SomeCommand is invalid:`,
					`      - message type ID is unavailable`,
					`      - message type is unavailable`,
					`  - route:records-event:SomeEvent is invalid:`,
					`      - message type ID is unavailable`,
					`      - message type is unavailable`,
				),
				Options: []ValidateOption{
					ForExecution(),
				},
				Component: configbuilder.Aggregate(
					func(b *configbuilder.AggregateBuilder) {
						b.TypeName("pkg.SomeAggregate")
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
								b.RouteType(RecordsEventRouteType)
								b.MessageTypeName("pkg.SomeEvent")
							},
						)
					},
				),
			},
			{
				Name: "nil aggregate",
				Error: multiline(
					`aggregate is invalid:`,
					`  - could not evaluate entire configuration`,
					`  - no identity`,
					`  - no handles-command routes`,
					`  - no records-event routes`,
				),
				Component: runtimeconfig.FromAggregate(nil),
			},
			{
				Name: "unconfigured aggregate",
				Error: multiline(
					`aggregate:AggregateMessageHandlerStub is invalid:`,
					`  - no identity`,
					`  - no handles-command routes`,
					`  - no records-event routes`,
				),
				Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{}),
			},
			{
				Name:  "aggregate identity must be valid",
				Error: `aggregate:AggregateMessageHandlerStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
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
				Name:  "aggregate must not have multiple identities",
				Error: `aggregate:AggregateMessageHandlerStub is invalid: multiple identities`,
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
				Name:  "aggregate must handle at least one command type",
				Error: `aggregate:AggregateMessageHandlerStub is invalid: no handles-command routes`,
				Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("handler", "d1e04684-ec56-44a7-8c7d-f111b2d6b2d2")
						c.Routes(
							// <-- MISSING HandlesCommand ROUTE
							dogma.RecordsEvent[EventStub[TypeA]](),
						)
					},
				}),
			},
			{
				Name:  "aggregate must record at least one event type",
				Error: `aggregate:AggregateMessageHandlerStub is invalid: no records-event routes`,
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
				Name:  "aggregate must not have multiple routes for the same command type",
				Error: `aggregate:AggregateMessageHandlerStub is invalid: multiple handles-command routes for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
				Name:  "aggregate must not have multiple routes for the same event type",
				Error: `aggregate:AggregateMessageHandlerStub is invalid: multiple records-event routes for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
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
			{
				Name:  "unsupported route type",
				Error: `aggregate:SomeAggregate is invalid: unsupported schedules-timeout route for pkg.SomeTimeout`,
				Component: configbuilder.Aggregate(
					func(b *configbuilder.AggregateBuilder) {
						b.TypeName("pkg.SomeAggregate")
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
								b.RouteType(RecordsEventRouteType)
								b.MessageTypeName("pkg.SomeEvent")
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
				String: `aggregate:AggregateMessageHandlerStub`,
				Description: multiline(
					`valid aggregate *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
					`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					`  - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				),
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
				Name:   "disabled",
				String: `aggregate:AggregateMessageHandlerStub`,
				Description: multiline(
					`valid aggregate *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
					`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					`  - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					`  - valid disabled flag, set to true`,
				),
				Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeA]](),
						)
						c.Disable()
					},
				}),
			},
			{
				Name:   "no runtime values",
				String: `aggregate:SomeAggregate`,
				Description: multiline(
					`valid aggregate pkg.SomeAggregate (value unavailable)`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
					`  - valid handles-command route for pkg.SomeCommand (type unavailable)`,
					`  - valid records-event route for pkg.SomeEvent (type unavailable)`,
				),
				Component: configbuilder.Aggregate(func(b *configbuilder.AggregateBuilder) {
					b.TypeName("pkg.SomeAggregate")
					b.Identity(func(b *configbuilder.IdentityBuilder) {
						b.Name("name")
						b.Key("19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					})
					b.Route(func(b *configbuilder.RouteBuilder) {
						b.RouteType(HandlesCommandRouteType)
						b.MessageTypeName("pkg.SomeCommand")
					})
					b.Route(func(b *configbuilder.RouteBuilder) {
						b.RouteType(RecordsEventRouteType)
						b.MessageTypeName("pkg.SomeEvent")
					})
				}),
			},
			{
				Name:   "empty",
				String: `aggregate`,
				Description: multiline(
					`invalid aggregate`,
					`  - no identity`,
					`  - no handles-command routes`,
					`  - no records-event routes`,
				),
				Component: &Aggregate{},
			},
			{
				Name:   "invalid",
				String: `aggregate:AggregateMessageHandlerStub`,
				Description: multiline(
					`invalid aggregate *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
					`  - no handles-command routes`,
					`  - no records-event routes`,
					`  - valid identity name/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				),
				Component: runtimeconfig.FromAggregate(&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					},
				}),
			},
			{
				Name:   "invalid sub-component",
				String: `aggregate:AggregateMessageHandlerStub`,
				Description: multiline(
					`valid aggregate *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
					`  - invalid identity name/non-uuid`,
					`      - invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
					`  - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					`  - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				),
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
		},
	)

	t.Run("func RouteSet()", func(t *testing.T) {
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

			handler := runtimeconfig.FromAggregate(h)

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
					"unsupported ExecutesCommand route",
					`unsupported executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					&Route{
						RouteType:       optional.Some(ExecutesCommandRouteType),
						MessageTypeID:   optional.Some(MessageTypeID[CommandStub[TypeA]]()),
						MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
						MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
					},
				},
				{
					"unsupported HandlesEvent route",
					`unsupported handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					&Route{
						RouteType:       optional.Some(HandlesEventRouteType),
						MessageTypeID:   optional.Some(MessageTypeID[EventStub[TypeA]]()),
						MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
						MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
					},
				},
				{
					"unsupported SchedulesTimeout route",
					`unsupported schedules-timeout route for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
					&Route{
						RouteType:       optional.Some(SchedulesTimeoutRouteType),
						MessageTypeID:   optional.Some(MessageTypeID[TimeoutStub[TypeA]]()),
						MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
						MessageType:     optional.Some(message.TypeFor[TimeoutStub[TypeA]]()),
					},
				},
			}

			for _, c := range cases {
				t.Run(c.Name, func(t *testing.T) {
					cfg := &Aggregate{
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
