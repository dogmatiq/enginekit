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
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func TestApplication_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		app := &ApplicationStub{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A")
			},
		}

		cfg := runtimeconfig.FromApplication(app)

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
			Name string
			Want string
			App  dogma.Application
		}{
			{
				"no identity",
				`application:ApplicationStub is invalid: no identity is configured`,
				&ApplicationStub{},
			},
			{
				"invalid identity",
				`application:ApplicationStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				&ApplicationStub{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("name", "non-uuid")
					},
				},
			},
			{
				"multiple identities",
				`application:ApplicationStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
				&ApplicationStub{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
						c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := runtimeconfig.FromApplication(c.App)

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

func TestApplication_RouteSet(t *testing.T) {
	t.Run("it returns the normalized routes", func(t *testing.T) {
		app := &ApplicationStub{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("app", "04c64a99-2b48-4cd5-a62a-85c6cb1d5e35")
				c.RegisterAggregate(&AggregateMessageHandlerStub{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("aggregate", "6a006c20-075f-4706-8230-4188b42b60aa")
						c.Routes(
							dogma.HandlesCommand[CommandStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeA]](),
							dogma.RecordsEvent[EventStub[TypeB]](),
						)
					},
				})
				c.RegisterProcess(&ProcessMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("process1", "3614c386-4d8d-4a1d-88fa-10f94313c803")
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeB]](),
							dogma.ExecutesCommand[CommandStub[TypeB]](),
							dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
						)
					},
				})
			},
		}

		cfg := runtimeconfig.FromApplication(app)

		Expect(
			t,
			"unexpected routes",
			cfg.RouteSet().MessageTypes(),
			map[message.Type]RouteDirection{
				message.TypeFor[CommandStub[TypeA]](): InboundDirection,
				message.TypeFor[EventStub[TypeA]]():   OutboundDirection,
				message.TypeFor[EventStub[TypeB]]():   InboundDirection | OutboundDirection,
				message.TypeFor[CommandStub[TypeB]](): OutboundDirection,
				message.TypeFor[TimeoutStub[TypeA]](): InboundDirection | OutboundDirection,
			},
		)
	})

	t.Run("it panics if the routes are invalid", func(t *testing.T) {
		cfg := &Application{
			X: XApplication{
				Handlers: []Handler{
					&Projection{
						X: XProjection{
							Routes: []*Route{
								{},
							},
						},
					},
				},
			},
		}

		ExpectPanic(
			t,
			`application is invalid: projection is invalid: route is invalid: could not evaluate entire configuration`,
			func() {
				cfg.RouteSet()
			},
		)
	})
}

func TestApplication_Interface(t *testing.T) {
	app := &ApplicationStub{
		ConfigureFunc: func(c dogma.ApplicationConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
		},
	}

	cfg := runtimeconfig.FromApplication(app)

	Expect(
		t,
		"unexpected result",
		cfg.Interface(),
		app,
	)
}

func TestApplication_HandlerByName(t *testing.T) {
	h := &AggregateMessageHandlerStub{
		ConfigureFunc: func(c dogma.AggregateConfigurer) {
			c.Identity("name", "40ddf2a2-f053-485c-8621-1fc8a58f8ddf")
			c.Routes(
				dogma.HandlesCommand[CommandStub[TypeA]](),
				dogma.RecordsEvent[EventStub[TypeA]](),
			)
		},
	}

	app := &ApplicationStub{
		ConfigureFunc: func(c dogma.ApplicationConfigurer) {
			c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
			c.RegisterAggregate(h)
		},
	}

	cfg := runtimeconfig.FromApplication(app)

	if got, ok := cfg.HandlerByName("name"); ok {
		want := runtimeconfig.FromAggregate(h)

		Expect(
			t,
			"unexpected handler",
			got,
			want,
		)
	} else {
		t.Fatal("expected handler to be found")
	}

	if _, ok := cfg.HandlerByName("unknown"); ok {
		t.Fatal("did not expect handler to be found")
	}
}

func TestApplication_render(t *testing.T) {
	cases := []renderTestCase{
		{
			Name:             "complete",
			ExpectDescriptor: `application:ApplicationStub`,
			ExpectDetails: multiline(
				`valid application *github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub`,
				`  - valid identity app/c85acb36-e47b-4ef6-b46b-64d847a853b7`,
				`  - valid aggregate *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
				`      - valid identity aggregate/8bb5eaf2-6b36-42bd-a1b3-90c27c9c80d4`,
				`      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`      - valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "c85acb36-e47b-4ef6-b46b-64d847a853b7")

					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "8bb5eaf2-6b36-42bd-a1b3-90c27c9c80d4")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:             "no runtime type information",
			ExpectDescriptor: `application:SomeApplication`,
			ExpectDetails: multiline(
				`valid application pkg.SomeApplication (runtime type unavailable)`,
				`  - valid identity app/c85acb36-e47b-4ef6-b46b-64d847a853b7`,
			),
			Component: configbuilder.Application(
				func(b *configbuilder.ApplicationBuilder) {
					b.SetSourceTypeName("pkg.SomeApplication")
					b.Identity(func(b *configbuilder.IdentityBuilder) {
						b.SetName("app")
						b.SetKey("c85acb36-e47b-4ef6-b46b-64d847a853b7")
					})
				},
			),
		},
		{
			Name:             "empty",
			ExpectDescriptor: `application`,
			ExpectDetails: multiline(
				`incomplete application`,
				`  - no identity is configured`,
			),
			Component: &Application{},
		},
		{
			Name:             "invalid",
			ExpectDescriptor: `application:ApplicationStub`,
			ExpectDetails: multiline(
				`invalid application *github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub`,
				`  - multiple identities are configured: identity:app/19cb98d5-dd17-4daf-ae00-1b413b7b899a and identity:app/89864744-89c5-4a80-a2bf-90aaebc467ba`,
				`  - valid identity app/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - valid identity app/89864744-89c5-4a80-a2bf-90aaebc467ba`,
			),
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
					c.Identity("app", "89864744-89c5-4a80-a2bf-90aaebc467ba")
				},
			}),
		},
		{
			Name:             "invalid sub-component",
			ExpectDescriptor: `application:ApplicationStub`,
			ExpectDetails: multiline(
				`valid application *github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub`,
				`  - valid identity app/19cb98d5-dd17-4daf-ae00-1b413b7b899a`,
				`  - invalid aggregate *github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
				`      - no "records-event" routes are configured`,
				`      - valid identity aggregate/8bb5eaf2-6b36-42bd-a1b3-90c27c9c80d4`,
				`      - valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			),
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")

					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "8bb5eaf2-6b36-42bd-a1b3-90c27c9c80d4")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
							)
						},
					})
				},
			}),
		},
	}

	runRenderTests(t, cases)
}

func TestApplication_validation(t *testing.T) {
	cases := []validationTestCase{
		{
			Name:   "application name may be shared with one of its handlers",
			Expect: "",
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("name", "14769f7f-87fe-48dd-916e-5bcab6ba6aca") // <-- SAME NAME
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("name", "40ddf2a2-f053-485c-8621-1fc8a58f8ddf") // <-- SAME NAME
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:   "multiple processes may schedule the same type of timeout message",
			Expect: "",
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("name", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterProcess(&ProcessMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProcessConfigurer) {
							c.Identity("process1", "3614c386-4d8d-4a1d-88fa-10f94313c803")
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
								dogma.ExecutesCommand[CommandStub[TypeA]](),
								dogma.SchedulesTimeout[TimeoutStub[TypeA]](), // <-- SAME MESSAGE TYPE
							)
						},
					})
					c.RegisterProcess(&ProcessMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProcessConfigurer) {
							c.Identity("process2", "f2c9acdd-93a8-4ca0-9014-56aaae16a3af")
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
								dogma.ExecutesCommand[CommandStub[TypeA]](),
								dogma.SchedulesTimeout[TimeoutStub[TypeA]](), // <-- SAME MESSAGE TYPE
							)
						},
					})
				},
			}),
		},
		{
			Name: "nil application",
			Expect: `application is invalid:` +
				"\n" + `- no identity is configured` +
				"\n" + `- could not evaluate entire configuration`,
			Component: runtimeconfig.FromApplication(nil),
		},
		{
			Name:      "unconfigured application",
			Expect:    `application:ApplicationStub is invalid: no identity is configured`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{}),
		},
		{
			Name:   "application identity must be valid",
			Expect: `application:ApplicationStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("name", "non-uuid")
				},
			}),
		},
		{
			Name:   "application must not have multiple identities",
			Expect: `application:ApplicationStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
				},
			}),
		},
		{
			Name:   "application must not contain invalid handlers",
			Expect: `application:ApplicationStub is invalid: aggregate:AggregateMessageHandlerStub is invalid: no "handles-command" routes are configured`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "fe78acbf-dfd4-490a-bf99-93b6acf9f891")
							c.Routes(
								// <-- MISSING "handles-command" ROUTE
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:   "application must not have the same identity key as one of its handlers",
			Expect: `application:ApplicationStub is invalid: entities have conflicting identities: the "14769f7f-87fe-48dd-916e-5bcab6ba6aca" key is shared by application:ApplicationStub and aggregate:AggregateMessageHandlerStub`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca") // <-- SAME IDENTITY KEY
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "14769F7F-87FE-48DD-916E-5BCAB6BA6ACA") // <-- SAME IDENTITY KEY (note: non-canonical UUID)
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:   "multiple handlers must not have the same identity name",
			Expect: `application:ApplicationStub is invalid: entities have conflicting identities: the "handler" name is shared by aggregate:AggregateMessageHandlerStub and integration:IntegrationMessageHandlerStub`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("handler", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db") // <-- SAME IDENTITY NAME
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("handler", "300a00e7-9d8f-4859-b67a-7eb36c7e3d83") // <-- SAME IDENTITY NAME
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeB]](),
								dogma.RecordsEvent[EventStub[TypeB]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:   "multiple handlers must not have the same identity key",
			Expect: `application:ApplicationStub is invalid: entities have conflicting identities: the "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db" key is shared by aggregate:AggregateMessageHandlerStub and integration:IntegrationMessageHandlerStub`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db") // <-- SAME IDENTITY KEY
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("integration", "4F2A6C38-0651-4CA5-B6A1-1EDF4B2624DB") // <-- SAME IDENTITY NAME (note: non-canonical UUID)
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeB]](),
								dogma.RecordsEvent[EventStub[TypeB]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:   "multiple handlers must not handle the same command type",
			Expect: `application:ApplicationStub is invalid: handlers have conflicting "handles-command" routes: github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] is handled by aggregate:AggregateMessageHandlerStub and integration:IntegrationMessageHandlerStub`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("integration", "1228c8a5-ad60-4e59-81b3-e236c31f12e2")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](), // <-- SAME MESSAGE TYPE
								dogma.RecordsEvent[EventStub[TypeB]](),
							)
						},
					})
				},
			}),
		},
		{
			Name:   "multiple handlers must not record the same event type",
			Expect: `application:ApplicationStub is invalid: handlers have conflicting "records-event" routes: github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] is recorded by aggregate:AggregateMessageHandlerStub and integration:IntegrationMessageHandlerStub`,
			Component: runtimeconfig.FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("integration", "1228c8a5-ad60-4e59-81b3-e236c31f12e2")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeB]](),
								dogma.RecordsEvent[EventStub[TypeA]](), // <-- SAME MESSAGE TYPE
							)
						},
					})
				},
			}),
		},
	}

	runValidationTests(t, cases)
}
