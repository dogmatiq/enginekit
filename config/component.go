package config

import (
	"fmt"
	"strings"

	"github.com/dogmatiq/enginekit/optional"
)

// A Component is some element of the configuration of a Dogma application.
type Component interface {
	fmt.Stringer

	normalize(*normalizationContext) Component
}

// Implementation contains information about the type and value that implements
// a [Component] of type T.
type Implementation[T any] struct {
	// TypeName is the fully-qualified name of the Go type that implements T.
	TypeName string

	// Source is the value that produced the configuration, if available.
	Source optional.Optional[T]
}

// ComponentError indicates that a [Component] is invalid.
type ComponentError struct {
	Component Component
	Causes    []error
}

func (e ComponentError) Error() string {
	var w strings.Builder

	w.WriteString(e.Component.String())
	w.WriteString(" is invalid")

	if len(e.Causes) == 1 {
		w.WriteString(": ")
		w.WriteString(e.Causes[0].Error())
	} else if len(e.Causes) > 1 {
		w.WriteString(":")

		for _, cause := range e.Causes {
			lines := strings.Split(cause.Error(), "\n")

			for i, line := range lines {
				w.WriteByte('\n')
				if i == 0 {
					w.WriteString("- ")
				} else {
					w.WriteString("  ")
				}
				w.WriteString(line)
			}
		}
	}

	return w.String()
}

func (e ComponentError) Unwrap() []error {
	return e.Causes
}

// MissingTypeNameError indicates that a component that refers to a Go type is
// missing the type name.
type MissingTypeNameError struct{}

func (e MissingTypeNameError) Error() string {
	return "missing type name"
}

// MissingImplementationError indicates that a component is invalid because it
// does not contain an implementation and the [WithImplementations] option was
// specified during normalization.
type MissingImplementationError struct{}

func (e MissingImplementationError) Error() string {
	return "missing implementation"
}

// ImplementationTypeNameMismatchError indicates that the type name of an
// implementation does not match the type name of the source value.
type ImplementationTypeNameMismatchError struct {
	ExpectedTypeName   string
	UnexpectedTypeName string
}

func (e ImplementationTypeNameMismatchError) Error() string {
	return fmt.Sprintf(
		"type name mismatch: implementation is %q, but the type was reported as %q",
		e.ExpectedTypeName,
		e.UnexpectedTypeName,
	)
}
