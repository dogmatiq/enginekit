package ssax

import (
	"iter"

	"github.com/dogmatiq/enginekit/optional"
	"golang.org/x/tools/go/ssa"
)

// Terminator returns the final "transfer of control" instruction in the given
// block.
//
// If the block does not contain any instructions (as is the case for external
// functions), or the terminator instruction is not of type T, ok is false.
//
// The instruction is always [ssa.If], [ssa.Jump], [ssa.Return], or [ssa.Panic].
func Terminator[T ssa.Instruction](b *ssa.BasicBlock) optional.Optional[T] {
	return optional.As[T](optional.Last(b.Instrs))
}

// InstructionsBefore yields all instructions in the block that precede the
// given instruction.
//
// It yields all instructions if inst is not in b.
func InstructionsBefore(b *ssa.BasicBlock, inst ssa.Instruction) iter.Seq[ssa.Instruction] {
	return func(yield func(ssa.Instruction) bool) {
		for _, x := range b.Instrs {
			if x == inst {
				return
			}

			if !yield(x) {
				return
			}
		}
	}
}
