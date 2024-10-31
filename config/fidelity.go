package config

// Fidelity is a bit-field that describes how faithfully a [Component] describes
// a complete configuration that can be used to execute an application.
//
// It is possible to have an [Immaculate] representation of an invalid
// configuration.
type Fidelity int

const (
	// Immaculate is the [Fidelity] value that indicates the configuration is an
	// exact match for the actual configuration that would be used at runtime.
	Immaculate Fidelity = 0

	// Incomplete is a [Fidelity] flag that indicates that the [Component] has
	// some configuration that could not be resolved accurately at configuration
	// time.
	//
	// Most commonly this is occurs during static analysis of code that uses
	// interfaces that cannot be followed statically.
	//
	// Its absence means that all of the _available_ configuration logic was
	// applied; it does not imply that all _mandatory_ configuration is present.
	Incomplete Fidelity = 1 << iota

	// Speculative is a [Fidelity] flag that indicates that the [Component] is
	// only present in the configuration under certain conditions, and that
	// those conditions could not be evaluated at configuration time.
	Speculative
)

// Has returns true if d is a superset of dir.
func (f Fidelity) Has(x Fidelity) bool {
	return f&x != 0
}

func validateFidelity(ctx *validateContext, f Fidelity) {
	failIfIncomplete(ctx, f)
	failIfSpeculative(ctx, f)
}

func failIfIncomplete(ctx *validateContext, f Fidelity) {
	if f.Has(Incomplete) {
		ctx.Invalid(IncompleteComponentError{})
	}
}

func failIfSpeculative(ctx *validateContext, f Fidelity) {
	if f.Has(Speculative) {
		ctx.Invalid(SpeculativeComponentError{})
	}
}

func describeFidelity(ctx *describeContext, f Fidelity) {
	if f&Incomplete != 0 {
		ctx.Print("incomplete ")
	} else if !ctx.options.ValidationResult.IsPresent() {
		ctx.Print("unvalidated ")
	} else if len(ctx.errors) == 0 {
		ctx.Print("valid ")
	} else {
		ctx.Print("invalid ")
	}

	if f&Speculative != 0 {
		ctx.Print("speculative ")
	}
}

// IncompleteComponentError indicates that a [Component] contains values that
// could not be resolved at the time the configuration was built.
//
// See [Fidelity] for more information.
type IncompleteComponentError struct{}

func (e IncompleteComponentError) Error() string {
	return "could not evaluate entire configuration"
}

// SpeculativeComponentError indicates that a [Component]'s inclusion in the
// configuration is subject to some condition that could not be evaluated at the
// time the configuration was built.
//
// See [Fidelity] for more information.
type SpeculativeComponentError struct{}

func (e SpeculativeComponentError) Error() string {
	return "conditions for the component's inclusion in the configuration could not be evaluated"
}
