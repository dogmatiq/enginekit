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
		Want    func(h dogma.ProcessMessageHandler) *config.Process
	}{
		{
			"nil handler",
			nil,
			func(dogma.ProcessMessageHandler) *config.Process {
				return &config.Process{
					X: config.XProcess{
						Fidelity: config.Incomplete,
					},
				}
			},
		},
		{
			"unconfigured handler",
			&ProcessMessageHandlerStub{},
			func(h dogma.ProcessMessageHandler) *config.Process {
				return &config.Process{
					X: config.XProcess{
						Source: config.Value[dogma.ProcessMessageHandler]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
							Value:    optional.Some(h),
						},
					},
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
			func(h dogma.ProcessMessageHandler) *config.Process {
				return &config.Process{
					X: config.XProcess{
						Source: config.Value[dogma.ProcessMessageHandler]{
							TypeName: optional.Some("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"),
							Value:    optional.Some(h),
						},
						Identities: []*config.Identity{
							{
								AsConfigured: config.IdentityProperties{
									Name: optional.Some("projection"),
									Key:  optional.Some("050415ad-ce90-496f-8987-40467e5415e0"),
								},
							},
						},
						Routes: []*config.Route{
							{
								XRoute: config.XRoute{
									RouteType:       optional.Some(config.HandlesEventRouteType),
									MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
									MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
								},
							},
							{
								XRoute: config.XRoute{
									RouteType:       optional.Some(config.ExecutesCommandRouteType),
									MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
									MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
								},
							},
							{
								XRoute: config.XRoute{
									RouteType:       optional.Some(config.SchedulesTimeoutRouteType),
									MessageTypeName: optional.Some("github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]"),
									MessageType:     optional.Some(message.TypeFor[TimeoutStub[TypeA]]()),
								},
							},
						},
						DisabledFlags: config.Flag[config.Disabled]{{}},
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
				FromProcess(c.Handler),
				c.Want(c.Handler),
			)
		})
	}
}