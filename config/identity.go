package config

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
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
	return RenderDescriptor(i)
}

func (i *Identity) renderDescriptor(ren *renderer.Renderer) {
	ren.Print("identity")

	name, nameOK := i.AsConfigured.Name.TryGet()
	key, keyOK := i.AsConfigured.Key.TryGet()

	if !nameOK && !keyOK {
		return
	}

	ren.Print(":")

	if !nameOK {
		ren.Print("?")
	} else if !isPrintableIdentifier(name) || strings.Contains(name, `"`) {
		ren.Print(strconv.Quote(name))
	} else {
		ren.Print(name)
	}

	ren.Print("/")

	if !keyOK {
		ren.Print("?")
	} else if uuid, err := uuidpb.Parse(key); err == nil {
		ren.Print(uuid.AsString())
	} else if !isPrintableIdentifier(key) || strings.Contains(key, `"`) {
		ren.Print(strconv.Quote(key))
	} else {
		ren.Print(key)
	}
}

func (i *Identity) renderDetails(ren *renderer.Renderer) {
	f, errs := validate(i)

	renderFidelity(ren, f, errs)
	ren.Print("identity ")

	if name, ok := i.AsConfigured.Name.TryGet(); !ok {
		ren.Print("?")
	} else if !isPrintableIdentifier(name) || strings.Contains(name, `"`) {
		ren.Print(strconv.Quote(name))
	} else {
		ren.Print(name)
	}

	ren.Print("/")

	if key, ok := i.AsConfigured.Key.TryGet(); !ok {
		ren.Print("?")
	} else if !isPrintableIdentifier(key) || strings.Contains(key, `"`) {
		ren.Print(strconv.Quote(key))
	} else {
		ren.Print(key)
	}

	if key, ok := i.AsConfigured.Key.TryGet(); ok {
		if uuid, err := uuidpb.Parse(key); err == nil {
			if uuid.AsString() != key {
				ren.Print(" (non-canonical)")
			}
		}
	}

	ren.Print("\n")
	renderErrors(ren, errs)
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
		i.AsConfigured.Fidelity |= Incomplete
	}

	if k, ok := i.AsConfigured.Key.TryGet(); ok {
		if id, err := uuidpb.Parse(k); err != nil {
			ctx.Fail(InvalidIdentityKeyError{k})
		} else {
			i.AsConfigured.Key = optional.Some(id.AsString())
		}
	} else {
		i.AsConfigured.Fidelity |= Incomplete
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
	normalizeChildren(ctx, identities)
	reportIdentityErrors(ctx, identities)

	id := identities[0].AsConfigured

	return &identitypb.Identity{
		Name: id.Name.Get(),
		Key:  uuidpb.MustParse(id.Key.Get()),
	}
}

func reportIdentityErrors(ctx *normalizationContext, identities []*Identity) {
	if len(identities) == 0 {
		ctx.Fail(MissingIdentityError{})
	} else if len(identities) > 1 {
		ctx.Fail(MultipleIdentitiesError{identities})
	}
}
