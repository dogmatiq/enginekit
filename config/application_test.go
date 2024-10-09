package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

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
			`application(?) is configured without an identity, Identity() must be called exactly once within Configure()`,
			nil,
		},
		{
			"unconfigured application",
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) is configured without an identity, Identity() must be called exactly once within Configure()`,
			&ApplicationStub{},
		},
		{
			"application must not have multiple identities",
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) is configured with multiple identities (foo/63bd2756-2397-4cae-b33b-96e809b384d8, foo/ee316cdb-894c-454e-91dd-ec0cc4531c42 and bar/ee316cdb-894c-454e-91dd-ec0cc4531c42), Identity() must be called exactly once within Configure()`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
				},
			},
		},
		{
			"application identity must be valid",
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) is configured with an invalid identity (name/"non-uuid"): invalid identity key ("non-uuid"): keys must be RFC 4122/9562 UUIDs: invalid UUID format, expected 36 characters`,
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("name", "non-uuid")
				},
			},
		},
		{
			"application must not contain invalid handlers",
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) contains an invalid handler: aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) must have at least one "HandlesCommand" route`,
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
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) and aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) have the same identity (app/14769f7f-87fe-48dd-916e-5bcab6ba6aca), which is not allowed`,
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
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) and aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) have the same identity key (14769f7f-87fe-48dd-916e-5bcab6ba6aca), which is not allowed`,
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
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub), integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) and projection(*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub) have the same identity (handler/4f2a6c38-0651-4ca5-b6a1-1edf4b2624db), which is not allowed`,
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
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have the same identity name (handler), which is not allowed`,
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
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have the same identity key (4f2a6c38-0651-4ca5-b6a1-1edf4b2624db), which is not allowed`,
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
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have "HandlesCommand" routes for the same command type (github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]), which is not allowed`,
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
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have "RecordsEvent" routes for the same event type (github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]), which is not allowed`,
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
				t.Fatal()
			}
		})
	}
}
