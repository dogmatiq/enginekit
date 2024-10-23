package ssax

import (
	"iter"

	"golang.org/x/tools/go/ssa"
)

// WalkFunc recursively yields all reachable blocks in the given function.
func WalkFunc(fn *ssa.Function) iter.Seq[*ssa.BasicBlock] {
	return func(yield func(*ssa.BasicBlock) bool) {
		if len(fn.Blocks) != 0 {
			for b := range WalkBlock(fn.Blocks[0]) {
				if !yield(b) {
					return
				}
			}
		}
	}
}

// WalkBlock recursively yields b and all reachable successor blocks of b.
//
// A block is considered reachable if there is a control flow path from b to
// that block that does not depend on a condition that is known to be false at
// compile-time.
func WalkBlock(b *ssa.BasicBlock) iter.Seq[*ssa.BasicBlock] {
	return walk(b, DirectSuccessors)
}

// DirectSuccessors yields the reachable direct successors of b.
func DirectSuccessors(b *ssa.BasicBlock) iter.Seq[*ssa.BasicBlock] {
	return func(yield func(*ssa.BasicBlock) bool) {
		successors := b.Succs

		if inst, ok := Terminator[*ssa.If](b).TryGet(); ok {
			if cond, ok := AsBool(inst.Cond).TryGet(); ok {
				if cond {
					successors = b.Succs[:1]
				} else {
					successors = b.Succs[1:]
				}
			}
		}

		for _, s := range successors {
			if !yield(s) {
				return
			}
		}
	}
}

// IsUnconditional returns true if all control-flow paths through the function
// containing b pass through b at some point.
func IsUnconditional(b *ssa.BasicBlock) bool {
	return UnconditionalPathExists(b.Parent().Blocks[0], b)
}

// PathExists returns true if, after dead-code elimnation, it's possible to
// traverse the control-flow graph from one specific node to another.
func PathExists(from, to *ssa.BasicBlock) bool {
	if from.Parent() != to.Parent() {
		panic("blocks are not in the same function")
	}

	for b := range WalkBlock(from) {
		if b == to {
			return true
		}
	}

	return false
}

// UnconditionalPathExists returns true if, after dead-code elimination, all
// control flow paths from one node always lead to another.
func UnconditionalPathExists(from, to *ssa.BasicBlock) bool {
	if from.Parent() != to.Parent() {
		panic("blocks are not in the same function")
	}

	seen := map[*ssa.BasicBlock]struct{}{}

	var exists func(*ssa.BasicBlock) bool
	exists = func(from *ssa.BasicBlock) bool {
		if _, ok := seen[from]; ok {
			return true
		}
		seen[from] = struct{}{}

		if from == to {
			return true
		}

		if len(from.Succs) == 0 {
			return false
		}

		for s := range DirectSuccessors(from) {
			if !exists(s) {
				return false
			}
		}

		return true
	}

	return exists(from)
}

// walk recursively yields b, and the blocks yielded by next(b). It stops
// recursing when a cycle is detected.
func walk(
	b *ssa.BasicBlock,
	next func(*ssa.BasicBlock) iter.Seq[*ssa.BasicBlock],
) iter.Seq[*ssa.BasicBlock] {
	return func(yield func(*ssa.BasicBlock) bool) {
		seen := map[*ssa.BasicBlock]struct{}{}

		var emit func(*ssa.BasicBlock) bool
		emit = func(b *ssa.BasicBlock) bool {
			if _, ok := seen[b]; ok {
				return true
			}

			seen[b] = struct{}{}
			if !yield(b) {
				return false
			}

			for n := range next(b) {
				if !emit(n) {
					return false
				}
			}

			return true
		}

		emit(b)
	}
}
