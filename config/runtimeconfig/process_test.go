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

func TestFromProcess(t *testing.T) {
	cases := []struct {
		Name    string
		Handler dogma.ProcessMessageHandler
		Want    func(h dogma.ProcessMessageHandler) config.Process
	}{
		{
			"nil handler",
			nil,
			func(dogma.ProcessMessageHandler) config.Process {
				return config.Process{}
			},
		},
		{
			"unconfigured handler",
			&ProcessMessageHandlerStub{},
			func(h dogma.ProcessMessageHandler) config.Process {
				return config.Process{
					Impl: optional.Some(
						config.Implementation[dogma.ProcessMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub",
							Source:   optional.Some(h),
						},
					),
					IsExhaustive: true,
				}
			},
		},
		{
			"configured handler",
			&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("projection", "050415ad-ce90-496f-8987-40467e5415e0")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					)
					c.Disable()
				},
			},
			func(h dogma.ProcessMessageHandler) config.Process {
				return config.Process{
					Impl: optional.Some(
						config.Implementation[dogma.ProcessMessageHandler]{
							TypeName: "*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub",
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
						{
							RouteType: optional.Some(config.ExecutesCommandRoute),
							MessageType: optional.Some(
								config.MessageType{
									TypeName: "github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]",
									Kind:     message.CommandKind,
									Type:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
								},
							),
						},
						{
							RouteType: optional.Some(config.SchedulesTimeoutRoute),
							MessageType: optional.Some(
								config.MessageType{
									TypeName: "github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]",
									Kind:     message.TimeoutKind,
									Type:     optional.Some(message.TypeFor[TimeoutStub[TypeA]]()),
								},
							),
						},
					},
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
				FromProcess(c.Handler),
				c.Want(c.Handler),
			)
		})
	}
}
