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
					ConfigurationSource: optional.Some(
						config.Source[dogma.ProjectionMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
							Value:    optional.Some(h),
						},
					),
					ConfigurationIsExhaustive: true,
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
					ConfigurationSource: optional.Some(
						config.Source[dogma.ProjectionMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProjectionMessageHandlerStub",
							Value:    optional.Some(h),
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
							RouteType:       optional.Some(config.HandlesEventRouteType),
							MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
							MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
						},
					},
					ConfiguredDeliveryPolicy: optional.Some(
						config.ProjectionDeliveryPolicy{
							TypeName:       optional.Some("github.com/dogmatiq/dogma.UnicastProjectionDeliveryPolicy"),
							Implementation: optional.Some[dogma.ProjectionDeliveryPolicy](dogma.UnicastProjectionDeliveryPolicy{}),
						},
					),
					ConfiguredAsDisabled:      true,
					ConfigurationIsExhaustive: true,
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