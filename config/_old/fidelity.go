package config

// Fidelity is a bit-field that describes how well a [Component] configuration
// represents the actual configuration that would be used at runtime.
//
// Importantly, it does not describe the validity of the configuration itself.
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
