package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/optional"
)

// Aggregate returns an [AggregateBuilder] that builds a new [config.Aggregate].
func Aggregate() *AggregateBuilder {
	return &AggregateBuilder{}
}

// AggregateBuilder constructs a [config.Aggregate].
type AggregateBuilder struct {
	target   config.Aggregate
	appendTo *[]config.Handler
}

// SetSourceTypeName sets the source of the configuration.
func (b *AggregateBuilder) SetSourceTypeName(typeName string) *AggregateBuilder {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
	return b
}

// SetSource sets the source of the configuration.
func (b *AggregateBuilder) SetSource(h dogma.AggregateMessageHandler) *AggregateBuilder {
	setSource(&b.target.AsConfigured.Source, h)
	return b
}

// AddIdentity returns an [IdentityBuilder] that adds a [config.Identity] to the
// handler.
func (b *AggregateBuilder) AddIdentity() *IdentityBuilder {
	return &IdentityBuilder{appendTo: &b.target.AsConfigured.Identities}
}

// BuildIdentity calls fn which configures a [config.Identity] that is added to
// the handler.
func (b *AggregateBuilder) BuildIdentity(fn func(*IdentityBuilder)) *AggregateBuilder {
	x := b.AddIdentity()
	fn(x)
	x.Done()
	return b
}

// SetDisabled sets whether the handler is disabled or not.
func (b *AggregateBuilder) SetDisabled(disabled bool) *AggregateBuilder {
	b.target.AsConfigured.IsDisabled = optional.Some(disabled)
	return b
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *AggregateBuilder) Edit(fn func(*config.AggregateAsConfigured)) *AggregateBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// UpdateFidelity merges f with the current fidelity of the handler.
func (b *AggregateBuilder) UpdateFidelity(f config.Fidelity) *AggregateBuilder {
	b.target.AsConfigured.Fidelity |= f
	return b
}

// Done completes the configuration of the handler.
func (b *AggregateBuilder) Done() *config.Aggregate {
	if b.target.AsConfigured.Fidelity&config.Incomplete == 0 {
		if !b.target.AsConfigured.Source.TypeName.IsPresent() {
			panic("handler must have a source or be marked as incomplete")
		}
		if !b.target.AsConfigured.IsDisabled.IsPresent() {
			panic("handler must be known to be enabled or disabled, or be marked as incomplete")
		}
	}

	if b.appendTo != nil {
		*b.appendTo = append(*b.appendTo, &b.target)
		b.appendTo = nil
	}

	return &b.target
}
