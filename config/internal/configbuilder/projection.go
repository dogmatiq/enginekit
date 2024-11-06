package configbuilder

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
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
func (b *ProjectionBuilder) Disabled(fn func(*FlagBuilder[config.Disabled])) {
	b.target.DisabledFlags = append(b.target.DisabledFlags, Flag(fn))
}

// DeliveryPolicy calls fn which configures a [config.ProjectionDeliveryPolicy]
// that is added to the handler.
func (b *ProjectionBuilder) DeliveryPolicy(fn func(*ProjectionDeliveryPolicyBuilder)) {
	b.target.DeliveryPolicyComponents = append(b.target.DeliveryPolicyComponents, ProjectionDeliveryPolicy(fn))
}

// Partial marks the compomnent as partially configured.
func (b *ProjectionBuilder) Partial(format string, args ...any) {
	b.target.IsPartialReasons = append(b.target.IsPartialReasons, fmt.Sprintf(format, args...))
}

// Speculative marks the component as speculative.
func (b *ProjectionBuilder) Speculative() {
	b.target.IsSpeculative = true
}

// Done completes the configuration of the handler.
func (b *ProjectionBuilder) Done() *config.Projection {
	return &b.target
}
