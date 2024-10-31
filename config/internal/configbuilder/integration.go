package configbuilder

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config"
)

// Integration returns a new [config.Integration] as configured by fn.
func Integration(fn func(*IntegrationBuilder)) *config.Integration {
	x := &IntegrationBuilder{}
	fn(x)
	return x.Done()
}

// IntegrationBuilder constructs a [config.Integration].
type IntegrationBuilder struct {
	target config.Integration
}

// TypeName sets the name of the concrete type that implements
// [dogma.IntegrationMessageHandler].
func (b *IntegrationBuilder) TypeName(n string) {
	setTypeName(&b.target.TypeName, &b.target.Source, n)
}

// Source sets the source value to h.
func (b *IntegrationBuilder) Source(h dogma.IntegrationMessageHandler) {
	setSource(&b.target.TypeName, &b.target.Source, h)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// handler.
func (b *IntegrationBuilder) Identity(fn func(*IdentityBuilder)) {
	b.target.IdentityComponents = append(b.target.IdentityComponents, Identity(fn))
}

// Route calls fn which configures a [config.Route] that is added to the
// handler.
func (b *IntegrationBuilder) Route(fn func(*RouteBuilder)) {
	b.target.RouteComponents = append(b.target.RouteComponents, Route(fn))
}

// Disabled calls fn which configures a [config.FlagModification] that is added
// to the handler's disabled flag.
func (b *IntegrationBuilder) Disabled(fn func(*FlagBuilder)) {
	b.target.DisabledFlag.Modifications = append(b.target.DisabledFlag.Modifications, Flag(fn))
}

// Done completes the configuration of the handler.
func (b *IntegrationBuilder) Done() *config.Integration {
	if !b.target.TypeName.IsPresent() {
		b.target.Fidelity |= config.Incomplete
	}
	return &b.target
}
