package config

import (
	"fmt"
)

// Component is the "top-level" interface for the individual elements that form
// a complete configuration of a Dogma application or handler.
type Component interface {
	fmt.Stringer

	// Fidelity reports how faithfully the [Component] describes a complete
	// configuration that can be used to execute an application.
	Fidelity() Fidelity

	// ComponentProperties returns the properties common to all [Component] types.
	ComponentProperties() *ComponentCommon

	validate(*validateContext)
	describe(*describeContext)
}

// ComponentCommon contains the properties common to all [Component] types.
type ComponentCommon struct {
	ComponentFidelity Fidelity
}

// Fidelity reports how faithfully the [Component] describes a complete
// configuration that can be used to execute an application.
func (p *ComponentCommon) Fidelity() Fidelity {
	return p.ComponentFidelity
}

// ComponentProperties returns the properties common to all [Component] types.
func (p *ComponentCommon) ComponentProperties() *ComponentCommon {
	return p
}

func validateComponent(
	ctx *validateContext,
	funcs ...func(*validateContext),
) {
	for _, fn := range funcs {
		fn(ctx)
	}

	if ctx.Component.Fidelity().Has(Incomplete) {
		ctx.Invalid(IncompleteComponentError{})
	}
}

var (
	_ Component = (*Identity)(nil)
	_ Component = (*FlagModification)(nil)
	_ Component = (*Route)(nil)

	_ Entity = (*Application)(nil)

	_ Handler   = (*Aggregate)(nil)
	_ Handler   = (*Process)(nil)
	_ Handler   = (*Integration)(nil)
	_ Handler   = (*Projection)(nil)
	_ Component = (*ProjectionDeliveryPolicy)(nil)
)
