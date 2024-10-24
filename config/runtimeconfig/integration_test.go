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
				return &config.Integration{
					AsConfigured: config.IntegrationAsConfigured{
						Fidelity: config.Incomplete,
					},
				}
			},
		},
		{
			"unconfigured handler",
			&IntegrationMessageHandlerStub{},
			func(h dogma.IntegrationMessageHandler) *config.Integration {
				return &config.Integration{
					AsConfigured: config.IntegrationAsConfigured{
						Source: config.Value[dogma.IntegrationMessageHandler]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
							Value:    optional.Some(h),
						},
						IsDisabled: optional.Some(false),
					},
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
					AsConfigured: config.IntegrationAsConfigured{
						Source: config.Value[dogma.IntegrationMessageHandler]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
							Value:    optional.Some(h),
						},
						Identities: []*config.Identity{
							{
								AsConfigured: config.IdentityAsConfigured{
									Name: optional.Some("integration"),
									Key:  optional.Some("51ffcb6f-171f-41a1-90e7-6fe1111649cd"),
								},
							},
						},
						Routes: []*config.Route{
							{
								AsConfigured: config.RouteAsConfigured{
									RouteType:       optional.Some(config.HandlesCommandRouteType),
									MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
									MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
								},
							},
							{
								AsConfigured: config.RouteAsConfigured{
									RouteType:       optional.Some(config.RecordsEventRouteType),
									MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
									MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
								},
							},
						},
						IsDisabled: optional.Some(true),
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
				FromIntegration(c.Handler),
				c.Want(c.Handler),
			)
		})
	}
}
