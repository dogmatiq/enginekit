package config

// Fidelity describes how well a [Component] configuration represents the actual
// configuration that would be when running an application.
type Fidelity struct {
	// IsExhaustive is true if all _available_ configuration logic was applied
	// when building the configuration.
	//
	// It does not imply that all _mandatory_ configuration is present.
	//
	// An exhaustive configuration is guaranteed to contain all of the
	// components that could possibly be included in the configuration at
	// runtime.
	IsExhaustive bool

	// IsSpeculative is true if the component is only included in the
	// configuration under certain conditions and those conditions could not be
	// evaluated at the time the configuration was built.
	IsSpeculative bool

	// IsUnresolved is true if any of the component's configuration values
	// could not be determined at the time the configuration was built.
	IsUnresolved bool

	// HasSpeculativeSubcomponents is true if the component has any speculative
	// sub-components.
	HasSpeculativeSubcomponents bool

	// HasUnresolvedSubcomponents is true if the component has any
	// sub-components that have unresolved configuration values.
	HasUnresolvedSubcomponents bool

	// HasMutuallyExclusiveSubcomponents is true if the component has
	// sub-components that are included under mutually exclusive conditions and
	// it was unknown at the time of configuration whether those conditions
	// would be met.
	HasMutuallyExclusiveSubcomponents bool
}
