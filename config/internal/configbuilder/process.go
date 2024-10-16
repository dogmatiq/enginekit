package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Process returns an [ProcessBuilder] that builds a new [config.Process].
func Process() *ProcessBuilder {
	return &ProcessBuilder{}
}

// ProcessBuilder constructs a [config.Process].
type ProcessBuilder struct {
	target   config.Process
	appendTo *[]config.Handler
}

// SetSourceTypeName sets the source of the configuration.
func (b *ProcessBuilder) SetSourceTypeName(typeName string) *ProcessBuilder {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
	return b
}

// SetSource sets the source of the configuration.
func (b *ProcessBuilder) SetSource(h dogma.ProcessMessageHandler) *ProcessBuilder {
	setSource(&b.target.AsConfigured.Source, h)
	return b
}

// AddIdentity returns an [IdentityBuilder] that adds a [config.Identity] to the
// handler.
func (b *ProcessBuilder) AddIdentity() *IdentityBuilder {
	return &IdentityBuilder{appendTo: &b.target.AsConfigured.Identities}
}

// BuildIdentity calls fn which configures a [config.Identity] that is added to
// the handler.
func (b *ProcessBuilder) BuildIdentity(
	fn func(*IdentityBuilder) config.Fidelity,
) *ProcessBuilder {
	x := b.AddIdentity()
	x.Done(fn(x))
	return b
}

// SetDisabled sets whether the handler is disabled or not.
func (b *ProcessBuilder) SetDisabled(disabled bool) *ProcessBuilder {
	b.target.AsConfigured.IsDisabled = optional.Some(disabled)
	return b
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *ProcessBuilder) Edit(fn func(*config.ProcessAsConfigured)) *ProcessBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// Done completes the configuration of the handler.
func (b *ProcessBuilder) Done(f config.Fidelity) *config.Process {
	if f&config.Incomplete == 0 {
		if !b.target.AsConfigured.Source.TypeName.IsPresent() {
			panic("handler must have a source or be marked as incomplete")
		}
		if !b.target.AsConfigured.IsDisabled.IsPresent() {
			panic("handler must be known to be enabled or disabled, or be marked as incomplete")
		}
	}

	b.target.AsConfigured.Fidelity = f

	if b.appendTo != nil {
		*b.appendTo = append(*b.appendTo, &b.target)
		b.appendTo = nil
	}

	return &b.target
}
