package config

import (
	"slices"
	"strconv"
	"unicode"

	"github.com/dogmatiq/enginekit/protobuf/identitypb"
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

func (i Identity) normalize(ctx *normalizeContext) Component {
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

func finalizeIdentity(ctx *normalizeContext, ent Entity) *identitypb.Identity {
	identities := normalizeIdentities(ctx, ent)

	return &identitypb.Identity{
		Name: identities[0].Name,
		Key:  uuidpb.MustParse(identities[0].Key),
	}
}

func normalizeIdentities(ctx *normalizeContext, ent Entity) []Identity {
	identities := slices.Clone(ent.identities())

	if len(identities) == 0 {
		ctx.Fail(MissingIdentityError{})
	} else if len(identities) > 1 {
		ctx.Fail(MultipleIdentitiesError{identities})
	}

	for i, id := range identities {
		identities[i] = normalize(ctx, id)
	}

	return identities
}
