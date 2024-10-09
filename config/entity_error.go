package config

import "fmt"

// NoIdentityError indicates that an [Entity] has been configured without an
// [Identity].
type NoIdentityError struct {
	Entity Entity
}

func (e NoIdentityError) Error() string {
	return fmt.Sprintf(
		"%s is configured without an identity, Identity() must be called exactly once within Configure()",
		e.Entity,
	)
}

// MultipleIdentitiesError indicates that an [Entity] has been configured with
// more than one [Identity].
type MultipleIdentitiesError struct {
	Entity     Entity
	Identities []Identity
}

func (e MultipleIdentitiesError) Error() string {
	return fmt.Sprintf(
		"%s is configured with multiple identities (%s), Identity() must be called exactly once within Configure()",
		e.Entity,
		renderList(e.Identities),
	)
}

// InvalidIdentityError indicates that an [Entity] has been configured with an
// invalid [Identity].
type InvalidIdentityError struct {
	Entity          Entity
	InvalidIdentity Identity
	Cause           error
}

func (e InvalidIdentityError) Error() string {
	return fmt.Sprintf(
		"%s is configured with an invalid identity (%s): %s",
		e.Entity,
		e.InvalidIdentity,
		e.Cause,
	)
}

func (e InvalidIdentityError) Unwrap() error {
	return e.Cause
}
