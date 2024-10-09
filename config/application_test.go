package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestApplication_validate(t *testing.T) {
	cases := []struct {
		Name string
		App  dogma.Application
		Want string
	}{
		{
			"nil",
			nil,
			`application(?) is configured without an identity, Identity() must be called exactly once within Configure()`,
		},
		{
			"unconfigured",
			&ApplicationStub{},
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) is configured without an identity, Identity() must be called exactly once within Configure()`,
		},
		{
			"multiple identity calls",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("foo", "63bd2756-2397-4cae-b33b-96e809b384d8")
					c.Identity("foo", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
					c.Identity("bar", "ee316cdb-894c-454e-91dd-ec0cc4531c42")
				},
			},
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) is configured with multiple identities (foo/63bd2756-2397-4cae-b33b-96e809b384d8, foo/ee316cdb-894c-454e-91dd-ec0cc4531c42 and bar/ee316cdb-894c-454e-91dd-ec0cc4531c42), Identity() must be called exactly once within Configure()`,
		},
		{
			"invalid identity",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("name", "non-uuid")
				},
			},
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) is configured with an invalid identity (name/"non-uuid"): invalid identity key ("non-uuid"): keys must be RFC 4122/9562 UUIDs: invalid UUID format, expected 36 characters`,
		},
		{
			"handler is invalid",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "non-uuid")
						},
					})
				},
			},
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) contains an invalid handler: aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) is configured with an invalid identity (aggregate/"non-uuid"): invalid identity key ("non-uuid"): keys must be RFC 4122/9562 UUIDs: invalid UUID format, expected 36 characters`,
		},
		{
			"handler has the same identity as the application",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
						},
					})
				},
			},
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) and aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) have the same identity (app/14769f7f-87fe-48dd-916e-5bcab6ba6aca), which is not allowed`,
		},
		{
			"handler has the same identity key as the application",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "14769F7F-87FE-48DD-916E-5BCAB6BA6ACA") // note: non-normalized UUID
						},
					})
				},
			},
			`application(*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub) and aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) have the same identity key (14769f7f-87fe-48dd-916e-5bcab6ba6aca), which is not allowed`,
		},
		{
			"multiple handlers have the same identity",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("handler", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("handler", "4F2A6C38-0651-4CA5-B6A1-1EDF4B2624DB") // note: non-normalized UUID
						},
					})
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("handler", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
						},
					})
				},
			},
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub), integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) and projection(*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub) have the same identity (handler/4f2a6c38-0651-4ca5-b6a1-1edf4b2624db), which is not allowed`,
		},
		{
			"multiple handlers have the same identity name",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("handler", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("handler", "300a00e7-9d8f-4859-b67a-7eb36c7e3d83")
						},
					})
				},
			},
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have the same identity name (handler), which is not allowed`,
		},
		{
			"multiple handlers have the same identity key",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("integration", "4F2A6C38-0651-4CA5-B6A1-1EDF4B2624DB") // note: non-normalized UUID
						},
					})
				},
			},
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have the same identity key (4f2a6c38-0651-4ca5-b6a1-1edf4b2624db), which is not allowed`,
		},
		{
			"multiple handlers have routes to handle the same command type",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](), // <-- PROBLEM
								dogma.RecordsEvent[EventStub[TypeA]](),
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("integration", "1228c8a5-ad60-4e59-81b3-e236c31f12e2")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](), // <-- PROBLEM
								dogma.RecordsEvent[EventStub[TypeB]](),
							)
						},
					})
				},
			},
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have "HandlesCommand" routes for the same command type (github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]), which is not allowed`,
		},
		{
			"multiple handlers have routes to record the same event type",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "14769f7f-87fe-48dd-916e-5bcab6ba6aca")
					c.RegisterAggregate(&AggregateMessageHandlerStub{
						ConfigureFunc: func(c dogma.AggregateConfigurer) {
							c.Identity("aggregate", "4f2a6c38-0651-4ca5-b6a1-1edf4b2624db")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeA]](),
								dogma.RecordsEvent[EventStub[TypeA]](), // <-- PROBLEM
							)
						},
					})
					c.RegisterIntegration(&IntegrationMessageHandlerStub{
						ConfigureFunc: func(c dogma.IntegrationConfigurer) {
							c.Identity("integration", "1228c8a5-ad60-4e59-81b3-e236c31f12e2")
							c.Routes(
								dogma.HandlesCommand[CommandStub[TypeB]](),
								dogma.RecordsEvent[EventStub[TypeA]](), // <-- PROBLEM
							)
						},
					})
				},
			},
			`aggregate(*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub) and integration(*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub) have "RecordsEvent" routes for the same event type (github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]), which is not allowed`,
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
