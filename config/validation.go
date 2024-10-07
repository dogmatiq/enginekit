package config

import "errors"

// ValidationOption is an option that changes the behavior of configuration
// validation.
type ValidationOption func(*validationOptions)

type validationOptions struct {
}

func newValidationOptions(options []ValidationOption) validationOptions {
	opts := validationOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

// Validatable is an interface that represents some configuration element that
// can be validated and normalized.
type Validatable[T any] interface {
	validate(validationOptions, *validationResult)
	normalize(validationOptions) T
}

// Validate returns an error if v is invalid.
func Validate[T Validatable[T]](v T, options ...ValidationOption) error {
	opts := newValidationOptions(options)
	res := &validationResult{}

	v.validate(opts, res)

	return res.Err()
}

// Normalize returns a normalized copy of v, or an error if v is invalid.
func Normalize[T Validatable[T]](v T, options ...ValidationOption) (T, error) {
	opts := newValidationOptions(options)
	res := &validationResult{}

	v.validate(opts, res)

	if err := res.Err(); err != nil {
		var zero T
		return zero, err
	}

	return v.normalize(opts), nil
}

type validationResult struct {
	Errors []error
}

func (r *validationResult) appendErr(err error) {
	if err != nil {
		r.Errors = append(r.Errors, err)
	}
}

func (r *validationResult) Err() error {
	return errors.Join(r.Errors...)
}
