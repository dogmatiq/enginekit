package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

func TestAggregate_RouteSet(t *testing.T) {
	t.Run("it returns the normalized routes", func(t *testing.T) {
		h := &AggregateMessageHandlerStub{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}

		cfg := runtimeconfig.FromAggregate(h)

		Expect(
			t,
			"unexpected routes",
			cfg.RouteSet().MessageTypes(),
			map[message.Type]RouteDirection{
				message.TypeFor[CommandStub[TypeA]](): InboundDirection,
				message.TypeFor[EventStub[TypeA]]():   OutboundDirection,
			},
		)
	})

	t.Run("it panics if the routes are invalid", func(t *testing.T) {
		cases := []struct {
			Name  string
			Want  string
			Route *Route
		}{
			{
				"empty route",
				`aggregate is invalid: route is invalid: could not evaluate entire configuration`,
				&Route{},
			},
			{
				"unexpected ExecutesCommand route",
				`aggregate is invalid: unexpected executes-command route for pkg.SomeCommandType`,
				&Route{
					RouteType:       optional.Some(ExecutesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommandType"),
				},
			},
			{
				"unexpected HandlesEvent route",
				`aggregate is invalid: unexpected handles-event route for pkg.SomeEventType`,
				&Route{
					RouteType:       optional.Some(HandlesEventRouteType),
					MessageTypeName: optional.Some("pkg.SomeEventType"),
				},
			},
			{
				"unexpected SchedulesTimeout route",
				`aggregate is invalid: unexpected schedules-timeout route for pkg.SomeTimeoutType`,
				&Route{
					RouteType:       optional.Some(SchedulesTimeoutRouteType),
					MessageTypeName: optional.Some("pkg.SomeTimeoutType"),
				},
			},
		}

		for _, c := range cases {
			t.Run(c.Name, func(t *testing.T) {
				cfg := &Aggregate{
					HandlerCommon: HandlerCommon{
						RouteComponents: []*Route{c.Route},
					},
				}

				ExpectPanic(
					t,
					c.Want,
					func() {
						cfg.RouteSet()
					},
				)
			})
		}
	})
}

func TestAggregate_IsDisabled(t *testing.T) {
	disable := false

	h := &AggregateMessageHandlerStub{
		ConfigureFunc: func(c dogma.AggregateConfigurer) {
			c.Identity("name", "19cb98d5-dd17-4daf-ae00-1b413b7b899a")
			c.Routes(
				dogma.HandlesCommand[CommandStub[TypeA]](),
				dogma.RecordsEvent[EventStub[TypeA]](),
			)
			if disable {
				c.Disable()
			}
		},
	}

	cfg := runtimeconfig.FromAggregate(h)

	if cfg.IsDisabled() {
		t.Fatal("did not expect handler to be disabled")
	}

	disable = true
	cfg = runtimeconfig.FromAggregate(h)

	if !cfg.IsDisabled() {
		t.Fatal("expected handler to be disabled")
	}
}
