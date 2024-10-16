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
					AsConfigured: config.ApplicationAsConfigured{
						Fidelity: config.Incomplete,
					},
				}
			},
		},
		{
			"unconfigured application",
			&ApplicationStub{},
			func(app dogma.Application) *config.Application {
				return &config.Application{
					AsConfigured: config.ApplicationAsConfigured{
						Source: config.Value[dogma.Application]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
							Value:    optional.Some(app),
						},
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
						Source: config.Value[dogma.Application]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
							Value:    optional.Some(app),
						},
						Identities: []*config.Identity{
							{
								AsConfigured: config.IdentityAsConfigured{
									Name: optional.Some("app"),
									Key:  optional.Some("bed53df8-bf22-4502-be4b-64d56532d8be"),
								},
							},
						},
						Handlers: []config.Handler{
							&config.Aggregate{
								AsConfigured: config.AggregateAsConfigured{
									Source: config.Value[dogma.AggregateMessageHandler]{
										TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"),
										Value:    optional.Some(aggregate),
									},
									IsDisabled: optional.Some(false),
								},
							},
							&config.Process{
								AsConfigured: config.ProcessAsConfigured{
									Source: config.Value[dogma.ProcessMessageHandler]{
										TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
										Value:    optional.Some(process),
									},
									IsDisabled: optional.Some(false),
								},
							},
							&config.Integration{
								AsConfigured: config.IntegrationAsConfigured{
									Source: config.Value[dogma.IntegrationMessageHandler]{
										TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
										Value:    optional.Some(integration),
									},
									IsDisabled: optional.Some(false),
								},
							},
							&config.Projection{
								AsConfigured: config.ProjectionAsConfigured{
									Source: config.Value[dogma.ProjectionMessageHandler]{
										TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub"),
										Value:    optional.Some(projection),
									},
									DeliveryPolicy: optional.Some(
										config.Value[dogma.ProjectionDeliveryPolicy]{
											TypeName: optional.Some("github.com/dogmatiq/dogma.UnicastProjectionDeliveryPolicy"),
											Value:    optional.Some[dogma.ProjectionDeliveryPolicy](dogma.UnicastProjectionDeliveryPolicy{}),
										},
									),
									IsDisabled: optional.Some(false),
								},
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
