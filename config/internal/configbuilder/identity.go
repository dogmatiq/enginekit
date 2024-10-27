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

// TargetComponent returns the component that is being configured.
func (b *IdentityBuilder) TargetComponent() config.Component {
	return &b.target
}

// Name sets the name element of the identity.
func (b *IdentityBuilder) Name(name string) {
	b.target.Name = optional.Some(name)
}

// Key sets the key element of the identity.
func (b *IdentityBuilder) Key(key string) {
	b.target.Key = optional.Some(key)
}

// Done completes the configuration of the identity.
func (b *IdentityBuilder) Done() *config.Identity {
	if !b.target.Name.IsPresent() || !b.target.Key.IsPresent() {
		b.target.ComponentFidelity |= config.Incomplete
	}
	return &b.target
}
