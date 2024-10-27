package config

import (
	"reflect"
	"strings"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
	"github.com/dogmatiq/enginekit/optional"
)

// A FlagSymbol is a type that uniquely identifies a specific [Flag].
type FlagSymbol interface {
	~struct{ flagSymbol }
}

// flagSymbol is a "marker" struct that is embedded in named types to declare a
// new [FlagSymbol].
type flagSymbol struct{}

func symbol[S FlagSymbol]() string {
	return strings.ToLower(reflect.TypeFor[S]().Name())
}

// A Flag is a [Component] that represents some boolean state of some other
// [Component].
//
// Each type of flag is uniquely identified by a [FlagSymbol].
type Flag[S FlagSymbol] struct {
	ComponentProperties
	Modifications []*FlagModification
}

func (f *Flag[S]) String() string {
	return RenderDescriptor(f)
}

func (f *Flag[S]) renderDescriptor(ren *renderer.Renderer) {
	ren.Print(symbol[S]())
}

func (f *Flag[S]) renderDetails(ren *renderer.Renderer) {
	fid, errs := validate(f)

	renderFidelity(ren, fid, errs)

	ren.Print(symbol[S]())
	ren.Print(" flag")
	ren.Print("\n")

	for _, i := range f.Modifications {
		ren.IndentBullet()
		i.renderDetails(ren)
		ren.Dedent()
	}

	renderErrors(ren, errs)
}

func (f *Flag[S]) normalize(*normalizationContext) {
	// TODO
	// // resolve reports whether or not the flag is set, if definitively known.
	// func (f *Flag[S]) resolve(c Component) bool {
	// 	ctx := strictContext(c)

	// 	if f&Incomplete != 0 {
	// 		return optional.None[bool]()
	// 	}

	// 	if len(f) == 0 {
	// 		return optional.Some(false)
	// 	}

	// 	for _, f := range f {
	// 		if !f.IsSpeculative {
	// 			return optional.Some(true)
	// 		}
	// 	}

	//		return optional.None[bool]()
	//	}
}

func (f *Flag[S]) clone() any {
	return &Flag[S]{
		clone(f.ComponentProperties),
		cloneSlice(f.Modifications),
	}
}

// A FlagModification is a [Component] that represents a specific point at which
// a flag is set or unset within the configuration.
type FlagModification struct {
	ComponentProperties
	Value optional.Optional[bool]
}

func (m *FlagModification) String() string {
	return RenderDescriptor(m)
}

func (m *FlagModification) renderDescriptor(ren *renderer.Renderer) {
	ren.Print("flag-modification:")
	if v, ok := m.Value.TryGet(); !ok {
		ren.Print("?")
	} else if v {
		ren.Print("set")
	} else {
		ren.Print("cleared")
	}
}

func (m *FlagModification) renderDetails(ren *renderer.Renderer) {
	fid, errs := validate(m)

	renderFidelity(ren, fid, errs)
	ren.Print("flag modification")

	if v, ok := m.Value.TryGet(); ok {
		if v {
			ren.Print(" (set)")
		} else {
			ren.Print(" (cleared)")
		}
	}

	ren.Print("\n")
	renderErrors(ren, errs)
}

func (m *FlagModification) normalize(*normalizationContext) {
	if !m.Value.IsPresent() {
		m.Fidelity |= Incomplete
	}
}

func (m *FlagModification) clone() any {
	return &FlagModification{
		clone(m.ComponentProperties),
		m.Value,
	}
}
