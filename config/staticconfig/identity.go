package staticconfig

import (
	"go/constant"
	"go/token"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"golang.org/x/tools/go/ssa"
)

func analyzeIdentityCall(
	b configbuilder.EntityBuilder,
	call configurerCall,
) {
	b.Identity(func(b *configbuilder.IdentityBuilder) {
		b.UpdateFidelity(call.Fidelity)

		if name, ok := resolveValue(call.Args[0]); ok {
			b.SetName(constant.StringVal(name))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}

		if key, ok := resolveValue(call.Args[1]); ok {
			b.SetKey(constant.StringVal(key))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}
	})
}

func resolveValue(
	v ssa.Value,
) (constant.Value, bool) {
	switch v := v.(type) {
	case *ssa.Const:
		return v.Value, true
	case ssa.Instruction:
		values := resolveExpr(v)
		switch len(values) {
		case 0:
			return nil, false
		case 1:
			return values[0], values[0] != nil
		default:
			panic("did not expect multiple values")
		}
	default:
		return nil, false
	}
}

func resolveExpr(
	inst ssa.Instruction,
) []constant.Value {
	switch inst := inst.(type) {
	case *ssa.Call:
		return resolveReturnValues(inst.Common())
	case *ssa.Extract:
		if expr, ok := inst.Tuple.(ssa.Instruction); ok {
			values := resolveExpr(expr)
			return values[inst.Index : inst.Index+1]
		}
		return nil
	default:
		return nil
	}
}

func resolveReturnValues(call *ssa.CallCommon) []constant.Value {
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
		inst, ok := b.Instrs[len(b.Instrs)-1].(*ssa.Return)
		if !ok {
			continue
		}

		var values []constant.Value

		for _, v := range inst.Results {
			if x, ok := resolveValue(v); ok {
				values = append(values, x)
			} else {
				values = append(values, nil)
			}
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
