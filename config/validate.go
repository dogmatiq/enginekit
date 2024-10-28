package config

import (
	"slices"
	"strings"
)

// Validate returns an error if the configuration is invalid.
func Validate(c Component, _ ...ValidateOption) error {
	ctx := &validationContext{
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

func (e InvalidComponentError) Unwrap() []error {
	return e.Causes
}

type validationOptions struct {
	Normalize bool
}

// validationContext carries the inputs and outputs of the component validation
// process.
//
// A nil pointer to a validationContext is a valid context that behaves in
// "strict mode", where errors are surfaced via a panic the moment they occur.
type validationContext struct {
	options   validationOptions
	component Component
	errors    []error
	children  []*validationContext
}

func (c *validationContext) ValidateChild(child Component) {
	var ctx *validationContext

	if c != nil {
		ctx = &validationContext{
			options:   c.options,
			component: child,
		}
		c.children = append(c.children, ctx)
	}

	child.validate(ctx)
}

func (c *validationContext) Options() validationOptions {
	if c != nil {
		return c.options
	}
	return validationOptions{}
}

func (c *validationContext) Fail(err error) {
	if c == nil {
		panic(err)
	}
	c.errors = append(c.errors, err)
}

func (c *validationContext) error() error {
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
