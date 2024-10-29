package config

import (
	"errors"
	"fmt"
	"iter"
	"slices"
	"strings"
)

// Validate returns an error if the configuration is invalid.
func Validate(c Component, _ ...ValidateOption) error {
	ctx := &validateContext{
		component: c,
	}

	c.validate(ctx)

	return ctx.error()
}

// ValidateOption changes the behavior of [Component.Validate].
type ValidateOption func(*validationOptions)

// Normalize is a [ValidateOption] that indicates causes the component to be
// normalized in-place during validation.
func Normalize() ValidateOption {
	return func(o *validationOptions) {
		o.Normalize = true
	}
}

// InvalidComponentError indicates that a [Component] is invalid.
type InvalidComponentError struct {
	Component Component
	Causes    []error
}

func (e InvalidComponentError) Error() string {
	var w strings.Builder

	fmt.Fprint(&w, e.Component)
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

func (e InvalidComponentError) Unwrap() []error {
	return e.Causes
}

// ErrorsByComponent returns the errors within err that are directly
// associated with the given component.
func ErrorsByComponent(c Component, err error) []error {
	var matches []error

	cerr, ok := err.(InvalidComponentError)
	if !ok {
		return nil
	}

	for _, err := range cerr.Unwrap() {
		if nested, ok := err.(InvalidComponentError); ok {
			matches = append(
				matches,
				ErrorsByComponent(c, nested)...,
			)
		} else if cerr.Component == c {
			matches = append(matches, err)
		}
	}

	return matches
}

func unwrap(err error) iter.Seq[error] {
	type many interface {
		Unwrap() []error
	}

	return func(yield func(error) bool) {
		if err := errors.Unwrap(err); err != nil {
			yield(err)
		} else if err, ok := err.(many); ok {
			for _, e := range err.Unwrap() {
				if !yield(e) {
					return
				}
			}
		}
	}
}

type validationOptions struct {
	Normalize bool
}

// validateContext carries the inputs and outputs of the component validation
// process.
//
// A nil pointer to a validateContext is a valid context that behaves in
// "strict mode", where errors are surfaced via a panic the moment they occur.
type validateContext struct {
	options   validationOptions
	component Component
	errors    []error
	children  []*validateContext
}

func (c *validateContext) ValidateChild(child Component) {
	var ctx *validateContext

	if c != nil {
		ctx = &validateContext{
			options:   c.options,
			component: child,
		}
		c.children = append(c.children, ctx)
	}

	child.validate(ctx)
}

func (c *validateContext) Options() validationOptions {
	if c != nil {
		return c.options
	}
	return validationOptions{}
}

func (c *validateContext) Fail(err error) {
	if c == nil {
		panic(err)
	}
	c.errors = append(c.errors, err)
}

func (c *validateContext) error() error {
	errors := slices.Clone(c.errors)

	for _, child := range c.children {
		if err := child.error(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return InvalidComponentError{
		Component: c.component,
		Causes:    errors,
	}
}
