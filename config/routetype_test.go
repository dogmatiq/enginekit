package config

import (
	"testing"

	"github.com/dogmatiq/enginekit/internal/test"
)

func TestRouteType(t *testing.T) {
	test.Enum(
		t,
		test.EnumSpec[RouteType]{
			Range:       RouteTypes,
			Switch:      SwitchByRouteType,
			MapToString: MapByRouteType[string],
		},
	)
}
