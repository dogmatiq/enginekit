package config

import "fmt"

// InvalidIdentityNameError indicates that the "name" element of an [Identity]
// is invalid.
type InvalidIdentityNameError struct {
	InvalidName string
}

func (e InvalidIdentityNameError) Error() string {
	return fmt.Sprintf("invalid name (%q), expected a non-empty, printable UTF-8 string with no whitespace", e.InvalidName)
}

// InvalidIdentityKeyError indicates that the "key" element of an [Identity]
// is invalid.
type InvalidIdentityKeyError struct {
	InvalidKey string
}

func (e InvalidIdentityKeyError) Error() string {
	return fmt.Sprintf("invalid key (%q), expected an RFC 4122/9562 UUID", e.InvalidKey)
}

// MissingIdentityError indicates that an [QEntity] has been configured without
// an [Identity].
type MissingIdentityError struct{}

func (e MissingIdentityError) Error() string {
	return "no identity is configured"
}

// MultipleIdentitiesError indicates that an [QEntity] has been configured with
// more than one [Identity].
type MultipleIdentitiesError struct {
	Identities []*Identity
}

func (e MultipleIdentitiesError) Error() string {
	return fmt.Sprintf(
		"multiple identities are configured: %s",
		renderList(e.Identities),
	)
}

// IdentityNameConflictError indicates that more than one [QEntity] within the
// same [Application] is shares the same "name" element of an [Identity].
type IdentityNameConflictError struct {
	Entities        []QEntity
	ConflictingName string
}

func (e IdentityNameConflictError) Error() string {
	return fmt.Sprintf(
		"entities have conflicting identities: the %q name is shared by %s",
		e.ConflictingName,
		renderList(e.Entities),
	)
}

// IdentityKeyConflictError indicates that more than one [QEntity] within the
// same [Application] is shares the same "key" element of an [Identity].
type IdentityKeyConflictError struct {
	Entities       []QEntity
	ConflictingKey string
}

func (e IdentityKeyConflictError) Error() string {
	return fmt.Sprintf(
		"entities have conflicting identities: the %q key is shared by %s",
		e.ConflictingKey,
		renderList(e.Entities),
	)
}
