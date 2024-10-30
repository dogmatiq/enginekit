package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
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

// TypeName sets the name of the concrete type that implements
// [dogma.ProjectionMessageHandler].
func (b *ProjectionBuilder) TypeName(n string) {
	setTypeName(&b.target.TypeName, &b.target.Source, n)
}

// Source sets the source value to h.
func (b *ProjectionBuilder) Source(h dogma.ProjectionMessageHandler) {
	setSource(&b.target.TypeName, &b.target.Source, h)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// handler.
func (b *ProjectionBuilder) Identity(fn func(*IdentityBuilder)) {
	b.target.IdentityComponents = append(b.target.IdentityComponents, Identity(fn))
}

// Route calls fn which configures a [config.Route] that is added to the
// handler.
func (b *ProjectionBuilder) Route(fn func(*RouteBuilder)) {
	b.target.RouteComponents = append(b.target.RouteComponents, Route(fn))
}

// Disabled calls fn which configures a [config.FlagModification] that is added
// to the handler's disabled flag.
func (b *ProjectionBuilder) Disabled(fn func(*FlagBuilder)) {
	b.target.DisabledFlag.Modifications = append(b.target.DisabledFlag.Modifications, Flag(fn))
}

// DeliveryPolicy calls fn which configures a [config.ProjectionDeliveryPolicy]
// that is added to the handler.
func (b *ProjectionBuilder) DeliveryPolicy(fn func(*ProjectionDeliveryPolicyBuilder)) {
	b.target.DeliveryPolicyComponents = append(b.target.DeliveryPolicyComponents, ProjectionDeliveryPolicy(fn))
}

// Done completes the configuration of the handler.
func (b *ProjectionBuilder) Done() *config.Projection {
	if !b.target.TypeName.IsPresent() {
		b.target.ComponentFidelity |= config.Incomplete
	}
	return &b.target
}

// ProjectionDeliveryPolicy returns a new [dogma.ProjectionDeliveryPolicy] as
// configured by fn.
func ProjectionDeliveryPolicy(fn func(*ProjectionDeliveryPolicyBuilder)) *config.ProjectionDeliveryPolicy {
	x := &ProjectionDeliveryPolicyBuilder{}
	fn(x)
	return x.Done()
}

// ProjectionDeliveryPolicyBuilder constructs a
// [config.ProjectionDeliveryPolicy].
type ProjectionDeliveryPolicyBuilder struct {
	target config.ProjectionDeliveryPolicy
}

// AsPerDeliveryPolicy configures the builder to use the same properties as p.
func (b *ProjectionDeliveryPolicyBuilder) AsPerDeliveryPolicy(p dogma.ProjectionDeliveryPolicy) {
	switch p := p.(type) {
	case dogma.UnicastProjectionDeliveryPolicy:
		b.target.DeliveryPolicyType = optional.Some(config.UnicastProjectionDeliveryPolicyType)
	case dogma.BroadcastProjectionDeliveryPolicy:
		b.target.DeliveryPolicyType = optional.Some(config.BroadcastProjectionDeliveryPolicyType)
		b.target.BroadcastToPrimaryFirst = optional.Some(p.PrimaryFirst)
	default:
		b.target.ComponentFidelity |= config.Incomplete
	}
}

// Done completes the configuration of the policy.
func (b *ProjectionDeliveryPolicyBuilder) Done() *config.ProjectionDeliveryPolicy {
	if t, ok := b.target.DeliveryPolicyType.TryGet(); !ok {
		b.target.ComponentFidelity |= config.Incomplete
	} else if t == config.BroadcastProjectionDeliveryPolicyType && !b.target.BroadcastToPrimaryFirst.IsPresent() {
		b.target.ComponentFidelity |= config.Incomplete
	}
	return &b.target
}
