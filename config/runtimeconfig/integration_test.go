package runtimeconfig_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

func TestFromIntegration(t *testing.T) {
	cases := []struct {
		Name    string
		Handler dogma.IntegrationMessageHandler
		Want    func(h dogma.IntegrationMessageHandler) *config.Integration
	}{
		{
			"nil handler",
			nil,
			func(dogma.IntegrationMessageHandler) *config.Integration {
				return &config.Integration{}
			},
		},
		{
			"unconfigured handler",
			&IntegrationMessageHandlerStub{},
			func(h dogma.IntegrationMessageHandler) *config.Integration {
				return &config.Integration{
					ConfigurationSource: optional.Some(
						config.Source[dogma.IntegrationMessageHandler]{
							TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub",
							Interface: optional.Some(h),
						},
					),
				}
			},
		},
		{
			"configured handler",
			&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("integration", "51ffcb6f-171f-41a1-90e7-6fe1111649cd")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
					c.Disable()
				},
			},
			func(h dogma.IntegrationMessageHandler) *config.Integration {
				return &config.Integration{
					ConfigurationSource: optional.Some(
						config.Source[dogma.IntegrationMessageHandler]{
							TypeName:  "*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub",
							Interface: optional.Some(h),
						},
					),
					ConfiguredIdentities: []config.Identity{
						{
							AsConfigured: config.IdentityProperties{
								Name: "integration",
								Key:  "51ffcb6f-171f-41a1-90e7-6fe1111649cd",
							},
						},
					},
					ConfiguredRoutes: []config.Route{
						{
							RouteType:       optional.Some(config.HandlesCommandRouteType),
							MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
						},
						{
							RouteType:       optional.Some(config.RecordsEventRouteType),
							MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
						},
					},
					ConfiguredAsDisabled: true,
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			Expect(
				t,
				"unexpected config",
				FromIntegration(c.Handler),
				c.Want(c.Handler),
			)
		})
	}
}
