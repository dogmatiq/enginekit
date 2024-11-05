package staticconfig

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
	"golang.org/x/tools/go/ssa"
)

type variadicConfigurerCallContext[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
	TC config.Component,
	BC configbuilder.ComponentBuilder[TC],
] struct {
	*configurerCallContext[T, E, B]

	BuildChild   func(func(BC))
	AnalyzeChild func(*configurerCallContext[T, E, B], BC, ssa.Value)

	seen map[ssa.Value]struct{}
}

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
	walkUpVariadic(
		&variadicConfigurerCallContext[T, E, B, TC, BC]{
			configurerCallContext: ctx,
			BuildChild:            buildChild,
			AnalyzeChild:          analyzeChild,
			seen:                  map[ssa.Value]struct{}{},
		},
		ctx.Args[len(ctx.Args)-1], // varadic slice is always the last argument
	)
}

func walkUpVariadic[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
	TC config.Component,
	BC configbuilder.ComponentBuilder[TC],
](
	ctx *variadicConfigurerCallContext[T, E, B, TC, BC],
	v ssa.Value,
) {
	if _, ok := ctx.seen[v]; ok {
		return
	}
	ctx.seen[v] = struct{}{}

	switch v := v.(type) {
	default:
		unimplementedAnalysis(ctx.Builder, v)

	case *ssa.Const:
		// We've found a nil slice.

	case *ssa.Phi:
		for _, edge := range v.Edges {
			walkUpVariadic(ctx, edge)
		}

	case *ssa.Alloc:
		walkDownVariadic(ctx, v)

	case *ssa.Slice:
		walkUpVariadic(ctx, v.X)

	case *ssa.Call:
		call := v.Common()

		if fn, ok := call.Value.(*ssa.Builtin); ok {
			if fn.Name() == "append" {
				for _, arg := range call.Args {
					walkUpVariadic(ctx, arg)
				}
			}
		}
	}
}

func walkDownVariadic[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
	TC config.Component,
	BC configbuilder.ComponentBuilder[TC],
](
	ctx *variadicConfigurerCallContext[T, E, B, TC, BC],
	alloc *ssa.Alloc,
) {
	for block := range ssax.WalkBlock(alloc.Block()) {
		unconditional := ssax.IsUnconditional(block)

		for inst := range ssax.InstructionsBefore(block, ctx.Instruction) {
			switch inst := inst.(type) {
			case *ssa.Store:
				addr, ok := inst.Addr.(*ssa.IndexAddr)
				if !ok {
					continue
				}

				if addr.X != alloc {
					continue
				}

				ctx.BuildChild(func(b BC) {
					ctx.Apply(b)

					if !unconditional {
						b.Speculative()
					}

					ctx.AnalyzeChild(ctx.configurerCallContext, b, inst.Val)
				})
			}
		}
	}
}
