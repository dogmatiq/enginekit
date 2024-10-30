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
		Component: c,
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
	Panic     bool
}

// validateContext carries the inputs and outputs of the component validation
// process.
type validateContext struct {
	Component Component
	Options   validationOptions

	errors   []error
	parent   *validateContext
	children []*validateContext
}

func newResolutionContext(c Component) *validateContext {
	return &validateContext{
		Component: c,
		Options: validationOptions{
			Panic: true,
		},
	}
}

func (c *validateContext) ForChild(child Component) *validateContext {
	ctx := &validateContext{
		Component: child,
		Options:   c.Options,
		parent:    c,
	}
	c.children = append(c.children, ctx)

	return ctx
}

func (c *validateContext) ValidateChild(child Component) {
	child.validate(c.ForChild(child))
}

func (c *validateContext) Invalid(err error) {
	if !c.Options.Panic {
		c.errors = append(c.errors, err)
		return
	}

	if c.parent == nil {
		panic(err)
	}

	err = InvalidComponentError{
		Component: c.Component,
		Causes:    []error{err},
	}

	c.parent.Invalid(err)
}

func (c *validateContext) Malformed(format string, args ...any) {
	panic(fmt.Sprintf(
		"malformed configuration representation: there is a problem with the %T value, not necessarily the configuration it represents %s",
		c.Component,
		fmt.Sprintf(format, args...),
	))
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
		Component: c.Component,
		Causes:    errors,
	}
}
