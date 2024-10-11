package config

import "slices"

// Normalize returns a normalized copy of the given component.
//
// It returns a non-nil error if the component is invalid, in which case the
// returned component is normalized as possible in light of the error.
func Normalize[T Component](component T, options ...NormalizeOption) (T, error) {
	ctx := &normalizeContext{
		Component: component,
	}

	for _, opt := range options {
		opt(&ctx.Options)
	}

	norm := component.normalize(ctx).(T)

	return norm, ctx.Err()
}

// MustNormalize returns a normalized copy of v, or panics if v is invalid.
func MustNormalize[T Component](component T, options ...NormalizeOption) T {
	norm, err := Normalize(component, options...)
	if err != nil {
		panic(err)
	}

	return norm
}

func normalize[T Component](ctx *normalizeContext, component T) T {
	return component.normalize(ctx.NewChild(component)).(T)
}

// normalizeContext is the context in which normalization occurs.
type normalizeContext struct {
	Component Component
	Options   normalizeOptions
	Errors    []error

	parent   *normalizeContext
	children []*normalizeContext
}

func newFinalizeContext(c Component) *normalizeContext {
	return &normalizeContext{
		Component: c,
		Options: normalizeOptions{
			PanicOnFailure: true,
		},
	}
}

// normalizeOptions is the result of applying a set of [NormalizeOption] values.
type normalizeOptions struct {
	PanicOnFailure         bool
	RequireImplementations bool
}

func (c *normalizeContext) NewChild(com Component) *normalizeContext {
	ctx := &normalizeContext{
		Component: com,
		Options:   c.Options,

		parent: c,
	}

	c.children = append(c.children, ctx)

	return ctx
}

func (c *normalizeContext) Fail(err error) {
	if c.Options.PanicOnFailure {
		for ctx := c; ctx != nil; ctx = ctx.parent {
			err = ComponentError{
				Component: ctx.Component,
				Causes:    []error{err},
			}
		}
		panic(err)
	}

	if err != nil {
		c.Errors = append(c.Errors, err)
	}
}

func (c *normalizeContext) Err() error {
	errors := slices.Clone(c.Errors)

	for _, child := range c.children {
		if err := child.Err(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return ComponentError{
		Component: c.Component,
		Causes:    errors,
	}
}
