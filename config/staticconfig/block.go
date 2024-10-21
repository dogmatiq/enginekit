package staticconfig

import (
	"go/constant"
	"iter"

	"golang.org/x/tools/go/ssa"
)

// walkReachable yields all blocks reachable from b.
func walkReachable(b *ssa.BasicBlock) iter.Seq[*ssa.BasicBlock] {
	return func(yield func(*ssa.BasicBlock) bool) {
		yielded := map[*ssa.BasicBlock]struct{}{}

		var walk func(*ssa.BasicBlock) bool
		walk = func(b *ssa.BasicBlock) bool {
			if _, ok := yielded[b]; ok {
				return true
			}

			yielded[b] = struct{}{}

			if !yield(b) {
				return false
			}

			for succ := range reachableSuccessors(b) {
				if !walk(succ) {
					return false
				}
			}

			return true
		}

		walk(b)
	}
}

// reachableSuccessors yields the successors of b that are actually reachable.
func reachableSuccessors(b *ssa.BasicBlock) iter.Seq[*ssa.BasicBlock] {
	return func(yield func(*ssa.BasicBlock) bool) {
		if inst, ok := transferOfControl[*ssa.If](b); ok {
			if v := staticValue(inst.Cond); v != nil {
				if constant.BoolVal(v) {
					yield(b.Succs[0])
				} else {
					yield(b.Succs[1])
				}

				return
			}
		}

		for _, succ := range b.Succs {
			if !yield(succ) {
				return
			}
		}
	}
}

// isConditional returns true if there is any control flow path through the
// function that does NOT pass through b.
func isConditional(b *ssa.BasicBlock) bool {
	return !isInevitable(b.Parent().Blocks[0], b)
}

// isInevitable returns true if all paths out of "from" pass through "to".
func isInevitable(from, to *ssa.BasicBlock) bool {
	if from == to {
		return true
	}

	if len(from.Succs) == 0 {
		return false
	}

	for succ := range reachableSuccessors(from) {
		if succ == from {
			continue
		}

		if !isInevitable(succ, to) {
			return false
		}
	}

	return true
}

// transferOfControl returns the last instruction in b, if it is of type T.
func transferOfControl[T ssa.Instruction](b *ssa.BasicBlock) (inst T, ok bool) {
	if len(b.Instrs) == 0 {
		return inst, false
	}

	inst, ok = b.Instrs[len(b.Instrs)-1].(T)
	return inst, ok
}
