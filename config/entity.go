package config

import "fmt"

// An Entity is a [Component] that represents the configuration of some
// configurable Dogma entity; that is, any type with a Configure() method that
// accepts one of the Dogma "configurer" interfaces.
type Entity interface {
	Component

	// Identity returns the entity's identity.
	//
	// It panics if no single valid identity is configured.
	Identity() Identity

	// IsExhaustive returns true if the complete configuration was loaded. It
	// may be false, for example, when attempting to load configuration using
	// static analysis, but the code depends on runtime type information.
	IsExhaustive() bool

	identities() []Identity
}

// MissingIdentityError indicates that an [Entity] has been configured without
// an [Identity].
type MissingIdentityError struct{}

func (e MissingIdentityError) Error() string {
	return "no identity is configured"
}

// MultipleIdentitiesError indicates that an [Entity] has been configured with
// more than one [Identity].
type MultipleIdentitiesError struct {
	Identities []Identity
}

func (e MultipleIdentitiesError) Error() string {
	return fmt.Sprintf(
		"multiple identities are configured: %s",
		renderList(e.Identities),
	)
}
