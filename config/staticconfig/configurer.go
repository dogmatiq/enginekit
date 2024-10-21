package staticconfig

import (
	"go/types"
	"iter"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"golang.org/x/tools/go/ssa"
)

// configureContext is a specialization of [context] that is used when analyzing
// a Configure() method.
type configureContext struct {
	*context

	Func              *ssa.Function
	Builder           configbuilder.EntityBuilder
	ConfigurerIndices []int
}

func (c *configureContext) IsConfigurer(v ssa.Value) bool {
	for _, i := range c.ConfigurerIndices {
		if v == c.Func.Params[i] {
			return true
		}
	}

	return false
}

type configurerCall struct {
	*ssa.CallCommon

	Fidelity config.Fidelity
}

// analyzeConfigurerCalls analyzes the calls to the "configurer" that is passed
// to t's "Configure()" method.
//
// Any calls that are not recognized are yielded.
func findConfigurerCalls(
	ctx *context,
	b configbuilder.EntityBuilder,
	t types.Type,
) iter.Seq[configurerCall] {
	configure := ctx.LookupMethod(t, "Configure")

	return func(yield func(configurerCall) bool) {
		emitConfigurerCallsInFunc(
			&configureContext{
				context:           ctx,
				Func:              configure,
				Builder:           b,
				ConfigurerIndices: []int{1},
			},
			configure,
			yield,
		)
	}
}

// emitConfigurerCallsInFunc yields all call to methods on the Dogma application
// or handler "configurer" within the given function.
//
// indices is a list of the positions of parameters to fn that are the
// configurer.
func emitConfigurerCallsInFunc(
	ctx *configureContext,
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
	ctx *configureContext,
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
	ctx *configureContext,
	call ssa.CallInstruction,
	yield func(configurerCall) bool,
) bool {
	com := call.Common()

	if com.IsInvoke() && ctx.IsConfigurer(com.Value) {
		// We've found a direct call to a method on the configurer.
		var f config.Fidelity
		if isConditional(call.Block()) {
			f |= config.Speculative
		}

		return yield(configurerCall{com, f})
	}

	// We've found a call to some function or method that does not belong to the
	// configurer. If any of the arguments are the configurer we analyze the
	// called function as well.
	//
	// This is an quite naive implementation. There are other ways that the
	// callee could gain access to the configurer. For example, it could be
	// passed inside a context, or assigned to a field within the entity struct.
	//
	// First, we build a list of the indices of arguments that are the
	// configurer. It doesn't make much sense, but the configurer could be
	// passed in multiple positions.
	var indices []int
	for i, arg := range com.Args {
		if ctx.IsConfigurer(arg) {
			indices = append(indices, i)
		}
	}

	// If none of the arguments are the configurer, we can skip analyzing the
	// callee. This prevents us from analyzing the entire program.
	if len(indices) == 0 {
		return true
	}

	// If we can't obtain the callee, this is a call to an interface method, or
	// some other un-analyzable function.
	fn := com.StaticCallee()
	if fn == nil {
		ctx.Builder.UpdateFidelity(config.Incomplete)
		return true
	}

	return emitConfigurerCallsInFunc(
		&configureContext{
			context:           ctx.context,
			Func:              fn,
			Builder:           ctx.Builder,
			ConfigurerIndices: indices,
		},
		fn,
		yield,
	)
}
