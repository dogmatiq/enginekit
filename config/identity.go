package config

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dogmatiq/enginekit/optional"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

// IdentityAsConfigured contains the raw unvalidated properties of an
// [Identity].
type IdentityAsConfigured struct {
	// Name is the human-readable name of the entity, if available.
	Name optional.Optional[string]

	// Key is the unique identifier for the entity, if available.
	Key optional.Optional[string]

	// Fidelity describes the configuration's accuracy in comparison to the
	// actual configuration that would be used at runtime.
	Fidelity Fidelity
}

// Identity represents the (potentially invalid) identity of an entity.
type Identity struct {
	AsConfigured IdentityAsConfigured
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (i *Identity) Fidelity() Fidelity {
	return i.AsConfigured.Fidelity
}

func (i *Identity) String() string {
	w := strings.Builder{}
	w.WriteString("identity")

	n, nOK := i.AsConfigured.Name.TryGet()
	k, kOK := i.AsConfigured.Key.TryGet()

	if !nOK && !kOK {
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

func (i *Identity) clone() Component {
	return &Identity{i.AsConfigured}
}

func (i *Identity) normalize(ctx *normalizationContext) {
	if n, ok := i.AsConfigured.Name.TryGet(); ok {
		if !isPrintableIdentifier(n) {
			ctx.Fail(InvalidIdentityNameError{n})
		}
	} else {
		i.AsConfigured.Fidelity.IsPartial = true
	}

	if k, ok := i.AsConfigured.Key.TryGet(); ok {
		if id, err := uuidpb.Parse(k); err != nil {
			ctx.Fail(InvalidIdentityKeyError{k})
		} else {
			i.AsConfigured.Key = optional.Some(id.AsString())
		}
	} else {
		i.AsConfigured.Fidelity.IsPartial = true
	}
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

func buildIdentity(ctx *normalizationContext, identities []*Identity) *identitypb.Identity {
	identities = clone(identities)
	normalizeIdentities(ctx, identities)

	id := identities[0].AsConfigured

	return &identitypb.Identity{
		Name: id.Name.Get(),
		Key:  uuidpb.MustParse(id.Key.Get()),
	}
}

func normalizeIdentities(ctx *normalizationContext, identities []*Identity) {
	normalize(ctx, identities...)

	if len(identities) == 0 {
		ctx.Fail(MissingIdentityError{})
	} else if len(identities) > 1 {
		ctx.Fail(MultipleIdentitiesError{identities})
	}
}
