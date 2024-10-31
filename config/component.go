package config

import (
	"fmt"
	"reflect"
)

// Component is the "top-level" interface for the individual elements that form
// a complete configuration of a Dogma application or handler.
type Component interface {
	fmt.Stringer

	validate(*validateContext)
	describe(*describeContext)
}

// ValueUnavailableError indicates that a [Component] cannot produce the actual
// value it represents because there is insufficient runtime type information
// available.
//
// See [ForExecution].
type ValueUnavailableError struct {
	Type reflect.Type
}

func (e ValueUnavailableError) Error() string {
	return fmt.Sprintf("%s value is unavailable", e.Type)
}

var (
	_ Component = (*Identity)(nil)
	_ Component = (*Flag[struct{ symbol }])(nil)
	_ Component = (*FlagModification)(nil)
	_ Component = (*Route)(nil)

	_ Entity = (*Application)(nil)

	_ Handler   = (*Aggregate)(nil)
	_ Handler   = (*Process)(nil)
	_ Handler   = (*Integration)(nil)
	_ Handler   = (*Projection)(nil)
	_ Component = (*ProjectionDeliveryPolicy)(nil)
)
