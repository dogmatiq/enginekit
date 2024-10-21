package staticconfig

import (
	"go/constant"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

// staticValue returns the constant value of v, if it's possible to obtain;
// otherwise, it returns nil.
func staticValue(v ssa.Value) constant.Value {
	switch v := v.(type) {
	case *ssa.Const:
		return v.Value
	case ssa.Instruction:
		values := staticValueOfInstruction(v)
		switch len(values) {
		case 0:
			return nil
		case 1:
			return values[0]
		default:
			panic("did not expect multiple values")
		}
	default:
		return nil
	}
}

// staticValueOfInstruction returns the constant value(s) of an instruction.
//
// If an individual value within the expression cannot be resolved, it is
// represented as a nil value in the returned slice.
//
// It returns an empty slice if the expression itself cannot be resolved.
func staticValueOfInstruction(inst ssa.Instruction) []constant.Value {
	switch inst := inst.(type) {
	case *ssa.Call:
		return staticReturnValues(inst.Common())
	case *ssa.Extract:
		if expr, ok := inst.Tuple.(ssa.Instruction); ok {
			values := staticValueOfInstruction(expr)
			return values[inst.Index : inst.Index+1]
		}
		return nil
	default:
		return nil
	}
}

// staticReturnValues returns the constant values returned by a function.
//
// If an invividual value cannot be resolved, it is represented as a nil value
// in the returned slice. The function must return the same value on all control
// paths for a value to be considered constant.
//
// It returns nil if the function itself cannot be resolved. For example, if it
// is a dynamic call to an interface method.
func staticReturnValues(call *ssa.CallCommon) []constant.Value {
	fn := call.StaticCallee()
	if fn == nil {
		return nil
	}

	if len(fn.Blocks) == 0 {
		return nil
	}

	if fn.Signature.Results().Len() == 0 {
		return nil
	}

	var results []constant.Value

	for b := range walkReachable(fn.Blocks[0]) {
		inst, ok := transferOfControl[*ssa.Return](b)
		if !ok {
			continue
		}

		var values []constant.Value

		for _, v := range inst.Results {
			values = append(values, staticValue(v))
		}

		if results == nil {
			results = values
		} else {
			for i, a := range values {
				b := results[i]
				if constant.Compare(a, token.NEQ, b) {
					results[i] = nil
				}
			}
		}
	}

	return results
}
