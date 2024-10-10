package config

// ValidationOption is an option that changes the behavior of configuration
// validation.
type ValidationOption func(*validationOptions)

// WithImplementations is a [ValidationOption] that requires all application,
// handler and message implementations to be available in order to consider the
// configuration valid.
func WithImplementations() ValidationOption {
	return func(o *validationOptions) {
		o.RequireImplementations = true
	}
}

type validationOptions struct {
	RequireImplementations bool
}

func newValidationOptions(options []ValidationOption) validationOptions {
	opts := validationOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

// Normalize returns a normalized copy of the given component.
//
// It returns a non-nil error if the component is invalid, in which case the
// returned component is normalized as possible in light of the error.
func Normalize[T Component](component T, options ...ValidationOption) (T, error) {
	ctx := &normalizationContext{
		Component: component,
		Options:   newValidationOptions(options),
	}

	norm := component.normalize(ctx).(T)

	return norm, ctx.Err()
}

// MustNormalize returns a normalized copy of v, or panics if v is invalid.
func MustNormalize[T Component](component T, options ...ValidationOption) T {
	norm, err := Normalize(component, options...)
	if err != nil {
		panic(err)
	}

	return norm
}

func normalize[T Component](ctx *normalizationContext, component T) T {
	child := &normalizationContext{
		Component: component,
		Options:   ctx.Options,
	}

	norm := component.normalize(child).(T)
	ctx.Fail(child.Err())

	return norm
}

type normalizationContext struct {
	Component Component
	Options   validationOptions
	Errors    []error
}

func (c *normalizationContext) Fail(err error) {
	if err != nil {
		c.Errors = append(c.Errors, err)
	}
}

func (c *normalizationContext) Err() error {
	if len(c.Errors) == 0 {
		return nil
	}

	return ComponentError{
		Component: c.Component,
		Causes:    c.Errors,
	}
}
