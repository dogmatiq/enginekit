package config

import (
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// IdentityProperties contains the raw unvalidated properties of an [Identity].
type IdentityProperties struct {
	Name     string
	Key      string
	Fidelity Fidelity
}

// Identity represents the (potentially invalid) identity of an entity.
type Identity struct {
	AsConfigured IdentityProperties
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (i Identity) Fidelity() Fidelity {
	return i.AsConfigured.Fidelity
}

func (i Identity) String() string {
	w := strings.Builder{}

	writeComponentPrefix(&w, "identity", i)

	n := i.AsConfigured.Name
	k := i.AsConfigured.Key

	if n == "" && k == "" {
		return w.String()
	}

	w.WriteByte(':')

	if isPrintableIdentifier(n) {
		w.WriteString(n)
	} else {
		w.WriteString(strconv.Quote(n))
	}

	w.WriteByte('/')

	if norm, err := uuidpb.Parse(k); err == nil {
		w.WriteString(norm.AsString())
	} else if isPrintableIdentifier(k) {
		w.WriteString(k)
	} else {
		w.WriteString(strconv.Quote(k))
	}

	return w.String()
}

func (i Identity) normalize(ctx *normalizeContext) Component {
	if !isPrintableIdentifier(i.AsConfigured.Name) {
		ctx.Fail(InvalidIdentityNameError{i.AsConfigured.Name})
	}

	if k, err := uuidpb.Parse(i.AsConfigured.Key); err != nil {
		ctx.Fail(InvalidIdentityKeyError{i.AsConfigured.Key})
	} else {
		i.AsConfigured.Key = k.AsString()
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
	id := normalizeIdentities(ctx, ent)[0].AsConfigured

	return &identitypb.Identity{
		Name: id.Name,
		Key:  uuidpb.MustParse(id.Key),
	}
}

func normalizeIdentities(ctx *normalizeContext, ent Entity) []Identity {
	identities := slices.Clone(ent.identitiesAsConfigured())

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
