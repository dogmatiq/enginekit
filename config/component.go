package config

import "fmt"

// Component is the "top-level" interface for the individual elements that form
// a complete configuration of a Dogma application or handler.
type Component interface {
	fmt.Stringer

	// Fidelity reports how faithfully the [Component] describes a complete
	// configuration that can be used to execute an application.
	Fidelity() Fidelity

	// Baggage is a collection of arbitrary data that is associated with the
	// [Component] by whatever system produced the configuration.
	Baggage() Baggage

	validate(ctx *validationContext)
}

// ValidateOption changes the behavior of [Component.Validate].
type ValidateOption func()

// ComponentCommon is a partial implementation of [Component].
type ComponentCommon struct {
	ComponentFidelity Fidelity
	ComponentBaggage  Baggage
}

// Fidelity reports how faithfully the [Component] describes a complete
// configuration that can be used to execute an application.
func (c *ComponentCommon) Fidelity() Fidelity {
	return c.ComponentFidelity
}

// Baggage returns a collection of arbitrary data that is associated with the
// [Component] by whatever system produced the configuration.
func (c *ComponentCommon) Baggage() Baggage {
	return c.ComponentBaggage
}

func (c *ComponentCommon) validate(ctx *validationContext) {
	if c.ComponentFidelity.Has(Incomplete) {
		ctx.Fail(IncompleteComponentError{})
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
