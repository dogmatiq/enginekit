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
		Want    func(h dogma.IntegrationMessageHandler) config.Integration
	}{
		{
			"nil handler",
			nil,
			func(dogma.IntegrationMessageHandler) config.Integration {
				return config.Integration{}
			},
		},
		{
			"unconfigured handler",
			&IntegrationMessageHandlerStub{},
			func(app dogma.IntegrationMessageHandler) config.Integration {
				return config.Integration{
					TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
					Implementation: optional.Some(app),
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
			func(app dogma.IntegrationMessageHandler) config.Integration {
				return config.Integration{
					TypeName:       optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"),
					Implementation: optional.Some(app),
					Identities: []config.Identity{
						{
							Name: "integration",
							Key:  "51ffcb6f-171f-41a1-90e7-6fe1111649cd",
						},
					},
					Routes: []config.Route{
						{
							Type: optional.Some(config.HandlesCommandRoute),
							MessageType: optional.Some(
								config.MessageType{
									TypeName: "github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]",
									Kind:     message.CommandKind,
									Type:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
								},
							),
						},
						{
							Type: optional.Some(config.RecordsEventRoute),
							MessageType: optional.Some(
								config.MessageType{
									TypeName: "github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]",
									Kind:     message.EventKind,
									Type:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
								},
							),
						},
					},
					IsDisabled: true,
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
