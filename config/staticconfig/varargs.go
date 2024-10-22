package staticconfig

import (
	"go/token"
	"go/types"

	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
	"golang.org/x/tools/go/ssa"
)

func findAllocation(v ssa.Value) (*ssa.Alloc, bool) {
	// fmt.Println("***", v, reflect.TypeOf(v))
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

func resolveVariadic(
	_ configbuilder.EntityBuilder,
	call configurerCall,
) ([]ssa.Value, bool) {
	n := len(call.Args) - 1
	variadics := call.Args[n]

	array, ok := findAllocation(variadics)
	if !ok {
		return nil, false
	}

	size := array.Type().Underlying().(*types.Pointer).Elem().(*types.Array).Len()
	result := make([]ssa.Value, size)

	for b := range ssax.WalkDown(array.Block()) {
		if !ssax.PathExists(b, call.Instruction.Block()) {
			continue
		}

		for inst := range ssax.InstructionsBefore(b, call.Instruction) {
			switch inst := inst.(type) {
			case *ssa.Store:
				if i, ok := isIndexOfArray(array, inst.Addr); ok {
					result[i] = inst.Val
				}
			}
		}
	}

	return result, true
}
