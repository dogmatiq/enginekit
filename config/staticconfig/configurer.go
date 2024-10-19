package staticconfig

import (
	"go/types"
	"iter"

	"github.com/dogmatiq/enginekit/config"
	"golang.org/x/tools/go/ssa"
)

type configurerCall struct {
	*ssa.CallCommon

	Fidelity config.Fidelity
}

// analyzeConfigurerCalls analyzes the calls to the "configurer" that is passed
// to t's "Configure()" method.
//
// Any calls that are not recognized are yielded.
func findConfigurerCalls(ctx *context, t types.Type) iter.Seq[configurerCall] {
	configure := ctx.LookupMethod(t, "Configure")
	return findConfigurerCallsInFunc(configure, 1)
}

// findConfigurerCallsInFunc yields all call to methods on the Dogma application
// or handler "configurer" within the given function.
//
// indices is a list of the positions of parameters to fn that are the
// configurer.
func findConfigurerCallsInFunc(
	fn *ssa.Function,
	indices ...int,
) iter.Seq[configurerCall] {
	isConfigurerCall := func(call *ssa.CallCommon) bool {
		for _, i := range indices {
			if call.Value == fn.Params[i] {
				return true
			}
		}
		return false
	}

	return func(yield func(configurerCall) bool) {
		for block := range walkReachable(fn.Blocks[0]) {
			var f config.Fidelity
			if isConditional(fn, block) {
				f |= config.Speculative
			}

			for _, inst := range block.Instrs {
				inst, ok := inst.(*ssa.Call)
				if !ok {
					continue
				}

				call := inst.Common()

				if isConfigurerCall(call) {
					// We've found a direct call to a method on the configurer.
					if !yield(configurerCall{call, f}) {
						return
					}
				}
			}
		}
	}
}

// func (e *entity) yieldIndirectCalls(
// 	call *ssa.CallCommon,
// 	yield func(*ssa.CallCommon) bool,
// ) bool {
// 	// com := call.Common()

// 	// var indices []int
// 	// for i, arg := range com.Args {
// 	// 	if _, ok := configurers[arg]; ok {
// 	// 		indices = append(indices, i)
// 	// 	}
// 	// }

// 	// if len(indices) == 0 {
// 	// 	return nil
// 	// }

// 	// if com.IsInvoke() {
// 	// 	t, ok := instantiatedTypes[com.Value.Type().String()]
// 	// 	if !ok {
// 	// 		// If we cannot find any instantiated types in mapping, most likely
// 	// 		// we hit the interface method and cannot analyze any further.
// 	// 		return nil
// 	// 	}

// 	// 	return findConfigurerCalls(
// 	// 		prog,
// 	// 		prog.LookupMethod(t, com.Method.Pkg(), com.Method.Name()),
// 	// 		instantiatedTypes,
// 	// 		// don't pass indices here, as we are already in the method.
// 	// 	)
// 	// }

// 	return e.yieldCalls(
// 		call.StaticCallee(),
// 		yield,
// 	)
// }
