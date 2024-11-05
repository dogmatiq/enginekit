package staticconfig

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/internal/typename"
	"golang.org/x/tools/go/ssa"
)

func analyzeRoutes[
	T config.Handler,
	H any,
	B configbuilder.HandlerBuilder[T, H],
](
	ctx *configurerCallContext[T, H, B],
) {
	analyzeVariadicArguments(
		ctx,
		ctx.Builder.Route,
		analyzeRoute,
	)
}

func analyzeRoute[
	T config.Handler,
	H any,
	B configbuilder.HandlerBuilder[T, H],
](
	ctx *configurerCallContext[T, H, B],
	b *configbuilder.RouteBuilder,
	r ssa.Value,
) {
	switch r := r.(type) {
	case *ssa.MakeInterface:
		analyzeRoute(ctx, b, r.X)

	case *ssa.Call:
		call := r.Common()
		fn := call.StaticCallee()

		if fn == nil {
			cannotAnalyzeNonStaticCall(b)
			return
		}

		switch fn.Object() {
		case ctx.Dogma.HandlesCommand:
			b.RouteType(config.HandlesCommandRouteType)
		case ctx.Dogma.HandlesEvent:
			b.RouteType(config.HandlesEventRouteType)
		case ctx.Dogma.ExecutesCommand:
			b.RouteType(config.ExecutesCommandRouteType)
		case ctx.Dogma.RecordsEvent:
			b.RouteType(config.RecordsEventRouteType)
		case ctx.Dogma.SchedulesTimeout:
			b.RouteType(config.SchedulesTimeoutRouteType)
		}

		b.MessageTypeName(typename.OfStatic(fn.TypeArgs()[0]))

	default:
		unimplementedAnalysis(b, r)
	}
}
