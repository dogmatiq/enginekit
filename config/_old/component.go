package config

import (
	"fmt"

	"github.com/dogmatiq/enginekit/config/internal/renderer"
)

// A Component is some element of the configuration of a Dogma application.
type Component interface {
	fmt.Stringer
	clonable

	// Component returns the (possibly invalid or incomplete) properties of the
	// component.
	CommonComponentProperties() *ComponentProperties

	renderDescriptor(*renderer.Renderer)
	renderDetails(*renderer.Renderer)
	normalize(*normalizationContext)
}

// ComponentProperties contains the properties common to all [Component]
// implementations.
type ComponentProperties struct {
	// Fidelity describes the configuration's accuracy and completeness in
	// comparison to the actual configuration that would be used at runtime.
	//
	// The fidelity must be [Immaculate] in order for the configuration to be
	// executed by an engine.
	Fidelity Fidelity
}

// CommonComponentProperties returns the (possibly invalid or incomplete)
// properties of the component.
func (p *ComponentProperties) CommonComponentProperties() *ComponentProperties {
	return p
}

func (p ComponentProperties) clone() any {
	return p
}

var (
	_ QEntity   = (*Application)(nil)
	_ Component = (*Identity)(nil)

	_ Handler   = (*Aggregate)(nil)
	_ Handler   = (*Process)(nil)
	_ Handler   = (*Integration)(nil)
	_ Handler   = (*Projection)(nil)
	_ Component = (*Route)(nil)
	_ Component = (*Flag[struct{ flagSymbol }])(nil)
	_ Component = (*FlagModification)(nil)
)
