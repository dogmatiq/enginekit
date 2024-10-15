package config

import (
	"fmt"
	"reflect"
)

// ImplementationUnavailableError indicates that a [Value] is invalid
// because it does not contain some runtime value or type information and the
// [WithImplementations] option was specified during normalization.
type ImplementationUnavailableError struct {
	MissingType reflect.Type
}

func (e ImplementationUnavailableError) Error() string {
	return fmt.Sprintf("missing implementation: %s value is not available", e.MissingType)
}

// TypeNameMismatchError indicates that the type name specified in some
// [Component] does not match the name of the actual type it refers to.
type TypeNameMismatchError struct {
	ExpectedTypeName   string
	UnexpectedTypeName string
}

func (e TypeNameMismatchError) Error() string {
	return fmt.Sprintf(
		"type name mismatch: implementation is %q, but the type name is reported as %q",
		e.ExpectedTypeName,
		e.UnexpectedTypeName,
	)
}
