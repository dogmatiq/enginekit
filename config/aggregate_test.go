package config_test

import (
	"testing"

	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

func TestAggregate(t *testing.T) {
	testHandler(t, configbuilder.Aggregate)
}
