package config

// Completeness indicates the level of completeness of an application or
// handler's configuration.
type Completeness int

const (
	// Unconfigured indicates that the configuration is empty; no methods were
	// called on the configurer.
	Unconfigured Completeness = iota

	// PartiallyConfigured indicates that some configuration has been applied,
	// but it is not sufficient to use the entity.
	PartiallyConfigured

	// FullyConfigured indicates a complete configuration has been applied, but
	// not necessarily that the configuration is valid.
	//
	// For example, a handler requires an identity and at least one route to be
	// fully configured, but its route may conflict with another handler.
	FullyConfigured
)
