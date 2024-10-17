package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/internal/typename"
	"github.com/dogmatiq/enginekit/optional"
)

// Projection returns a new [config.Projection] as configured by fn.
func Projection(fn func(*ProjectionBuilder)) *config.Projection {
	x := &ProjectionBuilder{}
	fn(x)
	return x.Done()
}

// ProjectionBuilder constructs a [config.Projection].
type ProjectionBuilder struct {
	target config.Projection
}

// SetSourceTypeName sets the source of the configuration.
func (b *ProjectionBuilder) SetSourceTypeName(typeName string) {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
}

// SetSource sets the source of the configuration.
func (b *ProjectionBuilder) SetSource(h dogma.ProjectionMessageHandler) {
	setSource(&b.target.AsConfigured.Source, h)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// handler.
func (b *ProjectionBuilder) Identity(fn func(*IdentityBuilder)) {
	x := &IdentityBuilder{}
	fn(x)
	b.target.AsConfigured.Identities = append(
		b.target.AsConfigured.Identities,
		x.Done(),
	)
}

// Route calls fn which configures a [config.Route] that is added to the
// handler.
func (b *ProjectionBuilder) Route(fn func(*RouteBuilder)) {
	x := &RouteBuilder{}
	fn(x)
	b.target.AsConfigured.Routes = append(
		b.target.AsConfigured.Routes,
		x.Done(),
	)
}

// SetDisabled sets whether the handler is disabled or not.
func (b *ProjectionBuilder) SetDisabled(disabled bool) {
	b.target.AsConfigured.IsDisabled = optional.Some(disabled)
}

// SetDeliveryPolicyTypeName sets the type name of the delivery policy.
func (b *ProjectionBuilder) SetDeliveryPolicyTypeName(typeName string) {
	if typeName == "" {
		panic("type name must not be empty")
	}

	b.target.AsConfigured.DeliveryPolicy = optional.Some(
		config.Value[dogma.ProjectionDeliveryPolicy]{
			TypeName: optional.Some(typeName),
		},
	)

}

// SetDeliveryPolicy sets the delivery policy for the handler.
func (b *ProjectionBuilder) SetDeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	if p == nil {
		panic("delivery policy must not be nil")
	}

	b.target.AsConfigured.DeliveryPolicy = optional.Some(
		config.Value[dogma.ProjectionDeliveryPolicy]{
			TypeName: optional.Some(typename.Of(p)),
			Value:    optional.Some(p),
		},
	)

}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *ProjectionBuilder) Edit(fn func(*config.ProjectionAsConfigured)) {
	fn(&b.target.AsConfigured)
}

// UpdateFidelity merges f with the current fidelity of the handler.
func (b *ProjectionBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.AsConfigured.Fidelity |= f
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

	return &b.target
}
