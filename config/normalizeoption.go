package config

// NormalizeOption is an option that changes the behavior of [Normalize] and
// [MustNormalize].
type NormalizeOption func(*normalizeOptions)

// WithImplementations is a [NormalizeOption] that requires all application,
// handler and message implementations to be available in order to consider the
// configuration valid.
func WithImplementations() NormalizeOption {
	return func(o *normalizeOptions) {
		o.RequireImplementations = true
	}
}
