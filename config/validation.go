package config

// ValidationOption is an option that changes the behavior of configuration
// validation.
type ValidationOption func(*validationOptions)

type validationOptions struct{}

func newValidationOptions(options []ValidationOption) validationOptions {
	opts := validationOptions{}

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}

// Normalize returns a normalized copy of v, or an error if v is invalid.
func Normalize[T Normalizable[U], U any](v T, options ...ValidationOption) (T, error) {
	opts := newValidationOptions(options)
	return normalize(opts, v)
}

// Normalizable is an interface that represents some configuration element that
// can be validated and normalized.
type Normalizable[T any] interface {
	normalize(validationOptions) (T, error)
}

func normalize[T Normalizable[U], U any](opts validationOptions, ent T) (T, error) {
	norm, err := ent.normalize(opts)
	return any(norm).(T), err
}
