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
