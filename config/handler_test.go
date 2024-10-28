package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

func testHandler[
	T Handler,
	B configbuilder.EntityBuilder[T],
](
	t *testing.T,
	build func(func(B)) T,
) {
	testEntity(t, build)
}
