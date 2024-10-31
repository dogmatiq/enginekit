package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/test"
)

func TestProjectionDeliveryPolicyType(t *testing.T) {
	test.Enum(
		t,
		test.EnumSpec[ProjectionDeliveryPolicyType]{
			Range:       ProjectionDeliveryPolicyTypes,
			Switch:      SwitchByProjectionDeliveryPolicyType,
			MapToString: MapByProjectionDeliveryPolicyType[string],
		},
	)
}
