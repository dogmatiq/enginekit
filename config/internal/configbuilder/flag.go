package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Flag returns a new [config.FlagModification] as configured by fn.
func Flag(fn func(*FlagBuilder)) *config.FlagModification {
	x := &FlagBuilder{}
	fn(x)
	return x.Done()
}

// FlagBuilder constructs a [config.FlagModification].
type FlagBuilder struct {
	target config.FlagModification
}

// Value sets the value of the target [config.FlagModification].
func (b *FlagBuilder) Value(v bool) {
	b.target.Value = optional.Some(v)
}

// Done completes the configuration of the flag.
func (b *FlagBuilder) Done() *config.FlagModification {
	if !b.target.Value.IsPresent() {
		b.target.Fidelity |= config.Incomplete
	}
	return &b.target
}
