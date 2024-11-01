package staticconfig

import (
	"go/types"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
)

// analyzeApplicationType analyzes t, which must be an implementation of
// [dogma.Application].
func analyzeApplicationType(ctx *context, t types.Type) {
	app := configbuilder.Application(
		func(b *configbuilder.ApplicationBuilder) {
			analyzeEntity(
				ctx,
				t,
				b,
				analyzeApplicationConfigurerCall,
			)
		},
	)

	ctx.Analysis.Applications = append(ctx.Analysis.Applications, app)
}

func analyzeApplicationConfigurerCall(
	ctx *configurerCallContext[*config.Application, dogma.Application, *configbuilder.ApplicationBuilder],
) {
	switch ctx.Method.Name() {
	case "RegisterAggregate":
		analyzeHandler(ctx, ctx.Builder.Aggregate, nil)
	case "RegisterProcess":
		analyzeHandler(ctx, ctx.Builder.Process, nil)
	case "RegisterIntegration":
		analyzeHandler(ctx, ctx.Builder.Integration, nil)
	case "RegisterProjection":
		analyzeHandler(ctx, ctx.Builder.Projection, analyzeProjectionConfigurerCall)
	default:
		ctx.Builder.Partial()
	}
}
