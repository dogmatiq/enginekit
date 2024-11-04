package staticconfig

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
	"golang.org/x/tools/go/ssa"
)

func analyzeVariadicArguments[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
	TChild config.Component,
	BChild configbuilder.ComponentBuilder[TChild],
](
	ctx *configurerCallContext[T, E, B],
	child func(func(BChild)),
	analyze func(*context, BChild, ssa.Value),
) {
	// The variadic slice parameter is always the last argument.
	varargs := ctx.Args[len(ctx.Args)-1]

	if ssax.IsZeroValue(varargs) {
		return
	}

	array, ok := findSliceArrayAllocation(varargs)
	if !ok {
		ctx.Builder.Partial()
		return
	}

	buildersByIndex := map[int][]BChild{}

	for block := range ssax.WalkBlock(array.Block()) {
		// If there's no path from this block to the call instruction, we can
		// safely ignore it, even if it modifies the underlying array.
		if !ssax.PathExists(block, ctx.Instruction.Block()) {
			continue
		}

		for inst := range ssax.InstructionsBefore(block, ctx.Instruction) {
			switch inst := inst.(type) {
			case *ssa.Store:
				if addr, ok := inst.Addr.(*ssa.IndexAddr); ok && addr.X == array {
					child(func(b BChild) {
						if index, ok := ssax.AsInt(addr.Index).TryGet(); ok {
							// If there are multiple writes to the same index,
							// we mark them all as speculative.
							//
							// TODO: Could we handle this more intelligently by
							// using the value of the store instruction closest
							// to the call instruction?
							conflicting := buildersByIndex[index]
							if len(conflicting) == 1 {
								conflicting[0].Speculative()
							}
							if len(conflicting) != 0 {
								b.Speculative()
							}
							buildersByIndex[index] = append(conflicting, b)
						} else {
							// If we can't resolve the index we assume the child
							// is speculative because we can't tell if it is
							// ever overwritten with a different value.
							b.Speculative()
						}

						if ctx.IsSpeculative {
							b.Speculative()
						}

						analyze(ctx.context, b, inst.Val)
					})
				}
			}
		}
	}
}

// findSliceArrayAllocation returns the underlying array allocation of a slice.
func findSliceArrayAllocation(v ssa.Value) (*ssa.Alloc, bool) {
	switch v := v.(type) {
	case *ssa.Alloc:
		return v, true
	case *ssa.Slice:
		return findSliceArrayAllocation(v.X)
	default:
		return nil, false
	}
}
