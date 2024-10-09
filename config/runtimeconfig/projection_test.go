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

func TestFromProjection(t *testing.T) {
	cases := []struct {
		Name    string
		Handler dogma.ProjectionMessageHandler
		Want    func(h dogma.ProjectionMessageHandler) config.Projection
	}{
		{
			"nil handler",
			nil,
			func(dogma.ProjectionMessageHandler) config.Projection {
				return config.Projection{}
			},
		},
		{
			"unconfigured handler",
			&ProjectionMessageHandlerStub{},
			func(h dogma.ProjectionMessageHandler) config.Projection {
				return config.Projection{
					Impl: optional.Some(
						config.Implementation[dogma.ProjectionMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
							Source:   optional.Some(h),
						},
					),
					IsExhaustive: true,
				}
			},
		},
		{
			"configured handler",
			&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("projection", "050415ad-ce90-496f-8987-40467e5415e0")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
					)
					c.DeliveryPolicy(
						dogma.UnicastProjectionDeliveryPolicy{},
					)
					c.Disable()
				},
			},
			func(h dogma.ProjectionMessageHandler) config.Projection {
				return config.Projection{
					Impl: optional.Some(
						config.Implementation[dogma.ProjectionMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
							Source:   optional.Some(h),
						},
					),
					ConfiguredIdentities: []config.Identity{
						{
							Name: "projection",
							Key:  "050415ad-ce90-496f-8987-40467e5415e0",
						},
					},
					ConfiguredRoutes: []config.Route{
						{
							RouteType: optional.Some(config.HandlesEventRoute),
							MessageType: optional.Some(
								config.MessageType{
									TypeName: "github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]",
									Kind:     message.EventKind,
									Type:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
								},
							),
						},
					},
					DeliveryPolicy: optional.Some(
						config.ProjectionDeliveryPolicy{
							TypeName:       optional.Some("github.com/dogmatiq/dogma.UnicastProjectionDeliveryPolicy"),
							Implementation: optional.Some[dogma.ProjectionDeliveryPolicy](dogma.UnicastProjectionDeliveryPolicy{}),
						},
					),
					IsDisabled:   true,
					IsExhaustive: true,
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			Expect(
				t,
				"unexpected config",
				FromProjection(c.Handler),
				c.Want(c.Handler),
			)
		})
	}
}
