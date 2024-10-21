package staticconfig

import (
	"fmt"
	"iter"
	"os"
	"reflect"
	"strings"

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
					analyzeRoutesCall(b, call)
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

func analyzeRoutesCall(
	b configbuilder.HandlerBuilder,
	call configurerCall,
) {
	// call.Value.Parent().WriteTo(os.Stderr)
	// routes := call.Args[0].(*ssa.Slice)
	// fmt.Println(routes.X, reflect.TypeOf(routes.X))

	for range resolveVariadic(b, call) {
		b.Route(func(b *configbuilder.RouteBuilder) {
			b.UpdateFidelity(call.Fidelity)
			b.UpdateFidelity(config.Incomplete)
		})
	}
}

func resolveVariadic(
	_ configbuilder.EntityBuilder,
	call configurerCall,
) []ssa.Value {
	call.Instruction.Parent().WriteTo(os.Stderr)
	// n := len(call.Args) - 1
	// variadics := call.Args[n].(*ssa.Slice)

	// for x := range walkPredecessors(call.Instruction.Block()) {
	// 	fmt.Println(x)
	// }

	walkReferrers(call.Instruction, 0, map[any]struct{}{})

	// b.UpdateFidelity(config.Incomplete)

	return nil
}

func walkReferrers(n any, depth int, seen map[any]struct{}) {
	if n == nil {
		return // why ?
	}

	fmt.Print(strings.Repeat("\t", depth), "- ", n)

	if _, ok := seen[n]; ok {
		fmt.Println(" (seen)")
		return
	}

	fmt.Print(" (", reflect.TypeOf(n), ")\n")
	seen[n] = struct{}{}

	if n, ok := n.(ssa.Value); ok {
		if refs := n.Referrers(); refs != nil {
			for _, ref := range *refs {
				walkReferrers(ref, depth+1, seen)
			}
		}
	}

	if n, ok := n.(ssa.Instruction); ok {
		for _, op := range n.Operands(nil) {
			walkReferrers(*op, depth+1, seen)
		}
	}
}

// // addMessagesFromRoutes analyzes the arguments in a call to a configurer's
// // Routes() method to populate the messages that are produced and consumed by
// // the handler.
// func addMessagesFromRoutes(
// 	messages configkit.EntityMessages[message.Name],
// 	args []ssa.Value,
// ) {
// 	var mii []*ssa.MakeInterface
// 	for _, arg := range args {
// 		recurseSSAValues(
// 			arg,
// 			&[]ssa.Value{},
// 			func(v ssa.Value) bool {
// 				if v, ok := v.(*ssa.Call); ok {
// 					// We don't want to recurse into the call to Routes() method
// 					// itself.
// 					if v.Common().IsInvoke() &&
// 						v.Common().Method.Name() == "Routes" {
// 						return true
// 					}
// 				}

// 				// We do want to collect all of the MakeInterface instructions
// 				// that can potentially indicate boxing into Dogma route
// 				// interfaces.
// 				if v, ok := v.(*ssa.MakeInterface); ok {
// 					mii = append(mii, v)
// 					return true
// 				}

// 				return false
// 			},
// 		)
// 	}

// 	for _, mi := range mii {
// 		// If this is the boxing to the following interfaces,
// 		// we need to analyze the concrete types:
// 		switch mi.X.Type().String() {
// 		case "github.com/dogmatiq/dogma.HandlesCommandRoute",
// 			"github.com/dogmatiq/dogma.HandlesEventRoute",
// 			"github.com/dogmatiq/dogma.ExecutesCommandRoute",
// 			"github.com/dogmatiq/dogma.SchedulesTimeoutRoute",
// 			"github.com/dogmatiq/dogma.RecordsEventRoute":

// 			// At this point we should expect that the interfaces above
// 			// are produced as a result of calls to following functions:
// 			// (At the time of writing this code, there is no other way
// 			// to produce these interfaces)
// 			//  `github.com/dogmatiq/dogma.HandlesCommand()
// 			//  `github.com/dogmatiq/dogma.HandlesEvent()`
// 			//  `github.com/dogmatiq/dogma.ExecutesCommand()`
// 			//  `github.com/dogmatiq/dogma.RecordsEvent()`
// 			//  `github.com/dogmatiq/dogma.SchedulesTimeout()`
// 			if f, ok := mi.X.(*ssa.Call).Common().Value.(*ssa.Function); ok {
// 				messages.Update(
// 					message.NameFromStaticType(f.TypeArgs()[0]),
// 					func(n message.Name, em *configkit.EntityMessage) {
// 						switch {
// 						case strings.HasPrefix(f.Name(), "HandlesCommand["):
// 							em.Kind = message.CommandKind
// 							em.IsConsumed = true

// 						case strings.HasPrefix(f.Name(), "HandlesEvent["):
// 							em.Kind = message.EventKind
// 							em.IsConsumed = true

// 						case strings.HasPrefix(f.Name(), "ExecutesCommand["):
// 							em.Kind = message.CommandKind
// 							em.IsProduced = true

// 						case strings.HasPrefix(f.Name(), "RecordsEvent["):
// 							em.Kind = message.EventKind
// 							em.IsProduced = true

// 						case strings.HasPrefix(f.Name(), "SchedulesTimeout["):
// 							em.Kind = message.TimeoutKind
// 							em.IsProduced = true
// 							em.IsConsumed = true
// 						}
// 					},
// 				)

// 			}
// 		}
// 	}
// }
