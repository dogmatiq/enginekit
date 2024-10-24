package configbuilder

import "github.com/dogmatiq/enginekit/config"

// FlagBuilder constructs a [config.Flag].
type FlagBuilder[L config.Label] struct {
	target config.Flag[L]
}

// Fidelity returns the fidelity of the configuration.
func (b *FlagBuilder[L]) Fidelity() config.Fidelity {
	return b.target.Fidelity()
}

// UpdateFidelity merges f into the fidelity of the configuration.
func (b *FlagBuilder[L]) UpdateFidelity(f config.Fidelity) {
	if f == config.Speculative {
		b.target.IsSpeculative = true
	}

	if f != config.Immaculate {
		panic("unsupported fidelity")
	}
}

// Done completes the configuration of the flag.
func (b *FlagBuilder[L]) Done() *config.Flag[L] {
	return &b.target
}
