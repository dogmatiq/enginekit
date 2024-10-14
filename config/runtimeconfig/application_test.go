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
				return &config.Application{
					ConfigurationFidelity: config.Fidelity{
						IsExhaustive: true,
					},
				}
			},
		},
		{
			"unconfigured application",
			&ApplicationStub{},
			func(app dogma.Application) *config.Application {
				return &config.Application{
					ConfigurationSource: optional.Some(
						config.Source[dogma.Application]{
							TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
							Interface: optional.Some(app),
						},
					),
					ConfigurationFidelity: config.Fidelity{
						IsExhaustive: true,
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
					ConfigurationSource: optional.Some(
						config.Source[dogma.Application]{
							TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub",
							Interface: optional.Some(app),
						},
					),
					ConfiguredIdentities: []config.Identity{
						{
							Name: "app",
							Key:  "bed53df8-bf22-4502-be4b-64d56532d8be",
						},
					},
					ConfiguredHandlers: []config.Handler{
						&config.Aggregate{
							ConfigurationSource: optional.Some(
								config.Source[dogma.AggregateMessageHandler]{
									TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub",
									Interface: optional.Some(aggregate),
								},
							),
							ConfigurationFidelity: config.Fidelity{
								IsExhaustive: true,
							},
						},
						&config.Process{
							ConfigurationSource: optional.Some(
								config.Source[dogma.ProcessMessageHandler]{
									TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub",
									Interface: optional.Some(process),
								},
							),
							ConfigurationFidelity: config.Fidelity{
								IsExhaustive: true,
							},
						},
						&config.Integration{
							ConfigurationSource: optional.Some(
								config.Source[dogma.IntegrationMessageHandler]{
									TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub",
									Interface: optional.Some(integration),
								},
							),
							ConfigurationFidelity: config.Fidelity{
								IsExhaustive: true,
							},
						},
						&config.Projection{
							ConfigurationSource: optional.Some(
								config.Source[dogma.ProjectionMessageHandler]{
									TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
									Interface: optional.Some(projection),
								},
							),
							ConfigurationFidelity: config.Fidelity{
								IsExhaustive: true,
							},
						},
					},
					ConfigurationFidelity: config.Fidelity{
						IsExhaustive: true,
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
