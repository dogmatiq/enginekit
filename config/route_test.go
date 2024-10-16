package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

func TestRoute_String(t *testing.T) {
	cases := []struct {
		Name  string
		Want  string
		Route *Route
	}{
		{
			"valid, handles command",
			`route:handles-command(pkg.SomeCommand)`,
			&Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommand"),
				},
			},
		},
		{
			"valid, handles event",
			`route:handles-event(pkg.SomeEvent)`,
			&Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(HandlesEventRouteType),
					MessageTypeName: optional.Some("pkg.SomeEvent"),
				},
			},
		},
		{
			"valid, executes command",
			`route:executes-command(pkg.SomeCommand)`,
			&Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(ExecutesCommandRouteType),
					MessageTypeName: optional.Some("pkg.SomeCommand"),
				},
			},
		},
		{
			"valid, records event",
			`route:records-event(pkg.SomeEvent)`,
			&Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(RecordsEventRouteType),
					MessageTypeName: optional.Some("pkg.SomeEvent"),
				},
			},
		},
		{
			"valid, schedules timeout",
			`route:schedules-timeout(pkg.SomeTimeout)`,
			&Route{
				AsConfigured: RouteAsConfigured{
					RouteType:       optional.Some(SchedulesTimeoutRouteType),
					MessageTypeName: optional.Some("pkg.SomeTimeout"),
				},
			},
		},
		{
			"missing route type",
			`route(pkg.SomeCommand)`,
			&Route{
				AsConfigured: RouteAsConfigured{
					MessageTypeName: optional.Some("pkg.SomeCommand"),
				},
			},
		},
		{
			"missing message type",
			`route:handles-command`,
			&Route{
				AsConfigured: RouteAsConfigured{
					RouteType: optional.Some(HandlesCommandRouteType),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			got := c.Route.String()

			if got != c.Want {
				t.Fatalf("unexpected string: got %q, want %q", got, c.Want)
			}
		})
	}
}
