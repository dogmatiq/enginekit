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
					TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
					Implementation: optional.Some(app),
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
					TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"),
					Implementation: optional.Some(app),
					Identities: []config.Identity{
						{
							Name: "app",
							Key:  "bed53df8-bf22-4502-be4b-64d56532d8be",
						},
					},
					Aggregates: []config.Aggregate{
						{
							TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"),
							Implementation: optional.Some(aggregate),
						},
					},
					Processes: []config.Process{
						{
							TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
							Implementation: optional.Some(process),
						},
					},
					Integrations: []config.Integration{
						{
							TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
							Implementation: optional.Some(integration),
						},
					},
					Projections: []config.Projection{
						{
							TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub"),
							Implementation: optional.Some(projection),
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
