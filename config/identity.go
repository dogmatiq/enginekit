package config

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"unicode"

	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// Identity represents the (potentially invalid) identity of an entity.
type Identity struct {
	Name string
	Key  string
}

func (i Identity) normalize(validationOptions) (_ Identity, errs error) {
	if !isValidIdentityName(i.Name) {
		errs = errors.Join(errs, InvalidIdentityNameError{i.Name})
	}

	if k, err := uuidpb.Parse(i.Key); err != nil {
		errs = errors.Join(errs, InvalidIdentityKeyError{i.Key, err})
	} else {
		i.Key = k.AsString()
	}

	return i, errs
}

func (i Identity) String() string {
	name := "?"
	if i.Name != "" {
		if isValidIdentityName(i.Name) {
			name = i.Name
		} else {
			name = strconv.Quote(i.Name)
		}
	}

	key := "?"
	if i.Key != "" {
		if normalized, err := uuidpb.Parse(i.Key); err == nil {
			key = normalized.AsString()
		} else {
			key = strconv.Quote(i.Key)
		}
	}

	return name + "/" + key
}

// isValidIdentityName returns true if n is a valid application or handler name.
//
// A valid name is a non-empty string consisting of Unicode printable
// characters, except whitespace.
func isValidIdentityName(n string) bool {
	if len(n) == 0 {
		return false
	}

	for _, r := range n {
		if unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}

// NoIdentityError is an error that occurs when an application or handler does
// not configure an identity.
type NoIdentityError struct {
	Entity Entity
}

func (e NoIdentityError) Error() string {
	return fmt.Sprintf(
		"%s is configured without an identity, Identity() must be called exactly once within Configure()",
		e.Entity,
	)
}

// MultipleIdentitiesError is an error that occurs when an application or
// handler configures multiple identities.
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

// InvalidIdentityError is an error that occurs when an application or handler
// is configured with an invalid identity.
type InvalidIdentityError struct {
	Entity   Entity
	Identity Identity
	Cause    error
}

func (e InvalidIdentityError) Error() string {
	return fmt.Sprintf(
		"%s is configured with an invalid identity (%s): %s",
		e.Entity,
		e.Identity,
		e.Cause,
	)
}

func (e InvalidIdentityError) Unwrap() error {
	return e.Cause
}

// IdentityConflictError is an error that occurs when multiple entities are
// configured with conflict identities.
type IdentityConflictError struct {
	Entities            []Entity
	ConflictingIdentity Identity
}

func (e IdentityConflictError) Error() string {
	return fmt.Sprintf(
		"%s have the same identity (%s), which is not allowed",
		renderList(e.Entities),
		e.ConflictingIdentity,
	)
}

// IdentityNameConflictError is an error that occurs when multiple entities are
// configured with identities that have conflicting names.
type IdentityNameConflictError struct {
	Entities        []Entity
	ConflictingName string
}

func (e IdentityNameConflictError) Error() string {
	return fmt.Sprintf(
		"%s have the same identity name (%s), which is not allowed",
		renderList(e.Entities),
		e.ConflictingName,
	)
}

// IdentityKeyConflictError is an error that occurs when multiple entities are
// configured with identities that have conflicting keys.
type IdentityKeyConflictError struct {
	Entities       []Entity
	ConflictingKey string
}

func (e IdentityKeyConflictError) Error() string {
	return fmt.Sprintf(
		"%s have the same identity key (%s), which is not allowed",
		renderList(e.Entities),
		e.ConflictingKey,
	)
}

// InvalidIdentityNameError is an error that occurs when an identity name is
// invalid.
type InvalidIdentityNameError struct {
	InvalidName string
}

func (e InvalidIdentityNameError) Error() string {
	return fmt.Sprintf("invalid identity name (%q): names must be non-empty, printable UTF-8 strings with no whitespace", e.InvalidName)
}

// InvalidIdentityKeyError is an error that occurs when an identity key is
// invalid.
type InvalidIdentityKeyError struct {
	InvalidKey string
	Cause      error
}

func (e InvalidIdentityKeyError) Error() string {
	return fmt.Sprintf("invalid identity key (%q): keys must be RFC 4122/9562 UUIDs: %s", e.InvalidKey, e.Cause)
}

func (e InvalidIdentityKeyError) Unwrap() error {
	return e.Cause
}

func renderList[T any](items []T) string {
	var s string

	for i, item := range items {
		if i == len(items)-1 {
			s += " and "
		} else if i > 0 {
			s += ", "
		}
		s += fmt.Sprint(item)
	}

	return s
}

func normalizedIdentity(ent Entity) Identity {
	identities := ent.configuredIdentities()

	if len(identities) == 0 {
		panic(NoIdentityError{Entity: ent})
	} else if len(identities) > 1 {
		panic(MultipleIdentitiesError{ent, identities})
	}

	id, err := identities[0].normalize(validationOptions{})
	if err != nil {
		panic(InvalidIdentityError{ent, id, err})
	}

	return id
}

func normalizeIdentitiesInPlace(
	opts validationOptions,
	ent Entity,
	errs *error,
	identities *[]Identity,
) {
	*identities = slices.Clone(*identities)

	if len(*identities) == 0 {
		*errs = errors.Join(*errs, NoIdentityError{Entity: ent})
	} else if len(*identities) > 1 {
		*errs = errors.Join(*errs, MultipleIdentitiesError{ent, *identities})
	}

	for i, id := range *identities {
		norm, err := id.normalize(opts)
		(*identities)[i] = norm

		if err != nil {
			*errs = errors.Join(*errs, InvalidIdentityError{ent, id, err})
		}
	}
}
