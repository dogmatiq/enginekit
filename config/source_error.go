package config

import (
	"fmt"
	"reflect"
)

// MissingTypeNameError indicates that a [Component] that refers to a Go type is
// missing the type name.
type MissingTypeNameError struct{}

func (e MissingTypeNameError) Error() string {
	return "missing type name"
}

// MissingImplementationError indicates that a [Component] is invalid because it
// does not contain some runtime value or type information and the
// [WithImplementations] option was specified during normalization.
type MissingImplementationError struct {
	MissingType reflect.Type
}

func (e MissingImplementationError) Error() string {
	return fmt.Sprintf("missing implementation: no %s value is available", e.MissingType)
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