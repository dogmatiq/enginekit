package staticconfig

import (
	"go/token"
	"iter"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
	"golang.org/x/tools/go/ssa"
)

func findAllocation(v ssa.Value) (*ssa.Alloc, bool) {
	switch v := v.(type) {
	case *ssa.Alloc:
		return v, true

	case *ssa.Slice:
		return findAllocation(v.X)

	case *ssa.UnOp:
		if v.Op == token.MUL { // pointer de-reference
			return findAllocation(v.X)
		}
		return nil, false

	default:
		return nil, false
	}
}

func isIndexOfArray(
	array *ssa.Alloc,
	v ssa.Value,
) (int, bool) {
	switch v := v.(type) {
	case *ssa.IndexAddr:
		if v.X != array {
			return 0, false
		}
		return ssax.AsInt(v.Index).TryGet()
	}
	return 0, false
}

func resolveVariadic[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
](
	b B,
	inst ssa.CallInstruction,
) iter.Seq[ssa.Value] {
	return func(yield func(ssa.Value) bool) {
		call := inst.Common()

		variadics := call.Args[len(call.Args)-1]
		if ssax.IsZeroValue(variadics) {
			return
		}

		array, ok := findAllocation(variadics)
		if !ok {
			b.Partial()
			return
		}

		for b := range ssax.WalkBlock(array.Block()) {
			if !ssax.PathExists(b, inst.Block()) {
				continue
			}

			for inst := range ssax.InstructionsBefore(b, inst) {
				switch inst := inst.(type) {
				case *ssa.Store:
					if _, ok := isIndexOfArray(array, inst.Addr); ok {
						if !yield(inst.Val) {
							return
						}
					}
				}
			}
		}
	}
}
