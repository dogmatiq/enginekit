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

func TestFromAggregate(t *testing.T) {
	cases := []struct {
		Name    string
		Handler dogma.AggregateMessageHandler
		Want    func(h dogma.AggregateMessageHandler) config.Aggregate
	}{
		{
			"nil handler",
			nil,
			func(dogma.AggregateMessageHandler) config.Aggregate {
				return config.Aggregate{}
			},
		},
		{
			"unconfigured handler",
			&AggregateMessageHandlerStub{},
			func(h dogma.AggregateMessageHandler) config.Aggregate {
				return config.Aggregate{
					ConfigurationSource: optional.Some(
						config.Source[dogma.AggregateMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub",
							Value:    optional.Some(h),
						},
					),
					ConfigurationIsExhaustive: true,
				}
			},
		},
		{
			"configured handler",
			&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("aggregate", "d9d75a75-7839-4b3e-a7e5-c8884b88ea57")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
					c.Disable()
				},
			},
			func(app dogma.AggregateMessageHandler) config.Aggregate {
				return config.Aggregate{
					ConfigurationSource: optional.Some(
						config.Source[dogma.AggregateMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub",
							Value:    optional.Some(app),
						},
					),
					ConfiguredIdentities: []config.Identity{
						{
							Name: "aggregate",
							Key:  "d9d75a75-7839-4b3e-a7e5-c8884b88ea57",
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
				FromAggregate(c.Handler),
				c.Want(c.Handler),
			)
		})
	}
}
