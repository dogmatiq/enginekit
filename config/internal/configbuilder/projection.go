package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// Projection returns an [ProjectionBuilder] that builds a new [config.Projection].
func Projection() *ProjectionBuilder {
	return &ProjectionBuilder{}
}

// ProjectionBuilder constructs a [config.Projection].
type ProjectionBuilder struct {
	target   config.Projection
	appendTo *[]config.Handler
}

// SetSourceTypeName sets the source of the configuration.
func (b *ProjectionBuilder) SetSourceTypeName(typeName string) *ProjectionBuilder {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
	return b
}

// SetSource sets the source of the configuration.
func (b *ProjectionBuilder) SetSource(h dogma.ProjectionMessageHandler) *ProjectionBuilder {
	setSource(&b.target.AsConfigured.Source, h)
	return b
}

// AddIdentity returns an [IdentityBuilder] that adds a [config.Identity] to the
// handler.
func (b *ProjectionBuilder) AddIdentity() *IdentityBuilder {
	return &IdentityBuilder{appendTo: &b.target.AsConfigured.Identities}
}

// BuildIdentity calls fn which configures a [config.Identity] that is added to
// the handler.
func (b *ProjectionBuilder) BuildIdentity(fn func(*IdentityBuilder)) *ProjectionBuilder {
	x := b.AddIdentity()
	fn(x)
	x.Done()
	return b
}

// SetDisabled sets whether the handler is disabled or not.
func (b *ProjectionBuilder) SetDisabled(disabled bool) *ProjectionBuilder {
	b.target.AsConfigured.IsDisabled = optional.Some(disabled)
	return b
}

// SetDeliveryPolicyTypeName sets the type name of the delivery policy.
func (b *ProjectionBuilder) SetDeliveryPolicyTypeName(typeName string) *ProjectionBuilder {
	if typeName == "" {
		panic("type name must not be empty")
	}

	b.target.AsConfigured.DeliveryPolicy = optional.Some(
		config.Value[dogma.ProjectionDeliveryPolicy]{
			TypeName: optional.Some(typeName),
		},
	)

	return b
}

// SetDeliveryPolicy sets the delivery policy for the handler.
func (b *ProjectionBuilder) SetDeliveryPolicy(p dogma.ProjectionDeliveryPolicy) *ProjectionBuilder {
	if p == nil {
		panic("delivery policy must not be nil")
	}

	b.target.AsConfigured.DeliveryPolicy = optional.Some(
		config.Value[dogma.ProjectionDeliveryPolicy]{
			TypeName: optional.Some(typename.Of(p)),
			Value:    optional.Some(p),
		},
	)

	return b
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *ProjectionBuilder) Edit(fn func(*config.ProjectionAsConfigured)) *ProjectionBuilder {
	fn(&b.target.AsConfigured)
	return b
}

// UpdateFidelity merges f with the current fidelity of the handler.
func (b *ProjectionBuilder) UpdateFidelity(f config.Fidelity) *ProjectionBuilder {
	b.target.AsConfigured.Fidelity |= f
	return b
}

// Done completes the configuration of the handler.
func (b *ProjectionBuilder) Done() *config.Projection {
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
