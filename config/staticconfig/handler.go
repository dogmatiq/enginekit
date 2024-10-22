package staticconfig

import (
	"iter"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/internal/typename"
	"golang.org/x/tools/go/ssa"
)

func analyzeHandler(
	ctx *context,
	b configbuilder.HandlerBuilder,
	h ssa.Value,
) iter.Seq[configurerCall] {
	return func(yield func(configurerCall) bool) {
		switch inst := h.(type) {
		default:
			b.UpdateFidelity(config.Incomplete)
		case *ssa.MakeInterface:
			t := inst.X.Type()
			b.SetSourceTypeName(typename.OfStatic(t))

			for call := range findConfigurerCalls(ctx, b, t) {
				switch call.Method.Name() {
				case "Identity":
					analyzeIdentityCall(b, call)
				case "Routes":
					// analyzeRoutesCall(ctx, b, call)
				case "Disable":
					b.SetDisabled(true)
				default:
					if !yield(call) {
						return
					}
				}
			}

			// If the handler wasn't disabled, and the configuration is NOT
			// incomplete, we know that the handler is enabled.
			if !b.IsDisabled().IsPresent() && b.Fidelity()&config.Incomplete == 0 {
				b.SetDisabled(false)
			}
		}
	}
}

func analyzeAggregate(
	ctx *context,
	b *configbuilder.AggregateBuilder,
	h ssa.Value,
) {
	for call := range analyzeHandler(ctx, b, h) {
		b.UpdateFidelity(call.Fidelity)
	}
}

func analyzeProcess(
	ctx *context,
	b *configbuilder.ProcessBuilder,
	h ssa.Value,
) {
	for call := range analyzeHandler(ctx, b, h) {
		b.UpdateFidelity(call.Fidelity)
	}
}

func analyzeIntegration(
	ctx *context,
	b *configbuilder.IntegrationBuilder,
	h ssa.Value,
) {
	for call := range analyzeHandler(ctx, b, h) {
		b.UpdateFidelity(call.Fidelity)
	}
}

func analyzeProjection(
	ctx *context,
	b *configbuilder.ProjectionBuilder,
	h ssa.Value,
) {
	for call := range analyzeHandler(ctx, b, h) {
		b.UpdateFidelity(call.Fidelity)
	}
}
