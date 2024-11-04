package ssax

import (
	"go/constant"
	"go/token"

	"github.com/dogmatiq/enginekit/optional"
	"golang.org/x/tools/go/ssa"
)

// StaticValue returns the singular value of v.
//
// If v cannot be resolved to a single value, it returns an empty optional.
func StaticValue(v ssa.Value) optional.Optional[ssa.Value] {
	values := staticValues(v)
	if len(values) > 1 {
		panic("did not expect multiple values")
	}

	if len(values) == 1 {
		return values[0]
	}

	return optional.None[ssa.Value]()
}

// staticValues returns the static value(s) that result from evaluating the
// given node.
//
// If an individual value within the expression cannot be resolved to a singular
// static value, it is represented as an empty optional in the returned slice.
//
// It returns an empty slice if the expression itself cannot be resolved.
func staticValues(v ssa.Value) []optional.Optional[ssa.Value] {
	switch v := v.(type) {
	case *ssa.Const:
		return optional.Slice[ssa.Value](v)

	case *ssa.Call:
		return staticValuesFromCall(v.Common())

	case *ssa.Extract:
		values := staticValues(v.Tuple)
		if len(values) <= v.Index {
			return nil
		}
		return values[v.Index : v.Index+1]

	case *ssa.MakeInterface:
		return staticValues(v.X)

	case *ssa.UnOp:
		if v.Op == token.MUL { // pointer de-reference
			return staticValues(v.X)
		}
	}

	// TODO(jmalloc): This implementation is incomplete.
	return nil
}

// staticValuesFromCall returns the static value(s) that result from evaluating
// a call to a function.
//
// If an individual value within the expression cannot be resolved to a singular
// static value, it is represented as an empty value in the returned slice.
//
// It returns an empty slice if the function itself cannot be resolved. For
// example, if it is a dynamic call to an interface method.
func staticValuesFromCall(call *ssa.CallCommon) []optional.Optional[ssa.Value] {
	// TODO: we could use StaticValue or some variant thereof to resolve the
	// callee in more cases.
	fn := call.StaticCallee()
	if fn == nil {
		// A call to an interface method.
		return nil
	}

	if len(fn.Blocks) == 0 {
		// Probably an external C function.
		return nil
	}

	n := fn.Signature.Results().Len()

	if n == 0 {
		// The function does not have any output parameters.
		return nil
	}

	outputs := make([]optional.Optional[ssa.Value], n)
	conflicting := make([]bool, n)

	for b := range WalkBlock(fn.Blocks[0]) {
		ret, ok := Terminator[*ssa.Return](b).TryGet()
		if !ok {
			continue
		}

		for i, v := range ret.Results {
			if conflicting[i] {
				continue
			}

			v := StaticValue(v)
			if !v.IsPresent() {
				continue
			}

			x := outputs[i]
			if !x.IsPresent() {
				outputs[i] = v
				continue
			}

			if !equal(x.Get(), v.Get()) {
				conflicting[i] = true
				outputs[i] = optional.None[ssa.Value]()
			}
		}
	}

	return outputs
}

// equal returns true if a and b refer to the same value, or are equal constant
// values.
func equal(a, b ssa.Value) bool {
	if a == b {
		return true
	}

	if a, ok := a.(*ssa.Const); ok {
		if b, ok := b.(*ssa.Const); ok {
			return constant.Compare(a.Value, token.EQL, b.Value)
		}
	}

	return false
}
