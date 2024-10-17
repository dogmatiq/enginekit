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
	ctx.Applications = append(
		ctx.Applications,
		configbuilder.Application(
			func(b *configbuilder.ApplicationBuilder) {
				b.SetSourceTypeName(typename.OfStatic(t))

				for call := range analyzeEntity(ctx, b, t) {
					switch call.Method.Name() {
					// // case "RegisterAggregate":
					// // 	analyzeRegisterAggregateCall(ctx, c)
					// // case "RegisterProcess":
					// // 	analyzeRegisterProcessCall(ctx, c)
					// // case "RegisterIntegration":
					// // 	analyzeRegisterIntegrationCall(ctx, c)
					// // case "RegisterProjection":
					// // 	analyzeRegisterProjectionCall(ctx, c)
					// // case "Handlers":
					// // 	panic("not implemented")
					default:
						b.UpdateFidelity(config.Incomplete)
					}
				}
			},
		),
	)

	// switch c.Common().Method.Name() {
	// 	case "Identity":
	// 		app.IdentityValue = analyzeIdentityCall(c)
	// 	case "RegisterAggregate":
	// 		addHandlerFromArguments(
	// 			prog,
	// 			dogmaPkg,
	// 			dogmaPkg.AggregateMessageHandler,
	// 			args,
	// 			app.HandlersValue,
	// 			configkit.AggregateHandlerType,
	// 		)
	// 	case "RegisterProcess":
	// 		addHandlerFromArguments(
	// 			prog,
	// 			dogmaPkg,
	// 			dogmaPkg.ProcessMessageHandler,
	// 			args,
	// 			app.HandlersValue,
	// 			configkit.ProcessHandlerType,
	// 		)
	// 	case "RegisterProjection":
	// 		addHandlerFromArguments(
	// 			prog,
	// 			dogmaPkg,
	// 			dogmaPkg.ProjectionMessageHandler,
	// 			args,
	// 			app.HandlersValue,
	// 			configkit.ProjectionHandlerType,
	// 		)
	// 	case "RegisterIntegration":
	// 		addHandlerFromArguments(
	// 			prog,
	// 			dogmaPkg,
	// 			dogmaPkg.IntegrationMessageHandler,
	// 			args,
	// 			app.HandlersValue,
	// 			configkit.IntegrationHandlerType,
	// 		)
	// 	}
	// }
}
