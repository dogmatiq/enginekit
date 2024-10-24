package config

import (
	"fmt"
	"slices"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
)

// A Component is some element of the configuration of a Dogma application.
type Component interface {
	fmt.Stringer

	// Fidelity returns information about how well the configuration represents
	// the actual configuration that would be used at runtime.
	Fidelity() Fidelity

	renderDescriptor(*renderer.Renderer)
	renderDetails(*renderer.Renderer)

	clone() Component
	normalize(*normalizationContext)
}

// ComponentTrait is a partial implementation of [Component].
type ComponentTrait struct {
	// F describes the configuration's accuracy in comparison to the actual
	// configuration that would be used at runtime.
	F Fidelity
}

// Fidelity returns information about how well the configuration represents
// the actual configuration that would be used at runtime.
func (c ComponentTrait) Fidelity() Fidelity {
	return c.F
}

var (
	_ Entity    = (*Application)(nil)
	_ Component = (*Identity)(nil)

	_ Handler   = (*Aggregate)(nil)
	_ Handler   = (*Process)(nil)
	_ Handler   = (*Integration)(nil)
	_ Handler   = (*Projection)(nil)
	_ Component = (*Route)(nil)
	_ Component = (*Flag[Label])(nil)
)

func clone[T Component](components []T) []T {
	clones := slices.Clone(components)

	for i, c := range components {
		clones[i] = c.clone().(T)
	}

	return clones
}

func cloneInPlace[T Component](components *[]T) {
	*components = clone(*components)
}
