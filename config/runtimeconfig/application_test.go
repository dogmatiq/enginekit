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
					EntityCommon: config.EntityCommon[dogma.Application]{
						ComponentCommon: config.ComponentCommon{
							ComponentFidelity: config.Incomplete,
						},
					},
				}
			},
		},
		{
			"unconfigured application",
			&ApplicationStub{},
			func(app dogma.Application) *config.Application {
				return &config.Application{
					EntityCommon: config.EntityCommon[dogma.Application]{
						SourceTypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
						Source:         optional.Some(app),
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
					EntityCommon: config.EntityCommon[dogma.Application]{
						ComponentCommon: config.ComponentCommon{},
						IdentityComponents: []*config.Identity{
							{
								Name: optional.Some("app"),
								Key:  optional.Some("bed53df8-bf22-4502-be4b-64d56532d8be"),
							},
						},
						SourceTypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
						Source:         optional.Some(app),
					},
					HandlerComponents: []config.Handler{
						&config.Aggregate{
							HandlerCommon: config.HandlerCommon[dogma.AggregateMessageHandler]{
								EntityCommon: config.EntityCommon[dogma.AggregateMessageHandler]{
									SourceTypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"),
									Source:         optional.Some(aggregate),
								},
							},
						},
						&config.Process{
							HandlerCommon: config.HandlerCommon[dogma.ProcessMessageHandler]{
								EntityCommon: config.EntityCommon[dogma.ProcessMessageHandler]{
									SourceTypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
									Source:         optional.Some(process),
								},
							},
						},
						&config.Integration{
							HandlerCommon: config.HandlerCommon[dogma.IntegrationMessageHandler]{
								EntityCommon: config.EntityCommon[dogma.IntegrationMessageHandler]{
									SourceTypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
									Source:         optional.Some(integration),
								},
							},
						},
						&config.Projection{
							HandlerCommon: config.HandlerCommon[dogma.ProjectionMessageHandler]{
								EntityCommon: config.EntityCommon[dogma.ProjectionMessageHandler]{
									SourceTypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub"),
									Source:         optional.Some(projection),
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
