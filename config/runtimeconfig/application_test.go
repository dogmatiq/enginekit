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
					ConfigurationSource: optional.Some(
						config.Source[dogma.Application]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
							Value:    optional.Some(app),
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
					ConfigurationSource: optional.Some(
						config.Source[dogma.Application]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
							Value:    optional.Some(app),
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
							ConfigurationSource: optional.Some(
								config.Source[dogma.AggregateMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub",
									Value:    optional.Some(aggregate),
								},
							),
							ConfigurationIsExhaustive: true,
						},
						config.Process{
							ConfigurationSource: optional.Some(
								config.Source[dogma.ProcessMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub",
									Value:    optional.Some(process),
								},
							),
							ConfigurationIsExhaustive: true,
						},
						config.Integration{
							ConfigurationSource: optional.Some(
								config.Source[dogma.IntegrationMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub",
									Value:    optional.Some(integration),
								},
							),
							ConfigurationIsExhaustive: true,
						},
						config.Projection{
							ConfigurationSource: optional.Some(
								config.Source[dogma.ProjectionMessageHandler]{
									TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
									Value:    optional.Some(projection),
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