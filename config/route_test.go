package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/message"
	"github.com/dogmatiq/enginekit/optional"
)

func TestRoute_render(t *testing.T) {
	cases := []renderTestCase{
		{
			Name:             "handles command",
			ExpectDescriptor: `route:handles-command:CommandStub[TypeA]`,
			ExpectDetails:    `valid handles-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some(typename.For[CommandStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "handles event",
			ExpectDescriptor: `route:handles-event:EventStub[TypeA]`,
			ExpectDetails:    `valid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesEventRouteType),
					MessageTypeName: optional.Some(typename.For[EventStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "executes command",
			ExpectDescriptor: `route:executes-command:CommandStub[TypeA]`,
			ExpectDetails:    `valid executes-command route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(ExecutesCommandRouteType),
					MessageTypeName: optional.Some(typename.For[CommandStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "records event",
			ExpectDescriptor: `route:records-event:EventStub[TypeA]`,
			ExpectDetails:    `valid records-event route for github.com/dogmatiq/enginekit/enginetest/stubs.EventStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(RecordsEventRouteType),
					MessageTypeName: optional.Some(typename.For[EventStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[EventStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "schedules timeout",
			ExpectDescriptor: `route:schedules-timeout:TimeoutStub[TypeA]`,
			ExpectDetails:    `valid schedules-timeout route for github.com/dogmatiq/enginekit/enginetest/stubs.TimeoutStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(SchedulesTimeoutRouteType),
					MessageTypeName: optional.Some(typename.For[TimeoutStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[TimeoutStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "no runtime type information",
			ExpectDescriptor: `route:handles-command:SomeCommand`,
			ExpectDetails:    `valid handles-command route for pkg.SomeCommand (runtime type unavailable)`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommand"),
				},
			},
		},
		{
			Name:             "empty",
			ExpectDescriptor: `route`,
			ExpectDetails:    `incomplete route`,
			Component:        &Route{},
		},
		{
			Name:             "missing route type",
			ExpectDescriptor: `route:CommandStub[TypeA]`,
			ExpectDetails:    `incomplete route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					MessageTypeName: optional.Some(typename.For[CommandStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "missing message type name",
			ExpectDescriptor: `route:handles-command`,
			ExpectDetails:    `incomplete handles-command route`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType: optional.Some(HandlesCommandRouteType),
				},
			},
		},
		{
			Name:             "mismatched message type name",
			ExpectDescriptor: `route:handles-command:SomeCommand`,
			ExpectDetails: multiline(
				`invalid handles-command route for pkg.SomeCommand`,
				`  - type name mismatch: pkg.SomeCommand does not match the runtime type (github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA])`,
			),
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommand"),
					MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "mismatched message kind",
			ExpectDescriptor: `route:handles-event:CommandStub[TypeA]`,
			ExpectDetails: multiline(
				`invalid handles-event route for github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA]`,
				`  - unexpected message kind: github.com/dogmatiq/enginekit/enginetest/stubs.CommandStub[github.com/dogmatiq/enginekit/enginetest/stubs.TypeA] is a command, expected an event`,
			),
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesEventRouteType),
					MessageTypeName: optional.Some(typename.For[CommandStub[TypeA]]()),
					MessageType:     optional.Some(message.TypeFor[CommandStub[TypeA]]()),
				},
			},
		},
		{
			Name:             "speculative",
			ExpectDescriptor: `route:handles-command:SomeCommand`,
			ExpectDetails:    `valid speculative handles-command route for pkg.SomeCommand (runtime type unavailable)`,
			Component: &Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommand"),
					Fidelity:        Speculative,
				},
			},
		},
	}

	runRenderTests(t, cases)
}
