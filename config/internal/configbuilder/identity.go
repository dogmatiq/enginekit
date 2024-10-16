package configbuilder

import (
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Identity returns an [IdentityBuilder] that builds a new [config.Identity].
func Identity() *IdentityBuilder {
	return &IdentityBuilder{}
}

// IdentityBuilder constructs a [config.Identity].
type IdentityBuilder struct {
	target   config.Identity
	appendTo *[]*config.Identity
}

// SetName sets the name element of the identity.
func (b *IdentityBuilder) SetName(name string) *IdentityBuilder {
	b.target.AsConfigured.Name = optional.Some(name)
	return b
}

// SetKey sets the key element of the identity.
func (b *IdentityBuilder) SetKey(key string) *IdentityBuilder {
	b.target.AsConfigured.Key = optional.Some(key)
	return b
}

// Edit calls fn, which can apply arbitrary changes to the identity.
func (b *IdentityBuilder) Edit(fn func(*config.IdentityAsConfigured)) *IdentityBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// UpdateFidelity merges f with the current fidelity of the identity.
func (b *IdentityBuilder) UpdateFidelity(f config.Fidelity) *IdentityBuilder {
	b.target.AsConfigured.Fidelity |= f
	return b
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

	if b.appendTo != nil {
		*b.appendTo = append(*b.appendTo, &b.target)
		b.appendTo = nil
	}

	return &b.target
}
