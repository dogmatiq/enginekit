package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/constraints"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

func testHandler[
	T interface {
		Handler
		Interface() H
	},
	B configbuilder.EntityBuilder[T, H],
	H constraints.Handler[C, R],
	C constraints.HandlerConfigurer[R],
	R dogma.Route,
](
	t *testing.T,
	build func(func(B)) T,
	runtime func(H) T,
	construct func(func(C)) H,
) {
	testEntity(
		t,
		build,
		runtime,
		construct,
	)

	t.Run("func RouteSet()", func(t *testing.T) {
	})

	t.Run("func IsDisabled()", func(t *testing.T) {
	})
}
