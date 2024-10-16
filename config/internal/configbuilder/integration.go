package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Integration returns an [IntegrationBuilder] that builds a new [config.Integration].
func Integration() *IntegrationBuilder {
	return &IntegrationBuilder{}
}

// IntegrationBuilder constructs a [config.Integration].
type IntegrationBuilder struct {
	target   config.Integration
	appendTo *[]config.Handler
}

// SetSourceTypeName sets the source of the configuration.
func (b *IntegrationBuilder) SetSourceTypeName(typeName string) *IntegrationBuilder {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
	return b
}

// SetSource sets the source of the configuration.
func (b *IntegrationBuilder) SetSource(h dogma.IntegrationMessageHandler) *IntegrationBuilder {
	setSource(&b.target.AsConfigured.Source, h)
	return b
}

// AddIdentity returns an [IdentityBuilder] that adds a [config.Identity] to the
// handler.
func (b *IntegrationBuilder) AddIdentity() *IdentityBuilder {
	return &IdentityBuilder{appendTo: &b.target.AsConfigured.Identities}
}

// BuildIdentity calls fn which configures a [config.Identity] that is added to
// the handler.
func (b *IntegrationBuilder) BuildIdentity(
	fn func(*IdentityBuilder) config.Fidelity,
) *IntegrationBuilder {
	x := b.AddIdentity()
	x.Done(fn(x))
	return b
}

// SetDisabled sets whether the handler is disabled or not.
func (b *IntegrationBuilder) SetDisabled(disabled bool) *IntegrationBuilder {
	b.target.AsConfigured.IsDisabled = optional.Some(disabled)
	return b
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *IntegrationBuilder) Edit(fn func(*config.IntegrationAsConfigured)) *IntegrationBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// Done completes the configuration of the handler.
func (b *IntegrationBuilder) Done(f config.Fidelity) *config.Integration {
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
