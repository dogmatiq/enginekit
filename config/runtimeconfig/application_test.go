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
		Want func(app dogma.Application) *config.Application
	}{
		{
			"nil application",
			nil,
			func(dogma.Application) *config.Application {
				return &config.Application{}
			},
		},
		{
			"unconfigured application",
			&ApplicationStub{},
			func(app dogma.Application) *config.Application {
				return &config.Application{
					AsConfigured: config.ApplicationAsConfigured{
						Source: optional.Some(
							config.Source[dogma.Application]{
								TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
								Interface: optional.Some(app),
							},
						),
					},
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
			func(app dogma.Application) *config.Application {
				return &config.Application{
					AsConfigured: config.ApplicationAsConfigured{
						Source: optional.Some(
							config.Source[dogma.Application]{
								TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
								Interface: optional.Some(app),
							},
						),
						Identities: []config.Identity{
							{
								AsConfigured: config.IdentityAsConfigured{
									Name: "app",
									Key:  "bed53df8-bf22-4502-be4b-64d56532d8be",
								},
							},
						},
						Handlers: []config.Handler{
							&config.Aggregate{
								AsConfigured: config.AggregateAsConfigured{
									Source: optional.Some(
										config.Source[dogma.AggregateMessageHandler]{
											TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub",
											Interface: optional.Some(aggregate),
										},
									),
									IsDisabled: optional.Some(false),
								},
							},
							&config.Process{
								AsConfigured: config.ProcessAsConfigured{
									Source: optional.Some(
										config.Source[dogma.ProcessMessageHandler]{
											TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub",
											Interface: optional.Some(process),
										},
									),
									IsDisabled: optional.Some(false),
								},
							},
							&config.Integration{
								AsConfigured: config.IntegrationAsConfigured{
									Source: optional.Some(
										config.Source[dogma.IntegrationMessageHandler]{
											TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub",
											Interface: optional.Some(integration),
										},
									),
									IsDisabled: optional.Some(false),
								},
							},
							&config.Projection{
								ConfigurationSource: optional.Some(
									config.Source[dogma.ProjectionMessageHandler]{
										TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
										Interface: optional.Some(projection),
									},
								),
							},
						},
					},
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
