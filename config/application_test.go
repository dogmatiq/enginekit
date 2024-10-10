package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
)

func TestApplication_Identity(t *testing.T) {
	t.Run("it returns the normalized identity", func(t *testing.T) {
		app := &ApplicationStub{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("name", "19CB98D5-DD17-4DAF-AE00-1B413B7B899A") // note: non-canonical UUID
			},
		}

		cfg := runtimeconfig.FromApplication(app)

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
			Name string
			Want string
			App  dogma.Application
		}{
			{
				"no identity",
				`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: no identity is configured`,
				&ApplicationStub{},
			},
			{
				"invalid identity",
				`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
				&ApplicationStub{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("name", "non-uuid")
					},
				},
			},
			{
				"multiple identities",
				`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8`,
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

func TestApplication_validation(t *testing.T) {
	cases := []struct {
		Name string
		Want string
		App  dogma.Application
	}{
		{
			`application name may be shared with one of its handlers`,
			``, // no error
			&ApplicationStub{
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
			},
		},
		{
			`multiple processes may schedule the same type of timeout message`,
			``, // no error
			&ApplicationStub{
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
			},
		},
		{
			"nil application",
			`partial application is invalid: no identity is configured`,
			nil,
		},
		{
			"unconfigured application",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: no identity is configured`,
			&ApplicationStub{},
		},
		{
			"application identity must be valid",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: identity:name/non-uuid is invalid: invalid key ("non-uuid"), expected an RFC 4122/9562 UUID`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("name", "non-uuid")
				},
			},
		},
		{
			"application must not have multiple identities",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: multiple identities are configured: identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8, identity:foo/63bd2756-2397-4cae-b33b-96e809b384d8 and identity:bar/ee316cdb-894c-454e-91dd-ec0cc4531c42`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
				},
			},
		},
		{
			"application must not contain invalid handlers",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub is invalid: expected at least one "HandlesCommand" route`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "fe78acbf-dfd4-490a-bf99-93b6acf9f891")
							c.Routes(
								// <-- MISSING "HandlesCommand" ROUTE
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			},
		},
		{
			"application must not have the same identity as one of its handlers",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: entities have conflicting identities: identity:app/14769f7f-87fe-48dd-916e-5bcab6ba6aca is shared by application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub and aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca") // <-- SAME IDENTITY
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca") // <-- SAME IDENTITY
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			},
		},
		{
			"application must not have the same identity key as one of its handlers",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: entities have conflicting identities: the "14769f7f-87fe-48dd-916e-5bcab6ba6aca" key is shared by application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub and aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub`,
			&ApplicationStub{
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
			},
		},
		{
			"multiple handlers must not have the same identity",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: entities have conflicting identities: identity:handler/4f2a6c38-0651-4ca5-b6a1-1edf4b2624db is shared by aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub, integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub and projection:github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("handler", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db") // <-- SAME IDENTITY
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("handler", "4F2A6C38-0651-4CA5-B6A1-1EDF4B2624DB") // <-- SAME IDENTITY (note: non-canonical UUID)
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeB]](),
								dogma.RecordsEvent[EventStub[TypeB]](),
							)
						},
					})
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("handler", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db") // <-- SAME IDENTITY
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			},
		},
		{
			"multiple handlers must not have the same identity name",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: entities have conflicting identities: the "handler" name is shared by aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub and integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
			&ApplicationStub{
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
			},
		},
		{
			"multiple handlers must not have the same identity key",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: entities have conflicting identities: the "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db" key is shared by aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub and integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
			&ApplicationStub{
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
			},
		},
		{
			"multiple handlers must not handle the same command type",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: handlers have conflicting "HandlesCommand" routes: github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] is handled by aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub and integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
			&ApplicationStub{
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
			},
		},
		{
			"multiple handlers must not record the same event type",
			`application:github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub is invalid: handlers have conflicting "RecordsEvent" routes: github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] is recorded by aggregate:github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub and integration:github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub`,
			&ApplicationStub{
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
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			cfg := runtimeconfig.FromApplication(c.App)

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
