package staticconfig

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/internal/typename"
	"golang.org/x/tools/go/ssa"
)

func analyzeRoutesCall(
	ctx *context,
	b configbuilder.HandlerBuilder,
	call configurerCall,
) {
	for r := range resolveVariadic(b, call) {
		b.Route(func(b *configbuilder.RouteBuilder) {
			b.UpdateFidelity(call.Fidelity)
			analyzeRoute(ctx, b, r)
		})
	}
}

func analyzeRoute(
	ctx *context,
	b *configbuilder.RouteBuilder,
	r ssa.Value,
) {
	call, ok := findRouteCall(ctx, r)
	if !ok {
		b.UpdateFidelity(config.Incomplete)
		return
	}

	fn := call.Common().StaticCallee()

	switch fn.Object() {
	case ctx.Dogma.HandlesCommand:
		b.SetRouteType(config.HandlesCommandRouteType)
	case ctx.Dogma.HandlesEvent:
		b.SetRouteType(config.HandlesEventRouteType)
	case ctx.Dogma.ExecutesCommand:
		b.SetRouteType(config.ExecutesCommandRouteType)
	case ctx.Dogma.RecordsEvent:
		b.SetRouteType(config.RecordsEventRouteType)
	case ctx.Dogma.SchedulesTimeout:
		b.SetRouteType(config.SchedulesTimeoutRouteType)
	}

	b.SetMessageTypeName(typename.OfStatic(fn.TypeArgs()[0]))
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
