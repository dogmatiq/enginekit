package staticconfig

import (
	"github.com/dogmatiq/enginekit/collections/sets"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
	"github.com/dogmatiq/enginekit/optional"
	"golang.org/x/tools/go/ssa"
)

// analyzeVariadicArguments analyzes the variadic arguments of a method call.
func analyzeVariadicArguments[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
	TC config.Component,
	BC configbuilder.ComponentBuilder[TC],
](
	ctx *configurerCallContext[T, E, B],
	buildChild func(func(BC)),
	analyzeChild func(*configurerCallContext[T, E, B], BC, ssa.Value),
) {
	allocs := collectVariadicAllocations(
		ctx.Builder,
		ctx.Args[len(ctx.Args)-1], // varadic slice is always the last argument
	)

	var isVarArg func(v ssa.Value) (optional.Optional[int], bool)

	isVarArg = func(v ssa.Value) (optional.Optional[int], bool) {
		if allocs.Has(v) {
			return optional.Some(0), true
		}

		switch v := v.(type) {
		case *ssa.Slice:
			if index, ok := isVarArg(v.X); ok {
				if v.Low == nil {
					return index, true
				}

				return optional.Sum(
					index,
					ssax.AsInt(v.Low),
				), true
			}
		case *ssa.IndexAddr:
			if index, ok := isVarArg(v.X); ok {
				return optional.Sum(
					index,
					ssax.AsInt(v.Index),
				), true
			}
		default:
			unimplementedAnalysis(ctx.Builder, v)
		}

		return optional.None[int](), false
	}

	indexCounts := map[int]int{}
	hasUnknownIndices := false
	var children []func(BC)

	for block := range ssax.WalkFunc(ctx.FunctionUnderAnalysis) {
		if !ssax.PathExists(block, ctx.Instruction.Block()) {
			continue
		}

		unconditional := ssax.IsUnconditional(block)

		for inst := range ssax.InstructionsBefore(block, ctx.Instruction) {
			if inst, ok := inst.(*ssa.Store); ok {
				index, ok := isVarArg(inst.Addr)
				if !ok {
					continue
				}

				if i, ok := index.TryGet(); ok {
					indexCounts[i]++
				} else {
					hasUnknownIndices = true
				}

				children = append(children, func(b BC) {
					ctx.Apply(b)

					if hasUnknownIndices || !unconditional {
						b.Speculative()
					} else if i, ok := index.TryGet(); ok && indexCounts[i] > 1 {
						b.Speculative()
					}

					analyzeChild(ctx, b, inst.Val)
				})
			}
		}
	}

	for _, child := range children {
		buildChild(child)
	}
}

func collectVariadicAllocations(
	b configbuilder.UntypedComponentBuilder,
	v ssa.Value,
) *sets.Set[ssa.Value] {
	allocs := sets.New[ssa.Value]()

	var collect func(v ssa.Value)
	collect = func(v ssa.Value) {
		switch v := v.(type) {
		case *ssa.Alloc:
			allocs.Add(v)

		case *ssa.Slice:
			collect(v.X)

		case *ssa.Const:
			// We've found a nil slice.

		case *ssa.Phi:
			for _, edge := range v.Edges {
				collect(edge)
			}

		case *ssa.Call:
			call := v.Common()

			if fn, ok := call.Value.(*ssa.Builtin); ok {
				if fn.Name() == "append" {
					for _, arg := range call.Args {
						collect(arg)
					}
				}
			}

		default:
			unimplementedAnalysis(b, v)
		}
	}

	collect(v)

	return allocs
}
