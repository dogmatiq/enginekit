package config

// NormalizeOption is an option that changes the behavior of [Normalize] and
// [MustNormalize].
type NormalizeOption func(*normalizationOptions)

// WithRuntimeTypes is a [NormalizeOption] that requires all application,
// handler and message implementations to be available in order to consider the
// configuration valid.
func WithRuntimeTypes() NormalizeOption {
	return func(o *normalizationOptions) {
		o.RequireValues = true
	}
}
