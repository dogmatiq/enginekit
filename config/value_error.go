package config

import (
	"fmt"
	"reflect"
)

// RuntimeValueUnavailableError indicates that a [Value] is invalid because it
// does not contain some runtime value or type information and the
// [WithRuntimeTypes] option was specified during normalization.
type RuntimeValueUnavailableError struct {
	MissingType reflect.Type
}

func (e RuntimeValueUnavailableError) Error() string {
	return fmt.Sprintf("%s value is not available", e.MissingType)
}

// TypeNameMismatchError indicates that the type name reported in some
// [Component] does not match the name of runtime type that it refers to.
type TypeNameMismatchError struct {
	ReportedTypeName string
	RuntimeTypeName  string
}

func (e TypeNameMismatchError) Error() string {
	return fmt.Sprintf(
		"type name mismatch: %s does not match the runtime type (%s)",
		e.ReportedTypeName,
		e.RuntimeTypeName,
	)
}
