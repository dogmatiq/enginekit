package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Identity returns a new [config.Identity] as configured by fn.
func Identity(fn func(*IdentityBuilder)) *config.Identity {
	x := &IdentityBuilder{}
	fn(x)
	return x.Done()
}

// IdentityBuilder constructs a [config.Identity].
type IdentityBuilder struct {
	target config.Identity
}

// SetName sets the name element of the identity.
func (b *IdentityBuilder) SetName(name string) {
	b.target.AsConfigured.Name = optional.Some(name)
}

// SetKey sets the key element of the identity.
func (b *IdentityBuilder) SetKey(key string) {
	b.target.AsConfigured.Key = optional.Some(key)
}

// Edit calls fn, which can apply arbitrary changes to the identity.
func (b *IdentityBuilder) Edit(fn func(*config.IdentityAsConfigured)) {
	fn(&b.target.AsConfigured)
}

// Fidelity returns the fidelity of the configuration.
func (b *IdentityBuilder) Fidelity() config.Fidelity {
	return b.target.AsConfigured.Fidelity
}

// UpdateFidelity merges f with the current fidelity of the configuration.
func (b *IdentityBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.AsConfigured.Fidelity |= f
}

// Done completes the configuration of the identity.
func (b *IdentityBuilder) Done() *config.Identity {
	if b.target.AsConfigured.Fidelity&config.Incomplete == 0 {
		if !b.target.AsConfigured.Name.IsPresent() {
			panic("identity must have a name or be marked as incomplete")
		}
		if !b.target.AsConfigured.Key.IsPresent() {
			panic("identity must have a key or be marked as incomplete")
		}
	}

	return &b.target
}
