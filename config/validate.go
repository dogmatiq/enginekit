package config

import (
	"slices"
	"strings"
)

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

type validationContext struct {
	Options struct {
		Normalize bool
	}

	component Component
	errors    []error
	children  []*validationContext
}

func newValidationContext(c Component, _ []ValidateOption) *validationContext {
	return &validationContext{
		component: c,
	}
}

func (c *validationContext) ValidateChild(child Component) {
	ctx := &validationContext{
		Options:   c.Options,
		component: child,
	}

	c.children = append(c.children, ctx)
	child.validate(ctx)
}

func (c *validationContext) Fail(err error) {
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

// Validate returns an error if the configuration is invalid.
func Validate(c Component, options ...ValidateOption) error {
	ctx := newValidationContext(c, options)
	c.validate(ctx)
	return ctx.error()
}
