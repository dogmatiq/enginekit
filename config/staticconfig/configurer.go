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

	ctx = ctx.NewChild(
		func(v ssa.Value) bool {
			return v == configure.Params[1]
		},
	)

	return func(yield func(configurerCall) bool) {
		emitConfigurerCallsInFunc(ctx, configure, yield)
	}
}

// emitConfigurerCallsInFunc yields all call to methods on the Dogma application
// or handler "configurer" within the given function.
//
// indices is a list of the positions of parameters to fn that are the
// configurer.
func emitConfigurerCallsInFunc(
	ctx *context,
	fn *ssa.Function,
	yield func(configurerCall) bool,
) bool {
	if len(fn.Blocks) == 0 {
		return true
	}

	for block := range walkReachable(fn.Blocks[0]) {
		for _, inst := range block.Instrs {
			if !emitConfigurerCallsInInstruction(ctx, inst, yield) {
				return false
			}
		}
	}

	return true
}

func emitConfigurerCallsInInstruction(
	ctx *context,
	inst ssa.Instruction,
	yield func(configurerCall) bool,
) bool {
	switch inst := inst.(type) {
	case ssa.CallInstruction:
		return emitConfigurerCallsInCallInstruction(ctx, inst, yield)
	default:
		return true
	}
}

func emitConfigurerCallsInCallInstruction(
	ctx *context,
	call ssa.CallInstruction,
	yield func(configurerCall) bool,
) bool {
	com := call.Common()

	if com.IsInvoke() {
		// We're invoking a method on an interface, that is, we don't know the
		// concrete type. If it's not a call to a method on the configurer,
		// there's nothing more we can analyze.
		if !ctx.IsConfigurer(com.Value) {
			return true
		}

		// We've found a direct call to a method on the configurer.
		var f config.Fidelity
		if isConditional(call.Block()) {
			f |= config.Speculative
		}

		return yield(configurerCall{com, f})
	}

	// We've found a call to some other function or method.
	//
	// If any of the parameters refer to the configurer, we need to analyze
	// _that_ function.
	//
	// This is an native implementation. There are other ways that this function
	// could gain access to the configurer. For example, it could be passed
	// inside a context, or assigned to a field within the entity struct.
	fn := com.StaticCallee()

	// Check at which argument indices the configurer is passed to the function.
	var indices []int
	for i, arg := range com.Args {
		if ctx.IsConfigurer(arg) {
			indices = append(indices, i)
		}
	}

	// Don't analyze fn if the configurer is not passed as an argument.
	if len(indices) == 0 {
		return true
	}

	return emitConfigurerCallsInFunc(
		ctx.NewChild(
			func(v ssa.Value) bool {
				for _, i := range indices {
					if v == fn.Params[i] {
						return true
					}
				}

				return false
			},
		),
		fn,
		yield,
	)
}
