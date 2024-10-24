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

// SetSourceTypeName sets the source of the configuration.
func (b *IntegrationBuilder) SetSourceTypeName(typeName string) {
	setSourceTypeName(&b.target.AsConfigured.Source, typeName)
}

// SetSource sets the source of the configuration.
func (b *IntegrationBuilder) SetSource(h dogma.IntegrationMessageHandler) {
	setSource(&b.target.AsConfigured.Source, h)
}

// Identity calls fn which configures a [config.Identity] that is added to the
// handler.
func (b *IntegrationBuilder) Identity(fn func(*IdentityBuilder)) {
	x := &IdentityBuilder{}
	fn(x)
	b.target.AsConfigured.Identities = append(
		b.target.AsConfigured.Identities,
		x.Done(),
	)
}

// Route calls fn which configures a [config.Route] that is added to the
// handler.
func (b *IntegrationBuilder) Route(fn func(*RouteBuilder)) {
	x := &RouteBuilder{}
	fn(x)
	b.target.AsConfigured.Routes = append(
		b.target.AsConfigured.Routes,
		x.Done(),
	)
}

// Disable calls fn which configures a [config.Flag] that indicates whether the
// handler is disabled.
func (b *IntegrationBuilder) Disable(fn func(*FlagBuilder[config.Disabled])) {
	x := &FlagBuilder[config.Disabled]{}
	fn(x)
	b.target.AsConfigured.DisabledFlags = append(
		b.target.AsConfigured.DisabledFlags,
		x.Done(),
	)
}

// Edit calls fn, which can apply arbitrary changes to the handler.
func (b *IntegrationBuilder) Edit(fn func(*config.IntegrationAsConfigured)) {
	fn(&b.target.AsConfigured)
}

// Fidelity returns the fidelity of the configuration.
func (b *IntegrationBuilder) Fidelity() config.Fidelity {
	return b.target.AsConfigured.Fidelity
}

// UpdateFidelity merges f with the current fidelity of the configuration.
func (b *IntegrationBuilder) UpdateFidelity(f config.Fidelity) {
	b.target.AsConfigured.Fidelity |= f
}

// Done completes the configuration of the handler.
func (b *IntegrationBuilder) Done() *config.Integration {
	if b.target.AsConfigured.Fidelity&config.Incomplete == 0 {
		if !b.target.AsConfigured.Source.TypeName.IsPresent() {
			panic("handler must have a source or be marked as incomplete")
		}
	}

	return &b.target
}
