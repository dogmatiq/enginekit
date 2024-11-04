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

func analyzeRoute(
	ctx *context,
	b *configbuilder.RouteBuilder,
	r ssa.Value,
) {
	call, ok := findRouteCall(ctx, r)
	if !ok {
		b.Partial()
		return
	}

	fn := call.Common().StaticCallee()

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
}

func findRouteCall(
	ctx *context,
	v ssa.Value,
) (*ssa.Call, bool) {
	switch v := v.(type) {
	case *ssa.Call:
		fn := v.Common().StaticCallee()
		if fn != nil {
			switch fn.Object() {
			case ctx.Dogma.HandlesCommand,
				ctx.Dogma.HandlesEvent,
				ctx.Dogma.ExecutesCommand,
				ctx.Dogma.RecordsEvent,
				ctx.Dogma.SchedulesTimeout:
				return v, true
			}
		}
	case *ssa.MakeInterface:
		return findRouteCall(ctx, v.X)
	}

	return nil, false
}
