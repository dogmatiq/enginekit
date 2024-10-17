package config

import (
	"slices"
)

// Normalize returns a normalized copy of the given component.
//
// It returns a non-nil error if the component is invalid, in which case the
// returned component is normalized as possible in light of the error.
func Normalize[T Component](c T, options ...NormalizeOption) (T, error) {
	ctx := &normalizationContext{
		Component: c,
	}

	for _, opt := range options {
		opt(&ctx.Options)
	}

	c = c.clone().(T)
	c.normalize(ctx)
	reportFidelityErrors(ctx, c)

	return c, ctx.Err()
}

// MustNormalize returns a normalized copy of v, or panics if v is invalid.
func MustNormalize[T Component](c T, options ...NormalizeOption) T {
	norm, err := Normalize(c, options...)
	if err != nil {
		panic(err)
	}

	return norm
}

func normalizeChildren[T Component](ctx *normalizationContext, components []T) {
	if ctx.Options.Shallow {
		return
	}

	for _, c := range components {
		childCtx := ctx.NewChild(c)
		c.normalize(childCtx)
		reportFidelityErrors(childCtx, c)
	}
}

func validate[T Component](c T) (Fidelity, []error) {
	ctx := shallowContext(c)
	c = c.clone().(T)
	c.normalize(ctx)
	return c.Fidelity(), ctx.Errors
}

func reportFidelityErrors(ctx *normalizationContext, c Component) {
	f := c.Fidelity()

	if f&Speculative != 0 {
		ctx.Fail(SpeculativeError{})
	}

	if f&Incomplete != 0 {
		ctx.Fail(IncompleteError{})
	}
}

// normalizationContext is the context in which normalization occurs.
type normalizationContext struct {
	Component Component
	Options   normalizationOptions
	Errors    []error

	parent   *normalizationContext
	children []*normalizationContext
}

func strictContext(c Component) *normalizationContext {
	return &normalizationContext{
		Component: c,
		Options: normalizationOptions{
			PanicOnFailure: true,
		},
	}
}

func shallowContext(c Component) *normalizationContext {
	return &normalizationContext{
		Component: c,
		Options: normalizationOptions{
			Shallow: true,
		},
	}
}

// normalizationOptions is the result of applying a set of [NormalizeOption]
// values.
type normalizationOptions struct {
	PanicOnFailure bool
	RequireValues  bool
	Shallow        bool
}

func (c *normalizationContext) NewChild(com Component) *normalizationContext {
	if c.Options.Shallow {
		panic("did not expect to descend into subcomponents")
	}

	ctx := &normalizationContext{
		Component: com,
		Options:   c.Options,

		parent: c,
	}

	c.children = append(c.children, ctx)

	return ctx
}

func (c *normalizationContext) Fail(err error) {
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

func (c *normalizationContext) Err() error {
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
