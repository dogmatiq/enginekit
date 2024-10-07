package config

// ValidationOption is an option that changes the behavior of configuration
// validation.
type ValidationOption func(*validationOptions)

type validationOptions struct {
}

func newValidationOptions(options []ValidationOption) validationOptions {
	var opts validationOptions

	for _, opt := range options {
		opt(&opts)
	}

	return opts
}
