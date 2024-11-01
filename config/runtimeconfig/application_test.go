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
					EntityCommon: config.EntityCommon{
						ComponentCommon: config.ComponentCommon{
							IsPartial: true,
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
					EntityCommon: config.EntityCommon{
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
					},
					Source: optional.Some(app),
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
					EntityCommon: config.EntityCommon{
						IdentityComponents: []*config.Identity{
							{
								Name: optional.Some("app"),
								Key:  optional.Some("bed53df8-bf22-4502-be4b-64d56532d8be"),
							},
						},
						TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
					},
					HandlerComponents: []config.Handler{
						&config.Aggregate{
							HandlerCommon: config.HandlerCommon{
								EntityCommon: config.EntityCommon{
									TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"),
								},
							},
							Source: optional.Some(aggregate),
						},
						&config.Process{
							HandlerCommon: config.HandlerCommon{
								EntityCommon: config.EntityCommon{
									TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
								},
							},
							Source: optional.Some(process),
						},
						&config.Integration{
							HandlerCommon: config.HandlerCommon{
								EntityCommon: config.EntityCommon{
									TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
								},
							},
							Source: optional.Some(integration),
						},
						&config.Projection{
							HandlerCommon: config.HandlerCommon{
								EntityCommon: config.EntityCommon{
									TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub"),
								},
							},
							Source: optional.Some(projection),
						},
					},
					Source: optional.Some(app),
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
