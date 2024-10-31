package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dogmatiq/enginekit/optional"
)

// A Symbol is a type that uniquely identifies a specific [Flag].
type Symbol interface {
	~struct{ symbol }
}

// symbol is a "marker" struct that is embedded in named types to declare a
// new [Symbol].
type symbol struct{}

// A Flag represents some boolean state of a [Component].
//
// Each type of flag is uniquely identified by a [Symbol].
type Flag[S Symbol] struct {
	ComponentCommon

	Modifications []*FlagModification
}

func (f *Flag[S]) String() string {
	return strings.ToLower(reflect.TypeFor[S]().Name()) + " flag"
}

// Get returns the definitive value of the flag, if possible.
func (f *Flag[S]) Get() optional.Optional[bool] {
	if len(f.Modifications) == 0 {
		return optional.None[bool]()
	}

	result := f.Modifications[0].Value

	for _, m := range f.Modifications[1:] {
		if !m.Fidelity().Has(Speculative) {
			return optional.None[bool]()
		}

		if result != m.Value {
			return optional.None[bool]()
		}
	}

	return result
}

func (f *Flag[S]) validate(ctx *validateContext) {
	validateComponent(ctx)

	for _, m := range f.Modifications {
		ctx.ValidateChild(m)
	}
}

func (f *Flag[S]) describe(ctx *describeContext) {
	if len(f.Modifications) == 0 {
		return
	}

	ctx.Print(f.String())

	if v, ok := f.Get().TryGet(); ok {
		ctx.Printf(" set to %t\n", v)
		return
	}

	ctx.Print("\n")

	for _, m := range f.Modifications {
		ctx.DescribeChild(m)
	}
}

// A FlagModification is a [Component] that represents a specific point at which
// a flag is set or unset within the configuration.
type FlagModification struct {
	ComponentCommon

	Value optional.Optional[bool]
}

func (m *FlagModification) String() string {
	if v, ok := m.Value.TryGet(); ok {
		return fmt.Sprintf("flag-modification:%t", v)
	}
	return "flag-modification:?"
}

func (m *FlagModification) validate(*validateContext) {
}

func (m *FlagModification) describe(ctx *describeContext) {
	ctx.DescribeFidelity()
	ctx.Print("flag modification")

	if v, ok := m.Value.TryGet(); ok {
		ctx.Printf(", set to %t", v)
	}
}
