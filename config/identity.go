package config

import (
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

func (i Identity) String() string {
	if i.Name == "" && i.Key == "" {
		return "identity"
	}

	name := i.Name
	key := i.Key

	if !isPrintableIdentifier(name) {
		name = strconv.Quote(name)
	}

	if norm, err := uuidpb.Parse(key); err == nil {
		key = norm.AsString()
	} else {
		if !isPrintableIdentifier(key) {
			key = strconv.Quote(key)
		}
	}

	return "identity:" + name + "/" + key
}

func (i Identity) normalize(ctx *normalizationContext) Component {
	if !isPrintableIdentifier(i.Name) {
		ctx.Fail(InvalidIdentityNameError{i.Name})
	}

	if k, err := uuidpb.Parse(i.Key); err != nil {
		ctx.Fail(InvalidIdentityKeyError{i.Key})
	} else {
		i.Key = k.AsString()
	}

	return i
}

// InvalidIdentityNameError indicates that the "name" component of an [Identity]
// is invalid.
type InvalidIdentityNameError struct {
	InvalidName string
}

func (e InvalidIdentityNameError) Error() string {
	return fmt.Sprintf("invalid name (%q), expected a non-empty, printable UTF-8 string with no whitespace", e.InvalidName)
}

// InvalidIdentityKeyError indicates that the "key" component of an [Identity]
// is invalid.
type InvalidIdentityKeyError struct {
	InvalidKey string
}

func (e InvalidIdentityKeyError) Error() string {
	return fmt.Sprintf("invalid key (%q), expected an RFC 4122/9562 UUID", e.InvalidKey)
}

// isPrintableIdentifier returns true if n contains only non-whitespace printable
// Unicode characters.
func isPrintableIdentifier(n string) bool {
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

func normalizedIdentity(ent Entity) Identity {
	identities := ent.identities()

	if len(identities) == 0 {
		panic(NoIdentityError{})
	} else if len(identities) > 1 {
		panic(MultipleIdentitiesError{identities})
	}

	return MustNormalize(identities[0])
}

func normalizeIdentities(ctx *normalizationContext, ent Entity) []Identity {
	identities := slices.Clone(ent.identities())

	if len(identities) == 0 {
		ctx.Fail(NoIdentityError{})
	} else if len(identities) > 1 {
		ctx.Fail(MultipleIdentitiesError{identities})
	}

	for i, id := range identities {
		identities[i] = normalize(ctx, id)
	}

	return identities
}
