package config_test

import (
	"testing"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/internal/constraints"
	"github.com/dogmatiq/enginekit/internal/test"
)

func testHandler[
	T interface {
		Handler
		Interface() H
	},
	B configbuilder.HandlerBuilder[T, H],
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

	t.Run("func IsDisabled()", func(t *testing.T) {
		t.Run("it returns false if the handler is not disabled", func(t *testing.T) {
			handler := build(func(b B) {})

			test.Expect(
				t,
				"unexpected flag value",
				handler.IsDisabled(),
				false,
			)
		})

		t.Run("it returns true if the handler is disabled", func(t *testing.T) {
			handler := build(func(b B) {
				b.Disabled(
					func(b *configbuilder.FlagBuilder[Disabled]) {
						b.Value(true)
					},
				)
			})

			test.Expect(
				t,
				"unexpected flag value",
				handler.IsDisabled(),
				true,
			)
		})

		t.Run("panics if the flag is partially configured", func(t *testing.T) {
			handler := build(func(b B) {
				b.Disabled(
					func(b *configbuilder.FlagBuilder[Disabled]) {
						b.Partial("<reason>")
					},
				)
			})

			test.ExpectPanic(
				t,
				`flag:disabled is invalid: could not evaluate entire configuration: <reason>`,
				func() {
					handler.IsDisabled()
				},
			)
		})

		t.Run("panics if the flag value is missing", func(t *testing.T) {
			handler := build(func(b B) {
				b.Disabled(
					func(b *configbuilder.FlagBuilder[Disabled]) {},
				)
			})

			test.ExpectPanic(
				t,
				`flag:disabled is invalid: value is unavailable`,
				func() {
					handler.IsDisabled()
				},
			)
		})
	})
}
