package staticconfig

import (
	"go/constant"
	"go/types"
	"iter"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"golang.org/x/tools/go/ssa"
)

// analyzeConfigurerCalls analyzes the calls to the "configurer" that is passed
// to t's "Configure()" method.
//
// Any calls that are not recognized are yielded.
func findConfigurerCalls(ctx *context, t types.Type) iter.Seq[*ssa.CallCommon] {
	configure := ctx.LookupMethod(t, "Configure")
	return findConfigurerCallsInFunc(configure, 1)
}

func analyzeIdentityCall(
	b configbuilder.EntityBuilder,
	call *ssa.CallCommon,
) {
	b.Identity(func(b *configbuilder.IdentityBuilder) {
		if name, ok := call.Args[0].(*ssa.Const); ok {
			b.SetName(constant.StringVal(name.Value))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}

		if key, ok := call.Args[1].(*ssa.Const); ok {
			b.SetKey(constant.StringVal(key.Value))
		} else {
			b.UpdateFidelity(config.Incomplete)
		}
	})
}

// indices refers to the positions of arguments that are the configurer. If none
// are provided it defaults to [1]. This accounts for the most common case where
// fn is the Configure() method on an application or handler. In this case the
// first parameter is the receiver, so the second parameter is the configurer
// itself.

// The instantiatedTypes map is used to store the types that have been
// instantiated in the function. This is necessary because the SSA
// representation of a function does not include type information for the
// arguments, so we need to track this information ourselves. The keys are the
// names of the type parameters and the values are the concrete types that have
// been instantiated.

// findConfigurerCallsInFunc yields all call to methods on the Dogma application or
// handler "configurer" within the given function.
func findConfigurerCallsInFunc(
	fn *ssa.Function,
	indices ...int,
) iter.Seq[*ssa.CallCommon] {
	isConfigurerCall := func(call *ssa.CallCommon) bool {
		for _, i := range indices {
			if call.Value == fn.Params[i] {
				return true
			}
		}
		return false
	}

	return func(yield func(*ssa.CallCommon) bool) {
		for _, block := range fn.Blocks {
			for _, inst := range block.Instrs {
				inst, ok := inst.(*ssa.Call)
				if !ok {
					continue
				}

				call := inst.Common()

				if isConfigurerCall(call) {
					// We've found a direct call to a method on the configurer.
					if !yield(call) {
						return
					}
					// } else {
					// 	// We've found a call to some other function or method. We need
					// 	// to analyse the instructions within *that* function to see if
					// 	// *it* makes any calls to the configurer.
					// 	if !e.yieldIndirectCalls(call, yield) {
					// 		return false
					// 	}
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
