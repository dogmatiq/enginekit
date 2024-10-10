package runtimeconfig_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/optional"
)

func TestFromApplication(t *testing.T) {
	var (
		aggregate   dogma.AggregateMessageHandler   = &AggregateMessageHandlerStub{}
		process     dogma.ProcessMessageHandler     = &ProcessMessageHandlerStub{}
		integration dogma.IntegrationMessageHandler = &IntegrationMessageHandlerStub{}
		projection  dogma.ProjectionMessageHandler  = &ProjectionMessageHandlerStub{}
	)

	cases := []struct {
		Name string
		App  dogma.Application
		Want func(app dogma.Application) config.Application
	}{
		{
			"nil application",
			nil,
			func(dogma.Application) config.Application {
				return config.Application{}
			},
		},
		{
			"unconfigured application",
			&ApplicationStub{},
			func(app dogma.Application) config.Application {
				return config.Application{
					Impl: optional.Some(
						config.Implementation[dogma.Application]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
							Source:   optional.Some(app),
						},
					),
					ConfigurationIsExhaustive: true,
				}
			},
		},
		{
			"configured application",
			&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("app", "bed53df8-bf22-4502-be4b-64d56532d8be")
					c.RegisterAggregate(aggregate)
					c.RegisterProcess(process)
					c.RegisterIntegration(integration)
					c.RegisterProjection(projection)
				},
			},
			func(app dogma.Application) config.Application {
				return config.Application{
					Impl: optional.Some(
						config.Implementation[dogma.Application]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
							Source:   optional.Some(app),
						},
					),
					ConfiguredIdentities: []config.Identity{
						{
							Name: "app",
							Key:  "bed53df8-bf22-4502-be4b-64d56532d8be",
						},
					},
					ConfiguredHandlers: []config.Handler{
						config.Aggregate{
							Impl: optional.Some(
								config.Implementation[dogma.AggregateMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub",
									Source:   optional.Some(aggregate),
								},
							),
							ConfigurationIsExhaustive: true,
						},
						config.Process{
							Impl: optional.Some(
								config.Implementation[dogma.ProcessMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub",
									Source:   optional.Some(process),
								},
							),
							ConfigurationIsExhaustive: true,
						},
						config.Integration{
							Impl: optional.Some(
								config.Implementation[dogma.IntegrationMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub",
									Source:   optional.Some(integration),
								},
							),
							ConfigurationIsExhaustive: true,
						},
						config.Projection{
							Impl: optional.Some(
								config.Implementation[dogma.ProjectionMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
									Source:   optional.Some(projection),
								},
							),
							ConfigurationIsExhaustive: true,
						},
					},
					ConfigurationIsExhaustive: true,
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			Expect(
				t,
				"unexpected config",
				FromApplication(c.App),
				c.Want(c.App),
			)
		})
	}
}
