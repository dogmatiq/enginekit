package staticconfig

import (
	"go/types"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/config/staticconfig/internal/ssax"
	"github.com/dogmatiq/enginekit/internal/typename"
	"golang.org/x/tools/go/ssa"
)

type entityContext[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
] struct {
	*context

	EntityType             types.Type
	Builder                B
	ConfigureMethod        *ssa.Function
	FunctionUnderAnalysis  *ssa.Function
	ConfigurerParamIndices []int
}

func (c *entityContext[T, E, B]) IsConfigurer(v ssa.Value) bool {
	for _, i := range c.ConfigurerParamIndices {
		if v == c.FunctionUnderAnalysis.Params[i] {
			return true
		}
	}
	return false
}

type configurerCallContext[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
] struct {
	*entityContext[T, E, B]
	*ssa.CallCommon

	Instruction   ssa.CallInstruction
	IsSpeculative bool
}

// configurerCallAnalyzer is a function that analyzes a call to a method on an
// entity's configurer.
type configurerCallAnalyzer[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
] func(*configurerCallContext[T, E, B])

// analyzeEntity analyzes the Configure() method of the type t, which must be a
// Dogma application or handler.
//
// It calls the analyze function for each call to a method on the configurer,
// other than Identity() which is handled the same in all cases.
func analyzeEntity[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
](
	ctx *context,
	t types.Type,
	builder B,
	analyze configurerCallAnalyzer[T, E, B],
) {
	builder.TypeName(typename.OfStatic(t))
	configure := ctx.LookupMethod(t, "Configure")

	ectx := &entityContext[T, E, B]{
		context:                ctx,
		EntityType:             t,
		Builder:                builder,
		ConfigureMethod:        configure,
		FunctionUnderAnalysis:  configure,
		ConfigurerParamIndices: []int{1},
	}

	fn := func(ctx *configurerCallContext[T, E, B]) {
		switch ctx.Method.Name() {
		case "Identity":
			analyzeIdentity(ctx)
		default:
			analyze(ctx)
		}
	}

	analyzeConfigurerCallsInFunc(
		ectx,
		fn,
	)
}

// analyzeConfigurerCallsInFunc analyzes calls to methods on the configurer in
// the function under analysis.
func analyzeConfigurerCallsInFunc[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
](
	ctx *entityContext[T, E, B],
	analyze configurerCallAnalyzer[T, E, B],
) {
	for b := range ssax.WalkFunc(ctx.FunctionUnderAnalysis) {
		for _, inst := range b.Instrs {
			if inst, ok := inst.(ssa.CallInstruction); ok {
				analyzeConfigurerCallsInInstruction(ctx, inst, analyze)
			}
		}
	}
}

// analyzeConfigurerCallsInInstruction analyzes calls to methods on the
// configurer in the given instruction.
func analyzeConfigurerCallsInInstruction[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
](
	ctx *entityContext[T, E, B],
	inst ssa.CallInstruction,
	analyze configurerCallAnalyzer[T, E, B],
) {
	com := inst.Common()

	if com.IsInvoke() && ctx.IsConfigurer(com.Value) {
		analyze(&configurerCallContext[T, E, B]{
			entityContext: ctx,
			CallCommon:    com,
			Instruction:   inst,
			IsSpeculative: !ssax.IsUnconditional(inst.Block()),
		})
		return
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

	// We don't analyze the callee if it is not passed the configurer.
	if len(indices) == 0 {
		return
	}

	// If we can't obtain the callee this is a call to an interface method or
	// some other un-analyzable function.
	fn := com.StaticCallee()
	if fn == nil {
		ctx.Builder.Partial()
		return
	}

	analyzeConfigurerCallsInFunc(
		&entityContext[T, E, B]{
			context:                ctx.context,
			EntityType:             ctx.EntityType,
			Builder:                ctx.Builder,
			ConfigureMethod:        ctx.ConfigureMethod,
			FunctionUnderAnalysis:  fn,
			ConfigurerParamIndices: indices,
		},
		analyze,
	)
}

func analyzeIdentity[
	T config.Entity,
	E any,
	B configbuilder.EntityBuilder[T, E],
](
	ctx *configurerCallContext[T, E, B],
) {
	ctx.
		Builder.
		Identity(func(b *configbuilder.IdentityBuilder) {
			if ctx.IsSpeculative {
				b.Speculative()
			}

			if name, ok := ssax.AsString(ctx.Args[0]).TryGet(); ok {
				b.Name(name)
			}

			if key, ok := ssax.AsString(ctx.Args[1]).TryGet(); ok {
				b.Key(key)
			}
		})
}
