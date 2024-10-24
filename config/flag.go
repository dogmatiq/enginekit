package config

import (
	"reflect"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
)

// A Flag is a [Component] that represents an optional boolean flag that may be
// set on some other [Component].
//
// Each flag type of flag is identified by L.
type Flag[L Label] struct {
	IsSpeculative bool
}

// Fidelity returns information about how well the configuration represents the
// actual configuration that would be used at runtime.
func (f *Flag[L]) Fidelity() Fidelity {
	if f.IsSpeculative {
		return Speculative
	}
	return Immaculate
}

func (f *Flag[L]) String() string {
	return RenderDescriptor(f)
}

func (f *Flag[L]) renderDescriptor(ren *renderer.Renderer) {
	ren.Print(labelAsString[L]())
}

func (f *Flag[L]) renderDetails(ren *renderer.Renderer) {
	fid, errs := validate(f)

	renderFidelity(ren, fid, errs)

	ren.Print(labelAsString[L]())
	ren.Print(" flag")
	ren.Print("\n")
	renderErrors(ren, errs)
}

func (f *Flag[L]) clone() Component {
	return &Flag[L]{f.IsSpeculative}
}

func (f *Flag[L]) normalize(*normalizationContext) {
}

// FlagSet is a set of flags of a certain type. It uses a separate flag to
// represent each time the flag is applied to the [Component].
type FlagSet[L Label] []*Flag[L]

// resolve reports whether or not the flag is set, if definitively known.
func (s FlagSet[L]) resolve(c Component) optional.Optional[bool] {
	if c.Fidelity()&Incomplete != 0 {
		return optional.None[bool]()
	}

	if len(s) == 0 {
		return optional.Some(false)
	}

	for _, f := range s {
		if !f.IsSpeculative {
			return optional.Some(true)
		}
	}

	return optional.None[bool]()
}

// Label is a type that is used as a unique identity for a [Flag].
type Label interface {
	isLabel()
}

func labelAsString[L Label]() string {
	return strings.ToLower(reflect.TypeFor[L]().Name())
}
