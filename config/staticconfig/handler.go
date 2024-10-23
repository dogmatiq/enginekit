package staticconfig

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"golang.org/x/tools/go/ssa"
)

func analyzeHandler[T configbuilder.HandlerBuilder](
	ctx *configurerCallContext[*configbuilder.ApplicationBuilder],
	build func(func(T)),
	analyze configurerCallAnalyzer[T],
) {
	build(func(b T) {
		b.UpdateFidelity(ctx.Fidelity)

		inst, ok := ctx.Args[0].(*ssa.MakeInterface)
		if !ok {
			b.UpdateFidelity(config.Incomplete)
			return
		}

		analyzeEntity(
			ctx.context,
			inst.X.Type(),
			b,
			func(ctx *configurerCallContext[T]) {
				switch ctx.Method.Name() {
				case "Routes":
					analyzeRoutes(ctx)

				case "Disable":
					// TODO(jmalloc): f is lost in this case, so any handler
					// that is _sometimes_ disabled will appear as always
					// disabled, which is a bit non-sensical.
					//
					// It probably needs similar treatment to
					// https://github.com/dogmatiq/enginekit/issues/55.
					ctx.Builder.SetDisabled(true)

				default:
					if analyze == nil {
						ctx.Builder.UpdateFidelity(config.Incomplete)
					} else {
						analyze(ctx)
					}
				}
			},
		)

		// If the handler wasn't disabled, and the configuration is NOT
		// incomplete, we know that the handler is enabled.
		if !b.IsDisabled().IsPresent() && b.Fidelity()&config.Incomplete == 0 {
			b.SetDisabled(false)
		}
	})
}

func analyzeProjectionConfigurerCall(
	ctx *configurerCallContext[*configbuilder.ProjectionBuilder],
) {
	switch ctx.Method.Name() {
	case "DeliveryPolicy":
		panic("not implemented") // TODO
	default:
		ctx.Builder.UpdateFidelity(config.Incomplete)
	}
}
