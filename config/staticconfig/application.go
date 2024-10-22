package staticconfig

import (
	"go/types"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/internal/typename"
)

// analyzeApplicationType analyzes t, which must be an implementation of
// [dogma.Application].
func analyzeApplicationType(ctx *context, t types.Type) {
	ctx.Analysis.Applications = append(
		ctx.Analysis.Applications,
		configbuilder.Application(
			func(b *configbuilder.ApplicationBuilder) {
				b.SetSourceTypeName(typename.OfStatic(t))

				for call := range findConfigurerCalls(ctx, b, t) {
					switch call.Method.Name() {
					case "Identity":
						analyzeIdentityCall(b, call)
					case "RegisterAggregate":
						b.Aggregate(func(b *configbuilder.AggregateBuilder) {
							b.UpdateFidelity(call.Fidelity)
							analyzeAggregate(ctx, b, call.Args[0])
						})
					case "RegisterProcess":
						b.Process(func(b *configbuilder.ProcessBuilder) {
							b.UpdateFidelity(call.Fidelity)
							analyzeProcess(ctx, b, call.Args[0])
						})
					case "RegisterIntegration":
						b.Integration(func(b *configbuilder.IntegrationBuilder) {
							b.UpdateFidelity(call.Fidelity)
							analyzeIntegration(ctx, b, call.Args[0])
						})
					case "RegisterProjection":
						b.Projection(func(b *configbuilder.ProjectionBuilder) {
							b.UpdateFidelity(call.Fidelity)
							analyzeProjection(ctx, b, call.Args[0])
						})
					default:
						b.UpdateFidelity(config.Incomplete)
					}
				}
			},
		),
	)
}
