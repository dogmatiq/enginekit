package staticconfig

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
)

func analyzeHandler[
	T config.Handler,
	H any,
	B configbuilder.HandlerBuilder[T, H],
](
	ctx *configurerCallContext[*config.Application, dogma.Application, *configbuilder.ApplicationBuilder],
	build func(func(B)),
	analyze configurerCallAnalyzer[T, H, B],
) {
	build(func(b B) {
		ctx.Apply(b)

		t := ssax.ConcreteType(ctx.Args[0])

		if !t.IsPresent() {
			b.Partial("the handler's type is unknown")
			return
		}

		analyzeEntity(
			ctx.context,
			t.Get(),
			b,
			func(ctx *configurerCallContext[T, H, B]) {
				switch ctx.Method.Name() {
				case "Routes":
					analyzeRoutes(ctx)

				case "Disable":
					ctx.Builder.Disabled(
						func(b *configbuilder.FlagBuilder[config.Disabled]) {
							ctx.Apply(b)
							b.Value(true)
						},
					)

				default:
					if analyze == nil {
						cannotAnalyzeUnrecognizedConfigurerMethod(ctx)
					} else {
						analyze(ctx)
					}
				}
			},
		)
	})
}

func analyzeProjectionConfigurerCall(
	ctx *configurerCallContext[*config.Projection, dogma.ProjectionMessageHandler, *configbuilder.ProjectionBuilder],
) {
	switch ctx.Method.Name() {
	case "DeliveryPolicy":
		panic("not implemented") // TODO
	default:
		cannotAnalyzeUnrecognizedConfigurerMethod(ctx)
	}
}
