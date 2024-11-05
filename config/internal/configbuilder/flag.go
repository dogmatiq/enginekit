package configbuilder

import (
	"fmt"

	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Flag returns a new [config.FlagModification] as configured by fn.
func Flag[S config.Symbol](fn func(*FlagBuilder[S])) *config.Flag[S] {
	x := &FlagBuilder[S]{}
	fn(x)
	return x.Done()
}

// FlagBuilder constructs a [config.FlagModification].
type FlagBuilder[S config.Symbol] struct {
	target config.Flag[S]
}

// Value sets the value of the target [config.FlagModification].
func (b *FlagBuilder[S]) Value(v bool) {
	b.target.Value = optional.Some(v)
}

// Partial marks the compomnent as partially configured.
func (b *FlagBuilder[S]) Partial(format string, args ...any) {
	b.target.IsPartialReasons = append(b.target.IsPartialReasons, fmt.Sprintf(format, args...))
}

// Speculative marks the component as speculative.
func (b *FlagBuilder[S]) Speculative() {
	b.target.IsSpeculative = true
}

// Done completes the configuration of the flag.
func (b *FlagBuilder[S]) Done() *config.Flag[S] {
	return &b.target
}
