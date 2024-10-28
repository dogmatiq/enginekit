package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestAggregate(t *testing.T) {
	testHandler(
		t,
		configbuilder.Aggregate,
		runtimeconfig.FromAggregate,
		func(fn func(dogma.AggregateConfigurer)) dogma.AggregateMessageHandler {
			return &AggregateMessageHandlerStub{ConfigureFunc: fn}
		},
	)
}
