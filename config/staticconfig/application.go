package staticconfig

import (
	"go/types"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/internal/typename"
)

// analyzeApplication analyzes t, which must be an implementation of
// [dogma.Application].
func analyzeApplication(ctx *context, t types.Type) {
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
						analyzeRegisterAggregateCall(b, call)
					case "RegisterProcess":
						analyzeRegisterProcessCall(b, call)
					case "RegisterIntegration":
						analyzeRegisterIntegrationCall(b, call)
					case "RegisterProjection":
						analyzeRegisterProjectionCall(b, call)
					default:
						b.UpdateFidelity(config.Incomplete)
					}
				}
			},
		),
	)
}

func analyzeRegisterAggregateCall(
	b *configbuilder.ApplicationBuilder,
	_ configurerCall,
) {
	b.Aggregate(func(b *configbuilder.AggregateBuilder) {
		b.UpdateFidelity(config.Incomplete)
	})
}

func analyzeRegisterProcessCall(
	b *configbuilder.ApplicationBuilder,
	_ configurerCall,
) {
	b.Process(func(b *configbuilder.ProcessBuilder) {
		b.UpdateFidelity(config.Incomplete)
	})
}

func analyzeRegisterIntegrationCall(
	b *configbuilder.ApplicationBuilder,
	_ configurerCall,
) {
	b.Integration(func(b *configbuilder.IntegrationBuilder) {
		b.UpdateFidelity(config.Incomplete)
	})
}

func analyzeRegisterProjectionCall(
	b *configbuilder.ApplicationBuilder,
	_ configurerCall,
) {
	b.Projection(func(b *configbuilder.ProjectionBuilder) {
		b.UpdateFidelity(config.Incomplete)
	})
}
